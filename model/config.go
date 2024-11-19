package model

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectName    string          `yaml:"project"`
	Middlewares    []Middleware    `yaml:"middlewares"`
	EndpointGroups []EndpointGroup `yaml:"endpointGroups"`
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

	return &cfg, nil
}
