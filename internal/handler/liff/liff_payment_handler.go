package liff

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/internal/handler"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type liffPaymentHandler struct {
	pmSvc  service.IPaymentMethodService
	cpmSvc service.ICustomerPaymentMethodService
}

func NewLiffPaymentHandler(pmSvc service.IPaymentMethodService, cpmSvc service.ICustomerPaymentMethodService) handler.ILiffPaymentHandler {
	return &liffPaymentHandler{pmSvc: pmSvc, cpmSvc: cpmSvc}
}

func (h *liffPaymentHandler) ListMethods(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	methods, err := h.pmSvc.List(c.Context(), *tenantID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, methods)
}

func (h *liffPaymentHandler) ListSavedMethods(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)
	methods, err := h.cpmSvc.List(c.Context(), *tenantID, customerID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, methods)
}

func (h *liffPaymentHandler) CreateSavedMethod(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)

	var req struct {
		PaymentMethodID uuid.UUID `json:"payment_method_id"`
		Provider        string    `json:"provider"`
		ProviderRef     string    `json:"provider_ref"`
		Last4           string    `json:"last4"`
		Brand           string    `json:"brand"`
		ExpMonth        int       `json:"exp_month"`
		ExpYear         int       `json:"exp_year"`
		IsDefault       bool      `json:"is_default"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	cpm := &domain.CustomerPaymentMethod{
		TenantID:        *tenantID,
		CustomerID:      customerID,
		PaymentMethodID: req.PaymentMethodID,
		Provider:        req.Provider,
		ProviderRef:     req.ProviderRef,
		Last4:           req.Last4,
		Brand:           req.Brand,
		ExpMonth:        req.ExpMonth,
		ExpYear:         req.ExpYear,
		IsDefault:       req.IsDefault,
	}

	if err := h.cpmSvc.Create(c.Context(), cpm, req.IsDefault); err != nil {
		return response.InternalError(c, err)
	}

	return response.Created(c, cpm)
}

func (h *liffPaymentHandler) UpdateSavedMethod(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Params("id"))

	var req struct {
		IsDefault bool `json:"is_default"`
	}

	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	cpm, err := h.cpmSvc.GetByID(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "payment method not found")
	}

	if cpm.CustomerID != customerID {
		return response.Forbidden(c)
	}

	cpm.IsDefault = req.IsDefault

	if err := h.cpmSvc.Update(c.Context(), cpm, req.IsDefault); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, cpm)
}

func (h *liffPaymentHandler) DeleteSavedMethod(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Params("id"))

	if err := h.cpmSvc.Delete(c.Context(), *tenantID, customerID, id); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, nil)
}

func (h *liffPaymentHandler) SetDefaultMethod(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	customerID := c.Locals("user_id").(uuid.UUID)
	id, _ := uuid.Parse(c.Params("id"))

	cpm, err := h.cpmSvc.GetByID(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "payment method not found")
	}

	if cpm.CustomerID != customerID {
		return response.Forbidden(c)
	}

	if err := h.cpmSvc.Update(c.Context(), cpm, true); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, nil)
}
