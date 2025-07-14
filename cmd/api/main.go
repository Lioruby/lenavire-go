package main

import (
	"fmt"
	addExpenseCommand "lenavire/internal/ledger/application/commands/add_expense"
	commands "lenavire/internal/ledger/application/commands/create_checkout"
	receivePaymentCommand "lenavire/internal/ledger/application/commands/receive_payment"
	queries "lenavire/internal/ledger/application/queries"
	"lenavire/internal/ledger/infrastructure/adapters"
	"lenavire/internal/ledger/infrastructure/api"
	"lenavire/internal/ledger/infrastructure/api/handlers"
	"lenavire/internal/ledger/infrastructure/database"
	"lenavire/internal/ledger/infrastructure/database/schema"
	websocket "lenavire/internal/ledger/infrastructure/websocket"

	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	database.ConnectDB()

	hub := websocket.NewLedgerActivityHub()
	go hub.Run()

	database.DB.AutoMigrate(&schema.PaymentModel{}, &schema.ExpenseModel{})

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}

		return c.Next()
	})

	/* Repositories */
	paymentRepository := adapters.NewPostgrePaymentRepository(database.DB)
	expenseRepository := adapters.NewPostgreExpenseRepository(database.DB)

	/* Providers */
	idProvider := adapters.NewUUIDIdProvider()
	dateProvider := adapters.NewRealDateProvider()

	/* Channels */
	ledgerActivityChannel := adapters.NewWebSocketLedgerActivityChannel(hub)

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

	paymentGateway := adapters.NewStripePaymentGateway(os.Getenv("STRIPE_SECRET_KEY"))
	addExpenseHandler := handlers.NewAddExpenseHandler(addExpenseCommandHandler)

	getLedgerQueryHandler := queries.NewGetLedgerQueryHandler(database.DB)
	getLedgerHandler := handlers.NewGetLedgerHandler(getLedgerQueryHandler)

	getContributorsRankQueryHandler := queries.NewGetContributorsRankQueryHandler(database.DB)
	getContributorsRankHandler := handlers.NewGetContributorsRankHandler(getContributorsRankQueryHandler)

	getExpensesHistoryQueryHandler := queries.NewGetExpensesHistoryQueryHandler(database.DB)
	getExpensesHistoryHandler := handlers.NewGetExpensesHistoryHandler(getExpensesHistoryQueryHandler)

	getPaymentsHistoryQueryHandler := queries.NewGetPaymentsHistoryQueryHandler(database.DB)
	getPaymentsHistoryHandler := handlers.NewGetPaymentsHistoryHandler(getPaymentsHistoryQueryHandler)

	addPaymentHandler := handlers.NewAddPaymentHandler(receivePaymentCommandHandler)

	createCheckoutCommandHandler := commands.NewCreateCheckoutCommandHandler(paymentGateway)
	createCheckoutHandler := handlers.NewCreateCheckoutHandler(createCheckoutCommandHandler)

	receiveStripePaymentHandler := handlers.NewReceiveStripePaymentHandler(receivePaymentCommandHandler)

	/* Routes */
	api.SetupRoutes(
		app,
		receiveStripePaymentHandler,
		receivePaymentHandler,
		addExpenseHandler,
		getLedgerHandler,
		getContributorsRankHandler,
		getExpensesHistoryHandler,
		getPaymentsHistoryHandler,
		addPaymentHandler,
		createCheckoutHandler,
		hub,
	)

	/* Start server */
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	app.Listen(":" + port)
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
	}
}
