package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// BuildUrl uses the DRONE_ADDR environment variable
// and the API_ENDPOINT constant value to create the base API URL.
// The DRONE_ADDR doesn't have a default value and the function
// will throw an error if the value is not set.
func BuildUrl() (string, error) {
	var listUrl strings.Builder
	backendAddr := viper.GetString(CONFIG_VALUE_ADDR)
	if backendAddr == "" {
		return "", ErrMissingAddr
	}
	listUrl.WriteString(backendAddr)

	if backendAddr[len(backendAddr)-1:] != "/" {
		listUrl.WriteString("/")
	}
	listUrl.WriteString(API_ENDPOINT)

	output := listUrl.String()

	log.Debug().Msgf("API ENDPOINT: %s", output)
	return output, nil
}

// ExecHttpRequest executes the HTTP Request.
// It configures the Authorization Header using the DRONE_TOKEN environment variable.
func ExecHttpRequest(httpClient HttpClientInterface, req *http.Request) (*http.Response, error) {
	authToken := viper.GetString(CONFIG_VALUE_TOKEN)
	if authToken == "" {
		return nil, ErrMissingToken
	}

	log.Debug().Msg("Configuring Authorization Header for the API request")

	req.Header.Set(AUTHORIZATION_HEADER, authToken)
	httpResponse, httpRespError := httpClient.Do(req)
	if httpRespError != nil {
		return nil, httpRespError
	}
	return httpResponse, nil
}

// ParseJsonRawResponse parses the HTTP Response
func ParseJsonRawResponse(httpResponseBody io.ReadCloser) (json.RawMessage, error) {
	log.Debug().Msg("Parsing JSON response...")

	body, bodyErr := io.ReadAll(httpResponseBody)
	if bodyErr != nil {
		return nil, bodyErr
	}

	var jsonRawResponse json.RawMessage
	jsonErr := json.Unmarshal(body, &jsonRawResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}

	log.Debug().Msg("JSON parsed successfully")
	return jsonRawResponse, nil
}
