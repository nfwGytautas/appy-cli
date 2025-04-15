package config

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"gopkg.in/yaml.v3"
)

type AppyConfig struct {
	Version   string          `yaml:"version"`
	Project   string          `yaml:"project"`
	Type      string          `yaml:"type"`
	Module    string          `yaml:"module"`
	Providers ProvidersConfig `yaml:"providers"`

	Workspace string `yaml:"-"`
}

func LoadConfig() (*AppyConfig, error) {
	cfg := &AppyConfig{}

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

func (c *AppyConfig) Save() error {
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

func (c *AppyConfig) GetDomainsRoot(domainName string) string {
	if c.Type == shared.ScaffoldHMD {
		return "domains/" + domainName
	}

	return "domain"
}

func (c *AppyConfig) Reconfigure() error {
	c.Providers.config = c
	err := c.Providers.StopWatchers()
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c.Workspace = cwd

	// Providers
	err = c.Providers.Configure(ProviderConfigureOpts{
		"ProjectName": c.Project,
		"Workspace":   cwd,
	})
	if err != nil {
		return err
	}

	fmt.Println("Reconfiguring main.go")

	// Adapt main.go
	f, err := os.OpenFile("main.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write template
	tmpl := template.Must(template.New("main.go").Parse(templates.MainWithProvidersGo))

	err = tmpl.Execute(f, map[string]any{
		"Providers": c.Providers.GetEnabledProviders(),
		"Module":    c.Module,
	})
	if err != nil {
		return err
	}

	// Save
	err = c.Save()
	if err != nil {
		return err
	}

	err = c.Providers.StartWatchers()
	if err != nil {
		return err
	}

	return nil
}

func (c *AppyConfig) RunHook(hookName string, data any) error {
	fmt.Println("Running hook", hookName)
	return c.Providers.RunHook(hookName, data)
}

func (c *AppyConfig) ApplyStringSubstitution(str string) string {
	str = strings.ReplaceAll(str, "${Workspace}", c.Workspace)
	str = strings.ReplaceAll(str, "${ProjectName}", c.Project)
	str = strings.ReplaceAll(str, "${Module}", c.Module)
	return str
}

func (c *AppyConfig) StartProviders() error {
	return c.Providers.StartWatchers()
}
