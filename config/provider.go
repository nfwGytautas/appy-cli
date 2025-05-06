package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nfwGytautas/appy-cli/plugins"
	"github.com/nfwGytautas/appy-cli/utils"
)

type Provider struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Enabled     bool   `yaml:"enabled"`

	repo   *Repository     `yaml:"-"`
	plugin *plugins.Plugin `yaml:"-"`
}

func (p *Provider) Configure(opts RepositoryConfigureOpts) error {
	if !p.Enabled {
		return nil
	}

	// Load the plugin
	var err error
	pluginPath := filepath.Join(p.Path, "plugin.lua")
	p.plugin, err = p.repo.config.Plugins.LoadPlugin(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to load plugin: %v", err)
	}

	if p.plugin != nil {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %v", err)
		}

		p.plugin.SetMetaFields(plugins.PluginMetaFields{
			ScriptRoot:   filepath.Join(cwd, p.Path) + "/",
			ProviderRoot: filepath.Join(cwd, "providers", p.Name) + "/",
		})
	}

	// Check if provider is already configured
	if p.IsConfigured() {
		return nil
	}

	utils.Console.DebugLn("Configuring provider `%s@%s`...", p.Name, p.Version)

	// Copy provider files from repository to providers directory
	sourceDir := filepath.Join(p.Path, "providers")
	destDir := filepath.Join("providers", p.Name)

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create provider directory: %v", err)
	}

	opts["ProviderRoot"] = p.Path

	// Copy all files from source to destination
	err = filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
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

		err = utils.TemplateAFile(path, destPath, opts)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy provider files: %v", err)
	}

	if p.plugin != nil {
		err = p.plugin.OnConfigure()
		if err != nil {
			return fmt.Errorf("failed to configure provider: %v", err)
		}
	}

	return nil
}

func (p *Provider) IsConfigured() bool {
	_, err := os.Stat(fmt.Sprintf("providers/%s", p.Name))
	return err == nil
}

func (p *Provider) DeleteConfiguration() error {
	err := os.RemoveAll(fmt.Sprintf("providers/%s", p.Name))
	if err != nil {
		return fmt.Errorf("failed to delete provider configuration: %v", err)
	}

	return nil
}
