package watchers_hmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/scaffolds"
	"github.com/nfwGytautas/appy-cli/utils"
	watchers_shared "github.com/nfwGytautas/appy-cli/watchers/shared"
)

var domainWatchers map[string]*utils.Watcher = make(map[string]*utils.Watcher)

func Watch() error {
	domainsWatcher, err := utils.NewWatcher("domains/", onDomainEvent)
	if err != nil {
		return err
	}

	domainsWatcher.Start()

	// Get all domains in domains/
	domains, err := os.ReadDir("domains/")
	if err != nil {
		return err
	}

	for _, domain := range domains {
		if domain.IsDir() {
			domainWatcher, err := watchers_shared.WatchDomain("domains/" + domain.Name())
			if err != nil {
				return err
			}

			domainWatchers[domain.Name()] = domainWatcher
		}
	}

	return nil
}

func onDomainEvent(event fsnotify.Event) {
	// Check if this is a file or directory
	fileInfo, err := os.Stat(event.Name)
	if err != nil {
		return
	}

	if !fileInfo.IsDir() {
		utils.ConsoleWarn("./domains/ should consist only of packages, got file")
		return
	}

	fmt.Println("Domain:", event.Name, event.Op)

	domain := filepath.Base(event.Name)

	if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
		// Disable watcher
		domainWatchers[domain].Stop()
		domainWatchers[domain] = nil
		return
	}

	if event.Op&fsnotify.Create == fsnotify.Create {
		cfg, err := config.LoadConfig()
		if err != nil {
			utils.ConsoleError("Failed to load config: %s (%v)", domain, err)
			return
		}

		// Create domain template
		err = scaffolds.ScaffoldDomain(cfg, domain)
		if err != nil {
			utils.ConsoleError("Failed to scaffold domain: %s (%v)", domain, err)
			return
		}

		// Add watcher
		domainWatcher, err := watchers_shared.WatchDomain("domains/" + domain)
		if err != nil {
			utils.ConsoleError("Failed to watch domain: %s (%v)", domain, err)
			return
		}

		domainWatchers[domain] = domainWatcher
	}
}
