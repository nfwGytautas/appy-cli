package config

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
	"gopkg.in/yaml.v3"
)

type AppyConfig struct {
	Version      string        `yaml:"version"`
	Project      string        `yaml:"project"`
	Type         string        `yaml:"type"`
	Module       string        `yaml:"module"`
	Repositories []*Repository `yaml:"repositories"`

	Workspace string `yaml:"-"`
	BuildDir  string `yaml:"-"`
}

func LoadConfig() (*AppyConfig, error) {
	cfg := &AppyConfig{}

	yamlFile, err := os.ReadFile("appy.yaml")
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
	yamlFile, err := os.OpenFile("appy.yaml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
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
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	c.Workspace = cwd

	c.BuildDir = filepath.Join(cwd, ".appy", "build")

	// Create build directory
	err = os.MkdirAll(c.BuildDir, 0755)
	if err != nil {
		return err
	}

	enabledProviders := []Provider{}

	for _, repository := range c.Repositories {
		repository.config = c
		err := repository.StopWatchers()
		if err != nil {
			return err
		}

		// Providers
		err = repository.Configure(RepositoryConfigureOpts{
			"ProjectName": c.Project,
			"Workspace":   cwd,
			"BuildDir":    c.BuildDir,
		})
		if err != nil {
			return err
		}

		enabledProviders = append(enabledProviders, repository.GetEnabledProviders()...)
	}

	utils.Console.DebugLn("Reconfiguring main.go")

	// Adapt main.go
	f, err := os.OpenFile("main.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write template
	tmpl := template.Must(template.New("main.go").Parse(templates.MainWithProvidersGo))

	err = tmpl.Execute(f, map[string]any{
		"Providers": enabledProviders,
		"Module":    c.Module,
	})
	if err != nil {
		return err
	}

	// Run go mod tidy
	err = utils.RunCommand(cwd, strings.Split(shared.ToolGoModTidy, " ")...)
	if err != nil {
		return err
	}

	// Save
	err = c.Save()
	if err != nil {
		return err
	}

	err = c.StartProviders()
	if err != nil {
		return err
	}

	return nil
}

func (c *AppyConfig) RunHook(hookName string, data any) error {
	utils.Console.InfoLn("Running hook %s", hookName)

	for _, repository := range c.Repositories {
		err := repository.RunHook(hookName, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *AppyConfig) ApplyStringSubstitution(str string) string {
	str = strings.ReplaceAll(str, "${Workspace}", c.Workspace)
	str = strings.ReplaceAll(str, "${ProjectName}", c.Project)
	str = strings.ReplaceAll(str, "${Module}", c.Module)
	str = strings.ReplaceAll(str, "${BuildDir}", c.BuildDir)
	return str
}

func (c *AppyConfig) StartProviders() error {
	for _, repository := range c.Repositories {
		err := repository.StartWatchers()
		if err != nil {
			return err
		}
	}

	return nil
}
