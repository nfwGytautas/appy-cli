package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nfwGytautas/appy-cli/utils"
)

type Hook struct {
	Name    string   `yaml:"name"`
	Actions []string `yaml:"actions"`

	provider *Provider `yaml:"-"`
}

func (h *Hook) Configure(opts RepositoryConfigureOpts) error {
	for it, action := range h.Actions {
		h.Actions[it] = h.provider.ApplyStringSubstitution(action)
	}

	return nil
}

func (h *Hook) Run(data any) error {
	// TODO: Validate hook
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, action := range h.Actions {
		if action == "appy copyTemplate" {
			return h.copyTemplate(data)
		}

		err := utils.RunCommand(cwd, action)
		if err != nil {
			return fmt.Errorf("failed to run hook action `%s`: %v", action, err)
		}
	}

	return nil
}

func (h *Hook) copyTemplate(data any) error {
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
	sourceDir := filepath.Join(h.provider.Path, "domain")
	destDir := filepath.Join("domains", domainName, "adapters", h.provider.Name)

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

		err = utils.TemplateAFile(path, destPath, data)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to copy provider domain files: %v", err)
	}

	return nil
}
