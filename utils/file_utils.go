package utils

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

var osReadFileFunc = os.ReadFile

// ReadJsonPayload receives a file path, reads it, and parses its content to
// a JSON RawMessage
func ReadJsonPayload(jsonFilePath string) (*json.RawMessage, error) {
	log.Debug().Msgf("File: %s", jsonFilePath)

	jsonFileContent, osErr := osReadFileFunc(jsonFilePath)
	if osErr != nil {
		return nil, osErr
	}

	var jsonRawData json.RawMessage
	jsonErr := json.Unmarshal(jsonFileContent, &jsonRawData)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &jsonRawData, nil
}

func setOsReadFileFunc(readFileFunc func(name string) ([]byte, error)) {
	osReadFileFunc = readFileFunc
}
