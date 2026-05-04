package superadmin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
	"github.com/retail/backend/pkg/validator"
)

type TenantHandler struct {
	tenantSvc service.ITenantService
	authSvc   service.IAuthService
}

func NewTenantHandler(tenantSvc service.ITenantService, authSvc service.IAuthService) handler.ITenantHandler {
	return &TenantHandler{
		tenantSvc: tenantSvc,
		authSvc:   authSvc,
	}
}

func (h *TenantHandler) List(c *fiber.Ctx) error {
	opts := model.ListOptions{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
		Search: c.Query("search"),
	}
	tenants, total, err := h.tenantSvc.List(c.Context(), opts)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.Paginated(c, tenants, total, opts.Page, opts.Limit)
}

func (h *TenantHandler) Create(c *fiber.Ctx) error {
	var tenant domain.Tenant
	if err := c.BodyParser(&tenant); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := validator.Struct(tenant); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.tenantSvc.Create(c.Context(), &tenant); err != nil {
		return response.InternalError(c, err)
	}
	return response.Created(c, tenant)
}

func (h *TenantHandler) Get(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}
	tenant, err := h.tenantSvc.GetByID(c.Context(), id)
	if err != nil {
		return response.NotFound(c, "tenant not found")
	}
	return response.OK(c, tenant)
}

func (h *TenantHandler) Update(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}

	var tenant domain.Tenant
	if err := c.BodyParser(&tenant); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	tenant.ID = id

	if err := validator.Struct(tenant); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.tenantSvc.Update(c.Context(), &tenant); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, tenant)
}

func (h *TenantHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}
	if err := h.tenantSvc.Delete(c.Context(), id); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"deleted": true})
}

func (h *TenantHandler) UpdateFeatures(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}
	var features domain.TenantFeature
	if err := c.BodyParser(&features); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := validator.Struct(features); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.tenantSvc.UpdateFeatures(c.Context(), id, features); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"updated": true})
}

func (h *TenantHandler) UpdateAppearance(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}
	var appearance domain.TenantAppearance
	if err := c.BodyParser(&appearance); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if err := validator.Struct(appearance); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.tenantSvc.UpdateAppearance(c.Context(), id, appearance); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"updated": true})
}

func (h *TenantHandler) CreateAdmin(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}
	var req model.RegisterAdminRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	req.TenantID = id

	if err := validator.Struct(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	user, err := h.authSvc.RegisterAdmin(c.Context(), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}
	return response.Created(c, user)
}

func (h *TenantHandler) ListAdmins(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}
	opts := model.ListOptions{
		Page:  c.QueryInt("page", 1),
		Limit: c.QueryInt("limit", 100),
	}
	admins, total, err := h.tenantSvc.ListAdmins(c.Context(), id, opts)
	if err != nil {
		return response.InternalError(c, err)
	}

	if admins == nil {
		admins = []domain.User{}
	}

	return response.Paginated(c, admins, total, opts.Page, opts.Limit)
}

func (h *TenantHandler) DeleteAdmin(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id")) // tenant id
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}
	adminId, err := uuid.Parse(c.Params("adminId"))
	if err != nil {
		return response.BadRequest(c, "invalid admin id")
	}
	if err := h.tenantSvc.DeleteAdmin(c.Context(), id, adminId); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"deleted": true})
}
