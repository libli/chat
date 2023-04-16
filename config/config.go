package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const configFilename = "config.yaml"

// Config is the config struct.
type Config struct {
	OpenAIKey string      `yaml:"OpenAIKey"`
	GinPort   string      `yaml:"GinPort"`
	Driver    string      `yaml:"Driver"`
	DBName    string      `yaml:"DBName"`
	MySQL     MySQLConfig `yaml:"MySQL"`
}

// MySQLConfig is the MySQL config struct.
type MySQLConfig struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
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

// MySQLDSN 生成 MySQL 的 DSN
func (c *Config) MySQLDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQL.Username, c.MySQL.Password, c.MySQL.Host, c.MySQL.Port, c.MySQL.DBName)
}
