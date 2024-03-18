package database

import (
	"database/sql"
	"fmt"

	"github.com/romacardozx/stori-card-challenge/internal/config"
	"github.com/romacardozx/stori-card-challenge/internal/transaction"

	_ "github.com/lib/pq"
)

func NewDatabase(cfg config.PostgresConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	err = createTables(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %v", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			date DATE NOT NULL,
			amount DECIMAL(10, 2) NOT NULL
		);

		CREATE TABLE IF NOT EXISTS summary (
			id SERIAL PRIMARY KEY,
			total_balance DECIMAL(10, 2) NOT NULL,
			total_transactions INTEGER NOT NULL,
			avg_debit DECIMAL(10, 2) NOT NULL,
			avg_credit DECIMAL(10, 2) NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}

func SaveTransactionsAndSummary(db *sql.DB, transactions []transaction.Transaction, summary transaction.Summary) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	for _, t := range transactions {
		_, err := tx.Exec(`
			INSERT INTO transactions (date, amount)
			VALUES ($1, $2)
		`, t.Date, t.Amount)
		if err != nil {
			return fmt.Errorf("failed to insert transaction: %v", err)
		}
	}

	_, err = tx.Exec(`
		INSERT INTO summary (total_balance, total_transactions, avg_debit, avg_credit)
		VALUES ($1, $2, $3, $4)
	`, summary.TotalBalance, summary.TotalTransactions, summary.AvgDebit, summary.AvgCredit)
	if err != nil {
		return fmt.Errorf("failed to insert summary: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}
