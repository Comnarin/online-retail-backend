package routes

import (
	"github.com/gofiber/fiber/v2"
	registry "github.com/retail/backend/cmd/app/register"
	jwtpkg "github.com/retail/backend/pkg/jwt"
)

func RegisterRoutes(
	app fiber.Router,
	jwtMgr *jwtpkg.Manager,
	handlers *registry.Handlers,
) {
	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	app.Get("/api/routes", func(c *fiber.Ctx) error {
		return c.JSON(c.App().Stack())
	})

	AuthRoutes(app, handlers)
	StorageRoutes(app, handlers)
	SuperAdminRoutes(app, jwtMgr, handlers)
	AdminRoutes(app, jwtMgr, handlers)
	LiffRoutes(app, jwtMgr, handlers)
}
