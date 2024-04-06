package data

import "github.com/BradMichel/stori/internal/transactions"

func GetTransactions(accountID string) []transactions.Transaction {
	return []transactions.Transaction{
		{
			ID:        "0",
			AccountID: accountID,
			Amount:    60.5,
			Date:      "7/15",
		},
		{
			ID:        "1",
			AccountID: accountID,
			Amount:    -10.3,
			Date:      "7/28",
		},
		{
			ID:        "2",
			AccountID: accountID,
			Amount:    -20.46,
			Date:      "8/2",
		},
		{
			ID:        "3",
			AccountID: accountID,
			Amount:    10,
			Date:      "8/13",
		},
	}
}
