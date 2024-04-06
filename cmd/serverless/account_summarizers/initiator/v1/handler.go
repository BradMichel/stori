//go:generate go get github.com/google/wire/cmd/wire@v0.6.0
//go:generate go run github.com/google/wire/cmd/wire gen

package v1

import (
	"context"

	"github.com/BradMichel/stori/pkg/time"
	"github.com/aws/aws-lambda-go/events"
)

type Initiator interface {
	Launch(context.Context, time.Time) error
}

type Time interface {
	Now() time.Time
}
type Logger interface {
	Error(...interface{})
}

type Handler struct {
	usecase     Initiator
	timeService Time
	logger      Logger
}

type Request events.CloudWatchEvent

func New(usecase Initiator, timeService Time, logger Logger) *Handler {
	return &Handler{
		usecase:     usecase,
		timeService: timeService,
		logger:      logger,
	}
}

func (h *Handler) Handle(ctx context.Context, _ Request) error {
	date := h.timeService.Now()
	err := h.usecase.Launch(ctx, date)
	if err != nil {
		h.logger.Error(err)
		return nil
	}

	return nil
}
