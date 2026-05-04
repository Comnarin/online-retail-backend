package main

import (
	"log"

	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/platform/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log.Println("Starting bilingual data migration...")

	// 1. Load config to get Database URL
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Connect to database
	db, err := gorm.Open(postgres.Open(cfg.Database.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 3. Ensure AutoMigrate creates name_th and name_en
	log.Println("Running AutoMigrate to ensure columns exist...")
	err = db.AutoMigrate(
		&domain.ProductCategory{},
		&domain.Product{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// 4. Run Idempotent SQL Migrations within a Transaction
	log.Println("Running SQL data migration...")
	err = db.Transaction(func(tx *gorm.DB) error {
		// Products name migration
		res := tx.Exec("UPDATE products SET name_th = name WHERE name_th IS NULL OR name_th = '';")
		if res.Error != nil {
			return res.Error
		}
		log.Printf("Migrated name to name_th in products. Rows affected: %d", res.RowsAffected)

		// Products description migration
		res = tx.Exec("UPDATE products SET description_th = description WHERE description_th IS NULL OR description_th = '';")
		if res.Error != nil {
			return res.Error
		}
		log.Printf("Migrated description to description_th in products. Rows affected: %d", res.RowsAffected)

		// Categories name migration
		res = tx.Exec("UPDATE product_categories SET name_th = name WHERE name_th IS NULL OR name_th = '';")
		if res.Error != nil {
			return res.Error
		}
		log.Printf("Migrated name to name_th in product_categories. Rows affected: %d", res.RowsAffected)

		return nil
	})

	if err != nil {
		log.Fatalf("Migration failed during SQL transaction: %v", err)
	}

	log.Println("Migration completed successfully!")
}
