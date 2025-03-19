package commands

import (
	"fmt"
	"lenavire/internal/ledger/application/ports"
	"lenavire/internal/ledger/domain/entities"
	"lenavire/internal/ledger/domain/valuesobjects"
)

type AddExpenseCommand struct {
	amount int
}

func NewAddExpenseCommand(amount int) AddExpenseCommand {
	return AddExpenseCommand{
		amount: amount,
	}
}

type AddExpenseCommandHandler struct {
	expenseRepository     ports.ExpenseRepository
	idProvider            ports.IdProvider
	dateProvider          ports.DateProvider
	ledgerActivityChannel ports.LedgerActivityChannel
}

func NewAddExpenseCommandHandler(
	expenseRepository ports.ExpenseRepository,
	idProvider ports.IdProvider,
	dateProvider ports.DateProvider,
	ledgerActivityChannel ports.LedgerActivityChannel,
) *AddExpenseCommandHandler {
	return &AddExpenseCommandHandler{
		expenseRepository:     expenseRepository,
		idProvider:            idProvider,
		dateProvider:          dateProvider,
		ledgerActivityChannel: ledgerActivityChannel,
	}
}

func (h *AddExpenseCommandHandler) Execute(addExpenseCommand AddExpenseCommand) error {

	amount, err := valuesobjects.NewAmount(addExpenseCommand.amount)

	if err != nil {
		return err
	}

	expense := &entities.Expense{
		Amount: amount,
		Date:   h.dateProvider.Now().Format("2006-01-02"),
		Id:     h.idProvider.Generate(),
	}

	err = h.expenseRepository.Create(expense)

	if err != nil {
		return fmt.Errorf("failed to create expense: %w", err)
	}

	err = h.ledgerActivityChannel.Send("expense-added")

	if err != nil {
		return fmt.Errorf("failed to send expense-added message: %w", err)
	}

	return nil
}
