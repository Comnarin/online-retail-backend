package database

import (
	"github.com/retail/backend/internal/domain"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Set to Info to see migration logs
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	log.Info().Msg("PostgreSQL connected")

	// Run AutoMigrate
	err = AutoMigrate(db)
	if err != nil {
		log.Error().Err(err).Msg("Database migration failed")
		return nil, err
	}

	// Seed data
	if err := Seed(db); err != nil {
		log.Error().Err(err).Msg("Database seeding failed")
	}

	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	log.Info().Msg("Running database migrations...")

	// Cleanup for relational migration and deadlock resolution
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// 1. Consolidate stock/inventory columns
	log.Info().Msg("Consolidating stock/inventory columns...")
	if err := db.Exec(`
		DO $$ 
		BEGIN 
			-- If stock exists, move data to inventory and drop stock
			IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='products' AND column_name='stock') THEN
				-- Update inventory with stock values where needed
				UPDATE products SET inventory = stock WHERE inventory = 0;
				-- Drop the old column
				ALTER TABLE products DROP COLUMN stock;
			END IF;
		END $$;
	`).Error; err != nil {
		log.Error().Err(err).Msg("Stock consolidation failed")
	}

	// DE-DUPLICATION: Keep only the most recent cart per customer before adding UNIQUE index
	log.Info().Msg("Cleaning up duplicate carts (Cascading)...")
	
	// 1. Delete cart items for non-latest active carts
	if err := db.Exec(`
		DELETE FROM cart_items 
		WHERE cart_id IN (
			SELECT id FROM carts 
			WHERE id NOT IN (
				SELECT id FROM (
					SELECT DISTINCT ON (tenant_id, customer_id) id 
					FROM carts 
					WHERE deleted_at IS NULL
					ORDER BY tenant_id, customer_id, created_at DESC
				) as t
			)
		)
	`).Error; err != nil {
		log.Warn().Err(err).Msg("Cart items cleanup failed (may be empty)")
	}

	// 2. Delete the duplicate active carts
	if err := db.Exec(`
		DELETE FROM carts 
		WHERE id NOT IN (
			SELECT id FROM (
				SELECT DISTINCT ON (tenant_id, customer_id) id 
				FROM carts 
				WHERE deleted_at IS NULL
				ORDER BY tenant_id, customer_id, created_at DESC
			) as t
		)
	`).Error; err != nil {
		log.Error().Err(err).Msg("Carts cleanup failed")
	}

	err := db.AutoMigrate(
		&domain.Tenant{},
		&domain.TenantFeature{},
		&domain.TenantAppearance{},
		&domain.User{},
		&domain.Customer{},
		&domain.ProductCategory{},
		&domain.Product{},
		&domain.ProductImage{},
		&domain.Order{},
		&domain.OrderAddress{},
		&domain.OrderItem{},
		&domain.Coupon{},
		&domain.MembershipTier{},
		&domain.TierBenefit{},
		&domain.CustomerMembership{},
		&domain.PointTransaction{},
		&domain.Transaction{},
		&domain.CustomerAddress{},
		&domain.PaymentMethod{},
		&domain.CustomerPaymentMethod{},
		&domain.Cart{},
		&domain.CartItem{},
	)
	if err != nil {
		return err
	}

	if err := MigrateUsersToCustomers(db); err != nil {
		return err
	}

	// 2. Re-run AutoMigrate for customer-related tables to establish new foreign keys
	// Now that data is moved and old keys are dropped, GORM can safely point them to 'customers'
	return db.AutoMigrate(
		&domain.Order{},
		&domain.CustomerMembership{},
		&domain.PointTransaction{},
		&domain.Transaction{},
		&domain.CustomerAddress{},
		&domain.CustomerPaymentMethod{},
	)
}

func MigrateUsersToCustomers(db *gorm.DB) error {
	// Check if we have users with 'customer' role
	var count int64
	db.Table("users").Where("role = ?", "customer").Count(&count)
	if count == 0 {
		return nil
	}

	log.Info().Int64("count", count).Msg("Migrating legacy customers to dedicated table...")

	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Drop known legacy constraints that point to 'users' table
		// This is necessary because these tables now point to 'customers' but might have old FKs
		constraints := []struct {
			table string
			name  string
		}{
			{"orders", "fk_orders_customer"},
			{"customer_memberships", "fk_customer_memberships_customer"},
			{"point_transactions", "fk_point_transactions_customer"},
			{"transactions", "fk_transactions_customer"},
			{"customer_addresses", "fk_customer_addresses_customer"},
			{"customer_payment_methods", "fk_customer_payment_methods_customer"},
		}

		for _, c := range constraints {
			// DROP CONSTRAINT IF EXISTS is idempotent and safe
			tx.Exec("ALTER TABLE " + c.table + " DROP CONSTRAINT IF EXISTS " + c.name)
		}

		// 2. Copy records to customers table
		// Note: We use raw SQL because the domain.User struct no longer has LineUserID/Phone/Avatar
		migrateQuery := `
			INSERT INTO customers (id, tenant_id, line_user_id, name, phone, avatar, is_active, created_at, updated_at)
			SELECT id, tenant_id, line_user_id, name, phone, avatar, is_active, created_at, updated_at 
			FROM users 
			WHERE role = 'customer'
			ON CONFLICT (id) DO NOTHING
		`
		if err := tx.Exec(migrateQuery).Error; err != nil {
			return err
		}

		// 3. Remove migrated records from users table
		if err := tx.Exec("DELETE FROM users WHERE role = 'customer'").Error; err != nil {
			return err
		}

		log.Info().Msg("Customer migration completed successfully")
		return nil
	})
}

func Seed(db *gorm.DB) error {
	var count int64
	email := "superadmin@retail.com"
	password := "password"

	db.Model(&domain.User{}).Where("email = ?", email).Count(&count)
	if count == 0 {
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		hashedStr := string(hashed)

		user := &domain.User{
			Email:    &email,
			Password: &hashedStr,
			Name:     "Super Admin",
			Role:     domain.RoleSuperAdmin,
			IsActive: true,
		}

		if err := db.Create(user).Error; err != nil {
			return err
		}
		log.Info().Msg("Superadmin user seeded successfully")
	}

	// Seed Payment Methods for all tenants if none exist
	var tenants []domain.Tenant
	if err := db.Find(&tenants).Error; err == nil {
		for _, tenant := range tenants {
			var pmCount int64
			db.Model(&domain.PaymentMethod{}).Where("tenant_id = ?", tenant.ID).Count(&pmCount)
			if pmCount == 0 {
				defaults := []domain.PaymentMethod{
					{
						TenantID:      tenant.ID,
						Name:          "PromptPay QR",
						Code:          "promptpay",
						Description:   "Scan QR to pay instantly",
						Icon:          "Zap",
						IsActive:      true,
						ExpiryMinutes: 5,
					},
					{
						TenantID:      tenant.ID,
						Name:          "Credit Card",
						Code:          "credit_card",
						Description:   "Secure payment via Omise/Stripe",
						Icon:          "CreditCard",
						IsActive:      true,
						ExpiryMinutes: 10,
					},
					{
						TenantID:      tenant.ID,
						Name:          "Cash on Delivery",
						Code:          "cod",
						Description:   "Pay when you receive the items",
						Icon:          "ShoppingBag",
						IsActive:      true,
						ExpiryMinutes: 1440, // 1 day
					},
				}
				for _, pm := range defaults {
					db.Create(&pm)
				}
				log.Info().Str("tenant", tenant.Name).Msg("Default payment methods seeded")
			}
		}
	}

	return nil
}
