package config

import (
	"chat/model"
	"fmt"
	"net/http"
	"os"
	"sync"

	roundrobin "github.com/thegeekyasian/round-robin-go"
	"github.com/ysmood/gop"
	"gopkg.in/yaml.v3"
)

const configFilename = "config.yaml"

var w sync.WaitGroup

type Config struct {
	OpenAIKey []*string    `yaml:"OpenAIKey"`
	DBName    string       `yaml:"DBName"`
	GinPort   string       `yaml:"GinPort"`
	InitUsers []model.User `yaml:"InitUsers"`
	rr        *roundrobin.RoundRobin[string]
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
	c.rr, _ = roundrobin.New(c.OpenAIKey...)
	gop.P(c)
	return c, nil
}

func (c *Config) RoundRobinKey(*http.Request) string {
	return *c.rr.Next()
}

// TODO
// func (c *Config) HashMapKey(*http.Request) string {
// }
