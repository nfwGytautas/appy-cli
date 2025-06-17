package variant_hmd

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
	variant_base "github.com/nfwGytautas/appy-cli/variants/base"
)

var ignoredDirs = []string{
	".git",
	".appy",
	".github",
	".vscode",
}

var lastConfigHash string

func (cfg *Config) Start(ctx context.Context) error {
	utils.Console.DebugLn("Starting HMD...")

	utils.Console.DebugLn("Configuring...")
	err := cfg.Reconfigure()
	if err != nil {
		return err
	}

	lastConfigHash, err = utils.CalculateFileHash(variant_base.YamlFilePath)
	if err != nil {
		return err
	}

	utils.Console.DebugLn("Starting watchers...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	utils.Console.DebugLn("Watching config")
	watcher.Add("appy.yaml")

	utils.Console.DebugLn("Watching providers")
	watcher.Add("providers/")
	filepath.Walk("providers", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			utils.Console.Error("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			return nil
		}

		utils.Console.DebugLn("Watching provider directory: %s", path)
		watcher.Add(path)

		return nil
	})

	utils.Console.DebugLn("Watching domains")
	watcher.Add("domains/")
	filepath.Walk("domains/", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			utils.Console.Error("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			return nil
		}

		utils.Console.DebugLn("Watching domain: %s", path)
		watcher.Add(path)
		return nil
	})

	utils.Console.DebugLn("Watching connectors")
	watcher.Add("connectors/")

	utils.Console.DebugLn("Watching interfaces")
	watcher.Add("interfaces/domains/")
	watcher.Add("interfaces/models/")

	go func() {
		defer watcher.Close()

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				cfg.handleWatcherEvent(watcher, event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.Console.ErrorLn("Error: %v", err)
			}
		}
	}()

	return nil
}

func (cfg *Config) handleWatcherEvent(watcher *fsnotify.Watcher, event fsnotify.Event) {
	if event.Op == fsnotify.Chmod {
		// Ignore chmod events
		return
	}
	parts := strings.Split(event.Name, string(os.PathSeparator))

	utils.Console.DebugLn("event: %v, parts: %v", event, parts)

	category := parts[0]
	finalPart := parts[len(parts)-1]

	if category == "domains" {
		// If part count > 2, its a domain event
		if len(parts) > 2 {
			domain := parts[1]

			// Usecase added or removed
			if event.Op&fsnotify.Create == fsnotify.Create {
				// New usecase added
				utils.Console.DebugLn("New usecase added: %s", finalPart)

				err := utils.TemplateAStringToFile(event.Name, templateDomainExampleUsecase, map[string]any{
					"DomainName":  domain,
					"UsecaseName": strings.TrimSuffix(finalPart, filepath.Ext(finalPart)),
				})
				if err != nil {
					utils.Console.ErrorLn("Error creating usecase file: %v", err)
					return
				}

				return
			}

			// Uninteresting

			return
		}

		// A new domain added?
		if event.Op&fsnotify.Create == fsnotify.Create {
			// New domain added
			utils.Console.DebugLn("New domain added: %s", finalPart)

			err := cfg.scaffoldDomain(finalPart)
			if err != nil {
				utils.Console.ErrorLn("Error scaffolding domain: %v", err)
				return
			}

			watcher.Add(filepath.Join("domains", finalPart))

			return
		}

		if event.Op&fsnotify.Remove == fsnotify.Remove {
			// Domain removed
			utils.Console.DebugLn("Domain removed: %s", finalPart)

			watcher.Remove(filepath.Join("domains", finalPart))

			return
		}

		// Uninteresting

		return
	}

	if category == "appy.yaml" {
		newConfigHash, err := utils.CalculateFileHash(variant_base.YamlFilePath)
		if err != nil {
			utils.Console.ErrorLn("Error calculating config hash: %v", err)
			return
		}

		if newConfigHash == lastConfigHash {
			// Config file changed but hash is the same, ignore
			utils.Console.DebugLn("Config file changed but hash is the same, ignoring...")
			return
		}

		// Config file changed
		utils.Console.InfoLn("Configuration change detected, reconfiguring...")
		err = cfg.LoadAndReconfigure()
		if err != nil {
			utils.Console.ErrorLn("Error reconfiguring: %v", err)
			return
		}

		lastConfigHash, err = utils.CalculateFileHash(variant_base.YamlFilePath)
		if err != nil {
			utils.Console.ErrorLn("Error calculating config hash: %v", err)
			return
		}

		return
	}
}

func (cfg *Config) scaffoldDomain(name string) error {
	tree := utils.GeneratedFileTree{}

	tree.SetPrefix("domains/" + name)
	tree.AddFile("domain.go", templateDomainExampleDomain, []string{shared.ToolGoFmt})
	tree.AddFile("ping.go", templateDomainExampleUsecase, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config":      cfg,
		"DomainName":  name,
		"ModelName":   name,
		"UsecaseName": "ping",
	})
	if err != nil {
		return err
	}

	return nil
}
