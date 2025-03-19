package handlers

import (
	"lenavire/internal/ledger/application/commands/receive_payment"
	"lenavire/internal/ledger/domain/valuesobjects"
	"lenavire/internal/ledger/infrastructure/api/dto"

	"github.com/gofiber/fiber/v2"
)

type ReceivePaymentHandler struct {
	CommandHandler *commands.ReceivedPaymentCommandHandler
}

func NewReceivePaymentHandler(handler *commands.ReceivedPaymentCommandHandler) *ReceivePaymentHandler {
	return &ReceivePaymentHandler{CommandHandler: handler}
}

func (h *ReceivePaymentHandler) ReceivePayment(c *fiber.Ctx) error {
	var req dto.StripeWebhookRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	paymentType := valuesobjects.OneTime

	if len(req.Data.Object.CustomFields) > 1 {
		if req.Data.Object.CustomFields[1].Dropdown.Value == "one-time" {
			paymentType = valuesobjects.OneTime
		} else if req.Data.Object.CustomFields[1].Dropdown.Value == "recurring" {
			paymentType = valuesobjects.Recurring
		}
	}

	command := commands.NewReceivePaymentCommand(
		req.Data.Object.AmountTotal/100,
		req.Data.Object.CustomerDetails.Name,
		req.Data.Object.CustomerDetails.Email,
		paymentType,
	)

	err := h.CommandHandler.Execute(command)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Payment received successfully"})
}
