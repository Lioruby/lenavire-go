package adapters

import (
	"fmt"
	"lenavire/internal/ledger/application/ports"
	"lenavire/internal/ledger/domain/entities"
	"os"

	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type StripePaymentGateway struct {
	secretKey string
}

func NewStripePaymentGateway(secretKey string) *StripePaymentGateway {
	stripe.Key = secretKey
	return &StripePaymentGateway{
		secretKey: secretKey,
	}
}

func (s *StripePaymentGateway) CreateCheckoutSession(amount int, subscriptionType entities.SubscriptionType) (*ports.CheckoutSession, error) {
	baseFrontendURL := os.Getenv("BASE_FRONTEND_URL")
	if baseFrontendURL == "" {
		baseFrontendURL = "http://localhost:5173"
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("eur"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Le Navire Contribution"),
					},
					UnitAmount: stripe.Int64(int64(amount)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:             stripe.String(s.getModeFromSubscriptionType(subscriptionType)),
		SuccessURL:       stripe.String(baseFrontendURL + "?payment=success"),
		CancelURL:        stripe.String(baseFrontendURL + "?payment=failed"),
		CustomerCreation: stripe.String("always"),
	}

	if subscriptionType == entities.Subscription {
		params.LineItems[0].PriceData.Recurring = &stripe.CheckoutSessionLineItemPriceDataRecurringParams{
			Interval: stripe.String("month"),
		}
	}

	sess, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create stripe checkout session: %w", err)
	}

	return &ports.CheckoutSession{
		ID:  sess.ID,
		URL: sess.URL,
	}, nil
}

func (s *StripePaymentGateway) getModeFromSubscriptionType(subscriptionType entities.SubscriptionType) string {
	switch subscriptionType {
	case entities.OneTime:
		return "payment"
	case entities.Subscription:
		return "subscription"
	default:
		return "payment"
	}
}
