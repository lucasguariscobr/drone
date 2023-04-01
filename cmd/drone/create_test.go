package drone

import (
	"bytes"
	"net/http"
	"superorbital/drone/utils"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const CREATE_JSON_FILE = "test.json"

var minimumDroneModel = `{"instructionIndex":0,"name":"Test Drone","plan":["land-drone"],"type":"quadcopter-small"}`

func TestMissingAddrCreateCmd(t *testing.T) {
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
		utils.MockOsReadFile(minimumDroneModel)
		buildMockHttpClient(utils.BuildNilTestResponse())
		cmdResponse, cmdErr := callCreateCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func TestMissingTokenCreateCmd(t *testing.T) {
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
		utils.MockOsReadFile(minimumDroneModel)
		buildMockHttpClient(utils.BuildNilTestResponse())
		cmdResponse, cmdErr := callCreateCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func TestValidationCreateCmd(t *testing.T) {
	cases := map[string]struct {
		jsonPayload string
		e           error
		output      string
	}{
		"invalidInstructionIndex": {
			jsonPayload: `{"instructionIndex":"text","name":"Test Drone","plan":["land-drone"],"type":"quadcopter-small"}`,
			e:           utils.ErrCreateDroneInstructionIndex,
			output:      "",
		},
		"missingPlan": {
			jsonPayload: `{"instructionIndex":1,"name":"Test","type":"quadcopter-small"}`,
			e:           utils.ErrCreateDronePlanLength,
			output:      "",
		},
		"emptyPlan": {
			jsonPayload: `{"instructionIndex":1,"name":"Test","plan":[],"type":"quadcopter-small"}`,
			e:           utils.ErrCreateDronePlanLength,
			output:      "",
		},
		"invalidLastInstruction": {
			jsonPayload: `{"instructionIndex":1,"name":"Test","plan":["up"],"type":"quadcopter-small"}`,
			e:           utils.ErrCreateDronePlanLastInstruction,
			output:      "",
		},
		"missingType": {
			jsonPayload: `{"instructionIndex":1,"name":"Test","plan":["land-drone"]}`,
			e:           utils.ErrCreateDroneMissingType,
			output:      "",
		},
		"invalidType": {
			jsonPayload: `{"instructionIndex":1,"name":"Test","plan":["land-drone"],"type":"plane-jumbo"}`,
			e:           utils.ErrCreateDroneType,
			output:      "",
		},
		"containsCost": {
			jsonPayload: `{"instructionIndex":1,"name":"Test","plan":["land-drone"],"type":"plane-jumbo","cost":{"amount":0,"amountDecimalShift":0,"currency":"USD"}}`,
			e:           utils.ErrCreateDroneCost,
			output:      "",
		},
		"containsStatus": {
			jsonPayload: `{"instructionIndex":1,"name":"Test","plan":["land-drone"],"type":"plane-jumbo","status":"starting-up"}`,
			e:           utils.ErrCreateDroneStatus,
			output:      "",
		},
	}

	for _, value := range cases {
		viper.Reset()
		viper.Set(utils.CONFIG_VALUE_ADDR, "ADDR")
		viper.Set(utils.CONFIG_VALUE_TOKEN, "TOKEN")
		utils.MockOsReadFile(value.jsonPayload)
		buildMockHttpClient(utils.BuildNilTestResponse())
		cmdResponse, cmdErr := callCreateCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func TestHttpErrorCreateCmd(t *testing.T) {
	cases := map[string]struct {
		jsonPayload    string
		httpStatusCode int
		e              error
		output         string
	}{
		"badRequest": {
			jsonPayload:    minimumDroneModel,
			httpStatusCode: http.StatusBadRequest,
			e:              utils.ErrBadRequest,
			output:         "",
		},
		"unauthorized": {
			jsonPayload:    minimumDroneModel,
			httpStatusCode: http.StatusUnauthorized,
			e:              utils.ErrUnauthorized,
			output:         "",
		},
		"tooManyRequests": {
			jsonPayload:    minimumDroneModel,
			httpStatusCode: http.StatusTooManyRequests,
			e:              utils.ErrTooManyRequests,
			output:         "",
		},
		"internalServerError": {
			jsonPayload:    minimumDroneModel,
			httpStatusCode: http.StatusInternalServerError,
			e:              utils.ErrInternalServer,
			output:         "",
		},
	}

	for _, value := range cases {
		viper.Reset()
		viper.Set(utils.CONFIG_VALUE_ADDR, "ADDR")
		viper.Set(utils.CONFIG_VALUE_TOKEN, "TOKEN")
		utils.MockOsReadFile(value.jsonPayload)
		buildMockHttpClient(utils.BuildTestResponse(value.httpStatusCode, value.output))
		cmdResponse, cmdErr := callCreateCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func TestCreateCmd(t *testing.T) {
	cases := map[string]struct {
		jsonPayload string
		e           error
		output      string
	}{
		"success": {
			jsonPayload: minimumDroneModel,
			e:           nil,
			output:      minimumDroneModel,
		},
	}

	for _, value := range cases {
		viper.Reset()
		viper.Set(utils.CONFIG_VALUE_ADDR, "ADDR")
		viper.Set(utils.CONFIG_VALUE_TOKEN, "TOKEN")
		utils.MockOsReadFile(value.jsonPayload)
		buildMockHttpClient(utils.BuildTestResponse(http.StatusCreated, value.output))
		cmdResponse, cmdErr := callCreateCmd()
		assert.Equal(t, value.e, cmdErr)
		assert.Equal(t, value.output, cmdResponse.String())
	}

}

func callCreateCmd() (*bytes.Buffer, error) {
	cmdResponse := new(bytes.Buffer)
	rootCmd.SetOut(cmdResponse)
	rootCmd.SetErr(cmdResponse)
	rootCmd.SetArgs([]string{"create", "--file", CREATE_JSON_FILE})
	cmdErr := rootCmd.Execute()
	return cmdResponse, cmdErr
}
