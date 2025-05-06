package project_hmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func onDomainUsecaseEvent(event fsnotify.Event) {
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

	domain := filepath.Base(filepath.Dir(event.Name))

	usecase := filepath.Base(event.Name)
	usecase = strings.TrimSuffix(usecase, filepath.Ext(usecase))
	usecase = strings.ToLower(usecase)

	usecaseName := strings.ReplaceAll(usecase, "_", " ")
	usecaseName = cases.Title(language.English).String(usecaseName)
	usecaseName = strings.ReplaceAll(usecaseName, " ", "")

	if event.Op&fsnotify.Create == fsnotify.Create {
		utils.Console.DebugLn("New usecase: %s", usecase)

		utils.Console.DebugLn("    + adding template")
		// Fill template
		f, err := os.OpenFile(event.Name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			utils.Console.ErrorLn("Failed to open usecase file: %s (%v)", usecase, err)
			return
		}

		tmpl := utils.NewTemplate(templates.DomainExampleUsecase)

		err = tmpl.Execute(f, map[string]string{
			"DomainName":  domain,
			"UsecaseName": usecaseName,
		})
		if err != nil {
			utils.Console.ErrorLn("Failed to execute usecase template: %s (%v)", usecase, err)
			return
		}

		err = f.Close()
		if err != nil {
			utils.Console.ErrorLn("Failed to close usecase file: %s (%v)", usecase, err)
			return
		}
	}
}

func onAdapterAddedEvent(event fsnotify.Event) {
	// always <domain>/adapters/<adapter>
	domain := filepath.Base(filepath.Dir(filepath.Dir(event.Name)))
	adapter := filepath.Base(event.Name)
	log.Println(event)
	if event.Op&fsnotify.Create == fsnotify.Create {
		for _, p := range config.GetConfig().Plugins.GetLoadedPlugins() {
			err := p.OnAdapterCreated(domain, adapter)
			if err != nil {
				utils.Console.Fatal(err)
			}
		}
	}
}
