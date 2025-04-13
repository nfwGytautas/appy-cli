package providers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nfwGytautas/appy-cli/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Repositories []Repository `yaml:"repositories"`
	Providers    []Provider   `yaml:"providers"`
}

type Repository struct {
	Url    string `yaml:"url"`
	Branch string `yaml:"branch"`
}

type Provider struct {
	Name           string `yaml:"name"`
	Path           string `yaml:"path"`
	Version        string `yaml:"version"`
	Description    string `yaml:"description"`
	ProviderAtRoot bool   `yaml:"atRootLevel"`
	Enabled        bool   `yaml:"enabled"`
	Hooks          []Hook `yaml:"hooks"`
}

type Hook struct {
	Name   string `yaml:"name"`
	Action string `yaml:"action"`
}

func (c *Config) GetEnabledProviders() []Provider {
	enabledProviders := []Provider{}
	for _, provider := range c.Providers {
		if provider.Enabled {
			enabledProviders = append(enabledProviders, provider)
		}
	}

	return enabledProviders
}

func (c *Config) Configure() error {
	// Check repositories directory exists
	if _, err := os.Stat(".appy/repositories"); os.IsNotExist(err) {
		err := os.MkdirAll(".appy/repositories", 0755)
		if err != nil {
			return fmt.Errorf("failed to create repositories directory: %v", err)
		}
	}

	// Resolve providers
	fmt.Println("Resolving providers...")
	for _, repository := range c.Repositories {
		repoName := strings.Split(repository.Url, "/")[len(strings.Split(repository.Url, "/"))-1]
		repositoryPath := fmt.Sprintf(".appy/repositories/%s", repoName)

		_, err := os.Stat(repositoryPath)
		if err == nil {
			fmt.Printf("Repository `%s` already pulled, skipping...\n", repository.Url)
			continue
		}

		if os.IsNotExist(err) {
			fmt.Printf("Repository `%s` does not exist, pulling...\n", repository.Url)

			// Pull repository
			cmd := exec.Command("git", "clone", "--branch", repository.Branch, repository.Url, repositoryPath)

			fmt.Println(cmd.String())

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("failed to pull repository: %v", err)
			}
		} else {
			return err
		}

		// Check root level for appy.yaml
		rootConfigPath := filepath.Join(repositoryPath, "appy.yaml")
		if _, err := os.Stat(rootConfigPath); err == nil {
			// Read root level config
			configData, err := os.ReadFile(rootConfigPath)
			if err != nil {
				return fmt.Errorf("failed to read root level config: %v", err)
			}

			var providerConfig Provider
			err = yaml.Unmarshal(configData, &providerConfig)
			if err != nil {
				return fmt.Errorf("failed to parse root level config: %v", err)
			}

			// Add providers from root config
			providerConfig.ProviderAtRoot = true
			providerConfig.Enabled = false
			providerConfig.Path = repositoryPath

			// Check if provider is already in list
			c.registerProvider(providerConfig)
			continue
		}

		// Check subdirectories for appy.yaml
		entries, err := os.ReadDir(repositoryPath)
		if err != nil {
			return fmt.Errorf("failed to read repository directory: %v", err)
		}

		for _, entry := range entries {
			if !entry.IsDir() || entry.Name()[0] == '.' {
				continue
			}

			subdirPath := filepath.Join(repositoryPath, entry.Name())
			configPath := filepath.Join(subdirPath, "appy.yaml")

			if _, err := os.Stat(configPath); err == nil {
				// Read subdirectory config
				configData, err := os.ReadFile(configPath)
				if err != nil {
					return fmt.Errorf("failed to read config in %s: %v", entry.Name(), err)
				}
				log.Println("Found provider at ", configPath)

				var providerConfig Provider
				err = yaml.Unmarshal(configData, &providerConfig)
				if err != nil {
					return fmt.Errorf("failed to parse config in %s: %v", entry.Name(), err)
				}

				// Add providers from root config
				providerConfig.ProviderAtRoot = false
				providerConfig.Enabled = false
				providerConfig.Path = subdirPath

				c.registerProvider(providerConfig)
			}
		}
	}

	// Configure providers
	for _, provider := range c.Providers {
		err := provider.Configure()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) RunHook(hookName string, data any) error {
	for _, provider := range c.Providers {
		for _, hook := range provider.Hooks {
			if hook.Name == hookName {
				return hook.Run(&provider, data)
			}
		}
	}

	return nil
}

func (c *Config) registerProvider(provider Provider) {
	log.Println("Registering provider", provider.Name)
	for _, existingProvider := range c.Providers {
		if existingProvider.Name == provider.Name {
			return
		}
	}

	fmt.Printf("Registered provider '%s'\n", provider.Name)
	c.Providers = append(c.Providers, provider)
}

func (p *Provider) Configure() error {
	if !p.Enabled {
		return nil
	}

	// Check if provider is already configured
	if p.IsConfigured() {
		return nil
	}

	fmt.Println("Configuring provider", p.Name)

	// Copy provider files from repository to providers directory
	sourceDir := filepath.Join(p.Path, "providers")
	destDir := filepath.Join("providers", p.Name)

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create provider directory: %v", err)
	}

	// Copy all files from source to destination
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself
		if path == sourceDir {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, info.Mode())
	})

	if err != nil {
		return fmt.Errorf("failed to copy provider files: %v", err)
	}

	return nil
}

func (p *Provider) IsConfigured() bool {
	_, err := os.Stat(fmt.Sprintf("providers/%s", p.Name))
	return err == nil
}

func (h *Hook) Run(provider *Provider, data any) error {
	if h.Action == "appy copyTemplate" {
		return h.copyTemplate(provider, data)
	}

	// TODO: CLI tools

	return nil
}

func (h *Hook) copyTemplate(provider *Provider, data any) error {
	// Get the domain name from the data
	domainData, ok := data.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid data type for domain creation hook")
	}

	domainName, ok := domainData["DomainName"].(string)
	if !ok {
		return fmt.Errorf("domain name not found in hook data")
	}

	// Source and destination paths
	sourceDir := filepath.Join(provider.Path, "domain")
	destDir := filepath.Join("domains", domainName, "adapters", provider.Name)

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create provider directory: %v", err)
	}

	// Copy the contents
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the source directory itself
		if path == sourceDir {
			return nil
		}

		// Calculate relative path
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if info.IsDir() {
			// Create directory
			return os.MkdirAll(destPath, info.Mode())
		}

		// Copy file
		templateData, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		file, err := os.Create(destPath)
		if err != nil {
			return err
		}

		// Write template
		tmpl := utils.NewTemplate(string(templateData))

		err = tmpl.Execute(file, data)
		if err != nil {
			file.Close()
			return err
		}
		file.Close()

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy provider domain files: %v", err)
	}

	return nil
}
