package breaker_test

import "github.com/stretchr/testify/mock"

type MockSingleResultWrapper struct {
	mock.Mock
}

func (m *MockSingleResultWrapper) Decode(v any) error {
	args := m.Called(v)
	return args.Error(0)
}

func (m *MockSingleResultWrapper) Err() error {
	args := m.Called()
	return args.Error(0)
}
