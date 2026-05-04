package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type tenantRepo struct{ db *gorm.DB }

func NewTenantRepository(db *gorm.DB) repository.ITenantRepository {
	return &tenantRepo{db: db}
}

func (r *tenantRepo) WithTx(tx *gorm.DB) repository.ITenantRepository {
	return &tenantRepo{db: tx}
}

func (r *tenantRepo) Create(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).Preload("Features").Preload("Appearance").
		Where("id = ?", id).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepo) GetBySlug(ctx context.Context, slug string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).Preload("Features").Preload("Appearance").
		Where("slug = ?", slug).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepo) List(ctx context.Context, opts model.ListOptions) ([]domain.Tenant, int64, error) {
	var tenants []domain.Tenant
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.Tenant{})
	if opts.Search != "" {
		q = q.Where("name ILIKE ?", "%"+opts.Search+"%")
	}
	q.Count(&total)

	offset := (opts.Page - 1) * opts.Limit
	err := q.Preload("Features").Preload("Appearance").Offset(offset).Limit(opts.Limit).Order("created_at DESC").Find(&tenants).Error
	return tenants, total, err
}

func (r *tenantRepo) Update(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *tenantRepo) UpdateFeatures(ctx context.Context, id uuid.UUID, features domain.TenantFeature) error {
	features.TenantID = id
	return r.db.WithContext(ctx).Save(&features).Error
}

func (r *tenantRepo) UpdateAppearance(ctx context.Context, id uuid.UUID, appearance domain.TenantAppearance) error {
	appearance.TenantID = id
	return r.db.WithContext(ctx).Save(&appearance).Error
}

func (r *tenantRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Tenant{}).Error
}
