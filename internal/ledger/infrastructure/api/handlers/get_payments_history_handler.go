package handlers

import (
	"lenavire/internal/ledger/application/queries"

	"github.com/gofiber/fiber/v2"
)

type GetPaymentsHistoryHandler struct {
	QueryHandler *queries.GetPaymentsHistoryQueryHandler
}

func NewGetPaymentsHistoryHandler(queryHandler *queries.GetPaymentsHistoryQueryHandler) *GetPaymentsHistoryHandler {
	return &GetPaymentsHistoryHandler{QueryHandler: queryHandler}
}

func (h *GetPaymentsHistoryHandler) GetPaymentsHistory(c *fiber.Ctx) error {
	query := queries.NewGetPaymentsHistoryQuery()
	result, err := h.QueryHandler.Execute(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
