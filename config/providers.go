package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/utils"
	"gopkg.in/yaml.v3"
)

type ProviderConfigureOpts map[string]any

type ProvidersConfig struct {
	Repositories []Repository `yaml:"repositories"`
	Providers    []Provider   `yaml:"providers"`

	config *AppyConfig `yaml:"-"`
}

func (c *ProvidersConfig) GetEnabledProviders() []Provider {
	enabledProviders := []Provider{}
	for _, provider := range c.Providers {
		if provider.Enabled {
			provider.config = c
			enabledProviders = append(enabledProviders, provider)
		}
	}

	return enabledProviders
}

func (c *ProvidersConfig) Configure(opts ProviderConfigureOpts) error {
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
			fmt.Printf("Repository `%s` exists, checking out branch `%s`...\n", repository.Url, repository.Branch)

			// Checkout branch
			cmd := exec.Command("git", "checkout", repository.Branch)

			fmt.Println(cmd.String())

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = repositoryPath
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("failed to checkout branch: %v", err)
			}

			// Pull repository
			cmd = exec.Command("git", "pull")

			fmt.Println(cmd.String())

			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = repositoryPath
			err = cmd.Run()
			if err != nil {
				return fmt.Errorf("failed to pull repository: %v", err)
			}
		} else {
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
				fmt.Println("Found provider at ", configPath)

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
		provider.config = c
		err := provider.Configure(opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ProvidersConfig) RunHook(hookName string, data any) error {
	for _, provider := range c.Providers {
		err := provider.RunLocalHook(hookName, data)
		if err != nil {
			fmt.Printf("failed to run hook `%s` for provider `%s`: %v\n", hookName, provider.Name, err)
			return err
		}
	}

	return nil
}

func (c *ProvidersConfig) RestartWatchers() error {
	for _, provider := range c.Providers {
		err := provider.RestartWatchers()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ProvidersConfig) StopWatchers() error {
	for _, provider := range c.Providers {
		err := provider.StopWatchers()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ProvidersConfig) StartWatchers() error {
	for _, provider := range c.Providers {
		err := provider.StartWatchers()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *ProvidersConfig) ApplyStringSubstitution(str string) string {
	return c.config.ApplyStringSubstitution(str)
}

func (c *ProvidersConfig) registerProvider(provider Provider) {
	for it, existingProvider := range c.Providers {
		if existingProvider.Name == provider.Name {
			if provider.Version != existingProvider.Version {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Provider `%s` version mismatch local: %s, remote: %s, update", provider.Name, existingProvider.Version, provider.Version),
					Default:   "Y",
					IsConfirm: true,
				}

				result, err := prompt.Run()
				if err != nil {
					fmt.Printf("failed to run prompt: %v", err)
					return
				}

				utils.ClearLines(1)

				if result == "" || result == "Y" || result == "y" {
					// Reconfigure provider
					// TODO: Migration
					provider.Enabled = c.Providers[it].Enabled
					c.Providers[it] = provider

					// Delete provider directory if exists
					// TODO: Versioning?? this is too aggressive
					err := provider.DeleteConfiguration()
					if err != nil {
						fmt.Printf("failed to delete provider configuration: %v", err)
						return
					}
				}
			}

			return
		}
	}

	fmt.Printf("Registered provider '%s@%s'\n", provider.Name, provider.Version)
	c.Providers = append(c.Providers, provider)
}
