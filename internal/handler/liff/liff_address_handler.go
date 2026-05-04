package liff

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/handler"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type liffAddressHandler struct {
	svc service.ICustomerAddressService
}

func NewLiffAddressHandler(svc service.ICustomerAddressService) handler.ILiffAddressHandler {
	return &liffAddressHandler{svc: svc}
}

func (h *liffAddressHandler) List(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)

	addrs, err := h.svc.List(c.Context(), *tenantID, customerID)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, addrs)
}

func (h *liffAddressHandler) Create(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)

	var req struct {
		Label     string `json:"label"`
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
		District  string `json:"district"`
		Province  string `json:"province"`
		ZipCode   string `json:"zip_code"`
		IsDefault bool   `json:"is_default"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	addr := &domain.CustomerAddress{
		TenantID:   *tenantID,
		CustomerID: customerID,
		Label:      req.Label,
		Name:       req.Name,
		Phone:      req.Phone,
		Address:    req.Address,
		District:   req.District,
		Province:   req.Province,
		ZipCode:    req.ZipCode,
		IsDefault:  req.IsDefault,
	}

	if err := h.svc.Create(c.Context(), addr, req.IsDefault); err != nil {
		return response.InternalError(c, err)
	}

	return response.Created(c, addr)
}

func (h *liffAddressHandler) Update(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Params("id"))

	var req struct {
		Label     string `json:"label"`
		Name      string `json:"name"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
		District  string `json:"district"`
		Province  string `json:"province"`
		ZipCode   string `json:"zip_code"`
		IsDefault bool   `json:"is_default"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	addr, err := h.svc.GetByID(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "address not found")
	}

	if addr.CustomerID != customerID {
		return response.Forbidden(c)
	}

	addr.Label = req.Label
	addr.Name = req.Name
	addr.Phone = req.Phone
	addr.Address = req.Address
	addr.District = req.District
	addr.Province = req.Province
	addr.ZipCode = req.ZipCode
	addr.IsDefault = req.IsDefault

	if err := h.svc.Update(c.Context(), addr, req.IsDefault); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, addr)
}

func (h *liffAddressHandler) Delete(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Params("id"))

	if err := h.svc.Delete(c.Context(), *tenantID, customerID, id); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, nil)
}
