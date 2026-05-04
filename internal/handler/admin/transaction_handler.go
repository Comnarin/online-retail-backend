package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/retail/backend/internal/handler"
	"github.com/retail/backend/internal/model"
	service "github.com/retail/backend/internal/services"
	"github.com/retail/backend/pkg/response"
)

type TransactionHandler struct {
	svc service.ITransactionService
}

func NewTransactionHandler(svc service.ITransactionService) handler.ITransactionHandler {
	return &TransactionHandler{svc: svc}
}

func (h *TransactionHandler) List(c *fiber.Ctx) error {
	tenantID := c.Locals("tenant_id").(*uuid.UUID)
	opts := model.ListOptions{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", 20),
		Search: c.Query("search"),
	}

	txs, total, err := h.svc.List(c.Context(), *tenantID, opts)
	if err != nil {
		return response.InternalError(c, err)
	}

	return response.Paginated(c, txs, total, opts.Page, opts.Limit)
}
