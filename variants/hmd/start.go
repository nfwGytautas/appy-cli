package variant_hmd

import (
	"context"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
	variant_base "github.com/nfwGytautas/appy-cli/variants/base"
)

var ignoredDirs = []string{
	".git",
	".appy",
}

var lastConfigHash string

func (cfg *Config) Start(ctx context.Context) error {
	utils.Console.DebugLn("Starting HMD...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	filepath.Walk(".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			utils.Console.Error("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() {
			return nil
		}

		// Ignore certain directories
		for _, ignoredDir := range ignoredDirs {
			if info.Name() == ignoredDir {
				utils.Console.DebugLn("Ignoring directory: %s", path)
				return filepath.SkipDir
			}
		}

		utils.Console.DebugLn("Watching directory: %s", path)
		watcher.Add(path)

		return nil
	})

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
				cfg.handleWatcherEvent(event)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.Console.ErrorLn("Error: %v", err)
			}
		}
	}()

	// Block until a signal is received
	<-stop
	utils.Console.ClearLines(1)
	utils.Console.DebugLn("Received signal, shutting down...")

	return nil
}

func (cfg *Config) handleWatcherEvent(event fsnotify.Event) {
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

			if parts[2] == "adapters" {
				// Adapter added or removed
				if event.Op&fsnotify.Create == fsnotify.Create {
					// New adapter added
					utils.Console.DebugLn("New adapter added: %s", finalPart)

					for _, p := range cfg.Plugins.GetLoadedPlugins() {
						err := p.OnAdapterCreated(domain, finalPart)
						if err != nil {
							utils.Console.ErrorLn("Error on plugin hook creating adapter: %v", err)
							return
						}
					}
				}

				return
			}

			// Usecase added or removed
			if event.Op&fsnotify.Create == fsnotify.Create {
				// New usecase added
				utils.Console.DebugLn("New usecase added: %s", finalPart)

				err := utils.TemplateAStringToFile(event.Name, templates.DomainExampleUsecase, map[string]any{
					"DomainName":  domain,
					"UsecaseName": strings.TrimSuffix(finalPart, filepath.Ext(finalPart)),
				})
				if err != nil {
					utils.Console.ErrorLn("Error creating usecase file: %v", err)
					return
				}

				return
			}

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
	tree.AddDirectory("adapters/")
	tree.AddDirectory("model/")
	tree.AddFile("domain.go", templates.DomainExampleDomain, []string{shared.ToolGoFmt})
	tree.AddFile("ping.go", templates.DomainExampleUsecase, []string{shared.ToolGoFmt})
	tree.AddFile("model/example.go", templates.DomainExampleModel, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config":      cfg,
		"DomainName":  name,
		"UsecaseName": "ping",
	})
	if err != nil {
		return err
	}

	return nil
}
