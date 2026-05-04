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

type MembershipHandler struct {
	membershipSvc service.IMembershipService
}

func NewMembershipHandler(membershipSvc service.IMembershipService) handler.IMembershipHandler {
	return &MembershipHandler{
		membershipSvc: membershipSvc,
	}
}

func (h *MembershipHandler) ListTiers(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	tiers, err := h.membershipSvc.ListTiers(c.Context(), *tenantID)
	if err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, tiers)
}

func (h *MembershipHandler) CreateTier(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	var tier domain.MembershipTier
	if err := c.BodyParser(&tier); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	tier.TenantID = *tenantID

	if err := validator.Struct(tier); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.membershipSvc.CreateTier(c.Context(), &tier); err != nil {
		return response.InternalError(c, err)
	}
	return response.Created(c, tier)
}

func (h *MembershipHandler) UpdateTier(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tier id")
	}
	var tier domain.MembershipTier
	if err := c.BodyParser(&tier); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	tier.ID = id

	if err := validator.Struct(tier); err != nil {
		return response.BadRequest(c, err.Error())
	}

	if err := h.membershipSvc.UpdateTier(c.Context(), *tenantID, &tier); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, tier)
}

func (h *MembershipHandler) DeleteTier(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return response.BadRequest(c, "invalid tier id")
	}
	if err := h.membershipSvc.DeleteTier(c.Context(), *tenantID, id); err != nil {
		return response.InternalError(c, err)
	}
	return response.OK(c, fiber.Map{"deleted": true})
}
func (h *MembershipHandler) ListTransactions(c *fiber.Ctx) error {
	tenantID, ok := c.Locals("tenant_id").(*uuid.UUID)
	if !ok || tenantID == nil {
		return response.BadRequest(c, "tenant context required")
	}

	userIDStr := c.Query("user_id")
	var userID *uuid.UUID
	if userIDStr != "" {
		uid, err := uuid.Parse(userIDStr)
		if err == nil {
			userID = &uid
		}
	}

	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	txs, total, err := h.membershipSvc.ListTransactions(c.Context(), *tenantID, userID, model.ListOptions{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.OK(c, fiber.Map{
		"data":  txs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
