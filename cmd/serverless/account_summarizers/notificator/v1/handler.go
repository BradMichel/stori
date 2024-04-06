//go:generate go get github.com/google/wire/cmd/wire@v0.6.0
//go:generate go run github.com/google/wire/cmd/wire gen

package v1

import (
	"context"
	"embed"

	"github.com/BradMichel/stori/pkg/aws/sqs"
	"github.com/aws/aws-lambda-go/events"

	"github.com/BradMichel/stori/internal/summarizers"
)

//go:embed request.json
var requestJSON embed.FS

type Notifier interface {
	Notify(ctx context.Context, notification summarizers.AccountNotification) error
}

type Validator interface {
	Validate(ctx context.Context, body []byte) error
}

type Logger interface {
	Error(...interface{})
}

type Handler struct {
	usecase   Notifier
	validator Validator
	logger    Logger
}

type Request events.SQSEvent

func New(usecase Notifier, validator Validator, logger Logger) *Handler {
	return &Handler{
		usecase:   usecase,
		validator: validator,
		logger:    logger,
	}
}

func (h *Handler) Handle(ctx context.Context, request Request) error {
	for _, record := range request.Records {
		var notification summarizers.AccountNotification
		body, err := sqs.Unmarshal(record.Body, &notification)
		if err != nil {
			h.logger.Error(err)
			return err
		}

		err = h.validator.Validate(ctx, []byte(body.Message))
		if err != nil {
			h.logger.Error(err)
			return err
		}

		err = h.usecase.Notify(ctx, notification)
		if err != nil {
			h.logger.Error(err)
			return err
		}
	}

	return nil
}
