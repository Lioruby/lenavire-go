package handlers

import (
	"lenavire/internal/ledger/application/queries"

	"github.com/gofiber/fiber/v2"
)

type GetContributorsRankHandler struct {
	QueryHandler *queries.GetContributorsRankQueryHandler
}

func NewGetContributorsRankHandler(queryHandler *queries.GetContributorsRankQueryHandler) *GetContributorsRankHandler {
	return &GetContributorsRankHandler{QueryHandler: queryHandler}
}

func (h *GetContributorsRankHandler) GetContributorsRank(c *fiber.Ctx) error {
	query := queries.NewGetContributorsRankQuery()
	result, err := h.QueryHandler.Execute(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
