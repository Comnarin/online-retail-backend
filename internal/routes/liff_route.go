package routes

import (
	"github.com/gofiber/fiber/v2"
	registry "github.com/retail/backend/cmd/app/register"
	"github.com/retail/backend/internal/handler/middleware"
	jwtpkg "github.com/retail/backend/pkg/jwt"
)

func LiffRoutes(app fiber.Router, jwtMgr *jwtpkg.Manager, handlers *registry.Handlers) {
	// Base for all LIFF APIs
	liff := app.Group("/api/v1/liff", 
		middleware.Auth(jwtMgr),
		middleware.RequireRole("customer", "tenant_admin", "superadmin"),
	)

	// All routes are now protected. 
	// Handlers will resolve tenant_id from the Auth session.
	{
		productH := handlers.LiffProduct
		liff.Get("/products", productH.List)
		liff.Get("/products/:id", productH.Get)
		liff.Get("/categories", productH.ListCategories)
		
		liff.Get("/appearance", handlers.Appearance.Get)

		orderH := handlers.LiffOrder
		liff.Post("/orders", orderH.Create)
		liff.Get("/orders/active-pending", orderH.GetActivePending)
		liff.Get("/orders", orderH.List)
		liff.Get("/orders/:id", orderH.Get)
		liff.Post("/orders/:id/upload-slip", orderH.UploadSlip)
		liff.Post("/orders/:id/cancel", orderH.Cancel)
		liff.Post("/orders/:id/payment-method", orderH.UpdatePaymentMethod)

		cartH := handlers.LiffCart
		liff.Get("/cart", cartH.Get)
		liff.Post("/cart/items", cartH.AddItem)
		liff.Patch("/cart/items/:id", cartH.UpdateItem)
		liff.Patch("/cart/items/:id/select", cartH.ToggleSelection)
		liff.Delete("/cart/items/:id", cartH.RemoveItem)
		liff.Delete("/cart/clear", cartH.Clear)

		addrH := handlers.LiffAddress
		liff.Get("/addresses", addrH.List)
		liff.Post("/addresses", addrH.Create)
		liff.Put("/addresses/:id", addrH.Update)
		liff.Delete("/addresses/:id", addrH.Delete)

		payH := handlers.LiffPayment
		liff.Get("/payment-methods", payH.ListMethods)
		liff.Get("/customer/payment-methods", payH.ListSavedMethods)
		liff.Post("/customer/payment-methods", payH.CreateSavedMethod)
		liff.Put("/customer/payment-methods/:id", payH.UpdateSavedMethod)
		liff.Delete("/customer/payment-methods/:id", payH.DeleteSavedMethod)
		liff.Post("/customer/payment-methods/:id/default", payH.SetDefaultMethod)

		liff.Get("/profile", handlers.LiffCustomer.GetProfile)
	}
}
