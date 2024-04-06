package v1_test

import (
	"context"
	"fmt"
	"testing"

	v1 "github.com/BradMichel/stori/cmd/serverless/account_summarizers/summarizer/v1"
	"github.com/BradMichel/stori/internal/accounts"
	"github.com/BradMichel/stori/internal/mockers"
	"github.com/BradMichel/stori/internal/summarizers"
	"github.com/BradMichel/stori/pkg/aws/sqs"
	"github.com/BradMichel/stori/pkg/logger"
	"github.com/BradMichel/stori/pkg/validators"
	"github.com/aws/aws-lambda-go/events"
)

func TestHandler_Handle(t *testing.T) {
	t.Parallel()

	type Fields struct {
		usecase   *mockers.SummarizerMocker
		validator *validators.Mocker
		logger    *logger.Mocker
	}

	type Args struct {
		ctx               context.Context
		req               v1.Request
		message           string
		accountExpected   accounts.Account
		sqsUnmarshalErr   error
		validateErr       error
		produceSummaryErr error
		expectedErr       error
	}

	Mock := func(fields Fields, args Args) {
		if args.expectedErr != nil {
			fields.logger.MockError(1, args.expectedErr)
		}

		if args.sqsUnmarshalErr != nil {
			return
		}

		fields.validator.MockValidate(args.ctx, []byte(args.message), args.validateErr, 1)
		if args.validateErr != nil {
			return

		}

		fields.usecase.MockProduceSummary(
			args.ctx,
			args.accountExpected,
			summarizers.AccountSummary{},
			args.produceSummaryErr,
			1,
		)
	}

	var GetFields = func() Fields {
		return Fields{
			usecase:   new(mockers.SummarizerMocker),
			validator: new(validators.Mocker),
			logger:    new(logger.Mocker),
		}

	}

	var unmarshalErr = sqs.UnmarshalErr
	var validateErr = fmt.Errorf("%w body: {} details: %s", validators.Err, "\n(root): id is required\n(root): email is required")
	var produceSummaryErr = fmt.Errorf("produce summary error")
	var message = `{"id":"1", "email":"test_email"}`
	var sqsBodyMessage = events.SQSMessage{Body: `{"Message": "{\"id\": \"1\", \"email\": \"test_email\"}"}`}
	var account = accounts.Account{ID: "1", Email: "test_email"}

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
				ctx:             context.Background(),
				req:             v1.Request{Records: []events.SQSMessage{sqsBodyMessage}},
				message:         message,
				accountExpected: account,
			},
			mock: Mock,
		},
		{
			name:   "error: produce summary is called and response with error",
			fields: GetFields(),
			args: Args{
				ctx:               context.Background(),
				req:               v1.Request{Records: []events.SQSMessage{sqsBodyMessage}},
				message:           message,
				accountExpected:   account,
				produceSummaryErr: produceSummaryErr,
				expectedErr:       produceSummaryErr,
			},
			mock: Mock,
		},
		{
			name:   "error: validate is called and response with error",
			fields: GetFields(),
			args: Args{
				ctx:             context.Background(),
				req:             v1.Request{Records: []events.SQSMessage{{Body: fmt.Sprintf(`{"Message": %q}`, "{}")}}},
				message:         message,
				accountExpected: account,
				validateErr:     validateErr,
				expectedErr:     validateErr,
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
				fields.validator.AssertExpectations(t)
			}
		})
	}

}
