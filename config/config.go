package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DatabaseURL string `yaml:"database_url"`
	MediaRoot   string `yaml:"media_root"`
	TargetFile  string `yaml:"target_file"` // full path to the target archive (tar gz file)
	AppVersion  string `yaml:"app_version"` // Papermerge DMS version this config is intended for
}

func ReadConfig(fileName string) (*Config, error) {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
