package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const configFilename = "config.yaml"

type Config struct {
	OpenAIKey string `yaml:"OpenAIKey"`
	DBName    string `yaml:"DBName"`
	GinPort   string `yaml:"GinPort"`
}

// Parse parses the config file and returns a Config struct.
func Parse() (*Config, error) {
	buf, err := os.ReadFile(configFilename)
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return c, nil
}
