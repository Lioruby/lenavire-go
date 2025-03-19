package main

import (
	"fmt"
	addExpenseCommands "lenavire/internal/ledger/application/commands/add_expense"
	receivePaymentCommands "lenavire/internal/ledger/application/commands/receive_payment"
	"lenavire/internal/ledger/infrastructure/adapters"
	"lenavire/internal/ledger/infrastructure/api"
	"lenavire/internal/ledger/infrastructure/api/handlers"
	"lenavire/internal/ledger/infrastructure/database"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	fmt.Println("🔄 Connexion à la base de données...")
	database.ConnectDB()

	fmt.Println("🔄 Création de l'application...")
	app := fiber.New()

	fmt.Println("🔄 Création des répositories...")
	/* Repositories */
	paymentRepository := adapters.NewPostgrePaymentRepository(database.DB)
	expenseRepository := adapters.NewPostgreExpenseRepository(database.DB)

	fmt.Println("🔄 Création des fournisseurs...")
	/* Providers */
	idProvider := adapters.NewUUIDIdProvider()
	dateProvider := adapters.NewRealDateProvider()

	fmt.Println("🔄 Création des canaux...")
	/* Channels */
	ledgerActivityChannel := adapters.NewFakeLedgerActivityChannel()

	fmt.Println("🔄 Création des handlers...")
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

	fmt.Println("🔄 Création des routes...")
	/* Routes */
	api.SetupRoutes(app, receivePaymentHandler, addExpenseHandler)

	fmt.Println("🔄 Démarrage du serveur...")
	/* Start server */
	app.Listen(":3000")
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Erreur lors du chargement du fichier .env :", err)
	}
}
