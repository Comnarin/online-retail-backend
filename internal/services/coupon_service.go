package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
)

type CouponService struct {
	repo repository.ICouponRepository
}

func NewCouponService(repo repository.ICouponRepository) ICouponService {
	return &CouponService{repo: repo}
}

func (u *CouponService) Create(ctx context.Context, coupon *domain.Coupon) error {
	return u.repo.Create(ctx, coupon)
}

func (u *CouponService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Coupon, error) {
	return u.repo.GetByID(ctx, tenantID, id)
}

func (u *CouponService) List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Coupon, int64, error) {
	return u.repo.List(ctx, tenantID, opts)
}

func (u *CouponService) Update(ctx context.Context, tenantID uuid.UUID, coupon *domain.Coupon) error {
	return u.repo.Update(ctx, tenantID, coupon)
}

func (u *CouponService) Toggle(ctx context.Context, tenantID, id uuid.UUID) (*domain.Coupon, error) {
	if err := u.repo.ToggleActive(ctx, tenantID, id); err != nil {
		return nil, err
	}
	return u.repo.GetByID(ctx, tenantID, id)
}

func (u *CouponService) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return u.repo.Delete(ctx, tenantID, id)
}
