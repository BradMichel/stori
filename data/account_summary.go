package data

import "github.com/BradMichel/stori/internal/summarizers"

func GetAccountSummary(accountID string) summarizers.AccountSummary {
	return summarizers.AccountSummary{
		AccountID: accountID,
		Total:     39.74,
		DebitAvg:  -15.38,
		CreditAvg: 35.25,
		Months: []summarizers.AccountMonthSummary{
			{
				Date:         "July",
				Debit:        -10.3,
				DebitCount:   1,
				Credit:       60.5,
				CreditCount:  1,
				Transactions: 2,
			},
			{
				Date:         "August",
				Debit:        -20.46,
				DebitCount:   1,
				Credit:       10,
				CreditCount:  1,
				Transactions: 2,
			},
		},
	}
}
