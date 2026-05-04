package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/retail/backend/internal/domain"
	"github.com/retail/backend/pkg/response"
)

// FeatureGate checks if a feature is enabled for the current tenant.
func RequireFeature(feature string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenant, ok := c.Locals("tenant").(*domain.Tenant)
		if !ok || tenant == nil {
			return response.Forbidden(c)
		}

		enabled := false
		switch feature {
		case "membership":
			enabled = tenant.Features.EnableMembership
		case "coupon":
			enabled = tenant.Features.EnableCoupons
		case "points":
			enabled = tenant.Features.EnablePoints
		case "reviews":
			enabled = tenant.Features.EnableReviews
		case "delivery":
			enabled = tenant.Features.EnableDelivery
		}

		if !enabled {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "feature '" + feature + "' is not enabled for this tenant",
			})
		}

		return c.Next()
	}
}
