package v1_test

import (
	"context"
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/BradMichel/stori/data"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	v1 "github.com/BradMichel/stori/cmd/serverless/account_summarizers/notificator/v1"
	"github.com/BradMichel/stori/internal/accounts"
	"github.com/BradMichel/stori/internal/mockers"
	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/pkg/aws/sqs"
	"github.com/BradMichel/stori/pkg/logger"
	"github.com/BradMichel/stori/pkg/validators"
)

//go:embed request.json
var requestJSON embed.FS

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	type Fields struct {
		usecase   *mockers.SummarizerMocker
		validator *validators.Service
		logger    *logger.Mocker
	}

	type Args struct {
		ctx                  context.Context
		req                  v1.Request
		message              string
		notificationExpected summarizers.AccountNotification
		sqsUnmarshalErr      error
		validateErr          error
		notifyErr            error
		expectedErr          error
	}

	Mock := func(fields Fields, args Args) {
		if args.expectedErr != nil {
			fields.logger.MockError(1, args.expectedErr)
		}

		if args.sqsUnmarshalErr != nil {
			return
		}

		if args.validateErr != nil {
			return

		}

		fields.usecase.MockNotify(
			args.ctx,
			args.notificationExpected,
			args.notifyErr,
			1,
		)
	}

	var GetFields = func() Fields {
		validator, err := validators.New(validators.Config{Source: "request.json"}, requestJSON)
		assert.NoError(t, err)
		return Fields{
			usecase:   new(mockers.SummarizerMocker),
			validator: validator,
			logger:    new(logger.Mocker),
		}

	}

	accountID := "1"
	var unmarshalErr = sqs.UnmarshalErr
	var validateErr = fmt.Errorf("%w body: {} details: %s", validators.Err, "\n(root): account is required\n(root): summary is required")
	var produceSummaryErr = fmt.Errorf("produce summary error")
	var message = `{"account":{"id":"1","email":"test_email","last_update":"0001-01-01T00:00:00Z"},"summary":{"account_id":"1","months":[{"date":"July","debit":-10.3,"debit_count":1,"credit":60.5,"credit_count":1,"transactions":2},{"date":"August","debit":-20.46,"debit_count":1,"credit":10,"credit_count":1,"transactions":2}],"credit_avg":35.25,"debit_avg":-15.38,"total":39.74}}`
	var sqsBodyMessage = events.SQSMessage{Body: fmt.Sprintf(`{"Message": %q}`, message)}
	var notification = summarizers.AccountNotification{
		Account: accounts.Account{ID: accountID, Email: "test_email"},
		Summary: data.GetAccountSummary(accountID),
	}

	notificationJSON, _ := json.Marshal(notification)
	fmt.Println(string(notificationJSON))

	tests := []struct {
		name   string
		fields Fields
		args   Args
		mock   func(Fields, Args)
	}{
		{
			name:   "success: produce summary is called and response is nil",
			fields: GetFields(),
			args: Args{
				ctx:                  context.Background(),
				req:                  v1.Request{Records: []events.SQSMessage{sqsBodyMessage}},
				message:              message,
				notificationExpected: notification,
			},
			mock: Mock,
		},
		{
			name:   "error: produce summary is called and response with error",
			fields: GetFields(),
			args: Args{
				ctx:                  context.Background(),
				req:                  v1.Request{Records: []events.SQSMessage{sqsBodyMessage}},
				message:              message,
				notificationExpected: notification,
				notifyErr:            produceSummaryErr,
				expectedErr:          produceSummaryErr,
			},
			mock: Mock,
		},
		{
			name:   "error: validate is called and response with error",
			fields: GetFields(),
			args: Args{
				ctx:                  context.Background(),
				req:                  v1.Request{Records: []events.SQSMessage{{Body: fmt.Sprintf(`{"Message": %q}`, "{}")}}},
				message:              message,
				notificationExpected: notification,
				validateErr:          validateErr,
				expectedErr:          validateErr,
			},
			mock: Mock,
		},
		{
			name:   "error: sqs unmarshal is called and response with error",
			fields: GetFields(),
			args: Args{
				ctx:             context.Background(),
				req:             v1.Request{Records: []events.SQSMessage{{Body: "invalid"}}},
				sqsUnmarshalErr: unmarshalErr,
				expectedErr:     fmt.Errorf("%w: %w", unmarshalErr, fmt.Errorf("invalid character 'i' looking for beginning of value")),
			},
			mock: Mock,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := tt.fields
			args := tt.args
			mock := tt.mock
			mock(fields, args)
			h := v1.New(fields.usecase, fields.validator, fields.logger)
			err := h.Handle(tt.args.ctx, tt.args.req)
			if err != nil {
				fields.logger.AssertExpectations(t)
				fields.usecase.AssertExpectations(t)
			}
		})
	}

}
