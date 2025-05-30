package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ListenPort int      `yaml:"listenPort"`
	Strategy   string   `yaml:"strategy"`
	LogLevel   string   `yaml:"logLevel"`
	Backends   []string `yaml:"backends"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
