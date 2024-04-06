package v1_test

import (
	"context"
	"fmt"
	"testing"

	v1 "github.com/BradMichel/stori/cmd/serverless/account_summarizers/initiator/v1"
	"github.com/BradMichel/stori/internal/mockers"

	"github.com/BradMichel/stori/pkg/logger"
	"github.com/BradMichel/stori/pkg/time"
)

type Fields struct {
	usecase     *mockers.SummarizerInitiatorMocker
	timeService *time.Mock
	logger      *logger.Mocker
}
type Args struct {
	ctx  context.Context
	req  v1.Request
	time time.Time
	err  error
}

func TestInitiator_Handle(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		fields Fields
		args   Args
		mock   func(Fields, Args)
	}{
		{
			name:   "success: launch is called and response is nil",
			fields: getFields(),
			args: Args{
				ctx:  context.Background(),
				req:  v1.Request{},
				time: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			mock: mockInitiatorHandle,
		},
		{
			name:   "error: launch is called and response with error",
			fields: getFields(),
			args: Args{
				ctx:  context.Background(),
				req:  v1.Request{},
				time: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC),
				err:  fmt.Errorf("launch error"),
			},
			mock: mockInitiatorHandle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := tt.fields
			args := tt.args
			mock := tt.mock
			mock(fields, args)
			h := v1.New(fields.usecase, fields.timeService, fields.logger)
			_ = h.Handle(tt.args.ctx, tt.args.req)
			fields.timeService.AssertExpectations(t)
			fields.usecase.AssertExpectations(t)
			fields.logger.AssertExpectations(t)
		})
	}
}

func getFields() Fields {
	return Fields{
		usecase:     new(mockers.SummarizerInitiatorMocker),
		timeService: new(time.Mock),
		logger:      new(logger.Mocker),
	}
}

func mockInitiatorHandle(fields Fields, args Args) {
	fields.timeService.MockNow(args.time, 1)
	fields.usecase.MockLaunch(args.ctx, args.time, args.err, 1)
	if args.err != nil {
		fields.logger.MockError(1, args.err)
	}
}
