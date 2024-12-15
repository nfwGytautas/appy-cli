package logic

import (
	"strings"

	"github.com/nfwGytautas/appy-cli/templates"
)

func Create(args []string) {
	createType := args[0]

	switch createType {
	case "endpoint":
		createEndpoint()
	default:
		panic("Unknown type: " + createType)
	}
}

func createEndpoint() {
	opts := struct {
		path   string
		method string
	}{}

	chainErrorChecks(
		getUserInput("Path: ", &opts.path),
		getUserInput("Method [GET|POST|PUT|DELETE]: ", &opts.method),
	)

	opts.method = strings.ToUpper(opts.method)

	if !stringValueIsOneOf(opts.method, []string{"GET", "POST", "PUT", "DELETE"}) {
		panic("Invalid method: " + opts.method)
	}

	opts.path = strings.ReplaceAll(opts.path, ":", "/")
	path := strings.Split(opts.path, "/")
	packageName := strings.Join(path, "_")
	file := "endpoints/" + opts.path + "/" + opts.method + ".go"

	endpointContent, err := templates.GenTemplate(templates.EndpointTemplate, templates.TemplateParams{
		"Package": packageName,
		"Name":    capitalizeFirstLetter(strings.ToLower(opts.method)),
		"Method":  opts.method,
	})

	chainErrorChecks(
		err,
		ensureDirectory("endpoints/"+opts.path),
		createFile(file, endpointContent),
		goPipeline(file),
	)
}
