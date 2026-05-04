package liff

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	"github.com/retail/backend/internal/platform/storage"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
	"github.com/retail/backend/pkg/validator"
)

type LiffOrderHandler struct {
	orderSvc service.IOrderService
	storage  *storage.MinIOClient
}

func NewLiffOrderHandler(orderSvc service.IOrderService, storage *storage.MinIOClient) handler.ILiffOrderHandler {
	return &LiffOrderHandler{orderSvc: orderSvc, storage: storage}
}

func (h *LiffOrderHandler) Create(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)

	var req service.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	req.TenantID = *tenantID
	req.CustomerID = userID

	if err := validator.Struct(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	order, err := h.orderSvc.CreateOrder(c.Context(), req)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Created(c, order)
}

func (h *LiffOrderHandler) List(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uuid.UUID)
	opts := model.ListOptions{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
		Status: c.Query("status"),
	}
	orders, total, err := h.orderSvc.ListByCustomer(c.Context(), userID, opts)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.Paginated(c, orders, total, opts.Page, opts.Limit)
}

func (h *LiffOrderHandler) Get(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
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

func (h *LiffOrderHandler) UploadSlip(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}

	file, err := c.FormFile("slip")
	if err != nil {
		return response.BadRequest(c, "slip file is required")
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
	objectName := fmt.Sprintf("slips/%s/%s/%s%s", tenantID.String(), id.String(), time.Now().Format("20060102150405"), ext)
	url, err := h.storage.UploadFile(c.Context(), objectName, src, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return response.InternalError(c, err)
	}

	// Update order with slip URL
	if err := h.orderSvc.UploadSlip(c.Context(), *tenantID, id, userID, url); err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true, "slip_image_url": url, "status": "pending_verification"})
}

func (h *LiffOrderHandler) GetActivePending(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)

	order, err := h.orderSvc.GetActivePendingOrder(c.Context(), *tenantID, userID)
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.OK(c, order)
}

func (h *LiffOrderHandler) Cancel(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}

	err = h.orderSvc.CancelOrder(c.Context(), *tenantID, id, userID)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{"success": true})
}

func (h *LiffOrderHandler) UpdatePaymentMethod(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid order id")
	}

	var req struct {
		PaymentMethodID uuid.UUID `json:"payment_method_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	order, err := h.orderSvc.UpdatePaymentMethod(c.Context(), *tenantID, id, userID, req.PaymentMethodID)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, order)
}

