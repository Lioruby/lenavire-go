package handlers

import (
	"encoding/json"
	"fmt"
	commands "lenavire/internal/ledger/application/commands/receive_payment"
	"lenavire/internal/ledger/domain/valuesobjects"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
)

type ReceiveStripePaymentHandler struct {
	CommandHandler *commands.ReceivedPaymentCommandHandler
}

func NewReceiveStripePaymentHandler(handler *commands.ReceivedPaymentCommandHandler) *ReceiveStripePaymentHandler {
	return &ReceiveStripePaymentHandler{CommandHandler: handler}
}

func (h *ReceiveStripePaymentHandler) Handle(c *fiber.Ctx) error {
	// Get raw body for signature verification
	body := c.Body()

	event := stripe.Event{}

	if err := json.Unmarshal(body, &event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse event"})
	}

	// Handle checkout.session.completed event
	if event.Type != "checkout.session.completed" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "event ignored"})
	}

	// Parse checkout session
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "failed to parse session"})
	}

	// Determine payment type based on mode
	paymentType := valuesobjects.OneTime
	if session.Mode == stripe.CheckoutSessionModeSubscription {
		paymentType = valuesobjects.Recurring
	}

	// Extract pseudo from custom fields
	pseudo := ""
	if len(session.CustomFields) > 0 {
		// Look for the pseudo field in custom fields
		for _, field := range session.CustomFields {
			if field.Key == "Pseudo" && field.Text != nil {
				pseudo = field.Text.Value
				break
			}
		}
	}

	// If pseudo not found in custom fields, use customer name as fallback
	if pseudo == "" && session.CustomerDetails != nil {
		pseudo = session.CustomerDetails.Name
	}

	// Create and execute command
	command := commands.NewReceivePaymentCommand(
		int(session.AmountTotal)/100, // Convert cents to euros
		pseudo,
		session.CustomerDetails.Email,
		paymentType,
	)

	if err := h.CommandHandler.Execute(command); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fmt.Printf("Payment processed: %s (%s) - %d EUR - Type: %s\n",
		pseudo, session.CustomerDetails.Email, session.AmountTotal/100, paymentType)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment received successfully"})
}
