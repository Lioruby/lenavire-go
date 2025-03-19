package entities

import "lenavire/internal/ledger/domain/valuesobjects"

type Payment struct {
	Id          string
	Amount      valuesobjects.Amount
	Name        string
	Email       string
	Date        string
	PaymentType valuesobjects.PaymentType
}

func NewPayment(
	id string,
	amount valuesobjects.Amount,
	name string,
	email string,
	date string,
	paymentType valuesobjects.PaymentType,
) *Payment {
	return &Payment{
		Id:          id,
		Amount:      amount,
		Name:        name,
		Email:       email,
		Date:        date,
		PaymentType: paymentType,
	}
}
