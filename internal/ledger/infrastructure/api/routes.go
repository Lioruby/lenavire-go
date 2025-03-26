package api

import (
	"lenavire/internal/ledger/infrastructure/api/handlers"

	ws "lenavire/internal/ledger/infrastructure/websocket"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(
	app *fiber.App,
	paymentHandler *handlers.ReceivePaymentHandler,
	addExpenseHandler *handlers.AddExpenseHandler,
	getLedgerHandler *handlers.GetLedgerHandler,
	getContributorsRankHandler *handlers.GetContributorsRankHandler,
	getExpensesHistoryHandler *handlers.GetExpensesHistoryHandler,
	getPaymentsHistoryHandler *handlers.GetPaymentsHistoryHandler,
	hub *ws.LedgerActivityHub,
) {

	app.Use("/ledger-activity", handlers.HandleWebSocket)

	app.Get("/ledger-activity", handlers.HandleLedgerActivity(hub))

	app.Post("/ledger/receive-payment", paymentHandler.ReceivePayment)
	app.Post("/ledger/add-expense", addExpenseHandler.AddExpense)
	app.Get("/ledger", getLedgerHandler.GetLedger)
	app.Get("/ledger/contributors-rank", getContributorsRankHandler.GetContributorsRank)
	app.Get("/ledger/expenses-history", getExpensesHistoryHandler.GetExpensesHistory)
	app.Get("/ledger/payments-history", getPaymentsHistoryHandler.GetPaymentsHistory)
}
