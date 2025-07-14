package handlers

import (
	"encoding/json"
	"fmt"
	commands "lenavire/internal/ledger/application/commands/receive_payment"
	"lenavire/internal/ledger/domain/valuesobjects"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
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

	// Get Stripe signature from header
	signature := c.Get("Stripe-Signature")

	// TODO: Replace with your webhook secret from Stripe dashboard
	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	// Verify webhook signature
	event, err := webhook.ConstructEvent(body, signature, webhookSecret)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid webhook signature"})
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

	// Extract pseudo from metadata or customer name
	pseudo := ""
	if session.Metadata != nil {
		if p, ok := session.Metadata["pseudo"]; ok {
			pseudo = p
		}
	}

	// If pseudo not in metadata, use customer name
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
