package tests_test

import (
	"context"
	"testing"
	"time"

	"github.com/BradMichel/stori/data"

	"github.com/BradMichel/stori/internal/accounts"
	"github.com/BradMichel/stori/internal/notifiers/account_summaries"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/stretchr/testify/assert"

	"github.com/BradMichel/stori/internal/notifiers"
	"github.com/BradMichel/stori/internal/summarizers"
	ses2 "github.com/BradMichel/stori/pkg/email/ses"
)

func TestEmailNotificatorUseCase_Notify(t *testing.T) {
	const fromEmail = "michel.en@hotmail.com"
	const accountID = "account_id_1"
	const clientEmail = "bradmichel10@gmail.com"
	config := notifiers.EmailNotificatorConfig{
		From: fromEmail,
	}
	notification := summarizers.AccountNotification{
		Account: accounts.Account{
			ID:     accountID,
			Email:  clientEmail,
			Period: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Summary: data.GetAccountSummary(accountID),
	}

	ctx := context.Background()
	awsConfig := GetAWSConfig()
	mySession := session.Must(session.NewSession(&awsConfig))
	sesClient := ses.New(mySession)
	emailSender := ses2.New(sesClient)
	builder := account_summaries.NewEmailBuilder()
	n := notifiers.NewEmailNotificator[summarizers.AccountNotification](config, builder, emailSender)
	err := n.Notify(ctx, notification)
	assert.ErrorIs(t, err, nil)

}
