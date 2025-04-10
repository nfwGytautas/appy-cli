package watchers_shared

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func WatchDomain(root string) (*utils.Watcher, error) {
	domain := filepath.Base(root)

	cfg, err := shared.LoadConfig()
	if err != nil {
		return nil, err
	}

	domainWatcher, err := utils.NewWatcher(root+"/usecase/", func(event fsnotify.Event) {
		onDomainUsecaseEvent(cfg, root, domain, event)
	})
	if err != nil {
		return nil, err
	}

	domainWatcher.Start()

	return domainWatcher, nil
}

func onDomainUsecaseEvent(cfg *shared.Config, root string, domain string, event fsnotify.Event) {
	// Check if this is a file or directory
	fileInfo, err := os.Stat(event.Name)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		return
	}

	if filepath.Ext(event.Name) != ".go" {
		return
	}

	usecase := filepath.Base(event.Name)
	usecase = strings.TrimSuffix(usecase, filepath.Ext(usecase))
	usecase = strings.ToLower(usecase)

	usecaseName := cases.Title(language.English).String(usecase)

	domainRoot := cfg.Module + "/" + cfg.GetDomainsRoot(domain)

	if event.Op&fsnotify.Create == fsnotify.Create {
		fmt.Println("New usecase:", usecase)

		fmt.Println("    + adding template")
		// Fill template
		f, err := os.OpenFile(event.Name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			utils.ConsoleError("Failed to open usecase file: %s (%v)", usecase, err)
			return
		}

		tmpl := template.Must(template.New(usecase).Parse(templates.DomainExampleUsecase))

		err = tmpl.Execute(f, map[string]string{
			"DomainName":  domain,
			"DomainRoot":  domainRoot,
			"UsecaseName": usecaseName,
		})
		if err != nil {
			utils.ConsoleError("Failed to execute usecase template: %s (%v)", usecase, err)
			return
		}

		err = f.Close()
		if err != nil {
			utils.ConsoleError("Failed to close usecase file: %s (%v)", usecase, err)
			return
		}

		fmt.Println("    + adding associated input port")

		// Create input port
		f, err = os.Create(root + "/ports/in/" + usecase + ".go")
		if err != nil {
			utils.ConsoleError("Failed to create input port file: %s (%v)", usecase, err)
			return
		}

		// Write template
		tmpl = template.Must(template.New(usecase).Parse(templates.DomainExampleInPort))

		err = tmpl.Execute(f, map[string]string{
			"DomainName":  domain,
			"DomainRoot":  domainRoot,
			"UsecaseName": usecaseName,
		})
		if err != nil {
			utils.ConsoleError("Failed to execute input port template: %s (%v)", usecase, err)
			return
		}

		err = f.Close()
		if err != nil {
			utils.ConsoleError("Failed to close input port file: %s (%v)", usecase, err)
			return
		}
	}
}
