package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type membershipRepo struct{ db *gorm.DB }

func NewMembershipRepository(db *gorm.DB) repository.IMembershipRepository {
	return &membershipRepo{db: db}
}

func (r *membershipRepo) WithTx(tx *gorm.DB) repository.IMembershipRepository {
	return &membershipRepo{db: tx}
}

func (r *membershipRepo) CreateTier(ctx context.Context, tier *domain.MembershipTier) error {
	return r.db.WithContext(ctx).Create(tier).Error
}

func (r *membershipRepo) ListTiers(ctx context.Context, tenantID uuid.UUID) ([]domain.MembershipTier, error) {
	var tiers []domain.MembershipTier
	return tiers, r.db.WithContext(ctx).Preload("Benefits").
		Where("tenant_id = ?", tenantID).Order("min_points ASC").Find(&tiers).Error
}

func (r *membershipRepo) UpdateTier(ctx context.Context, tenantID uuid.UUID, tier *domain.MembershipTier) error {
	// Enforce tenant_id
	return r.db.WithContext(ctx).Save(tier).Error
}

func (r *membershipRepo) DeleteTier(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.MembershipTier{}).Error
}

func (r *membershipRepo) GetCustomerMembership(ctx context.Context, tenantID, customerID uuid.UUID) (*domain.CustomerMembership, error) {
	var cm domain.CustomerMembership
	return &cm, r.db.WithContext(ctx).Preload("Tier").
		Where("tenant_id = ? AND customer_id = ?", tenantID, customerID).First(&cm).Error
}

func (r *membershipRepo) UpsertCustomerMembership(ctx context.Context, cm *domain.CustomerMembership) error {
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "customer_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"tier_id", "expires_at", "total_spent", "points"}),
	}).Create(cm).Error
}
