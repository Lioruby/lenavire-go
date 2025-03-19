package mappers

import (
	"lenavire/internal/ledger/domain/entities"
	"lenavire/internal/ledger/domain/valuesobjects"
	"lenavire/internal/ledger/infrastructure/database/schema"
	"log"
	"time"
)

type ExpenseMapper struct{}

func NewExpenseMapper() *ExpenseMapper {
	return &ExpenseMapper{}
}

func (m *ExpenseMapper) ToDomain(expense *schema.ExpenseModel) *entities.Expense {
	amount, err := valuesobjects.NewAmount(expense.Amount)
	if err != nil {
		log.Fatal(err)
	}

	return &entities.Expense{
		Id:     expense.ID,
		Amount: amount,
		Date:   expense.Date.Format("2006-01-02"),
	}
}

func (m *ExpenseMapper) ToPersistence(expense *entities.Expense) *schema.ExpenseModel {
	date, err := time.Parse("2006-01-02", expense.Date)
	if err != nil {
		log.Fatal(err)
	}

	return &schema.ExpenseModel{
		ID:     expense.Id,
		Amount: expense.Amount.Value,
		Date:   date,
	}
}
