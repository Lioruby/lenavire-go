package adapters

import "lenavire/internal/ledger/domain/entities"

type InMemoryExpenseRepository struct {
	Expenses []*entities.Expense
}

func NewInMemoryExpenseRepository() *InMemoryExpenseRepository {
	return &InMemoryExpenseRepository{
		Expenses: []*entities.Expense{},
	}
}

func (r *InMemoryExpenseRepository) Create(expense *entities.Expense) error {
	r.Expenses = append(r.Expenses, expense)
	return nil
}
