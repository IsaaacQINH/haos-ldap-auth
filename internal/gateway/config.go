package gateway

import (
	"log"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type BindUser struct {
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	Server        string            `yaml:"server"`
	TLS           bool              `yaml:"tls"`
	BaseDN        string            `yaml:"basedn"`
	Bind          BindUser          `yaml:"bind"`
	Groups        []string          `yaml:"groups"`
	Mappings      map[string]string `yaml:"mappings"`
	UserAttribute string            `yaml:"user_attribute"`
	Attributes    []string          `yaml:"attributes"`
	Timeout       int               `yaml:"timeout"`
}

func (c *Config) GetConf() *Config {
	ex, err := os.Executable()

	if err != nil {
		log.Printf("Error getting executable path: %v", err)
	}

	yml, err := os.ReadFile(path.Join(path.Dir(ex), "config.yml"))

	if err != nil {
		log.Printf("Error reading config file: %v", err)
	}

	err = yaml.Unmarshal(yml, c)

	if err != nil {
		log.Printf("Error parsing config file: %v", err)
	}

	return c
}
