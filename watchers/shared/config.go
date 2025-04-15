package watchers_shared

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/utils"
)

var lastConfigHash string
var configLock sync.Mutex

func WatchConfig() error {
	var err error

	lastConfigHash, err = utils.CalculateFileHash(".appy/appy.yaml")
	if err != nil {
		return err
	}

	watcher, err := utils.NewWatcher(".appy/appy.yaml", func(event fsnotify.Event) {
		onConfigChange()
	})
	if err != nil {
		return err
	}

	watcher.Start()

	return nil
}

func onConfigChange() {
	configLock.Lock()
	defer configLock.Unlock()

	// Check if hash changed providers changed
	currentConfigHash, err := utils.CalculateFileHash(".appy/appy.yaml")
	if err != nil {
		fmt.Println("Error calculating config hash:", err)
		return
	}

	if currentConfigHash != lastConfigHash {
		fmt.Println("Config changed")

		cfg, err := config.LoadConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
			return
		}

		err = cfg.Reconfigure()
		if err != nil {
			fmt.Println("Error reconfiguring:", err)
			panic(err)
		}

		err = cfg.Save()
		if err != nil {
			fmt.Println("Error saving config:", err)
			panic(err)
		}

		// Avoid double update
		currentConfigHash, err = utils.CalculateFileHash(".appy/appy.yaml")
		if err != nil {
			fmt.Println("Error calculating config hash:", err)
			return
		}

		lastConfigHash = currentConfigHash
	}
}
