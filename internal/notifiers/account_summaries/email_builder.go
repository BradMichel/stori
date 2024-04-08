package account_summaries

import (
	"context"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/kelseyhightower/envconfig"
	"github.com/mcuadros/go-defaults"

	"github.com/BradMichel/stori/internal/accounts"

	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/pkg/email"
)

const templatePath = "template.html"

type BlueAccountTemplate string
type BlackCardTemplate string
type GreenCardTemplate string

type EmailBuilderConfig struct {
	BlueAccountTemplatePath string `envconfig:"BLUE_ACCOUNT_TEMPLATE_PATH" default:"stori_account_summary_template.html" required:"true"`
	BlackCardTemplatePath   string `envconfig:"BLACK_CARD_TEMPLATE_PATH" default:"stori_black_summary_template.html" required:"true"`
	GreenCardTemplatePath   string `envconfig:"GREEN_CARD_TEMPLATE_PATH" default:"stori_green_summary_template.html" required:"true"`
	Headers                 EmailHeaders
	Resources               FootImageByAccountType
}

type EmailBuilder struct {
	config                 EmailBuilderConfig
	templatesByAccountType map[string]string
	footImageByAccountType map[string]string
}

//go:embed template.html
var emailTemplate embed.FS

func NewEmailBuilder(
	config EmailBuilderConfig,
	blueAccountTemplate BlueAccountTemplate,
	blackCardTemplate BlackCardTemplate,
	greenCardTemplate GreenCardTemplate,
) *EmailBuilder {
	var templatesByAccountType = map[string]string{
		accounts.BlueAccountKey: string(blueAccountTemplate),
		accounts.BlackCardKey:   string(blackCardTemplate),
		accounts.GreenCardKey:   string(greenCardTemplate),
	}
	var footImageByAccountType = map[string]string{
		accounts.BlueAccountKey: config.Resources.BlueAccount,
		accounts.BlackCardKey:   config.Resources.BlackCard,
		accounts.GreenCardKey:   config.Resources.GreenCard,
	}
	return &EmailBuilder{
		config:                 config,
		templatesByAccountType: templatesByAccountType,
		footImageByAccountType: footImageByAccountType,
	}
}

func NewEmailBuilderConfig() (EmailBuilderConfig, error) {
	config := EmailBuilderConfig{}
	defaults.SetDefaults(&config)
	err := envconfig.Process("", &config)
	return config, err
}

func (b *EmailBuilder) Build(
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

func (b *EmailBuilder) getContent(notification summarizers.AccountNotification) (string, error) {
	accountSummary := b.getAccountSummary(notification)
	htmlTemplate, err := b.getTemplateByAccountType(notification.Account.Type)
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

func (b *EmailBuilder) getTemplateByAccountType(accountType string) (*template.Template, error) {
	templateString, ok := b.templatesByAccountType[accountType]
	if ok {
		return template.New(accountType).Parse(templateString)
	}

	return template.ParseFS(emailTemplate, templatePath)
}

func (b *EmailBuilder) getAccountSummary(notification summarizers.AccountNotification) Email {
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

	footImage, _ := b.footImageByAccountType[notification.Account.Type]
	return Email{
		Headers: b.config.Headers,
		Data:    accountSummary,
		Resources: AccountSummaryEmailResources{
			FootImage: footImage,
		},
	}
}
