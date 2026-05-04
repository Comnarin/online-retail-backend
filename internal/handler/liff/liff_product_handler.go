package liff

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type LiffProductHandler struct {
	productSvc service.IProductService
}

func NewLiffProductHandler(productSvc service.IProductService) handler.ILiffProductHandler {
	return &LiffProductHandler{productSvc: productSvc}
}

func (h *LiffProductHandler) List(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(*uuid.UUID)
	opts := model.ListOptions{
		Page:     c.QueryInt("page", 1),
		Limit:    c.QueryInt("limit", 20),
		Search:   c.Query("search"),
		Category: c.Query("category"),
	}
	products, total, err := h.productSvc.List(c.Context(), *tenantID, opts)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Paginated(c, products, total, opts.Page, opts.Limit)
}

func (h *LiffProductHandler) Get(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(*uuid.UUID)
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

func (h *LiffProductHandler) ListCategories(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(*uuid.UUID)
	cats, err := h.productSvc.ListCategories(c.Context(), *tenantID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, cats)
}
