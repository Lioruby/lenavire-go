package handlers

import (
	"lenavire/internal/ledger/application/commands/create_checkout"
	"lenavire/internal/ledger/infrastructure/api/dto"

	"github.com/gofiber/fiber/v2"
)

type CreateCheckoutHandler struct {
	commandHandler *commands.CreateCheckoutCommandHandler
}

func NewCreateCheckoutHandler(commandHandler *commands.CreateCheckoutCommandHandler) *CreateCheckoutHandler {
	return &CreateCheckoutHandler{
		commandHandler: commandHandler,
	}
}

func (h *CreateCheckoutHandler) Handle(c *fiber.Ctx) error {
	var request dto.CreateCheckoutRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	command := commands.NewCreateCheckoutCommand(request.Amount, request.SubscriptionType)
	result, err := h.commandHandler.Execute(command)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"checkout_url": result.CheckoutURL,
		"session_id":   result.SessionId,
	})
}