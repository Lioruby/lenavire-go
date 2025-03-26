package queries

import (
	"encoding/json"

	"gorm.io/gorm"
)

type GetPaymentsHistoryQuery struct{}

func NewGetPaymentsHistoryQuery() GetPaymentsHistoryQuery {
	return GetPaymentsHistoryQuery{}
}

type GetPaymentsHistoryQueryHandler struct {
	db *gorm.DB
}

type GetPaymentsHistoryQueryResult struct {
	Payments []Payment `json:"payments"`
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
			ORDER BY date ASC
		) p
	`

	err := h.db.Raw(sqlQuery).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	var payments []Payment
	if err := json.Unmarshal(result.Payments, &payments); err != nil {
		return nil, err
	}

	return &GetPaymentsHistoryQueryResult{Payments: payments}, nil
}
