package liff

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type LiffCustomerHandler struct {
	customerSvc   service.ICustomerService
	membershipSvc service.IMembershipService
}

func NewLiffCustomerHandler(
	customerSvc service.ICustomerService,
	membershipSvc service.IMembershipService,
) handler.ILiffCustomerHandler {
	return &LiffCustomerHandler{
		customerSvc:   customerSvc,
		membershipSvc: membershipSvc,
	}
}

func (h *LiffCustomerHandler) GetProfile(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	userID := c.Locals("user_id").(uuid.UUID)

	customer, err := h.customerSvc.GetCustomer(c.Context(), *tenantID, userID)
	if err != nil {
		return response.NotFound(c, "customer not found")
	}

	membership, _ := h.membershipSvc.GetCustomerMembership(c.Context(), *tenantID, userID)

	// Build profile response matching frontend expectations
	profile := fiber.Map{
		"id":                 customer.ID,
		"name":               customer.Name,
		"email":              customer.Email,
		"avatar_url":         customer.Avatar,
		"membership_tier":    "Member", // Default
		"points_balance":     0,        // Default
		"voucher_count":      0,        // Mock
		"favorites_count":    0,        // Mock
	}

	if membership != nil {
		profile["membership_tier"] = membership.Tier.Name
		profile["points_balance"] = membership.Points
	}

	return response.OK(c, profile)
}
