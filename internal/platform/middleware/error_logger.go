package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func ErrorLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue to next handler
		err := c.Next()

		// Check if there was an error response (status >= 400)
		status := c.Response().StatusCode()
		if status >= 400 {
			apiErr := c.Locals("api_error")
			if apiErr != nil {
				log.Error().
					Str("method", c.Method()).
					Str("path", c.Path()).
					Int("status", status).
					Interface("error", apiErr).
					Msg("API Error Response")
			} else if err != nil {
				// Fallback for unhandled Fiber errors
				log.Error().
					Str("method", c.Method()).
					Str("path", c.Path()).
					Int("status", status).
					Err(err).
					Msg("Unhandled API Error")
			}
		}

		return err
	}
}
