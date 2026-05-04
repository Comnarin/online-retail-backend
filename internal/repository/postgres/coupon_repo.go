package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type couponRepo struct{ db *gorm.DB }

func NewCouponRepository(db *gorm.DB) repository.ICouponRepository {
	return &couponRepo{db: db}
}

func (r *couponRepo) Create(ctx context.Context, coupon *domain.Coupon) error {
	return r.db.WithContext(ctx).Create(coupon).Error
}

func (r *couponRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Coupon, error) {
	var c domain.Coupon
	if err := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *couponRepo) GetByCode(ctx context.Context, tenantID uuid.UUID, code string) (*domain.Coupon, error) {
	var c domain.Coupon
	if err := r.db.WithContext(ctx).Where("tenant_id = ? AND code = ?", tenantID, code).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *couponRepo) List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Coupon, int64, error) {
	var coupons []domain.Coupon
	var total int64
	q := r.db.WithContext(ctx).Model(&domain.Coupon{}).Where("tenant_id = ?", tenantID)
	q.Count(&total)
	offset := (opts.Page - 1) * opts.Limit
	err := q.Offset(offset).Limit(opts.Limit).Order("created_at DESC").Find(&coupons).Error
	return coupons, total, err
}

func (r *couponRepo) Update(ctx context.Context, tenantID uuid.UUID, coupon *domain.Coupon) error {
	return r.db.WithContext(ctx).Model(&domain.Coupon{}).
		Where("id = ? AND tenant_id = ?", coupon.ID, tenantID).
		Updates(coupon).Error
}

func (r *couponRepo) IncrementUsage(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Coupon{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

func (r *couponRepo) ToggleActive(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&domain.Coupon{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		UpdateColumn("is_active", gorm.Expr("NOT is_active")).Error
}

func (r *couponRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.Coupon{}).Error
}

func (r *couponRepo) WithTx(tx *gorm.DB) repository.ICouponRepository {
	return &couponRepo{db: tx}
}