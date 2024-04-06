package validators

import (
	"context"

	"github.com/stretchr/testify/mock"
)

const validateMethod = "Validate"

type Mocker struct {
	mock.Mock
}

func (m *Mocker) MockValidate(ctx context.Context, body []byte, err error, times int) {
	m.On(validateMethod, ctx, body).Return(err).Times(times)
}

func (m *Mocker) Validate(ctx context.Context, body []byte) error {
	args := m.Called(ctx, body)
	return args.Error(0)
}
