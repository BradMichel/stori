package logger

import (
	"fmt"
	"github.com/stretchr/testify/mock"
)

const errorMethod = "Error"

type Mocker struct {
	mock.Mock
}

func (m *Mocker) MockError(times int, args ...interface{}) {
	m.On(errorMethod, fmt.Sprint(args)).Times(times)
}

func (m *Mocker) Error(args ...interface{}) {
	m.Called(fmt.Sprint(args))
}
