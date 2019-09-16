package bot

import (
	"encoding/json"
	"io/ioutil"
)

func getAnswers(path string) (
	answers map[string]string,
	err error,
) {
	answersMappingJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(answersMappingJSON, &answers)
	if err != nil {
		return
	}

	return
}
