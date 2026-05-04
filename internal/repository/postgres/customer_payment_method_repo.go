package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type customerPaymentMethodRepo struct{ db *gorm.DB }

func NewCustomerPaymentMethodRepository(db *gorm.DB) repository.ICustomerPaymentMethodRepository {
	return &customerPaymentMethodRepo{db: db}
}

func (r *customerPaymentMethodRepo) List(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerPaymentMethod, error) {
	var methods []domain.CustomerPaymentMethod
	err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Order("is_default DESC, created_at DESC").
		Find(&methods).Error
	return methods, err
}

func (r *customerPaymentMethodRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomerPaymentMethod, error) {
	var method domain.CustomerPaymentMethod
	err := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&method).Error
	return &method, err
}

func (r *customerPaymentMethodRepo) Create(ctx context.Context, cpm *domain.CustomerPaymentMethod) error {
	return r.db.WithContext(ctx).Create(cpm).Error
}

func (r *customerPaymentMethodRepo) Update(ctx context.Context, cpm *domain.CustomerPaymentMethod) error {
	return r.db.WithContext(ctx).Save(cpm).Error
}

func (r *customerPaymentMethodRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.CustomerPaymentMethod{}).Error
}

func (r *customerPaymentMethodRepo) SetDefault(ctx context.Context, tenantID, customerID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Reset all for this customer
		if err := tx.Model(&domain.CustomerPaymentMethod{}).
			Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Set default
		return tx.Model(&domain.CustomerPaymentMethod{}).
			Where("tenant_id = ? AND id = ?", tenantID, id).
			Update("is_default", true).Error
	})
}
