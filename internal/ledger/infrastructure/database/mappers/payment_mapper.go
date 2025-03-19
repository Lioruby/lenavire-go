package mappers

import (
	"lenavire/internal/ledger/domain/entities"
	"lenavire/internal/ledger/domain/valuesobjects"
	"lenavire/internal/ledger/infrastructure/database/schema"
	"log"
	"time"
)

type PaymentMapper struct{}

func NewPaymentMapper() *PaymentMapper {
	return &PaymentMapper{}
}

func (m *PaymentMapper) ToDomain(payment *schema.PaymentModel) *entities.Payment {
	paymentType := valuesobjects.OneTime
	if payment.PaymentType == "one-time" {
		paymentType = valuesobjects.OneTime
	} else if payment.PaymentType == "recurring" {
		paymentType = valuesobjects.Recurring
	}

	amount, err := valuesobjects.NewAmount(payment.Amount)
	if err != nil {
		log.Fatal(err)
	}

	return entities.NewPayment(
		payment.ID,
		amount,
		payment.Name,
		payment.Email,
		payment.Date.Format("2006-01-02"),
		paymentType,
	)
}

func (m *PaymentMapper) ToPersistence(payment *entities.Payment) *schema.PaymentModel {
	date, err := time.Parse("2006-01-02", payment.Date)
	if err != nil {
		log.Fatal(err)
	}

	return &schema.PaymentModel{
		ID:          payment.Id,
		Amount:      payment.Amount.Value,
		Date:        date,
		Name:        payment.Name,
		Email:       payment.Email,
		PaymentType: string(payment.PaymentType),
	}
}
