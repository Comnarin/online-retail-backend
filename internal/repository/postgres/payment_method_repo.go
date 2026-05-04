package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type paymentMethodRepo struct{ db *gorm.DB }

func NewPaymentMethodRepository(db *gorm.DB) repository.IPaymentMethodRepository {
	return &paymentMethodRepo{db: db}
}

func (r *paymentMethodRepo) List(ctx context.Context, tenantID uuid.UUID) ([]domain.PaymentMethod, error) {
	var methods []domain.PaymentMethod
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND is_active = ?", tenantID, true).Find(&methods).Error
	return methods, err
}

func (r *paymentMethodRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.PaymentMethod, error) {
	var method domain.PaymentMethod
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).First(&method).Error
	return &method, err
}

func (r *paymentMethodRepo) Create(ctx context.Context, pm *domain.PaymentMethod) error {
	return r.db.WithContext(ctx).Create(pm).Error
}

func (r *paymentMethodRepo) Update(ctx context.Context, tenantID uuid.UUID, pm *domain.PaymentMethod) error {
	return r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Save(pm).Error
}

func (r *paymentMethodRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("tenant_id = ? AND id = ?", tenantID, id).Delete(&domain.PaymentMethod{}).Error
}
