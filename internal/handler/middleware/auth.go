package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	jwtpkg "github.com/retail/backend/pkg/jwt"
	"github.com/retail/backend/pkg/response"
)

func Auth(jwtMgr *jwtpkg.Manager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string
		
		// 1. Check Authorization Header
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 2. Check Bridge Header
		if token == "" {
			token = c.Get("x-liff-bridge")
		}

		// 3. Cookie Fallback
		if token == "" {
			token = c.Cookies("liff_auth_token")
		}
		if token == "" {
			token = c.Cookies("auth_token")
		}

		if token == "" {
			return response.Unauthorized(c, "Authentication required (Header or Cookie)")
		}

		claims, err := jwtMgr.Parse(token)
		if err != nil {
			return response.Unauthorized(c, "Invalid or Expired Token: "+err.Error())
		}

		// Check Redis blacklist
		if jwtMgr.IsBlacklisted(c.Context(), claims.ID) {
			return response.Unauthorized(c)
		}

		// Inject into locals
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("jti", claims.ID)

		// Tenant ID Handling
		tenantID := claims.TenantID

		// If superadmin provides X-Tenant-ID, use it to impersonate
		if claims.Role == "superadmin" {
			if tHeader := c.Get("X-Tenant-ID"); tHeader != "" {
				if tid, err := uuid.Parse(tHeader); err == nil {
					tenantID = &tid
				}
			}
		}
		c.Locals("tenant_id", tenantID)

		return c.Next()
	}
}

// RequireRole gates a route to specific roles
func RequireRole(roles ...string) fiber.Handler {
	roleSet := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		roleSet[r] = struct{}{}
	}
	return func(c *fiber.Ctx) error {
		role, ok := c.Locals("role").(string)
		if !ok {
			return response.Unauthorized(c)
		}
		if _, allowed := roleSet[role]; !allowed {
			return response.Forbidden(c)
		}
		return c.Next()
	}
}
