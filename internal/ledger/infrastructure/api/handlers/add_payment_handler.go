package handlers

import (
	commands "lenavire/internal/ledger/application/commands/receive_payment"
	"lenavire/internal/ledger/domain/valuesobjects"
	"lenavire/internal/ledger/infrastructure/api/dto"

	"github.com/gofiber/fiber/v2"
)

type AddPaymentHandler struct {
	CommandHandler *commands.ReceivedPaymentCommandHandler
}

func NewAddPaymentHandler(commandHandler *commands.ReceivedPaymentCommandHandler) *AddPaymentHandler {
	return &AddPaymentHandler{CommandHandler: commandHandler}
}

func (h *AddPaymentHandler) AddPayment(c *fiber.Ctx) error {
	var req dto.AddPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	paymentType := valuesobjects.OneTime

	if req.PaymentType == "recurring" {
		paymentType = valuesobjects.Recurring
	}

	command := commands.NewReceivePaymentCommand(req.Amount, req.Name, req.Email, paymentType)
	err := h.CommandHandler.Execute(command)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "payment added successfully"})
}
