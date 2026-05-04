package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
)

type DashboardService struct {
	dashboardRepo repository.IDashboardRepository
	tenantRepo    repository.ITenantRepository
}

func NewDashboardService(dashboardRepo repository.IDashboardRepository, tenantRepo repository.ITenantRepository) IDashboardService {
	return &DashboardService{
		dashboardRepo: dashboardRepo,
		tenantRepo:    tenantRepo,
	}
}

func (s *DashboardService) GetStats(ctx context.Context, tenantID uuid.UUID) (*model.DashboardStats, error) {
	return s.dashboardRepo.GetStats(ctx, tenantID)
}

func (s *DashboardService) GetRecentOrders(ctx context.Context, tenantID uuid.UUID, limit int) ([]model.RecentOrderRow, error) {
	return s.dashboardRepo.GetRecentOrders(ctx, tenantID, limit)
}

func (s *DashboardService) GetAppearance(ctx context.Context, tenantID uuid.UUID) (*domain.Tenant, error) {
	return s.tenantRepo.GetByID(ctx, tenantID)
}

func (s *DashboardService) UpdateAppearance(ctx context.Context, tenantID uuid.UUID, name, orderCode string, appearance domain.TenantAppearance, features domain.TenantFeature) error {
	// Update Tenant Name and OrderCode
	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}
	tenant.Name = name
	tenant.OrderCode = orderCode
	if err := s.tenantRepo.Update(ctx, tenant); err != nil {
		return err
	}

	// Update Features
	features.TenantID = tenantID
	if err := s.tenantRepo.UpdateFeatures(ctx, tenantID, features); err != nil {
		return err
	}

	// Update Appearance
	appearance.TenantID = tenantID
	return s.tenantRepo.UpdateAppearance(ctx, tenantID, appearance)
}
