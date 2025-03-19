package queries

import (
	"encoding/json"

	"gorm.io/gorm"
)

type GetLedgerQuery struct{}

func NewGetLedgerQuery() GetLedgerQuery {
	return GetLedgerQuery{}
}

type GetLedgerQueryHandler struct {
	db *gorm.DB
}

type GetLedgerQueryResult struct {
	TotalRevenue    int              `json:"totalRevenue"`
	TotalExpenses   int              `json:"totalExpenses"`
	TotalReceived   int              `json:"totalReceived"`
	Payments        []Payment        `json:"payments"`
	TopContributors []TopContributor `json:"topContributors"`
}

type Payment struct {
	Amount      int    `json:"amount"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PaymentType string `json:"paymentType"`
}

type TopContributor struct {
	Amount int    `json:"amount"`
	Name   string `json:"name"`
}

func NewGetLedgerQueryHandler(db *gorm.DB) *GetLedgerQueryHandler {
	return &GetLedgerQueryHandler{db: db}
}

func (h *GetLedgerQueryHandler) Execute(query GetLedgerQuery) (*GetLedgerQueryResult, error) {
	var result struct {
		TotalExpenses   int             `json:"total_expenses"`
		TotalReceived   int             `json:"total_received"`
		Payments        json.RawMessage `json:"payments"`
		TopContributors json.RawMessage `json:"top_contributors"`
	}

	sqlQuery := `
		WITH payment_stats AS (
			SELECT 
				COALESCE(SUM(amount), 0) as total_received,
				COALESCE(
					(
						SELECT json_agg(
							json_build_object(
								'amount', amount,
								'name', name,
								'email', email,
								'payment_type', payment_type
							)
						)
						FROM (
							SELECT *
							FROM payments
							ORDER BY date DESC
							LIMIT 3
						) recent
					),
					'[]'::json
				) as payments,
				COALESCE(
					(
						SELECT json_agg(
							json_build_object(
								'amount', total_amount,
								'name', name
							)
						)
						FROM (
							SELECT 
								email,
								SUM(amount) as total_amount,
								MAX(name) as name
							FROM payments
							WHERE date_trunc('month', date) = date_trunc('month', CURRENT_DATE)
							GROUP BY email
							ORDER BY total_amount DESC
							LIMIT 20
						) top
					),
					'[]'::json
				) as top_contributors
			FROM payments
		),
		expense_stats AS (
			SELECT COALESCE(SUM(amount), 0) as total_expenses
			FROM expenses
		)
		SELECT 
			e.total_expenses,
			p.total_received,
			p.payments,
			p.top_contributors
		FROM payment_stats p
		CROSS JOIN expense_stats e`

	err := h.db.Raw(sqlQuery).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	var payments []Payment
	var topContributors []TopContributor

	if len(result.Payments) == 0 {
		payments = []Payment{}
	} else {
		if err := json.Unmarshal(result.Payments, &payments); err != nil {
			return nil, err
		}
	}

	if len(result.TopContributors) == 0 {
		topContributors = []TopContributor{}
	} else {
		if err := json.Unmarshal(result.TopContributors, &topContributors); err != nil {
			return nil, err
		}
	}

	return &GetLedgerQueryResult{
		TotalRevenue:    result.TotalReceived - result.TotalExpenses,
		TotalExpenses:   result.TotalExpenses,
		TotalReceived:   result.TotalReceived,
		Payments:        payments,
		TopContributors: topContributors,
	}, nil
}
