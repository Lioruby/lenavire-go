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
				'id', id,
				'amount', amount,
				'date', date,
				'payment_type', payment_type,
				'name', name,
				'email', email
			)
		) as payments
		FROM payments
		ORDER BY date DESC
		LIMIT 100
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
