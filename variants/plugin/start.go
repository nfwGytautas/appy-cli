package variant_plugin

import (
	"context"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/plugins"
	"github.com/nfwGytautas/appy-cli/utils"
)

func (cfg *Config) Start(ctx context.Context) error {
	utils.Console.DebugLn("Starting Plugin...")

	cfg.engine = plugins.NewPluginEngine(map[string]any{
		"module":  "AppyPluginTestModule",
		"project": "AppyPluginTestProject",
	})

	utils.Console.DebugLn("Starting watchers...")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	watcher.Add("plugin.lua")

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

	return nil
}

func (cfg *Config) handleWatcherEvent(event fsnotify.Event) {
	if event.Op == fsnotify.Chmod {
		// Ignore chmod events
		return
	}

	parts := strings.Split(event.Name, string(os.PathSeparator))

	finalPart := parts[len(parts)-1]

	if finalPart == "plugin.lua" {
		utils.Console.DebugLn("Plugin file changed: %s", event.Name)

		plugin, err := cfg.engine.LoadPlugin(event.Name)
		if err != nil {
			utils.Console.ErrorLn("Error loading plugin: %s", event.Name)
			return
		}

		utils.Console.Info("Plugin loaded: %v", plugin)

		return
	}
}
