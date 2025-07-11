package commands

import (
	"fmt"
	"lenavire/internal/ledger/application/ports"
	"lenavire/internal/ledger/domain/entities"
	"lenavire/internal/ledger/domain/valuesobjects"
)

type CreateCheckoutCommand struct {
	Amount           int
	SubscriptionType entities.SubscriptionType
}

func NewCreateCheckoutCommand(amount int, subscriptionType entities.SubscriptionType) CreateCheckoutCommand {
	return CreateCheckoutCommand{
		Amount:           amount,
		SubscriptionType: subscriptionType,
	}
}

type CreateCheckoutCommandHandler struct {
	paymentGateway ports.PaymentGateway
}

func NewCreateCheckoutCommandHandler(
	paymentGateway ports.PaymentGateway,
) *CreateCheckoutCommandHandler {
	return &CreateCheckoutCommandHandler{
		paymentGateway: paymentGateway,
	}
}

type CreateCheckoutResult struct {
	CheckoutURL string
	SessionId   string
}

func (h *CreateCheckoutCommandHandler) Execute(command CreateCheckoutCommand) (*CreateCheckoutResult, error) {
	_, err := valuesobjects.NewAmount(command.Amount)
	if err != nil {
		return nil, err
	}

	session, err := h.paymentGateway.CreateCheckoutSession(command.Amount, command.SubscriptionType)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return &CreateCheckoutResult{
		CheckoutURL: session.URL,
		SessionId:   session.ID,
	}, nil
}