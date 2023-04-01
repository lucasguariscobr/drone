package drone

import (
	"net/http"

	"superorbital/drone/utils"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:           "list",
	Aliases:       []string{"l"},
	Short:         "List all drones in your collection",
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          List,
}

// List will return all existing models using a JSON format.
func List(cmd *cobra.Command, args []string) error {
	setLogOutput()

	listUrl, urlError := utils.BuildUrl()
	if urlError != nil {
		return urlError
	}

	httpResponse, httpError := execListHttpRequest(listUrl)
	if httpError != nil {
		return httpError
	}

	if httpResponse.StatusCode != http.StatusOK {
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

func execListHttpRequest(listUrl string) (*http.Response, error) {
	req, httpReqError := http.NewRequest(http.MethodGet, listUrl, nil)
	if httpReqError != nil {
		return nil, httpReqError
	}

	return utils.ExecHttpRequest(httpClient, req)
}

func init() {
	rootCmd.AddCommand(listCmd)
}
