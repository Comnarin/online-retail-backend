package routes

import (
	"github.com/gofiber/fiber/v2"
	registry "github.com/retail/backend/cmd/app/register"
	"github.com/retail/backend/internal/handler/middleware"
	jwtpkg "github.com/retail/backend/pkg/jwt"
)

func SuperAdminRoutes(app fiber.Router, jwtMgr *jwtpkg.Manager, handlers *registry.Handlers) {
	saGroup := app.Group("/api/v1/superadmin",
		middleware.Auth(jwtMgr),
		middleware.RequireRole("superadmin"),
	)
	
	h := handlers.Tenant
	
	r := saGroup.Group("/tenants")
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/:id/admins", h.ListAdmins)
	r.Post("/:id/admins", h.CreateAdmin)
	r.Delete("/:id/admins/:adminId", h.DeleteAdmin)
	r.Put("/:id/features", h.UpdateFeatures)
	r.Put("/:id/appearance", h.UpdateAppearance)
	r.Get("/:id", h.Get)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}
