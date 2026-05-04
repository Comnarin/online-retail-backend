package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func OK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Data:    data,
	})
}

func Paginated(c *fiber.Ctx, data interface{}, total int64, page, limit int) error {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
		Success:    true,
		Data:       data,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	})
}

func BadRequest(c *fiber.Ctx, msg string, err ...error) error {
	var detailedErr error
	if len(err) > 0 {
		detailedErr = err[0]
	}

	// Store for middleware logging
	c.Locals("api_error", msg)
	if detailedErr != nil {
		c.Locals("api_error", msg+" | underlying: "+detailedErr.Error())
	}

	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Error:   msg,
	})
}

func Unauthorized(c *fiber.Ctx, msg ...string) error {
	errMsg := "unauthorized"
	if len(msg) > 0 {
		errMsg = msg[0]
	}
	c.Locals("api_error", errMsg)
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Error:   errMsg,
	})
}

func Forbidden(c *fiber.Ctx) error {
	c.Locals("api_error", "forbidden: insufficient permissions")
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success: false,
		Error:   "forbidden: insufficient permissions",
	})
}

func NotFound(c *fiber.Ctx, msg string) error {
	c.Locals("api_error", "not found: "+msg)
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Error:   msg,
	})
}

func InternalError(c *fiber.Ctx, err error) error {
	if err != nil {
		c.Locals("api_error", "internal server error | underlying: "+err.Error())
	} else {
		c.Locals("api_error", "internal server error")
	}

	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Error:   "internal server error",
	})
}
