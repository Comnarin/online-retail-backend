package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
)

type MembershipService struct {
	repo      repository.IMembershipRepository
	pointRepo repository.IPointRepository
}

func NewMembershipService(repo repository.IMembershipRepository, pointRepo repository.IPointRepository) IMembershipService {
	return &MembershipService{repo: repo, pointRepo: pointRepo}
}

func (u *MembershipService) ListTransactions(ctx context.Context, tenantID uuid.UUID, customerID *uuid.UUID, opts model.ListOptions) ([]domain.PointTransaction, int64, error) {
	if customerID != nil {
		return u.pointRepo.ListTransactions(ctx, tenantID, *customerID, opts)
	}
	// If no customerID, we need a tenant-wide list.
	return u.pointRepo.ListTransactions(ctx, tenantID, uuid.Nil, opts)
}

func (u *MembershipService) CreateTier(ctx context.Context, tier *domain.MembershipTier) error {
	return u.repo.CreateTier(ctx, tier)
}

func (u *MembershipService) ListTiers(ctx context.Context, tenantID uuid.UUID) ([]domain.MembershipTier, error) {
	return u.repo.ListTiers(ctx, tenantID)
}

func (u *MembershipService) UpdateTier(ctx context.Context, tenantID uuid.UUID, tier *domain.MembershipTier) error {
	return u.repo.UpdateTier(ctx, tenantID, tier)
}

func (u *MembershipService) DeleteTier(ctx context.Context, tenantID, id uuid.UUID) error {
	return u.repo.DeleteTier(ctx, tenantID, id)
}

func (u *MembershipService) GetCustomerMembership(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.CustomerMembership, error) {
	return u.repo.GetCustomerMembership(ctx, tenantID, customerID)
}
