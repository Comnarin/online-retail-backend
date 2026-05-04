package admin

import (
	"errors"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/repository"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/internal/platform/storage"
	"github.com/retail/backend/pkg/response"
	"github.com/retail/backend/pkg/validator"
)

type ProductHandler struct {
	productSvc service.IProductService
	storage    *storage.MinIOClient
}

func NewProductHandler(productSvc service.IProductService, storage *storage.MinIOClient) handler.IProductHandler {
	return &ProductHandler{productSvc: productSvc, storage: storage}
}

func (h *ProductHandler) List(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	opts := model.ListOptions{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
		Search: c.Query("search"),
	}
	products, total, err := h.productSvc.List(c.Context(), *tenantID, opts)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Paginated(c, products, total, opts.Page, opts.Limit)
}

func (h *ProductHandler) Get(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}
	product, err := h.productSvc.GetByID(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "product not found")
	}

	return response.OK(c, product)
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	var req model.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body", err)
	}

	// Validation
	if err := validator.Struct(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	product, err := h.productSvc.Create(c.Context(), *tenantID, &req)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.Created(c, product)
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}

	var req model.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body", err)
	}

	// Validation
	if err := validator.Struct(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	product, err := h.productSvc.Update(c.Context(), *tenantID, id, &req)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, product)
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}
	if err := h.productSvc.Delete(c.Context(), *tenantID, id); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"deleted": true})
}

func (h *ProductHandler) ListCategories(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	cats, err := h.productSvc.ListCategories(c.Context(), *tenantID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, cats)
}

func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	var cat domain.ProductCategory
	if err := c.BodyParser(&cat); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	cat.TenantID = *tenantID

	if err := validator.Struct(cat); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.productSvc.CreateCategory(c.Context(), &cat); err != nil {
		return response.InternalError(c, err)
	}
	return response.Created(c, cat)
}

func (h *ProductHandler) UpdateCategory(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid category id")
	}

	var cat domain.ProductCategory
	if err := c.BodyParser(&cat); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.productSvc.UpdateCategory(c.Context(), *tenantID, id, &cat); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, cat)
}

func (h *ProductHandler) DeleteCategory(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid category id")
	}

	if err := h.productSvc.DeleteCategory(c.Context(), *tenantID, id); err != nil {
		if errors.Is(err, repository.ErrNotEmpty) {
			return response.BadRequest(c, "cannot delete category with products assigned")
		}
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"deleted": true})
}

func (h *ProductHandler) UploadImage(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}

	// Ensure product exists
	_, err = h.productSvc.GetByID(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "product not found")
	}

	file, err := c.FormFile("image")
	if err != nil {
		return response.BadRequest(c, "image file is required")
	}

	src, err := file.Open()
	if err != nil {
		return response.InternalError(c, err)
	}
	defer src.Close()

	// Upload to MinIO
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return response.BadRequest(c, "only jpg, png, webp images are allowed")
	}

	objectName := "products/" + tenantID.String() + "/" + uuid.New().String() + ext
	url, err := h.storage.UploadFile(c.Context(), objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return response.InternalError(c, err)
	}

	// Create ProductImage
	productImage := &domain.ProductImage{
		ProductID: id,
		URL:       url,
	}

	if err := h.productSvc.AddImage(c.Context(), productImage); err != nil {
		return response.InternalError(c, err)
	}

	return response.Created(c, productImage)
}

func (h *ProductHandler) DeleteImage(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	productID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid product id")
	}

	imageID, err := uuid.Parse(c.Params("image_id"))
	if err != nil {
		return response.BadRequest(c, "invalid image id")
	}

	if err := h.productSvc.RemoveImage(c.Context(), *tenantID, productID, imageID); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"deleted": true})
}
