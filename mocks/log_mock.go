package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Infow(msg string, keysAndValues ...interface{}) {
	m.Called(append([]interface{}{msg}, keysAndValues...)...)
}

func (m *MockLogger) Errorw(msg string, keysAndValues ...interface{}) {
	m.Called(append([]interface{}{msg}, keysAndValues...)...)
}

func (m *MockLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	m.Called(append([]interface{}{msg}, keysAndValues...)...)
}

func (m *MockLogger) Warnw(msg string, keysAndValues ...interface{}) {
	m.Called(append([]interface{}{msg}, keysAndValues...)...)
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}
