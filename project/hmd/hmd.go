package project_hmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
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
			domainWatcher, err := watchDomain("domains/" + domain.Name())
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

	log.Println("Event:", event.Name, event.Op)

	if !fileInfo.IsDir() && event.Name != "domains.go" {
		utils.Console.WarnLn("./domains/ should consist only of packages, got file")
		return
	}

	utils.Console.DebugLn("Domain: %s, %s", event.Name, event.Op)

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
			utils.Console.ErrorLn("Failed to load config: %s (%v)", domain, err)
			return
		}

		// Create domain template
		err = scaffoldDomain(cfg, domain)
		if err != nil {
			utils.Console.ErrorLn("Failed to scaffold domain: %s (%v)", domain, err)
			return
		}

		// Create test folder
		err = os.MkdirAll("tests/"+domain, 0755)
		if err != nil {
			utils.Console.ErrorLn("Failed to create test folder: %s (%v)", domain, err)
			return
		}

		// Add watcher
		domainWatcher, err := watchDomain("domains/" + domain)
		if err != nil {
			utils.Console.ErrorLn("Failed to watch domain: %s (%v)", domain, err)
			return
		}

		domainWatchers[domain] = domainWatcher
	}
}

func scaffoldDomain(cfg *config.AppyConfig, name string) error {
	tree := utils.GeneratedFileTree{}

	tree.SetPrefix("domains/" + name)
	tree.AddDirectory("adapters/")
	tree.AddDirectory("connectors/")
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

	err = cfg.RunHook(shared.HookOnDomainCreated, map[string]any{
		"DomainName": name,
		"Module":     cfg.Module,
	})
	if err != nil {
		return err
	}

	return nil
}
