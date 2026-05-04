package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
)

type ProductService struct {
	productRepo repository.IProductRepository
}

func NewProductService(productRepo repository.IProductRepository) IProductService {
	return &ProductService{productRepo: productRepo}
}

func (u *ProductService) Create(ctx context.Context, tenantID uuid.UUID, req *model.CreateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		TenantID:      tenantID,
		CategoryID:    req.CategoryID,
		NameTh:        req.NameTh,
		NameEn:        req.NameEn,
		DescriptionTh: req.DescriptionTh,
		DescriptionEn: req.DescriptionEn,
		Price:         req.Price,
		Inventory:     req.Inventory,
		Status:        req.Status,
	}

	// Normalize SKU: empty string should be NULL in database
	if req.SKU != nil && *req.SKU != "" {
		product.SKU = req.SKU
	}

	if product.Status == "" {
		product.Status = "active"
	}

	if err := u.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

func (u *ProductService) GetByID(ctx context.Context, tenantID, id uuid.UUID) (*domain.Product, error) {
	return u.productRepo.GetByID(ctx, tenantID, id)
}

func (u *ProductService) List(ctx context.Context, tenantID uuid.UUID, opts model.ListOptions) ([]domain.Product, int64, error) {
	return u.productRepo.List(ctx, tenantID, opts)
}

func (u *ProductService) Update(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, req *model.UpdateProductRequest) (*domain.Product, error) {
	product := &domain.Product{
		ID:         id,
		TenantID:   tenantID,
		CategoryID: req.CategoryID,
		NameTh:     req.NameTh,
		NameEn:     req.NameEn,
		Status:     req.Status,
	}

	// Price and Inventory are handled via direct checks since they are pointers or float/int
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Inventory != nil {
		product.Inventory = *req.Inventory
	}
	if req.DescriptionTh != "" {
		product.DescriptionTh = req.DescriptionTh
	}
	if req.DescriptionEn != "" {
		product.DescriptionEn = req.DescriptionEn
	}

	// SKU handling:
	// If "" or nil is passed for SKU, we want NULL in the DB.
	// We use a map for Updates to ensure GORM doesn't skip the NULL/nil field
	updateData := make(map[string]interface{})
	if req.NameTh != "" {
		updateData["name_th"] = req.NameTh
	}
	if req.NameEn != "" {
		updateData["name_en"] = req.NameEn
	}
	if req.CategoryID != nil {
		updateData["category_id"] = req.CategoryID
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}
	if req.Price != nil {
		updateData["price"] = *req.Price
	}
	if req.Inventory != nil {
		updateData["inventory"] = *req.Inventory
	}
	if req.DescriptionTh != "" {
		updateData["description_th"] = req.DescriptionTh
	}
	if req.DescriptionEn != "" {
		updateData["description_en"] = req.DescriptionEn
	}
	
	// Handle SKU clearing
	if req.SKU != nil {
		if *req.SKU == "" {
			updateData["sku"] = nil
		} else {
			updateData["sku"] = *req.SKU
		}
	}

	// Update with map to ensure nullable fields like SKU are correctly handled
	if err := u.productRepo.UpdateFields(ctx, tenantID, id, updateData); err != nil {
		return nil, err
	}
	return product, nil
}

func (u *ProductService) Delete(ctx context.Context, tenantID, id uuid.UUID) error {
	return u.productRepo.Delete(ctx, tenantID, id)
}

func (u *ProductService) CreateCategory(ctx context.Context, cat *domain.ProductCategory) error {
	return u.productRepo.CreateCategory(ctx, cat)
}

func (u *ProductService) ListCategories(ctx context.Context, tenantID uuid.UUID) ([]domain.ProductCategory, error) {
	return u.productRepo.ListCategories(ctx, tenantID)
}

func (u *ProductService) UpdateCategory(ctx context.Context, tenantID uuid.UUID, id uuid.UUID, cat *domain.ProductCategory) error {
	cat.ID = id
	return u.productRepo.UpdateCategory(ctx, tenantID, cat)
}

func (u *ProductService) DeleteCategory(ctx context.Context, tenantID, id uuid.UUID) error {
	return u.productRepo.DeleteCategory(ctx, tenantID, id)
}

func (u *ProductService) AddImage(ctx context.Context, image *domain.ProductImage) error {
	// Let repo handle adding the image
	return u.productRepo.AddImage(ctx, image)
}

func (u *ProductService) RemoveImage(ctx context.Context, tenantID, productID, imageID uuid.UUID) error {
	return u.productRepo.RemoveImage(ctx, tenantID, productID, imageID)
}
