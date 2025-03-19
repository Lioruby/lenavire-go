package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		fmt.Println("Error: DATABASE_URL not defined")
		os.Exit(1)
	}

	databaseUrl = databaseUrl + "?sslmode=require"

	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("🚀 Connected to the database!")
	fmt.Println("Type 'help' to see the available commands")
	fmt.Println("Type 'exit' to quit")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		command := strings.TrimSpace(scanner.Text())
		if command == "exit" {
			break
		}

		handleCommand(db, command)
	}
}

func handleCommand(db *gorm.DB, command string) {
	switch command {
	case "help":
		fmt.Println("Available commands:")
		fmt.Println("  payments     - Show the latest payments")
		fmt.Println("  expenses     - Show the latest expenses")
		fmt.Println("  stats        - Show the general statistics")
		fmt.Println("  contributors - Show the top contributors of the month")
		fmt.Println("  help         - Show this help")
		fmt.Println("  exit         - Quit the program")

	case "payments":
		var payments []struct {
			Amount      int
			Name        string
			Email       string
			PaymentType string `gorm:"column:payment_type"`
			Date        string
		}
		db.Raw("SELECT amount, name, email, payment_type, date FROM payments ORDER BY date DESC LIMIT 5").Scan(&payments)

		fmt.Println("\nLatest payments:")
		for _, p := range payments {
			fmt.Printf("- %s (%s) paid %d€ on %s via %s\n",
				p.Name, p.Email, p.Amount, p.Date, p.PaymentType)
		}

	case "expenses":
		var expenses []struct {
			Amount      int
			Description string
			Date        string
		}
		db.Raw("SELECT amount, description, date FROM expenses ORDER BY date DESC LIMIT 5").Scan(&expenses)

		fmt.Println("\nLatest expenses:")
		for _, e := range expenses {
			fmt.Printf("- %d€ : %s (%s)\n", e.Amount, e.Description, e.Date)
		}

	case "stats":
		var stats struct {
			TotalReceived int
			TotalExpenses int
		}
		db.Raw(`
			SELECT 
				(SELECT COALESCE(SUM(amount), 0) FROM payments) as total_received,
				(SELECT COALESCE(SUM(amount), 0) FROM expenses) as total_expenses
		`).Scan(&stats)

		fmt.Printf("\nStatistics:\n")
		fmt.Printf("Total received: %d€\n", stats.TotalReceived)
		fmt.Printf("Total expenses: %d€\n", stats.TotalExpenses)
		fmt.Printf("Balance: %d€\n", stats.TotalReceived-stats.TotalExpenses)

	case "contributors":
		var contributors []struct {
			Name   string
			Amount int
		}
		db.Raw(`
			SELECT 
				MAX(name) as name,
				SUM(amount) as amount
			FROM payments
			WHERE date_trunc('month', date) = date_trunc('month', CURRENT_DATE)
			GROUP BY email
			ORDER BY amount DESC
			LIMIT 5
		`).Scan(&contributors)

		fmt.Println("\nTop contributors of the month:")
		for _, c := range contributors {
			fmt.Printf("- %s: %d€\n", c.Name, c.Amount)
		}

	default:
		fmt.Println("Unknown command. Type 'help' to see the available commands")
	}
}
