package logic

import (
	"github.com/nfwGytautas/appy-cli/model"
	"github.com/nfwGytautas/appy-cli/templates"
)

func Scaffold() {
	cfg, err := model.ReadConfig()
	errorCheckAndPanic(err)

	// Create appy folder if not exists
	err = ensureDirectory(".appy")
	errorCheckAndPanic(err)

	// Variables
	variablesContent, err := templates.GenTemplate(templates.VariablesAutogen, templates.TemplateParams{
		"Variables": cfg.Variables,
	})
	errorCheckAndPanic(err)

	// Create middleware
	middlewareContent, err := templates.GenTemplate(templates.Middleware, templates.TemplateParams{
		"Middlewares": cfg.Middlewares,
	})
	errorCheckAndPanic(err)

	// Create endpoint groups
	groupsContent, err := templates.GenTemplate(templates.EndpointGroup, templates.TemplateParams{
		"Groups": cfg.EndpointGroups,
	})
	errorCheckAndPanic(err)

	// Endpoints
	endpointsContent, err := writeEndpoints(cfg.Endpoints)
	errorCheckAndPanic(err)

	// Serve
	serveContent, err := templates.GenTemplate(templates.Serve, templates.TemplateParams{
		"ServePoints": cfg.ServePoints,
	})
	errorCheckAndPanic(err)

	// Final http autogen
	httpAutogenContent, err := templates.GenTemplate(templates.HttpAutogen, templates.TemplateParams{
		"Middleware": middlewareContent,
		"Groups":     groupsContent,
		"Endpoints":  endpointsContent,
		"Serve":      serveContent,
	})
	errorCheckAndPanic(err)
	httpAutogenContent = beutifyContent(httpAutogenContent)

	// Create files
	err = createFile(".appy/variables_autogen.go", variablesContent)
	errorCheckAndPanic(err)

	err = createFile(".appy/http_autogen.go", cfg.ReplaceWithVariables(httpAutogenContent))
	errorCheckAndPanic(err)

	// Format files
	err = goPipeline(".appy/http_autogen.go")
	errorCheckAndPanic(err)

	err = goPipeline(".appy/variables_autogen.go")
	errorCheckAndPanic(err)
}

func writeEndpoints(endpoints []model.Endpoint) (string, error) {
	var content string
	for _, e := range endpoints {
		childrenContent, err := writeEndpoints(e.Children)
		if err != nil {
			return "", err
		}

		endpointContent, err := templates.GenTemplate(templates.Endpoint, templates.TemplateParams{
			"Endpoint":        e,
			"ChildrenContent": childrenContent,
		})
		if err != nil {
			return "", err
		}

		content += endpointContent
	}

	return content, nil
}
