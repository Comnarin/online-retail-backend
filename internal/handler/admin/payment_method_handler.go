package admin

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/internal/platform/storage"
	"github.com/retail/backend/pkg/response"
)

type PaymentMethodHandler struct {
	pmSvc   service.IPaymentMethodService
	storage *storage.MinIOClient
}

func NewPaymentMethodHandler(pmSvc service.IPaymentMethodService, storage *storage.MinIOClient) handler.IAdminPaymentMethodHandler {
	return &PaymentMethodHandler{pmSvc: pmSvc, storage: storage}
}

func (h *PaymentMethodHandler) List(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	methods, err := h.pmSvc.List(c.Context(), *tenantID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, methods)
}

func (h *PaymentMethodHandler) Update(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid payment method id")
	}

	pm, err := h.pmSvc.GetByID(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "payment method not found")
	}

	var req struct {
		Name          *string `json:"name"`
		Description   *string `json:"description"`
		ExpiryMinutes *int    `json:"expiry_minutes"`
		IsActive      *bool   `json:"is_active"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if req.Name != nil {
		pm.Name = *req.Name
	}
	if req.Description != nil {
		pm.Description = *req.Description
	}
	if req.ExpiryMinutes != nil {
		pm.ExpiryMinutes = *req.ExpiryMinutes
	}
	if req.IsActive != nil {
		pm.IsActive = *req.IsActive
	}

	if err := h.pmSvc.Update(c.Context(), *tenantID, pm); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, pm)
}

func (h *PaymentMethodHandler) UploadQRCode(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid payment method id")
	}

	pm, err := h.pmSvc.GetByID(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "payment method not found")
	}

	file, err := c.FormFile("qr_code")
	if err != nil {
		return response.BadRequest(c, "qr_code file is required")
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		return response.BadRequest(c, "only jpg, png, webp images are allowed")
	}

	// Open file
	src, err := file.Open()
	if err != nil {
		return response.InternalError(c, err)
	}
	defer src.Close()

	// Upload to MinIO
	objectName := fmt.Sprintf("qr-codes/%s/%s%s", tenantID.String(), time.Now().Format("20060102150405"), ext)
	url, err := h.storage.UploadFile(c.Context(), objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return response.InternalError(c, err)
	}

	// Update payment method
	pm.QRCodeURL = url
	if err := h.pmSvc.Update(c.Context(), *tenantID, pm); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"qr_code_url": url})
}
