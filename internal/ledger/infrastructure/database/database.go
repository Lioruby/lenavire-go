package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"lenavire/internal/ledger/infrastructure/database/schema"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_PORT"),
		)
	} else {

		databaseUrl = databaseUrl + "?sslmode=require"
	}

	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	DB = db
	SyncLastDB(db)
	fmt.Println("✅ Database connection successful!")
}

type RawPayment struct {
	ID          string  `json:"id"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	PaymentType string  `json:"paymentType"`
}

type RawExpense struct {
	ID     string `json:"id"`
	Amount int    `json:"amount"`
	Date   string `json:"date"`
}

func SyncLastDB(db *gorm.DB) {
	rawPaymentsData, err := os.ReadFile("app/internal/ledger/infrastructure/database/payments.json")
	if err != nil {
		fmt.Println("❌ Erreur lecture fichier payments.json:", err)
		return
	}

	var rawPayments []RawPayment
	if err := json.Unmarshal(rawPaymentsData, &rawPayments); err != nil {
		fmt.Println("❌ Erreur parsing JSON:", err)
		return
	}

	var payments []schema.PaymentModel
	for _, raw := range rawPayments {
		date, _ := time.Parse(time.RFC3339, raw.Date)

		payments = append(payments, schema.PaymentModel{
			ID:          raw.ID,
			Date:        date,
			Amount:      int(raw.Amount),
			Name:        raw.Name,
			Email:       raw.Email,
			PaymentType: raw.PaymentType,
		})
	}

	rawExpensesData, err := os.ReadFile("app/internal/ledger/infrastructure/database/expenses.json")
	if err != nil {
		fmt.Println("❌ Erreur lecture fichier expenses.json:", err)
		return
	}

	var rawExpenses []RawExpense
	if err := json.Unmarshal(rawExpensesData, &rawExpenses); err != nil {
		fmt.Println("❌ Erreur parsing JSON:", err)
		return
	}

	var expenses []schema.ExpenseModel

	for _, raw := range rawExpenses {
		date, _ := time.Parse(time.RFC3339, raw.Date)

		expenses = append(expenses, schema.ExpenseModel{
			ID:     raw.ID,
			Amount: raw.Amount,
			Date:   date,
		})
	}

	if err := db.Create(&payments).Error; err != nil {
		fmt.Println("❌ Erreur insertion payments:", err)
		return
	}

	if err := db.Create(&expenses).Error; err != nil {
		fmt.Println("❌ Erreur insertion expenses:", err)
		return
	}

	fmt.Println("✅ Paiements importés avec succès.")
}
