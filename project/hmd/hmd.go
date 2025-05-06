package project_hmd

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
)

type watchers struct {
	domainWatchers  *fsnotify.Watcher
	domainsWatcher  *fsnotify.Watcher
	adapterWatchers *fsnotify.Watcher
}

func Watch(ctx context.Context, wg *sync.WaitGroup) error {
	var err error

	watchers := &watchers{}
	watchers.adapterWatchers, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	watchers.domainsWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	watchers.domainWatchers, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watchers.domainsWatcher.Add("domains/")
	if err != nil {
		return err
	}

	// Get all domains in domains/
	domains, err := os.ReadDir("domains/")
	if err != nil {
		return err
	}

	for _, domain := range domains {
		if domain.IsDir() {
			watchers.domainWatchers.Add("domains/" + domain.Name())
			watchers.adapterWatchers.Add("domains/" + domain.Name() + "/adapters/")
		}
	}

	wg.Add(1)
	go func() {
		defer watchers.domainsWatcher.Close()
		defer watchers.domainWatchers.Close()
		defer watchers.adapterWatchers.Close()
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watchers.domainsWatcher.Events:
				if !ok {
					return
				}
				onDomainEvent(event, watchers)
			case event, ok := <-watchers.domainWatchers.Events:
				if !ok {
					return
				}
				onDomainUsecaseEvent(event)
			case event, ok := <-watchers.adapterWatchers.Events:
				if !ok {
					return
				}
				onAdapterAddedEvent(event)
			case err, ok := <-watchers.domainWatchers.Errors:
				if !ok {
					return
				}
				utils.Console.ErrorLn("Error: %v", err)
			case err, ok := <-watchers.adapterWatchers.Errors:
				if !ok {
					return
				}
				utils.Console.ErrorLn("Error: %v", err)
			case err, ok := <-watchers.domainsWatcher.Errors:
				if !ok {
					return
				}
				utils.Console.ErrorLn("Error: %v", err)
			}
		}
	}()

	return nil
}

func onDomainEvent(event fsnotify.Event, watchers *watchers) {
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
		watchers.domainWatchers.Remove("domains/" + domain)
		watchers.adapterWatchers.Remove("domains/" + domain + "/adapters/")
		return
	}

	if event.Op&fsnotify.Create == fsnotify.Create {
		// Create domain template
		err = scaffoldDomain(domain)
		if err != nil {
			utils.Console.ErrorLn("Failed to scaffold domain: %s (%v)", domain, err)
			return
		}

		// Add watcher
		watchers.domainWatchers.Add("domains/" + domain)
		watchers.adapterWatchers.Add("domains/" + domain + "/adapters/")

		for _, p := range config.GetConfig().Plugins.GetLoadedPlugins() {
			err = p.OnDomainCreated(domain)
			if err != nil {
				utils.Console.Fatal(err)
			}
		}
	}
}

func scaffoldDomain(name string) error {
	tree := utils.GeneratedFileTree{}

	tree.SetPrefix("domains/" + name)
	tree.AddDirectory("adapters/")
	tree.AddDirectory("model/")
	tree.AddFile("domain.go", templates.DomainExampleDomain, []string{shared.ToolGoFmt})
	tree.AddFile("ping.go", templates.DomainExampleUsecase, []string{shared.ToolGoFmt})
	tree.AddFile("model/example.go", templates.DomainExampleModel, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config":      config.GetConfig(),
		"DomainName":  name,
		"UsecaseName": "ping",
	})
	if err != nil {
		return err
	}

	return nil
}
