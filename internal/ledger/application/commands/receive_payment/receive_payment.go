package commands

import (
	"fmt"
	"lenavire/internal/ledger/application/ports"
	"lenavire/internal/ledger/domain/entities"
	"lenavire/internal/ledger/domain/valuesobjects"
)

type ReceivePaymentCommand struct {
	amount      int
	name        string
	email       string
	paymentType valuesobjects.PaymentType
}

func NewReceivePaymentCommand(amount int, name string, email string, paymentType valuesobjects.PaymentType) ReceivePaymentCommand {
	return ReceivePaymentCommand{
		amount:      amount,
		name:        name,
		email:       email,
		paymentType: paymentType,
	}
}

type ReceivedPaymentCommandHandler struct {
	paymentRepository     ports.PaymentRepository
	idProvider            ports.IdProvider
	dateProvider          ports.DateProvider
	expenseRepository     ports.ExpenseRepository
	ledgerActivityChannel ports.LedgerActivityChannel
}

func NewReceivedPaymentCommandHandler(
	paymentRepository ports.PaymentRepository,
	idProvider ports.IdProvider,
	dateProvider ports.DateProvider,
	expenseRepository ports.ExpenseRepository,
	ledgerActivityChannel ports.LedgerActivityChannel,
) *ReceivedPaymentCommandHandler {
	return &ReceivedPaymentCommandHandler{
		paymentRepository:     paymentRepository,
		idProvider:            idProvider,
		dateProvider:          dateProvider,
		expenseRepository:     expenseRepository,
		ledgerActivityChannel: ledgerActivityChannel,
	}
}

func (h *ReceivedPaymentCommandHandler) Execute(command ReceivePaymentCommand) error {
	amount, err := valuesobjects.NewAmount(command.amount)

	if err != nil {
		return err
	}

	payment := entities.NewPayment(
		h.idProvider.Generate(),
		amount,
		command.name,
		command.email,
		h.dateProvider.Now().Format("2006-01-02"),
		command.paymentType,
	)

	err = h.paymentRepository.Create(payment)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	err = h.createTVAExpense(amount.Value)

	if err != nil {
		return fmt.Errorf("failed to create tva expense: %w", err)
	}

	err = h.sendLedgerActivityNotification()

	if err != nil {
		return fmt.Errorf("failed to send ledger activity notification: %w", err)
	}

	return nil
}

func (h *ReceivedPaymentCommandHandler) createTVAExpense(amount int) error {
	tvaValue := (amount * 20) / 100
	tvaAmount, err := valuesobjects.NewAmount(tvaValue)

	if err != nil {
		return err
	}

	expense := entities.NewExpense(
		h.idProvider.Generate(),
		tvaAmount,
		h.dateProvider.Now().Format("2006-01-02"),
	)

	err = h.expenseRepository.Create(expense)

	if err != nil {
		return err
	}

	return nil
}

func (h *ReceivedPaymentCommandHandler) sendLedgerActivityNotification() error {
	err := h.ledgerActivityChannel.Send("payment-received")

	if err != nil {
		return err
	}

	return nil
}
