package watchers_shared

import (
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/utils"
)

var lastConfigHash string
var configLock sync.Mutex

func WatchConfig() error {
	var err error

	lastConfigHash, err = utils.CalculateFileHash("appy.yaml")
	if err != nil {
		return err
	}

	watcher, err := utils.NewWatcher("appy.yaml", func(event fsnotify.Event) {
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

		cfg, err := config.LoadConfig()
		if err != nil {
			utils.Console.Fatal(err)
		}

		err = cfg.Reconfigure()
		if err != nil {
			utils.Console.Fatal(err)
		}

		err = cfg.Save()
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
