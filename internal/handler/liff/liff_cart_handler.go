package liff

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type LiffCartHandler struct {
	cartSvc service.ICartService
}

func NewLiffCartHandler(cartSvc service.ICartService) handler.ILiffCartHandler {
	return &LiffCartHandler{cartSvc: cartSvc}
}

func (h *LiffCartHandler) Get(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)

	cart, err := h.cartSvc.GetCart(c.Context(), *tenantID, userID)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, cart)
}

func (h *LiffCartHandler) AddItem(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)

	var req struct {
		ProductID uuid.UUID `json:"product_id"`
		Quantity  int       `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.cartSvc.AddItem(c.Context(), *tenantID, userID, req.ProductID, req.Quantity); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true})
}

func (h *LiffCartHandler) UpdateItem(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)
	itemID, _ := uuid.Parse(c.Params("id"))

	var req struct {
		Quantity int `json:"quantity"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.cartSvc.UpdateItemQty(c.Context(), *tenantID, userID, itemID, req.Quantity); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true})
}

func (h *LiffCartHandler) RemoveItem(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)
	itemID, _ := uuid.Parse(c.Params("id"))

	if err := h.cartSvc.RemoveItem(c.Context(), *tenantID, userID, itemID); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true})
}

func (h *LiffCartHandler) ToggleSelection(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)
	itemID, _ := uuid.Parse(c.Params("id"))

	var req struct {
		Selected bool `json:"selected"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.cartSvc.ToggleSelection(c.Context(), *tenantID, userID, itemID, req.Selected); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true})
}

func (h *LiffCartHandler) Clear(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)

	if err := h.cartSvc.ClearCart(c.Context(), *tenantID, userID); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true})
}
