package entities

import "lenavire/internal/ledger/domain/valuesobjects"

type Expense struct {
	Id     string
	Amount valuesobjects.Amount
	Date   string
}

func NewExpense(id string, amount valuesobjects.Amount, date string) *Expense {
	return &Expense{
		Id:     id,
		Amount: amount,
		Date:   date,
	}
}
