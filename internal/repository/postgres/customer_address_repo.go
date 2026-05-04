package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type customerAddressRepo struct{ db *gorm.DB }

func NewCustomerAddressRepository(db *gorm.DB) repository.ICustomerAddressRepository {
	return &customerAddressRepo{db: db}
}

func (r *customerAddressRepo) List(ctx context.Context, tenantID, customerID uuid.UUID) ([]domain.CustomerAddress, error) {
	var addrs []domain.CustomerAddress
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
		Order("is_default DESC, created_at DESC").
		Find(&addrs).Error
	return addrs, err
}

func (r *customerAddressRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.CustomerAddress, error) {
	var addr domain.CustomerAddress
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&addr).Error
	return &addr, err
}

func (r *customerAddressRepo) Create(ctx context.Context, addr *domain.CustomerAddress) error {
	return r.db.WithContext(ctx).Create(addr).Error
}

func (r *customerAddressRepo) Update(ctx context.Context, addr *domain.CustomerAddress) error {
	return r.db.WithContext(ctx).Save(addr).Error
}

func (r *customerAddressRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&domain.CustomerAddress{}).Error
}

func (r *customerAddressRepo) SetDefault(ctx context.Context, tenantID, customerID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Reset all addresses for this customer to not default
		if err := tx.Model(&domain.CustomerAddress{}).
			Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).
			Update("is_default", false).Error; err != nil {
			return err
		}

		// Set the chosen one to default
		return tx.Model(&domain.CustomerAddress{}).
			Where("tenant_id = ? AND id = ?", tenantID, id).
			Update("is_default", true).Error
	})
}
