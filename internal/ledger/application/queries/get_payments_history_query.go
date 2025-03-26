package queries

import (
	"encoding/json"

	"gorm.io/gorm"
)

type GetPaymentsHistoryQuery struct{}

type PaymentResponse struct {
	ID          string `json:"id"`
	Amount      int    `json:"amount"`
	Date        string `json:"date"`
	PaymentType string `json:"payment_type"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}

func NewGetPaymentsHistoryQuery() GetPaymentsHistoryQuery {
	return GetPaymentsHistoryQuery{}
}

type GetPaymentsHistoryQueryHandler struct {
	db *gorm.DB
}

type GetPaymentsHistoryQueryResult struct {
	Payments []PaymentResponse `json:"payments"`
}

func NewGetPaymentsHistoryQueryHandler(db *gorm.DB) *GetPaymentsHistoryQueryHandler {
	return &GetPaymentsHistoryQueryHandler{db: db}
}

func (h *GetPaymentsHistoryQueryHandler) Execute(query GetPaymentsHistoryQuery) (*GetPaymentsHistoryQueryResult, error) {
	var result struct {
		Payments json.RawMessage `json:"payments"`
	}

	sqlQuery := `
		SELECT json_agg(
			json_build_object(
				'id', p.id,
				'amount', p.amount,
				'date', p.date,
				'payment_type', p.payment_type,
				'name', p.name,
				'email', p.email
			)
		) as payments
		FROM (
			SELECT id, amount, date, payment_type, name, email
			FROM payments
			ORDER BY date DESC
		) p
	`

	err := h.db.Raw(sqlQuery).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	var payments []PaymentResponse
	if err := json.Unmarshal(result.Payments, &payments); err != nil {
		return nil, err
	}

	return &GetPaymentsHistoryQueryResult{Payments: payments}, nil
}
