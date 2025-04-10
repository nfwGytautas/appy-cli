package shared

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version string `yaml:"version"`
	Project string `yaml:"project"`
	Type    string `yaml:"type"`
	Module  string `yaml:"module"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	yamlFile, err := os.ReadFile(".appy/appy.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Save() error {
	yamlFile, err := os.OpenFile(".appy/appy.yaml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	err = yaml.NewEncoder(yamlFile).Encode(c)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) GetDomainsRoot(domainName string) string {
	if c.Type == ScaffoldHMD {
		return "domains/" + domainName
	}

	return "domain"
}
