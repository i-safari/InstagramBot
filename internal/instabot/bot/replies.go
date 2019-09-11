package bot

import (
	"InstaFollower/internal/instabot/utils"
	"encoding/json"
	"io/ioutil"
)

var (
	answers map[string]string
)

func getAnswers(cfg *utils.Config) error {
	answersMappingJSON, err := ioutil.ReadFile(cfg.PathToDialogs)
	if err != nil {
		return err
	}
	err = json.Unmarshal(answersMappingJSON, &answers)
	if err != nil {
		return err
	}

	return nil
}
