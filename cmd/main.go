package main

import (
	"github.com/romacardozx/stori-card-challenge/internal/email"
	"log"

	"github.com/romacardozx/stori-card-challenge/internal/config"
	"github.com/romacardozx/stori-card-challenge/internal/database"
	"github.com/romacardozx/stori-card-challenge/internal/file"
	"github.com/romacardozx/stori-card-challenge/internal/transaction"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	transactions, err := file.ReadTransactionsFromCSV(cfg.CSVFilePath)
	if err != nil {
		log.Fatalf("Failed to read transactions from CSV: %v", err)
	}

	// Process transactions
	summary, err := transaction.ProcessTransactions(transactions)
	if err != nil {
		log.Fatalf("Failed to process transactions: %v", err)
	}

	// Connect to the database
	db, err := database.NewDatabase(cfg.PostgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	// Save transactions and summary to the database
	err = database.SaveTransactionsAndSummary(db, transactions, summary)
	if err != nil {
		log.Fatalf("Failed to save transactions and summary to the database: %v", err)
	}

	err = email.SendSummaryEmail(cfg.SMTPConfig, summary)
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	log.Println("Process completed successfully")
}
