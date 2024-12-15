package templates

const EndpointTemplate = `
package {{ .Package }}

import (
	appy_http "github.com/nfwGytautas/appy-go/http"
)

type {{ .Name }}Endpoint struct {
	appy_http.Endpoint

	// TODO: Add your endpoint data here
}

func (e *{{ .Name }}Endpoint) Options() appy_http.EndpointOptions {
	// TODO: Specify the options for the endpoint
	return appy_http.EndpointOptions{}
}

func (e *{{ .Name }}Endpoint) Parse(pr *appy_http.ParamReader) {
	// TODO: Add your request parsing logic here
}

func (e *{{ .Name }}Endpoint) Validate() *appy_http.ValidationResult {
	// TODO: Check if the request is valid
	return appy_http.Valid()
}

func (e *{{ .Name }}Endpoint) Handle() appy_http.EndpointResult {
	// TODO: Add logic for handling the request
	return appy_http.Ok()
}
`
