package ports

import "lenavire/internal/ledger/domain/entities"

type ExpenseRepository interface {
	Create(expense *entities.Expense) error
}
