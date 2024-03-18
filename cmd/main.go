package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/romacardozx/stori-card-challenge/internal/config"
	"github.com/romacardozx/stori-card-challenge/internal/database"
	"github.com/romacardozx/stori-card-challenge/internal/email"
	"github.com/romacardozx/stori-card-challenge/internal/file"
	"github.com/romacardozx/stori-card-challenge/internal/transaction"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.NewDatabase(cfg.PostgresConfig)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/process-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, "Method not allowed")
			return
		}

		transactions, err := file.ReadTransactionsFromCSV(cfg.CSVFilePath)
		if err != nil {
			log.Printf("Failed to read transactions from CSV: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to read transactions from CSV")
			return
		}

		summary, err := transaction.ProcessTransactions(transactions)
		if err != nil {
			log.Printf("Failed to process transactions: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to process transactions")
			return
		}

		err = database.SaveTransactionsAndSummary(db, transactions, summary)
		if err != nil {
			log.Printf("Failed to save transactions and summary to the database: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to save transactions and summary to the database")
			return
		}

		err = email.SendSummaryEmail(cfg.SMTPConfig, summary)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Failed to send email")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Process completed successfully")
	})

	log.Printf("Server started on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}
