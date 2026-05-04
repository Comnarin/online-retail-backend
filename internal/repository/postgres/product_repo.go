package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	"gorm.io/gorm"
)

type productRepo struct{ db *gorm.DB }

func NewProductRepository(db *gorm.DB) repository.IProductRepository {
	return &productRepo{db: db}
}

func (r *productRepo) Create(ctx context.Context, product *domain.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *productRepo) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Product, error) {
	var p domain.Product
	if err := r.db.WithContext(ctx).Preload("Category").Preload("Images").
		Where("id = ? AND tenant_id = ?", id, tenantID).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepo) List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Product, int64, error) {
	var products []domain.Product
	var total int64

	q := r.db.WithContext(ctx).Model(&domain.Product{}).Where("products.tenant_id = ?", tenantID)
	if opts.Search != "" {
		q = q.Where("products.name_th ILIKE ? OR products.name_en ILIKE ? OR products.sku ILIKE ?", "%"+opts.Search+"%", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}
	if opts.Category != "" {
		if catID, err := uuid.Parse(opts.Category); err == nil {
			q = q.Where("products.category_id = ?", catID)
		} else {
			q = q.Joins("Category").Where("\"Category\".name_th = ? OR \"Category\".name_en = ?", opts.Category, opts.Category)
		}
	}
	q.Count(&total)
	offset := (opts.Page - 1) * opts.Limit
	err := q.Preload("Category").Preload("Images").Offset(offset).Limit(opts.Limit).Order("products.sort_order ASC, products.created_at DESC").Find(&products).Error
	return products, total, err
}

func (r *productRepo) Update(ctx context.Context, tenantID uuid.UUID, product *domain.Product) error {
	// Enforce tenant_id in the where clause to prevent cross-tenant updates
	return r.db.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ? AND tenant_id = ?", product.ID, tenantID).
		Updates(product).Error
}

func (r *productRepo) UpdateFields(ctx context.Context, tenantID, id uuid.UUID, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.Product{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Updates(fields).Error
}

func (r *productRepo) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.Product{}).Error
}

func (r *productRepo) CreateCategory(ctx context.Context, cat *domain.ProductCategory) error {
	return r.db.WithContext(ctx).Create(cat).Error
}

func (r *productRepo) ListCategories(ctx context.Context, tenantID uuid.UUID) ([]domain.ProductCategory, error) {
	var cats []domain.ProductCategory
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Order("sort_order ASC").Find(&cats).Error
	return cats, err
}

func (r *productRepo) UpdateCategory(ctx context.Context, tenantID uuid.UUID, cat *domain.ProductCategory) error {
	return r.db.WithContext(ctx).Model(&domain.ProductCategory{}).
		Where("id = ? AND tenant_id = ?", cat.ID, tenantID).
		Updates(cat).Error
}

func (r *productRepo) DeleteCategory(ctx context.Context, tenantID, id uuid.UUID) error {
	// Check if categories have products
	var count int64
	r.db.WithContext(ctx).Model(&domain.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return repository.ErrNotEmpty
	}
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).
		Delete(&domain.ProductCategory{}).Error
}

func (r *productRepo) WithTx(tx *gorm.DB) repository.IProductRepository {
	return &productRepo{db: tx}
}

func (r *productRepo) AddImage(ctx context.Context, image *domain.ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *productRepo) RemoveImage(ctx context.Context, tenantID, productID, imageID uuid.UUID) error {
	// First ensure product belongs to tenant
	var count int64
	r.db.WithContext(ctx).Model(&domain.Product{}).Where("id = ? AND tenant_id = ?", productID, tenantID).Count(&count)
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.WithContext(ctx).Where("id = ? AND product_id = ?", imageID, productID).Delete(&domain.ProductImage{}).Error
}