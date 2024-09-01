package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CacheSize int    `yaml:"cacheSize"`
	Unit      string `yaml:"unit"`
	Bucket    string `yaml:"bucket"`
	Port      int    `yaml:"port"`
}

func New(configPath string) (*Config, error) {
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := Config{}
	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
