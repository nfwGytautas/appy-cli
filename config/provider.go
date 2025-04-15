package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
)

type Provider struct {
	Name           string    `yaml:"name"`
	Path           string    `yaml:"path"`
	Version        string    `yaml:"version"`
	Description    string    `yaml:"description"`
	ProviderAtRoot bool      `yaml:"atRootLevel"`
	Enabled        bool      `yaml:"enabled"`
	Hooks          []Hook    `yaml:"hooks"`
	Watchers       []Watcher `yaml:"watchers"`

	config *ProvidersConfig `yaml:"-"`
}

func (p *Provider) Configure(opts ProviderConfigureOpts) error {
	for _, hook := range p.Hooks {
		hook.provider = p
		err := hook.Configure(opts)
		if err != nil {
			return fmt.Errorf("failed to configure hook `%s`: %v", hook.Name, err)
		}
	}

	for _, watcher := range p.Watchers {
		watcher.provider = p
		err := watcher.Configure(opts)
		if err != nil {
			return fmt.Errorf("failed to configure watcher: %v", err)
		}
	}

	if !p.Enabled {
		return nil
	}

	// Check if provider is already configured
	if p.IsConfigured() {
		return nil
	}

	fmt.Printf("Configuring provider `%s@%s`...\n", p.Name, p.Version)

	// Copy provider files from repository to providers directory
	sourceDir := filepath.Join(p.Path, "providers")
	destDir := filepath.Join("providers", p.Name)

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create provider directory: %v", err)
	}

	opts["ProviderRoot"] = p.Path

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

		err = utils.TemplateAFile(path, destPath, opts)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy provider files: %v", err)
	}

	err = p.RunLocalHook(shared.HookOnProviderConfigured, opts)
	if err != nil {
		return fmt.Errorf("failed to run provider configured hook: %v", err)
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

func (p *Provider) RunLocalHook(hookName string, data any) error {
	for _, hook := range p.Hooks {
		hook.provider = p
		if hook.Name == hookName {
			return hook.Run(data)
		}
	}

	return nil
}

func (p *Provider) RestartWatchers() error {
	for _, watcher := range p.Watchers {
		watcher.provider = p
		err := watcher.Restart()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) StopWatchers() error {
	if !p.Enabled {
		return nil
	}

	for _, watcher := range p.Watchers {
		watcher.provider = p
		err := watcher.Stop()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) StartWatchers() error {
	if !p.Enabled {
		return nil
	}

	for _, watcher := range p.Watchers {
		watcher.provider = p
		err := watcher.Start()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) ApplyStringSubstitution(str string) string {
	str = strings.ReplaceAll(str, "${ProviderRoot}", "${Workspace}/"+p.Path)
	str = p.config.ApplyStringSubstitution(str)
	return str
}
