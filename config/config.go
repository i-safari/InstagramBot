package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config ...
type Config struct {
	Username string `json:"username"`
	Password string `json:"password"`
	LogPath  string `json:"logpath"`
	Vault    string `json:"vault"`
	Token    string `json:"token"`
}

// CreateConfig ...
func CreateConfig(path string) (*Config, error) {
	configJSON, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(configJSON, config); err != nil {
		return nil, err
	}

	return config, nil
}
