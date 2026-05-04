package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
)

type TenantService struct {
	tenantRepo repository.ITenantRepository
	userRepo   repository.IUserRepository
}

func NewTenantService(tenantRepo repository.ITenantRepository, userRepo repository.IUserRepository) ITenantService {
	return &TenantService{tenantRepo: tenantRepo, userRepo: userRepo}
}

func (u *TenantService) Create(ctx context.Context, tenant *domain.Tenant) error {
	return u.tenantRepo.Create(ctx, tenant)
}

func (u *TenantService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	return u.tenantRepo.GetByID(ctx, id)
}

func (u *TenantService) List(ctx context.Context, opts model.ListOptions) ([]domain.Tenant, int64, error) {
	return u.tenantRepo.List(ctx, opts)
}

func (u *TenantService) Update(ctx context.Context, tenant *domain.Tenant) error {
	return u.tenantRepo.Update(ctx, tenant)
}

func (u *TenantService) UpdateFeatures(ctx context.Context, id uuid.UUID, features domain.TenantFeature) error {
	return u.tenantRepo.UpdateFeatures(ctx, id, features)
}

func (u *TenantService) UpdateAppearance(ctx context.Context, id uuid.UUID, appearance domain.TenantAppearance) error {
	return u.tenantRepo.UpdateAppearance(ctx, id, appearance)
}

func (u *TenantService) Delete(ctx context.Context, id uuid.UUID) error {
	return u.tenantRepo.Delete(ctx, id)
}

func (u *TenantService) ListAdmins(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.User, int64, error) {
	return u.userRepo.ListByTenant(ctx, tenantID, opts)
}

func (u *TenantService) DeleteAdmin(ctx context.Context, tenantID, adminID uuid.UUID) error {
	// First ensure the admin belongs to the correct tenant
	if _, err := u.userRepo.GetByID(ctx, tenantID, adminID); err != nil {
		return err
	}
	return u.userRepo.Delete(ctx, tenantID, adminID)
}
