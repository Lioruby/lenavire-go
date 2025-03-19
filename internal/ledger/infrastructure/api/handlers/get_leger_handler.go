package handlers

import (
	"lenavire/internal/ledger/application/queries"

	"github.com/gofiber/fiber/v2"
)

type GetLedgerHandler struct {
	QueryHandler *queries.GetLedgerQueryHandler
}

func NewGetLedgerHandler(queryHandler *queries.GetLedgerQueryHandler) *GetLedgerHandler {
	return &GetLedgerHandler{QueryHandler: queryHandler}
}

func (h *GetLedgerHandler) GetLedger(c *fiber.Ctx) error {
	query := queries.NewGetLedgerQuery()
	result, err := h.QueryHandler.Execute(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
