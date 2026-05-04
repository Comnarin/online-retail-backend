package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
	"github.com/retail/backend/pkg/validator"
)

type AuthHandler struct {
	authService service.IAuthService
}

func NewAuthHandler(authService service.IAuthService) handler.IAuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body")
	}

	if err := validator.Struct(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	result, err := h.authService.LoginWithEmail(c.Context(), req)
	if err != nil {
		return response.Unauthorized(c, err.Error())
	}

	// Set cookies for secure persistent session
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "user_role",
		Value:    string(result.User.Role),
		Path:     "/",
		HTTPOnly: false, // Accessible by proxy.ts
		Secure:   true,
		SameSite: "Lax",
	})

	return response.OK(c, result)
}

func (h *AuthHandler) LineLogin(c *fiber.Ctx) error {
	var req model.LineLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "invalid body")
	}

	if err := validator.Struct(req); err != nil {
		return response.BadRequest(c, err.Error())
	}

	result, err := h.authService.LoginWithLine(c.Context(), req)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Set cookies for secure persistent session (Required for Server Components)
	c.Cookie(&fiber.Cookie{
		Name:     "liff_auth_token",
		Value:    result.Token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	// Compatibility cookie
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "liff_user_role",
		Value:    "customer",
		Path:     "/",
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "user_role",
		Value:    "customer",
		Path:     "/",
		HTTPOnly: false,
		Secure:   true,
		SameSite: "Lax",
	})

	return response.OK(c, result)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Clear all possible auth cookies
	cookieNames := []string{"auth_token", "liff_auth_token", "user_role", "liff_user_role"}
	
	for _, name := range cookieNames {
		c.Cookie(&fiber.Cookie{
			Name:     name,
			Value:    "",
			Path:     "/",
			Expires:  time.Now().Add(-1 * time.Hour), // Expire immediately
			HTTPOnly: true,
			Secure:   true,
			SameSite: "Lax",
		})
	}

	return response.OK(c, fiber.Map{"success": true})
}

func (h *AuthHandler) DevLiffLogin(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("tenantId"))
	if err != nil {
		return response.BadRequest(c, "invalid tenant id")
	}
	result, err := h.authService.DevLiffLogin(c.Context(), id)
	if err != nil {
		return response.BadRequest(c, err.Error())
	}

	// Set cookie for browser-based testing convenience
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    result.Token,
		Path:     "/",
		HTTPOnly: false,
		SameSite: "Lax",
	})

	return response.OK(c, result)
}
