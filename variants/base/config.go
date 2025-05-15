package variant_base

import (
	"os"

	"gopkg.in/yaml.v3"
)

const YamlFilePath = "appy.yaml"

type Config struct {
	Version string `yaml:"version"`
	Type    string `yaml:"type"`
}

func GetType() (string, error) {
	base := &Config{}

	yamlFileContents, err := os.ReadFile(YamlFilePath)
	if err != nil {
		return "", err
	}

	err = yaml.Unmarshal(yamlFileContents, base)
	if err != nil {
		return "", err
	}

	return base.Type, nil
}
