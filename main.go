package main

import (
	"fmt"
	addExpenseCommand "lenavire/internal/ledger/application/commands/add_expense"
	receivePaymentCommand "lenavire/internal/ledger/application/commands/receive_payment"
	getLedgerQuery "lenavire/internal/ledger/application/queries"
	"lenavire/internal/ledger/infrastructure/adapters"
	"lenavire/internal/ledger/infrastructure/api"
	"lenavire/internal/ledger/infrastructure/api/handlers"
	"lenavire/internal/ledger/infrastructure/database"
	"lenavire/internal/ledger/infrastructure/database/schema"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	database.ConnectDB()

	database.DB.AutoMigrate(&schema.PaymentModel{}, &schema.ExpenseModel{})

	app := fiber.New()

	/* Repositories */
	paymentRepository := adapters.NewPostgrePaymentRepository(database.DB)
	expenseRepository := adapters.NewPostgreExpenseRepository(database.DB)

	/* Providers */
	idProvider := adapters.NewUUIDIdProvider()
	dateProvider := adapters.NewRealDateProvider()

	/* Channels */
	ledgerActivityChannel := adapters.NewFakeLedgerActivityChannel()

	/* Handlers */
	receivePaymentCommandHandler := receivePaymentCommand.NewReceivedPaymentCommandHandler(
		paymentRepository,
		idProvider,
		dateProvider,
		expenseRepository,
		ledgerActivityChannel,
	)
	receivePaymentHandler := handlers.NewReceivePaymentHandler(receivePaymentCommandHandler)

	addExpenseCommandHandler := addExpenseCommand.NewAddExpenseCommandHandler(
		expenseRepository,
		idProvider,
		dateProvider,
		ledgerActivityChannel,
	)
	addExpenseHandler := handlers.NewAddExpenseHandler(addExpenseCommandHandler)

	getLedgerQueryHandler := getLedgerQuery.NewGetLedgerQueryHandler(database.DB)
	getLedgerHandler := handlers.NewGetLedgerHandler(getLedgerQueryHandler)

	/* Routes */
	api.SetupRoutes(
		app,
		receivePaymentHandler,
		addExpenseHandler,
		getLedgerHandler,
	)

	/* Start server */
	app.Listen(":3000")
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erreur lors du chargement du fichier .env :", err)
	}
}
