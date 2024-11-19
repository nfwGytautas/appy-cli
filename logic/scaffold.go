package logic

import (
	"github.com/nfwGytautas/appy-cli/model"
	"github.com/nfwGytautas/appy-cli/templates"
)

func Scaffold() {
	cfg, err := model.ReadConfig()
	if err != nil {
		panic(err)
	}

	// Create appy folder if not exists
	err = ensureDirectory(".appy")
	if err != nil {
		panic(err)
	}

	// Create endpoint groups
	groupsContent, err := templates.GenTemplate(templates.EndpointGroup, templates.TemplateParams{
		"Groups": cfg.EndpointGroups,
	})
	if err != nil {
		panic(err)
	}

	httpAutogenContent, err := templates.GenTemplate(templates.HttpAutogen, templates.TemplateParams{
		"Groups": groupsContent,
	})
	if err != nil {
		panic(err)
	}

	err = createFile(".appy/http_autogen.go", httpAutogenContent)
	if err != nil {
		panic(err)
	}

	err = goFmtFile(".appy/http_autogen.go")
	if err != nil {
		panic(err)
	}

}
