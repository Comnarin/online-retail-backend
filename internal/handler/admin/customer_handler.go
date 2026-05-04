package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type CustomerHandler struct {
	customerSvc service.ICustomerService
}

func NewCustomerHandler(customerSvc service.ICustomerService) handler.ICustomerHandler {
	return &CustomerHandler{
		customerSvc: customerSvc,
	}
}

func (h *CustomerHandler) List(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	opts := model.ListOptions{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
		Search: c.Query("search"),
	}
	customers, total, err := h.customerSvc.ListCustomers(c.Context(), *tenantID, opts)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.Paginated(c, customers, total, opts.Page, opts.Limit)
}

func (h *CustomerHandler) Get(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid customer id")
	}
	customer, err := h.customerSvc.GetCustomer(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "customer not found")
	}
	return response.OK(c, customer)
}
