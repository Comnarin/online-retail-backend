package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type cartRepo struct{ db *gorm.DB }

func NewCartRepository(db *gorm.DB) repository.ICartRepository {
	return &cartRepo{db: db}
}

func (r *cartRepo) GetByCustomerID(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.Cart, error) {
	var cart domain.Cart
	// Find or Create the cart for the customer.
	// The unique index idx_cart_tenant_customer prevents duplicate carts.
	err := r.db.WithContext(ctx).
		Preload("Items").
		Preload("Items.Product").
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		FirstOrCreate(&cart, domain.Cart{
			TenantID:   tenantID,
			CustomerID: customerID,
		}).Error
	if err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *cartRepo) UpsertItem(ctx context.Context, item *domain.CartItem) error {
	// Use GORM's OnConflict for upsert based on (cart_id, product_id)
	// We need a unique constraint on (cart_id, product_id) for this to work perfectly, 
	// but let's do it manually for safety if the constraint isn't there yet.
	var existing domain.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", item.CartID, item.ProductID).
		First(&existing).Error

	if err == nil {
		// Update existing
		existing.Quantity += item.Quantity
		existing.Selected = true // Re-select if it was unselected
		return r.db.WithContext(ctx).Save(&existing).Error
	}

	// Create new
	if item.ID == uuid.Nil {
		item.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *cartRepo) UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error {
	return r.db.WithContext(ctx).Model(&domain.CartItem{}).
		Where("id = ?", itemID).
		Update("quantity", quantity).Error
}

func (r *cartRepo) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&domain.CartItem{}, itemID).Error
}

func (r *cartRepo) ToggleSelectItem(ctx context.Context, itemID uuid.UUID, selected bool) error {
	return r.db.WithContext(ctx).Model(&domain.CartItem{}).
		Where("id = ?", itemID).
		Update("selected", selected).Error
}

func (r *cartRepo) ClearSelectedItems(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("cart_id = ? AND selected = ?", cartID, true).
		Delete(&domain.CartItem{}).Error
}

func (r *cartRepo) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("cart_id = ?", cartID).
		Delete(&domain.CartItem{}).Error
}

func (r *cartRepo) WithTx(tx *gorm.DB) repository.ICartRepository {
	return &cartRepo{db: tx}
}
