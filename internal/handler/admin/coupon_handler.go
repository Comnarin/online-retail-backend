package admin

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

type CouponHandler struct {
	couponSvc service.ICouponService
}

func NewCouponHandler(couponSvc service.ICouponService) handler.ICouponHandler {
	return &CouponHandler{couponSvc: couponSvc}
}

func (h *CouponHandler) List(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	opts := model.ListOptions{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 50),
		Search: c.Query("search"),
	}
	coupons, total, err := h.couponSvc.List(c.Context(), *tenantID, opts)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.Paginated(c, coupons, total, opts.Page, opts.Limit)
}

func (h *CouponHandler) Get(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid coupon id")
	}
	coupon, err := h.couponSvc.GetByID(c.Context(), *tenantID, id)
	if err != nil {
		return response.NotFound(c, "coupon not found")
	}
	return response.OK(c, coupon)
}

func (h *CouponHandler) Create(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	var coupon domain.Coupon
	if err := c.BodyParser(&coupon); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	coupon.TenantID = *tenantID

	if err := validator.Struct(coupon); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.couponSvc.Create(c.Context(), &coupon); err != nil {
		return response.InternalError(c, err)
	}
	return response.Created(c, coupon)
}

func (h *CouponHandler) Update(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid coupon id")
	}
	var coupon domain.Coupon
	if err := c.BodyParser(&coupon); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	coupon.ID = id

	if err := validator.Struct(coupon); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.couponSvc.Update(c.Context(), *tenantID, &coupon); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, coupon)
}

func (h *CouponHandler) Toggle(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid coupon id")
	}
	coupon, err := h.couponSvc.Toggle(c.Context(), *tenantID, id)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, coupon)
}

func (h *CouponHandler) Delete(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid coupon id")
	}
	if err := h.couponSvc.Delete(c.Context(), *tenantID, id); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"deleted": true})
}
