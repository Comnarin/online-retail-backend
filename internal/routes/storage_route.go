package routes

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	registry "github.com/retail/backend/cmd/app/register"
)

func StorageRoutes(app fiber.Router, handlers *registry.Handlers) {
	// Public Storage Proxy (to support tunnels/external access)
	storage := app.Group("/products")
	storage.Get("/*", func(c *fiber.Ctx) error {
		path := c.Params("*")
		fmt.Printf("🚀 [StorageProxy] Request for path: %s\n", path)
		
		if path == "" {
			return c.Status(404).SendString("Not Found")
		}

		if handlers.Storage == nil {
			return c.Status(500).SendString("Storage not initialized")
		}

		// Prepend 'products/' if the path doesn't already start with it, 
		// but since we want the bucket structure to be products/tenantID/filename,
		// and the route is /products/*, the '*' will be 'tenantID/filename'.
		// We need to fetch 'products/tenantID/filename' from MinIO.
		fullPath := "products/" + path
		obj, info, err := handlers.Storage.GetObject(c.Context(), fullPath)
		if err != nil {
			return c.Status(404).SendString("File not found")
		}

		c.Set("Content-Type", info.ContentType)
		c.Set("Content-Length", fmt.Sprintf("%d", info.Size))
		c.Set("Cache-Control", "public, max-age=1800") // 30 minutes cache

		return c.SendStream(obj)
	})
}
