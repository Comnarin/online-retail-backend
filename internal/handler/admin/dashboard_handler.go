package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type DashboardHandler struct {
	dashboardSvc service.IDashboardService
}

func NewDashboardHandler(dashboardSvc service.IDashboardService) handler.IDashboardHandler {
	return &DashboardHandler{dashboardSvc: dashboardSvc}
}

func (h *DashboardHandler) Stats(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	stats, err := h.dashboardSvc.GetStats(c.Context(), *tenantID)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, stats)
}

func (h *DashboardHandler) RecentOrders(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	limit := c.QueryInt("limit", 10)

	orders, err := h.dashboardSvc.GetRecentOrders(c.Context(), *tenantID, limit)
	if err != nil {
		return response.InternalError(c, err)
	}

	_ = model.ListOptions{}
	return response.OK(c, orders)
}
