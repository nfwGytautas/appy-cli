package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/utils"
	"gopkg.in/yaml.v3"
)

type RepositoryConfigureOpts map[string]any

type Repository struct {
	Url    string `yaml:"url"`
	Branch string `yaml:"branch"`

	Providers []Provider  `yaml:"providers"`
	config    *AppyConfig `yaml:"-"`
}

func (r *Repository) GetEnabledProviders() []Provider {
	enabledProviders := []Provider{}
	for _, provider := range r.Providers {
		if provider.Enabled {
			provider.repo = r
			enabledProviders = append(enabledProviders, provider)
		}
	}

	return enabledProviders
}

func (r *Repository) Configure(opts RepositoryConfigureOpts) error {
	// Check repositories directory exists
	if _, err := os.Stat(".appy/repositories"); os.IsNotExist(err) {
		err := os.MkdirAll(".appy/repositories", 0755)
		if err != nil {
			return fmt.Errorf("failed to create repositories directory: %v", err)
		}
	}

	// Resolve providers
	utils.Console.DebugLn("Resolving providers...")

	repoName := strings.Split(r.Url, "/")[len(strings.Split(r.Url, "/"))-1]
	repositoryPath := fmt.Sprintf(".appy/repositories/%s", repoName)

	_, err := os.Stat(repositoryPath)
	if err == nil {
		utils.Console.DebugLn("Repository `%s` exists, checking out branch `%s`...", r.Url, r.Branch)

		// Checkout branch
		err = utils.RunCommand(repositoryPath, "git checkout "+r.Branch)
		if err != nil {
			return fmt.Errorf("failed to checkout branch: (%v)", err)
		}

		// Pull repository
		err = utils.RunCommand(repositoryPath, "git pull")
		if err != nil {
			return fmt.Errorf("failed to pull repository: (%v)", err)
		}
	} else {
		utils.Console.InfoLn("Repository `%s` does not exist, pulling...", r.Url)

		// Pull repository
		err = utils.RunCommand("", "git clone --branch "+r.Branch+" "+r.Url+" "+repositoryPath)
		if err != nil {
			return fmt.Errorf("failed to pull repository: (%v)", err)
		}
	}

	utils.Console.DebugLn("Parsing `%s` providers", r.Url)

	// Check root level for appy.yaml
	rootConfigPath := filepath.Join(repositoryPath, "appy.yaml")
	if _, err := os.Stat(rootConfigPath); err == nil {
		// Read root level config
		configData, err := os.ReadFile(rootConfigPath)
		if err != nil {
			return fmt.Errorf("failed to read root level config: %v", err)
		}

		// Template the contents of config
		filledData, err := utils.TemplateAString(string(configData), opts)
		if err != nil {
			return fmt.Errorf("failed to template config: %v", err)
		}

		var providerConfig Provider
		err = yaml.Unmarshal([]byte(filledData), &providerConfig)
		if err != nil {
			return fmt.Errorf("failed to parse root level config: %v", err)
		}

		// Add providers from root config
		providerConfig.ProviderAtRoot = true
		providerConfig.Enabled = false
		providerConfig.Path = repositoryPath

		// Check if provider is already in list
		r.registerProvider(providerConfig)
	} else {
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
				utils.Console.DebugLn("Found provider at %s", configPath)

				// Template the contents of config
				filledData, err := utils.TemplateAString(string(configData), opts)
				if err != nil {
					return fmt.Errorf("failed to template config in %s", entry.Name())
				}

				var providerConfig Provider
				err = yaml.Unmarshal([]byte(filledData), &providerConfig)
				if err != nil {
					return fmt.Errorf("failed to parse config in %s: %v", entry.Name(), err)
				}

				// Add providers from root config
				providerConfig.ProviderAtRoot = false
				providerConfig.Enabled = false
				providerConfig.Path = subdirPath

				r.registerProvider(providerConfig)
			}
		}
	}

	// Configure providers
	utils.Console.InfoLn("Configuring providers...")
	for _, provider := range r.Providers {
		provider.repo = r
		utils.Console.InfoLn("   + `%s@%s`", provider.Name, provider.Version)
		err := provider.Configure(opts)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) RunHook(hookName string, data any) error {
	for _, provider := range r.Providers {
		err := provider.RunLocalHook(hookName, data)
		if err != nil {
			utils.Console.ErrorLn("failed to run hook `%s` for provider `%s`: %v\n", hookName, provider.Name, err)
			return err
		}
	}

	return nil
}

func (r *Repository) RestartWatchers() error {
	for _, provider := range r.Providers {
		err := provider.RestartWatchers()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) StopWatchers() error {
	for _, provider := range r.Providers {
		err := provider.StopWatchers()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) StartWatchers() error {
	for _, provider := range r.Providers {
		err := provider.StartWatchers()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) ApplyStringSubstitution(str string) string {
	return r.config.ApplyStringSubstitution(str)
}

func (r *Repository) registerProvider(provider Provider) {
	for it, existingProvider := range r.Providers {
		if existingProvider.Name == provider.Name {
			if provider.Version != existingProvider.Version {
				prompt := promptui.Prompt{
					Label:     fmt.Sprintf("Provider `%s` version mismatch local: %s, remote: %s, update", provider.Name, existingProvider.Version, provider.Version),
					Default:   "Y",
					IsConfirm: true,
				}

				result, err := prompt.Run()
				if err != nil {
					utils.Console.ErrorLn("failed to run prompt: %v", err)
					return
				}

				utils.Console.ClearLines(1)

				if result == "" || result == "Y" || result == "y" {
					// Reconfigure provider
					// TODO: Migration
					provider.Enabled = r.Providers[it].Enabled
					r.Providers[it] = provider

					// Delete provider directory if exists
					// TODO: Versioning?? this is too aggressive
					err := provider.DeleteConfiguration()
					if err != nil {
						utils.Console.ErrorLn("failed to delete provider configuration: %v", err)
						return
					}
				}
			}

			return
		}
	}

	utils.Console.InfoLn("Registered provider '%s@%s'", provider.Name, provider.Version)
	r.Providers = append(r.Providers, provider)
}
