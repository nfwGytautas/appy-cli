package logic

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/nfwGytautas/appy-cli/templates"
)

func Scaffold() {
	var err error

	type endpointOpts struct {
		Method   string
		Path     string
		FullType string
	}

	imports := []string{}
	endpoints := []endpointOpts{}

	// Iterate over 'endpoints/' and check what we have
	err = walkFilesInDirectory("endpoints", func(path string, filename string) {
		// Read the first line to get package name
		lines, err := readFileLines(path)
		errorCheckAndPanic(err)

		if !strings.Contains(lines[0], "package") {
			panic("Package not found in file: " + path)
		}

		// Get package name
		packageName := strings.Split(lines[0], " ")[1]

		imports = append(imports, fmt.Sprintf("%s \"%s/%s\"", packageName, getGoModule(), filepath.Dir(path)))
		endpoints = append(endpoints, endpointOpts{
			Method:   filename,
			Path:     filepath.Dir(strings.Split(path, "endpoints/")[1]),
			FullType: fmt.Sprintf("%s.%sEndpoint", packageName, capitalizeFirstLetter(strings.ToLower(filename))),
		})
	})
	errorCheckAndPanic(err)

	// Unique imports
	imports = uniqueStrings(imports)

	endpointsContent, err := templates.GenTemplate(templates.EndpointsScaffold, templates.TemplateParams{
		"Imports":   imports,
		"Endpoints": endpoints,
	})
	errorCheckAndPanic(err)

	// Generate the scaffold
	chainErrorChecks(
		createFile(".appy/endpoints.go", endpointsContent),
	)
}
