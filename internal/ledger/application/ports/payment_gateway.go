package ports

import "lenavire/internal/ledger/domain/entities"

type CheckoutSession struct {
	ID  string
	URL string
}

type PaymentGateway interface {
	CreateCheckoutSession(amount int, subscriptionType entities.SubscriptionType) (*CheckoutSession, error)
}