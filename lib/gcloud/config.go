package gcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/imdario/mergo"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Debug        bool   `json:"debug,omitempty" env:"CWS_DEBUG"`
	ExtID        string `json:"extension_id" env:"CWS_EXTENSION_ID"`
	ID           string `json:"client_id" env:"CWS_CLIENT_ID"`
	Secret       string `json:"client_secret" env:"CWS_CLIENT_SECRET"`
	RefreshToken string `json:"refresh_token" env:"CWS_REFRESH_TOKEN"`
}

func loadConfig(configPath string) (*Config, error) {
	envConf := Config{}
	if err := envconfig.Process(context.Background(), &envConf); err != nil {
		return nil, err
	}
	if _, err := os.Stat(configPath); err == nil {
		data, err := ioutil.ReadFile(configPath)
		if err != nil {
			return nil, err
		}
		fileConf := Config{}
		if err := json.Unmarshal(data, &fileConf); err != nil {
			return nil, err
		}
		if err := mergo.Merge(&envConf, fileConf); err != nil {
			return nil, err
		}
	}
	return &envConf, envConf.validate()
}

func (conf *Config) validate() error {
	missingVals := []string{}
	if conf.ExtID == "" {
		missingVals = append(missingVals, "extension_id")
	}
	if conf.ID == "" {
		missingVals = append(missingVals, "client_id")
	}
	if conf.Secret == "" {
		missingVals = append(missingVals, "secret_id")
	}
	if conf.RefreshToken == "" {
		missingVals = append(missingVals, "refresh_token")
	}
	if len(missingVals) > 0 {
		return fmt.Errorf("Configuration is missing %v which are required for cws to run", strings.Join(missingVals, ", "))
	}
	return nil
}
