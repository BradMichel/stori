package notifiers

import (
	"context"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/pkg/email"
)

const accountSummaryTemplatePath = "account_summarize_template.html"

//go:embed account_summarize_template.html
var accountSummaryTemplate embed.FS

type AccountSummaryEmailBuilder struct{}

func NewAccountSummaryEmailBuilder() *AccountSummaryEmailBuilder {
	return &AccountSummaryEmailBuilder{}
}

func (b *AccountSummaryEmailBuilder) Build(
	_ context.Context,
	from string,
	notification summarizers.AccountNotification,
) (email.Email, error) {
	content, err := b.getContent(notification)
	if err != nil {
		return email.Email{}, err
	}

	period := notification.Account.Period.Format("2006-01-02")
	return email.Email{
		From:    from,
		To:      notification.Account.Email,
		Subject: fmt.Sprintf("Account Summary before %s", period),
		Content: content,
	}, nil
}

func (b *AccountSummaryEmailBuilder) getContent(notification summarizers.AccountNotification) (string, error) {
	accountSummary := b.getAccountSummary(notification)
	htmlTemplate, err := template.ParseFS(accountSummaryTemplate, accountSummaryTemplatePath)
	if err != nil {
		return "", err
	}

	w := new(strings.Builder)
	err = htmlTemplate.Execute(w, accountSummary)
	if err != nil {
		return "", err
	}

	return w.String(), nil
}

func (b *AccountSummaryEmailBuilder) getAccountSummary(notification summarizers.AccountNotification) AccountSummary {
	accountSummary := AccountSummary{
		ID:                  notification.Account.ID,
		TotalBalance:        notification.Summary.Total,
		AverageDebitAmount:  notification.Summary.DebitAvg,
		AverageCreditAmount: notification.Summary.CreditAvg,
		Months:              make([]AccountMonthSummary, len(notification.Summary.Months)),
	}

	for i, month := range notification.Summary.Months {
		accountSummary.Months[i] = AccountMonthSummary{
			Month:         month.Date,
			NTransactions: month.Transactions,
		}
	}

	return accountSummary
}
