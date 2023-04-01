package drone

import (
	"bytes"
	"net/http"
	"superorbital/drone/utils"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var SUCCESS_RESPONSE string = `{"name": "rubyred","type": "quadcopter-large","plan": ["land-drone"],"code": "en-route","instructionIndex": 0,"cost": {"amount": 1000,"currency": "USD","amountDecimalShift": -2}}`

func TestMissingAddrListCmd(t *testing.T) {
	cases := map[string]struct {
		e      error
		output string
	}{
		"missingAddr": {
			e:      utils.ErrMissingAddr,
			output: "",
		},
	}

	for _, value := range cases {
		viper.Reset()
		buildMockHttpClient(utils.BuildNilTestResponse())
		cmdResponse, cmdErr := callListCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func TestMissingTokenListCmd(t *testing.T) {
	cases := map[string]struct {
		e      error
		output string
	}{
		"missingToken": {
			e:      utils.ErrMissingToken,
			output: "",
		},
	}

	for _, value := range cases {
		viper.Reset()
		viper.Set(utils.CONFIG_VALUE_ADDR, "ADDR")
		buildMockHttpClient(utils.BuildNilTestResponse())
		cmdResponse, cmdErr := callListCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func TestHttpErrorListCmd(t *testing.T) {
	cases := map[string]struct {
		httpStatusCode int
		e              error
		output         string
	}{
		"badRequest": {
			httpStatusCode: http.StatusBadRequest,
			e:              utils.ErrBadRequest,
			output:         "",
		},
		"unauthorized": {
			httpStatusCode: http.StatusUnauthorized,
			e:              utils.ErrUnauthorized,
			output:         "",
		},
		"tooManyRequests": {
			httpStatusCode: http.StatusTooManyRequests,
			e:              utils.ErrTooManyRequests,
			output:         "",
		},
		"internalServerError": {
			httpStatusCode: http.StatusInternalServerError,
			e:              utils.ErrInternalServer,
			output:         "",
		},
	}

	for _, value := range cases {
		viper.Reset()
		viper.Set(utils.CONFIG_VALUE_ADDR, "ADDR")
		viper.Set(utils.CONFIG_VALUE_TOKEN, "TOKEN")
		buildMockHttpClient(utils.BuildTestResponse(value.httpStatusCode, value.output))
		cmdResponse, cmdErr := callListCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func TestListCmd(t *testing.T) {
	cases := map[string]struct {
		e      error
		output string
	}{
		"success": {
			e:      nil,
			output: SUCCESS_RESPONSE,
		},
	}

	for _, value := range cases {
		viper.Reset()
		viper.Set(utils.CONFIG_VALUE_ADDR, "ADDR")
		viper.Set(utils.CONFIG_VALUE_TOKEN, "TOKEN")
		buildMockHttpClient(utils.BuildTestResponse(http.StatusOK, SUCCESS_RESPONSE))
		cmdResponse, cmdErr := callListCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func callListCmd() (*bytes.Buffer, error) {
	cmdResponse := new(bytes.Buffer)
	rootCmd.SetOut(cmdResponse)
	rootCmd.SetErr(cmdResponse)
	rootCmd.SetArgs([]string{"list"})
	cmdErr := rootCmd.Execute()
	return cmdResponse, cmdErr
}

func buildMockHttpClient(responseFunc func(req *http.Request) (*http.Response, error)) {
	mock := utils.MockHttpClient{}
	mock.MockDo = responseFunc
	httpClient = mock
}
