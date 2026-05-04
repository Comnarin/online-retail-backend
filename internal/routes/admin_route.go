package routes

import (
	"github.com/gofiber/fiber/v2"
	registry "github.com/retail/backend/cmd/app/register"
	"github.com/retail/backend/internal/handler/middleware"
	jwtpkg "github.com/retail/backend/pkg/jwt"
)

func AdminRoutes(app fiber.Router, jwtMgr *jwtpkg.Manager, handlers *registry.Handlers) {
	adminGroup := app.Group("/api/v1/admin",
		middleware.Auth(jwtMgr),
		middleware.RequireRole("tenant_admin", "superadmin"),
	)
	
	// Dashboard & Appearance
	dashboardH := handlers.Dashboard
	adminGroup.Get("/dashboard/stats", dashboardH.Stats)
	adminGroup.Get("/dashboard/recent-orders", dashboardH.RecentOrders)
	
	appearanceH := handlers.Appearance
	adminGroup.Get("/appearance", appearanceH.Get)
	adminGroup.Put("/appearance", appearanceH.Update)

	// Products
	productH := handlers.Product
	pg := adminGroup.Group("/products")
	pg.Get("/", productH.List)
	pg.Post("/", productH.Create)
	pg.Get("/categories", productH.ListCategories)
	pg.Post("/categories", productH.CreateCategory)
	pg.Put("/categories/:id", productH.UpdateCategory)
	pg.Delete("/categories/:id", productH.DeleteCategory)
	pg.Get("/:id", productH.Get)
	pg.Put("/:id", productH.Update)
	pg.Delete("/:id", productH.Delete)
	pg.Post("/:id/images", productH.UploadImage)
	pg.Delete("/:id/images/:image_id", productH.DeleteImage)

	// Orders
	orderH := handlers.Order
	og := adminGroup.Group("/orders")
	og.Get("/", orderH.List)
	og.Get("/:id", orderH.Get)
	og.Patch("/:id/status", orderH.UpdateStatus)
	og.Post("/:id/confirm-payment", orderH.ConfirmPayment)
	og.Post("/:id/reject-payment", orderH.RejectPayment)

	// Payment Methods
	pmH := handlers.AdminPaymentMethod
	pmg := adminGroup.Group("/payment-methods")
	pmg.Get("/", pmH.List)
	pmg.Put("/:id", pmH.Update)
	pmg.Post("/:id/qr-code", pmH.UploadQRCode)

	// Customers
	customerH := handlers.Customer
	cg := adminGroup.Group("/customers")
	cg.Get("/", customerH.List)
	cg.Get("/:id", customerH.Get)

	// Coupons
	couponH := handlers.Coupon
	couponG := adminGroup.Group("/coupons")
	couponG.Get("/", couponH.List)
	couponG.Post("/", couponH.Create)
	couponG.Get("/:id", couponH.Get)
	couponG.Put("/:id", couponH.Update)
	couponG.Patch("/:id/toggle", couponH.Toggle)
	couponG.Delete("/:id", couponH.Delete)

	// Membership
	membershipH := handlers.Membership
	mg := adminGroup.Group("/membership")
	mg.Get("/tiers", membershipH.ListTiers)
	mg.Post("/tiers", membershipH.CreateTier)
	mg.Put("/tiers/:id", membershipH.UpdateTier)
	mg.Delete("/tiers/:id", membershipH.DeleteTier)
	mg.Get("/transactions", membershipH.ListTransactions)

	// Transactions
	txH := handlers.Transaction
	adminGroup.Get("/transactions", txH.List)
}
