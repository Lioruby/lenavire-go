package handlers

import (
	"lenavire/internal/ledger/application/queries"

	"github.com/gofiber/fiber/v2"
)

type GetExpensesHistoryHandler struct {
	QueryHandler *queries.GetExpensesHistoryQueryHandler
}

func NewGetExpensesHistoryHandler(queryHandler *queries.GetExpensesHistoryQueryHandler) *GetExpensesHistoryHandler {
	return &GetExpensesHistoryHandler{QueryHandler: queryHandler}
}

func (h *GetExpensesHistoryHandler) GetExpensesHistory(c *fiber.Ctx) error {
	query := queries.NewGetExpensesHistoryQuery()
	result, err := h.QueryHandler.Execute(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
