package database

import (
	"database/sql"
	"fmt"
	"github.com/romacardozx/stori-card-challenge/internal/config"
	"github.com/romacardozx/stori-card-challenge/internal/transaction"
)

func NewDatabase(config config.PostgresConfig) (*sql.DB, error) {
	// Construir la cadena de conexión a PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func SaveTransactionsAndSummary(db *sql.DB, transactions []transaction.Transaction, summary transaction.Summary) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, t := range transactions {
		_, err = tx.Exec("INSERT INTO transactions (date, amount, description) VALUES ($1, $2, $3)", t.Date, t.Amount, t.Description)
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec("INSERT INTO summary (total_balance, total_transactions, avg_debit, avg_credit) VALUES ($1, $2, $3, $4)",
		summary.TotalBalance, summary.TotalTransactions, summary.AvgDebit, summary.AvgCredit)
	if err != nil {
		return err
	}

	return tx.Commit()
}
