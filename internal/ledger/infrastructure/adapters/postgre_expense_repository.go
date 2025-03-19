package adapters

import (
	"lenavire/internal/ledger/domain/entities"
	"lenavire/internal/ledger/infrastructure/database/mappers"

	"gorm.io/gorm"
)

type PostgreExpenseRepository struct {
	db     *gorm.DB
	mapper *mappers.ExpenseMapper
}

func NewPostgreExpenseRepository(db *gorm.DB) *PostgreExpenseRepository {
	mapper := mappers.NewExpenseMapper()

	return &PostgreExpenseRepository{db: db, mapper: mapper}
}

func (r *PostgreExpenseRepository) Create(expense *entities.Expense) error {
	expenseModel := r.mapper.ToPersistence(expense)

	return r.db.Create(expenseModel).Error
}
