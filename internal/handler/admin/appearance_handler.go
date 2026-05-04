package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/handler"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type AppearanceHandler struct {
	dashboardSvc service.IDashboardService
}

func NewAppearanceHandler(dashboardSvc service.IDashboardService) handler.IAppearanceHandler {
	return &AppearanceHandler{dashboardSvc: dashboardSvc}
}

func (h *AppearanceHandler) Get(c *fiber.Ctx) error {
	tenantID, _ := c.Locals("tenant_id").(*uuid.UUID)

	appearance, err := h.dashboardSvc.GetAppearance(c.Context(), *tenantID)
	if err != nil {
		return response.NotFound(c, "appearance not found")
	}

	return response.OK(c, appearance)
}

func (h *AppearanceHandler) Update(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	var body struct {
		Name       string                  `json:"name"`
		OrderCode  string                  `json:"order_code"`
		Appearance domain.TenantAppearance `json:"appearance"`
		Features   domain.TenantFeature    `json:"features"`
	}

	if err := c.BodyParser(&body); err != nil {
		return response.BadRequest(c, "invalid body")
	}

	if err := h.dashboardSvc.UpdateAppearance(c.Context(), *tenantID, body.Name, body.OrderCode, body.Appearance, body.Features); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"name":       body.Name,
		"order_code": body.OrderCode,
		"appearance": body.Appearance,
		"features":   body.Features,
	})
}
