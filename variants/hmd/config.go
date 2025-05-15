package variant_hmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nfwGytautas/appy-cli/plugins"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
	variant_base "github.com/nfwGytautas/appy-cli/variants/base"
	"gopkg.in/yaml.v3"
)

const VariantName = "HMD - Hexagonal Multi Domain"

type Config struct {
	variant_base.Config `yaml:",inline"`

	Project      string        `yaml:"project"`
	Module       string        `yaml:"module"`
	Repositories []*Repository `yaml:"repositories"`

	Workspace string                `yaml:"-"`
	BuildDir  string                `yaml:"-"`
	Plugins   *plugins.PluginEngine `yaml:"-"`
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

func (c *Config) Reconfigure() error {
	if c.Plugins != nil {
		c.Plugins.Shutdown()
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c.Workspace = cwd
	c.BuildDir = filepath.Join(cwd, ".appy", "build")
	c.Plugins = plugins.NewPluginEngine(map[string]any{
		"module":  c.Module,
		"project": c.Project,
	})

	// Create build directory
	err = os.MkdirAll(c.BuildDir, 0755)
	if err != nil {
		return err
	}

	enabledProviders := []Provider{}

	for _, repository := range c.Repositories {
		repository.config = c

		// Providers
		err = repository.Configure(RepositoryConfigureOpts{
			"Config":    c,
			"Workspace": cwd,
			"BuildDir":  c.BuildDir,
		})
		if err != nil {
			return err
		}

		enabledProviders = append(enabledProviders, repository.GetEnabledProviders()...)
	}

	utils.Console.DebugLn("Reconfiguring providers")

	// Adapt providers
	f, err := os.OpenFile("providers/providers.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	// Write template
	tmpl := utils.NewTemplate(templates.ProvidersGo)

	err = tmpl.Execute(f, map[string]any{
		"Providers": enabledProviders,
		"Module":    c.Module,
	})
	if err != nil {
		f.Close()
		return err
	}
	f.Close()

	err = utils.RunTools("providers/providers.go", []string{
		shared.ToolGoImports,
		shared.ToolGoFmt,
	})
	if err != nil {
		return err
	}

	// Run go mod tidy
	err = utils.RunCommand(cwd, shared.ToolGoModTidy)
	if err != nil {
		return err
	}

	// Load plugin hooks
	for _, p := range c.Plugins.GetLoadedPlugins() {
		err = p.OnLoad()
		if err != nil {
			return fmt.Errorf("failed to run 'onLoad' hook: %v", err)
		}
	}

	// Save
	err = c.Save()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) LoadAndReconfigure() error {
	yamlFileContents, err := os.ReadFile(variant_base.YamlFilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFileContents, c)
	if err != nil {
		return err
	}

	err = c.Reconfigure()
	if err != nil {
		return err
	}

	return nil
}
