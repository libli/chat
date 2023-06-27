package config

import (
	"chat/model"
	"fmt"
	"log"
	"net/http"
	"os"

	roundrobin "github.com/thegeekyasian/round-robin-go"
	"github.com/tidwall/gjson"
	"github.com/ysmood/gop"
	"sigs.k8s.io/yaml"
)

const configFilename = "config.yaml"
const GJSON_KEY_OPENAIKEY = `OpenAIKey`

type Config struct {
	DBName    string       `yaml:"DBName"`
	GinPort   string       `yaml:"GinPort"`
	InitUsers []model.User `yaml:"InitUsers"`

	openAIKeys []*string
	rr         *roundrobin.RoundRobin[string]
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

	jsonBuf, _ := yaml.YAMLToJSON(buf)
	for _, v := range gjson.GetBytes(jsonBuf, GJSON_KEY_OPENAIKEY).Array() {
		vStr := v.String()
		c.openAIKeys = append(c.openAIKeys, &vStr)
	}

	if len(c.openAIKeys) < 1 {
		log.Fatalf("Plz Set at least one `%s`\n", GJSON_KEY_OPENAIKEY)
	}
	c.rr, _ = roundrobin.New(c.openAIKeys...)

	gop.P(c)
	return c, nil
}

func (c *Config) AllKeys() []string {
	var keys []string
	for k := range c.openAIKeys {
		keys = append(keys, *c.openAIKeys[k])
	}
	return keys
}

func (c *Config) RoundRobinKey(*http.Request) string {
	return *c.rr.Next()
}

// TODO
// func (c *Config) HashMapKey(*http.Request) string {
// }
