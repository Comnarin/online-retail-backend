package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type OrderHandler struct {
	orderSvc service.IOrderService
}

func NewOrderHandler(orderSvc service.IOrderService) handler.IOrderHandler {
	return &OrderHandler{orderSvc: orderSvc}
}

func (h *OrderHandler) List(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	opts := model.ListOptions{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
		Search: c.Query("search"),
	}
	orders, total, err := h.orderSvc.ListByTenant(c.Context(), *tenantID, opts)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.Paginated(c, orders, total, opts.Page, opts.Limit)
}

func (h *OrderHandler) Get(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}
	order, err := h.orderSvc.GetOrder(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "order not found")
	}
	return response.OK(c, order)
}

func (h *OrderHandler) UpdateStatus(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}

	var req struct {
		Status domain.OrderStatus `json:"status" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := h.orderSvc.UpdateStatus(c.Context(), *tenantID, id, req.Status); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true})
}

func (h *OrderHandler) ConfirmPayment(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}

	if err := h.orderSvc.ConfirmPayment(c.Context(), *tenantID, id); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true, "status": "confirmed"})
}

func (h *OrderHandler) RejectPayment(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}

	if err := h.orderSvc.RejectPayment(c.Context(), *tenantID, id); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true, "status": "cancelled"})
}

