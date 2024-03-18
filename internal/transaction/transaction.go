package transaction

import (
	"time"
)

type Transaction struct {
	Date        time.Time
	Amount      float64
	Description string
}

type Summary struct {
	TotalBalance      float64
	TotalTransactions int
	AvgDebit          float64
	AvgCredit         float64
}

func ProcessTransactions(transactions []Transaction) (Summary, error) {
	var summary Summary
	var totalDebit float64
	var totalCredit float64
	var debitCount int
	var creditCount int

	for _, t := range transactions {
		if t.Amount < 0 {
			totalDebit += t.Amount
			debitCount++
		} else {
			totalCredit += t.Amount
			creditCount++
		}
	}

	summary.TotalBalance = totalCredit + totalDebit
	summary.TotalTransactions = len(transactions)

	if debitCount > 0 {
		summary.AvgDebit = totalDebit / float64(debitCount)
	}

	if creditCount > 0 {
		summary.AvgCredit = totalCredit / float64(creditCount)
	}

	return summary, nil
}
