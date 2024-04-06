package notifiers_test

import (
	"context"
	"testing"
	"time"

	"github.com/BradMichel/stori/data"

	"github.com/BradMichel/stori/internal/accounts"
	"github.com/stretchr/testify/assert"

	"github.com/BradMichel/stori/internal/mockers"
	"github.com/BradMichel/stori/internal/notifiers"
	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/pkg/email"
)

func TestEmailNotificatorUseCase_Notify(t *testing.T) {
	t.Parallel()

	type Fields struct {
		emailSender *mockers.EmailSenderMocker
		config      notifiers.EmailNotificatorConfig
	}

	type Args struct {
		ctx           context.Context
		notification  summarizers.AccountNotification
		emailExpected email.Email
		sendErr       error
		expectedErr   error
	}

	Mock := func(fields Fields, args Args) {
		fields.emailSender.MockSend(args.ctx, args.emailExpected, args.sendErr, 1)
	}

	const fromEmail = "bradmichel10@gmail.com"
	const accountID = "account_id_1"
	const clientEmail = "michel.en@hotmail.com"

	var GetFields = func() Fields {
		return Fields{
			emailSender: new(mockers.EmailSenderMocker),
			config: notifiers.EmailNotificatorConfig{
				From: fromEmail,
				//TemplatePath: "./account_summarize_template.html",
			},
		}
	}

	tests := []struct {
		name   string
		fields Fields
		args   Args
		mock   func(fields Fields, args Args)
	}{
		{
			name:   "success: email sent successfully",
			fields: GetFields(),
			args: Args{
				ctx: context.Background(),
				notification: summarizers.AccountNotification{
					Account: accounts.Account{
						ID:     accountID,
						Email:  clientEmail,
						Period: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					Summary: data.GetAccountSummary(accountID),
				},
				emailExpected: email.Email{
					From:    fromEmail,
					Subject: "Account Summary before 2020-01-01",
					To:      clientEmail,
					Content: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Account Summary</title>\n</head>\n<body>\n    <h1>Account Summary for account_id_1</h1>\n    <p>Total Balance: 39.74</p>\n    <p>Average Debit Amount: -15.38</p>\n    <p>Average Credit Amount: 35.25</p>\n    <h2>Monthly Summary</h2>\n    \n    <h3>Month: July</h3>\n    <p>Number of Transactions: 2</p>\n    \n    <h3>Month: August</h3>\n    <p>Number of Transactions: 2</p>\n    \n    <div><img data-imagetype=\"External\" src=\"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExbTFvM2QzZjZ0NmR4bzFoNHVicXV3c2tvdmxpbTI3bmk0a25iaTdhNCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/VXWN05HTttWxoCnXdP/source.gif\" width=\"200\" height=\"110\" alt=\"stori\"><br>\n    </div>\n</body>\n</html>",
				},
			},
			mock: Mock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mock(tt.fields, tt.args)
			builder := notifiers.NewAccountSummaryEmailBuilder()
			n := notifiers.NewEmailNotificator[summarizers.AccountNotification](tt.fields.config, builder, tt.fields.emailSender)
			err := n.Notify(tt.args.ctx, tt.args.notification)
			assert.ErrorIs(t, err, tt.args.expectedErr)
			tt.fields.emailSender.AssertExpectations(t)
		})
	}

}
