package model

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectName    string          `yaml:"project"`
	Variables      []Variable      `yaml:"variables"`
	Middlewares    []Middleware    `yaml:"middlewares"`
	EndpointGroups []EndpointGroup `yaml:"endpointGroups"`
	Endpoints      []Endpoint      `yaml:"endpoints"`
}

func ReadConfig() (*Config, error) {
	cfg := Config{}

	// Check if appy.yaml exists
	if _, err := os.Stat("appy.yaml"); os.IsNotExist(err) {
		return nil, errors.New("appy.yaml not found")
	}

	contents, err := os.ReadFile("appy.yaml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(contents, &cfg)
	if err != nil {
		return nil, err
	}

	cfg.postLoad()

	return &cfg, nil
}

func (c Config) ReplaceWithVariables(in string) string {
	for _, v := range c.Variables {
		in = strings.ReplaceAll(in, fmt.Sprintf("{{ .%s }}", v.Name), fmt.Sprintf("appyVar_%s", v.Name))
	}

	return in
}

func (c *Config) postLoad() {
	for i := range c.Endpoints {
		c.Endpoints[i].ResolveChildren()
	}
}
