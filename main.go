package main

import (
	addExpenseCommands "lenavire/internal/ledger/application/commands/add_expense"
	receivePaymentCommands "lenavire/internal/ledger/application/commands/receive_payment"
	"lenavire/internal/ledger/infrastructure/adapters"
	"lenavire/internal/ledger/infrastructure/api"
	"lenavire/internal/ledger/infrastructure/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	/* Repositories */
	paymentRepository := adapters.NewInMemoryPaymentRepository()
	expenseRepository := adapters.NewInMemoryExpenseRepository()

	/* Providers */
	idProvider := adapters.NewStubIdProvider("xxxxx")
	dateProvider := adapters.NewStubDateProvider("2021-01-04")

	/* Channels */
	ledgerActivityChannel := adapters.NewFakeLedgerActivityChannel()

	/* Handlers */
	receivePaymentCommandHandler := receivePaymentCommands.NewReceivedPaymentCommandHandler(
		paymentRepository,
		idProvider,
		dateProvider,
		expenseRepository,
		ledgerActivityChannel,
	)
	receivePaymentHandler := handlers.NewReceivePaymentHandler(receivePaymentCommandHandler)

	addExpenseCommandHandler := addExpenseCommands.NewAddExpenseCommandHandler(
		expenseRepository,
		idProvider,
		dateProvider,
		ledgerActivityChannel,
	)
	addExpenseHandler := handlers.NewAddExpenseHandler(addExpenseCommandHandler)

	/* Routes */
	api.SetupRoutes(app, receivePaymentHandler, addExpenseHandler)

	/* Start server */
	app.Listen(":3000")
}
