package utils

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockHttpClient struct {
	mock.Mock
	MockDo func(req *http.Request) (*http.Response, error)
}

// Do is a mocking method to help testing HTTP Requests.
// Its used to replace the default HTTP Client implementation
func (m MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}
