package time

import (
	"time"

	"github.com/stretchr/testify/mock"
)

const nowMethod = "Now"

type Mock struct {
	mock.Mock
}

func (m *Mock) MockNow(t time.Time, times int) {
	m.On(nowMethod).Return(t).Times(times)
}

func (m *Mock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}
