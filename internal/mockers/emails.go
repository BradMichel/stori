package mockers

import (
	"context"
	"github.com/BradMichel/stori/pkg/email"
	"github.com/stretchr/testify/mock"
)

const (
	sendMethod = "Send"
)

type EmailSenderMocker struct {
	mock.Mock
}

func (m *EmailSenderMocker) MockSend(ctx context.Context, e email.Email, err error, times int) {
	m.On(sendMethod, ctx, e).Return(err).Times(times)
}

func (m *EmailSenderMocker) Send(ctx context.Context, e email.Email) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}
