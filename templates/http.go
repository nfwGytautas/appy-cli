package templates

const HttpAutogen = `
package appy_generated

import (
	"github.com/nfwGytautas/appy-go"
)

func AppySetupHttp() {
	// Setup root group
	root := appy.HTTP().Root()

	//
	// Middleware
	//
	{{ .Middleware }}

	//
	// Endpoint groups
	//
	{{ .Groups }}

	//
	// Endpoints
	//
	{{ .Endpoints }}
}

`

const EndpointGroup = `
{{ range .Groups }}
{{ .Name }} = {{ if .Parent }} {{ .Parent }}.Group("{{ .Path }}") {{ else }} root.Group("{{ .Path }}") {{ end }}
{{ if .Middleware }} {{ .Name }}.Use({{ fnCommaSeperated .Middleware }}) {{ end }}
{{ end }}
`

const Middleware = `
{{ range .Middlewares }}
{{ fnSanitizeVar .Name }} := {{ .Provider }}({{ fnUpackMapValues .Params }}).Provide()
{{ end }}
`

const Endpoint = `
// {{ .Endpoint.Name }}
{{ if .Endpoint.Impl }}
{{ .Endpoint.Group }}.{{ .Endpoint.Method }}("{{ .Endpoint.Name }}", appy.AppyHttpBootstrap({{ .Endpoint.Impl }}))
{{ end }}
{{ if .ChildrenContent }}
{
{{ .ChildrenContent }}
}
{{ end }}
`

const Serve = `
{{ range .ServePoints }}
{{ .Group }}.StaticFS("{{ .On }}", gin.Dir({{ .Dir }}, false))
{{ end }}
`
