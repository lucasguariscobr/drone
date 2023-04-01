package utils

import (
	"io"
	"net/http"
	"strings"
)

func BuildTestResponse(httpStatusCode int, mockBody string) func(req *http.Request) (*http.Response, error) {
	r := io.NopCloser(strings.NewReader(mockBody))
	return func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: httpStatusCode,
			Body:       r,
		}, nil
	}
}

func BuildNilTestResponse() func(req *http.Request) (*http.Response, error) {
	return func(req *http.Request) (*http.Response, error) {
		return nil, nil
	}
}

func MockOsReadFile(jsonPayload string) {
	readFileFunc := func(_ string) ([]byte, error) {
		return []byte(jsonPayload), nil
	}
	setOsReadFileFunc(readFileFunc)
}
