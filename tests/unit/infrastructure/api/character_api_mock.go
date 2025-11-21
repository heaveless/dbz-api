package api_test

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockExternalClient struct {
	mock.Mock
}

func (m *MockExternalClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
