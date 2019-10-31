package utils

import (
	"encoding/json"
	"io/ioutil"
)

// GetMapFromJSON gets a map[string]string from JSON file
func GetMapFromJSON(pathToJSONFile string) (
	result map[string]string,
	err error,
) {
	mappingJSON, err := ioutil.ReadFile(pathToJSONFile)
	if err != nil {
		return
	}
	err = json.Unmarshal(mappingJSON, &result)
	if err != nil {
		return
	}

	return
}
