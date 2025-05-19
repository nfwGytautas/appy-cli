package variants

import (
	"context"
	"errors"
	"maps"
	"os"
	"slices"

	variant_base "github.com/nfwGytautas/appy-cli/variants/base"
	variant_hmd "github.com/nfwGytautas/appy-cli/variants/hmd"
	variant_plugin "github.com/nfwGytautas/appy-cli/variants/plugin"
	"gopkg.in/yaml.v3"
)

type VariantFactory func() Variant

var variantTypes = map[string]VariantFactory{
	variant_hmd.VariantName:    func() Variant { return &variant_hmd.Config{} },
	variant_plugin.VariantName: func() Variant { return &variant_plugin.Config{} },
}

// Variant is an interface that defines the methods that all variants must implement
type Variant interface {
	Scaffold() error
	Start(context.Context) error
}

// Loads a variant from current directory, returns nil if empty directory
func Load() (Variant, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	// If empty directory, return nil
	if len(entries) == 0 || (len(entries) == 1 && entries[0].Name() == ".git") {
		return nil, nil
	}

	// Load the base variant to figure out the type
	variantType, err := variant_base.GetType()
	if err != nil {
		return nil, err
	}

	// Check if the variant type is supported
	variantFactory, ok := variantTypes[variantType]
	if !ok {
		return nil, errors.New("unsupported variant type " + variantType)
	}

	variant := variantFactory()

	// Load the variant config
	yamlFileContents, err := os.ReadFile(variant_base.YamlFilePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFileContents, variant)
	if err != nil {
		return nil, err
	}

	return variant, nil
}

// GetVariantType returns all possible variant types
func GetVariantTypes() []string {
	return slices.Collect(maps.Keys(variantTypes))
}

// CreateEmptyVariant creates an empty variant of the given type
func CreateEmptyVariant(variantType string) (Variant, error) {
	variantFactory, ok := variantTypes[variantType]
	if !ok {
		return nil, errors.New("unsupported variant type " + variantType)
	}

	variant := variantFactory()

	return variant, nil
}
