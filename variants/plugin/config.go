package variant_plugin

import (
	"os"

	"github.com/nfwGytautas/appy-cli/plugins"
	variant_base "github.com/nfwGytautas/appy-cli/variants/base"
	"gopkg.in/yaml.v3"
)

const VariantName = "Plugin"

type Config struct {
	variant_base.Config `yaml:"-,inline"`

	Name        string `yaml:"name"`
	Description string `yaml:"description"`

	engine *plugins.PluginEngine `yaml:"-"`
}

func (c *Config) Save() error {
	yamlFile, err := os.OpenFile(variant_base.YamlFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	err = yaml.NewEncoder(yamlFile).Encode(c)
	if err != nil {
		return err
	}

	return nil
}
