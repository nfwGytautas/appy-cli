package project_shared

import (
	"context"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/utils"
)

var lastConfigHash string
var configLock sync.Mutex

func WatchConfig(ctx context.Context, wg *sync.WaitGroup) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add("appy.yaml")
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer watcher.Close()
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					onConfigChange()
				}
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

func onConfigChange() {
	configLock.Lock()
	defer configLock.Unlock()

	// Check if hash changed providers changed
	currentConfigHash, err := utils.CalculateFileHash("appy.yaml")
	if err != nil {
		utils.Console.ErrorLn("Error calculating config hash: %v", err)
		return
	}

	if currentConfigHash != lastConfigHash {
		if !utils.Verbose {
			s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
			s.Suffix = " Reconfiguring...\n"
			s.FinalMSG = "Done!"
			s.Start()
			defer s.Stop()
		}

		utils.Console.DebugLn("Config changed")
		cfg := config.GetConfig()

		err = cfg.Reconfigure()
		if err != nil {
			utils.Console.Fatal(err)
		}

		// Avoid double update
		currentConfigHash, err = utils.CalculateFileHash("appy.yaml")
		if err != nil {
			utils.Console.Fatal(err)
		}

		lastConfigHash = currentConfigHash

		return
	}
}
