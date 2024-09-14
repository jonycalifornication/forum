package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	GithubClientID     string `json:"client_id_github"`
	GithubClientSecret string `json:"client_secret_github"`
	GoogleClientID     string `json:"client_id_google"`
	GoogleClientSecret string `json:"client_secret_google"`
}

var Cfg Config

func LoadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal([]byte(data), &Cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &Cfg, nil
}
