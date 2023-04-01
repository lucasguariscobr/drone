package drone

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"superorbital/drone/utils"

	"github.com/spf13/cobra"
)

var jsonFilePath string

var createCmd = &cobra.Command{
	Use:           "create",
	Aliases:       []string{"c"},
	Short:         "Creates a new drone resource",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          Create,
}

// Create will receive a file path to a JSON file and
// will send a POST HTTP request to the API.
// It validates the input JSON and returns a JSON with the contents of the new drone model.
func Create(cmd *cobra.Command, args []string) error {
	setLogOutput()

	jsonRawData, jsonErr := utils.ReadJsonPayload(jsonFilePath)
	if jsonErr != nil {
		return jsonErr
	}

	validationErr := utils.ValidateDroneModelJson(jsonRawData)
	if validationErr != nil {
		return validationErr
	}

	createUrl, urlError := utils.BuildUrl()
	if urlError != nil {
		return urlError
	}

	httpResponse, httpError := execCreateHttpRequest(createUrl, jsonRawData)
	if httpError != nil {
		return httpError
	}

	if httpResponse.StatusCode != http.StatusCreated {
		return utils.ErrorBuilder(httpResponse.StatusCode)
	}

	defer httpResponse.Body.Close()

	jsonRawResponse, jsonErr := utils.ParseJsonRawResponse(httpResponse.Body)
	if jsonErr != nil {
		return jsonErr
	}

	cmd.Print(string(jsonRawResponse))
	return nil
}

func execCreateHttpRequest(createUrl string, jsonPayload *json.RawMessage) (*http.Response, error) {
	req, httpReqError := http.NewRequest(http.MethodPost, createUrl, nil)
	if httpReqError != nil {
		return nil, httpReqError
	}

	req.Body = io.NopCloser(bytes.NewReader(*jsonPayload))
	return utils.ExecHttpRequest(httpClient, req)
}

func init() {
	createCmd.Flags().StringVarP(&jsonFilePath, "file", "f", "", "JSON file that contains the payload")
	createCmd.MarkFlagRequired("file")
	rootCmd.AddCommand(createCmd)
}
