package api

import (
	"lenavire/internal/ledger/infrastructure/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, paymentHandler *handlers.ReceivePaymentHandler, addExpenseHandler *handlers.AddExpenseHandler) {
	app.Post("/receive-payment", paymentHandler.ReceivePayment)
	app.Post("/add-expense", addExpenseHandler.AddExpense)
}
