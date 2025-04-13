package shared

import (
	"fmt"
	"os"
	"text/template"

	"github.com/nfwGytautas/appy-cli/providers"
	"github.com/nfwGytautas/appy-cli/templates"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Version   string           `yaml:"version"`
	Project   string           `yaml:"project"`
	Type      string           `yaml:"type"`
	Module    string           `yaml:"module"`
	Providers providers.Config `yaml:"providers"`
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

func (c *Config) Reconfigure() error {
	// Providers
	err := c.Providers.Configure()
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

	return nil
}

func (c *Config) RunHook(hookName string, data any) error {
	fmt.Println("Running hook", hookName)
	return c.Providers.RunHook(hookName, data)
}
