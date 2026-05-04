package routes

import (
	"github.com/gofiber/fiber/v2"
	registry "github.com/retail/backend/cmd/app/register"
)

func AuthRoutes(app fiber.Router,handlers *registry.Handlers) {
	authGroup := app.Group("/api/v1/auth")
	h := handlers.Auth
	authGroup.Post("/login", h.Login)
	authGroup.Post("/line/login", h.LineLogin)
	authGroup.Post("/logout", h.Logout)
	authGroup.Get("/dev-liff/:tenantId", h.DevLiffLogin)
}
