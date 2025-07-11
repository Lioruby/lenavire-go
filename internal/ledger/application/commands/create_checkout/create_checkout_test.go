package commands

import (
	"errors"
	"lenavire/internal/ledger/application/ports"
	"lenavire/internal/ledger/domain/entities"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type FakePaymentGateway struct {
	sessionToReturn *ports.CheckoutSession
	errorToReturn   error
}

func (f *FakePaymentGateway) CreateCheckoutSession(amount int, subscriptionType entities.SubscriptionType) (*ports.CheckoutSession, error) {
	if f.errorToReturn != nil {
		return nil, f.errorToReturn
	}
	return f.sessionToReturn, nil
}

func NewFakePaymentGateway() *FakePaymentGateway {
	return &FakePaymentGateway{
		sessionToReturn: &ports.CheckoutSession{
			ID:  "session_123",
			URL: "https://checkout.example.com/session_123",
		},
	}
}

var _ = Describe("CreateCheckoutCommandHandler", func() {
	var (
		handler        *CreateCheckoutCommandHandler
		paymentGateway *FakePaymentGateway
	)

	BeforeEach(func() {
		paymentGateway = NewFakePaymentGateway()
		handler = NewCreateCheckoutCommandHandler(paymentGateway)
	})

	Describe("Execute", func() {
		Context("when creating a one-time checkout", func() {
			It("should create a checkout session successfully", func() {
				command := NewCreateCheckoutCommand(1000, entities.OneTime)

				result, err := handler.Execute(command)

				Expect(err).To(BeNil())
				Expect(result.SessionId).To(Equal("session_123"))
				Expect(result.CheckoutURL).To(Equal("https://checkout.example.com/session_123"))
			})
		})

		Context("when creating a subscription checkout", func() {
			It("should create a checkout session successfully", func() {
				command := NewCreateCheckoutCommand(2000, entities.Subscription)

				result, err := handler.Execute(command)

				Expect(err).To(BeNil())
				Expect(result.SessionId).To(Equal("session_123"))
				Expect(result.CheckoutURL).To(Equal("https://checkout.example.com/session_123"))
			})
		})

		Context("when amount is invalid", func() {
			It("should return an error", func() {
				command := NewCreateCheckoutCommand(-100, entities.OneTime)

				result, err := handler.Execute(command)

				Expect(result).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})

		Context("when payment gateway fails", func() {
			It("should return an error", func() {
				paymentGateway.errorToReturn = errors.New("payment gateway error")
				command := NewCreateCheckoutCommand(1000, entities.OneTime)

				result, err := handler.Execute(command)

				Expect(result).To(BeNil())
				Expect(err).To(MatchError(ContainSubstring("failed to create checkout session")))
			})
		})
	})
})

func TestCommands(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Create Checkout Command Suite")
}
