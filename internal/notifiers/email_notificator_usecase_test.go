package notifiers_test

import (
	"context"
	"testing"
	"time"

	"github.com/BradMichel/stori/internal/notifiers/account_summaries"

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
					Content: "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <title>Account Summary</title>\n</head>\n<body>\n<div class=\"container mt-5\">\n    <div class=\"card\">\n        <div class=\"card-header\">\n            <h1 class=\"text-center\"> account_id_1</h1>\n        </div>\n        <div class=\"card-body\">\n            <h5 class=\"card-title\"> 39.74</h5>\n            <p class=\"card-text\"> -15.38</p>\n            <p class=\"card-text\"> 35.25</p>\n            <h2></h2>\n            <div>\n                \n                <div>\n                    <h3><strong></strong> July</h3>\n                    <h3><strong></strong> 2</h3>\n                </div>\n                \n                <div>\n                    <h3><strong></strong> August</h3>\n                    <h3><strong></strong> 2</h3>\n                </div>\n                \n            </div>\n            <div class=\"mt-4\">\n                <img data-imagetype=\"External\" src=\"https://media.giphy.com/media/v1.Y2lkPTc5MGI3NjExbTFvM2QzZjZ0NmR4bzFoNHVicXV3c2tvdmxpbTI3bmk0a25iaTdhNCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/VXWN05HTttWxoCnXdP/source.gif\" width=\"200\" height=\"110\" alt=\"stori\">\n            </div>\n        </div>\n    </div>\n</div>\n</body>\n</html>",
				},
			},
			mock: Mock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.mock(tt.fields, tt.args)
			builder := account_summaries.NewEmailBuilder(account_summaries.EmailBuilderConfig{}, "", "", "")
			n := notifiers.NewEmailNotificator[summarizers.AccountNotification](tt.fields.config, builder, tt.fields.emailSender)
			err := n.Notify(tt.args.ctx, tt.args.notification)
			assert.ErrorIs(t, err, tt.args.expectedErr)
			tt.fields.emailSender.AssertExpectations(t)
		})
	}

}
