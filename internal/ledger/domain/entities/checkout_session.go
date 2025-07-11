package entities

import "lenavire/internal/ledger/domain/valuesobjects"

type SubscriptionType string

const (
	OneTime      SubscriptionType = "one_time"
	Subscription SubscriptionType = "subscription"
)

type CheckoutSession struct {
	Id               string
	Amount           valuesobjects.Amount
	SubscriptionType SubscriptionType
	StripeSessionId  string
	Status           string
	CreatedAt        string
}

func NewCheckoutSession(id string, amount valuesobjects.Amount, subscriptionType SubscriptionType, stripeSessionId string, status string, createdAt string) *CheckoutSession {
	return &CheckoutSession{
		Id:               id,
		Amount:           amount,
		SubscriptionType: subscriptionType,
		StripeSessionId:  stripeSessionId,
		Status:           status,
		CreatedAt:        createdAt,
	}
}