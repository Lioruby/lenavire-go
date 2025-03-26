package queries

import (
	"encoding/json"

	"gorm.io/gorm"
)

type GetExpensesHistoryQuery struct{}

type Expense struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	Date   string `json:"date"`
}

type GetExpensesHistoryQueryResult struct {
	Expenses []Expense `json:"expenses"`
}

func NewGetExpensesHistoryQuery() GetExpensesHistoryQuery {
	return GetExpensesHistoryQuery{}
}

type GetExpensesHistoryQueryHandler struct {
	db *gorm.DB
}

func NewGetExpensesHistoryQueryHandler(db *gorm.DB) *GetExpensesHistoryQueryHandler {
	return &GetExpensesHistoryQueryHandler{db: db}
}

func (h *GetExpensesHistoryQueryHandler) Execute(query GetExpensesHistoryQuery) (*GetExpensesHistoryQueryResult, error) {
	var result struct {
		Expenses json.RawMessage `json:"expenses"`
	}

	sqlQuery := `
		SELECT json_agg(
			json_build_object(
				'id', id,
				'amount', amount,
				'date', date
			)
		) as expenses
		FROM expenses
		ORDER BY date DESC
	`

	err := h.db.Raw(sqlQuery).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	var expenses []Expense
	if err := json.Unmarshal(result.Expenses, &expenses); err != nil {
		return nil, err
	}

	return &GetExpensesHistoryQueryResult{Expenses: expenses}, nil
}
