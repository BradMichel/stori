package summarizers

import (
	"context"
	"time"

	"github.com/BradMichel/stori/internal/transactions"
)

type AccountService struct {
}

type SerConfig struct {
}

func NewAccountService() *AccountService {
	return &AccountService{}
}

func (s *AccountService) BuildSummary(
	_ context.Context,
	accountID string,
	transactionList []transactions.Transaction,
) (AccountSummary, error) {
	accountSummary, err := s.buildFromTransactions(transactionList)
	if err != nil {
		return AccountSummary{}, err
	}
	accountSummary.AccountID = accountID
	accountSummary.DebitAvg = accountSummary.totalDebit / float64(accountSummary.countDebit)
	accountSummary.CreditAvg = accountSummary.totalCredit / float64(accountSummary.countCredit)
	return accountSummary, nil
}

func (s *AccountService) buildFromTransactions(transactionList []transactions.Transaction) (AccountSummary, error) {
	const layout = "1/2"
	accountSummary := AccountSummary{}
	monthSummary := [12]AccountMonthSummary{}
	for _, transaction := range transactionList {
		dateTime, err := time.Parse(layout, transaction.Date)
		if err != nil {
			return AccountSummary{}, err
		}

		month := dateTime.Month()
		summary := monthSummary[month]
		summary.Date = month.String()
		accountSummary, summary = s.summarizeMonth(accountSummary, summary, transaction)
		monthSummary[month] = summary
	}

	accountSummary.Months = make([]AccountMonthSummary, 0)
	for _, summary := range monthSummary {
		if summary == (AccountMonthSummary{}) {
			continue
		}

		accountSummary.Months = append(accountSummary.Months, summary)
	}

	return accountSummary, nil
}

func (s *AccountService) summarizeMonth(
	accountSummary AccountSummary,
	summary AccountMonthSummary,
	transaction transactions.Transaction,
) (AccountSummary, AccountMonthSummary) {
	accountSummary.Total += transaction.Amount
	summary.Transactions++

	if transaction.Amount < 0 {
		accountSummary.totalDebit += transaction.Amount
		accountSummary.countDebit++
		summary.Debit += transaction.Amount
		summary.DebitCount++

		return accountSummary, summary
	}

	accountSummary.totalCredit += transaction.Amount
	accountSummary.countCredit++
	summary.Credit += transaction.Amount
	summary.CreditCount++

	return accountSummary, summary
}
