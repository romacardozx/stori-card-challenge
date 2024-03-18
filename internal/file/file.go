package file

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/romacardozx/stori-card-challenge/internal/transaction"
)

func ReadTransactionsFromCSV(filePath string) ([]transaction.Transaction, error) {
	// Abrir el archivo CSV
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()

	// Crear un nuevo lector de CSV
	reader := csv.NewReader(file)

	// Leer todos los registros del archivo CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV records: %v", err)
	}

	var transactions []transaction.Transaction

	// Iterar sobre los registros y crear las transacciones
	for _, record := range records {
		// Parsear los campos del registro CSV
		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse date: %v", err)
		}

		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse amount: %v", err)
		}

		description := record[2]

		// Crear una transacción y agregarla al slice de transacciones
		transaction := transaction.Transaction{
			Date:        date,
			Amount:      amount,
			Description: description,
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
