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
	d := fmt.Sprintf("%v", args)
	m.On(errorMethod, d).Times(times)
}

func (m *Mocker) Error(args ...interface{}) {
	d := fmt.Sprintf("%v", args)
	m.Called(d)
}
