package dto

type AddExpenseRequest struct {
	Amount        int    `json:"amount"`
	OperationType string `json:"operation_type"`
}
