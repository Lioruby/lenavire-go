package dto

type AddPaymentRequest struct {
	Amount      int    `json:"amount"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PaymentType string `json:"payment_type"`
}
