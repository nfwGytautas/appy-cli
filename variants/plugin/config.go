package variant_plugin

import variant_base "github.com/nfwGytautas/appy-cli/variants/base"

const VariantName = "Plugin"

type Config struct {
	variant_base.Config

	Name       string `yaml:"name"`
	Version    string `yaml:"version"`
	Descriptor string `yaml:"descriptor"`
}
