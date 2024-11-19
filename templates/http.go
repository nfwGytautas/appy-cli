package templates

const HttpAutogen = `
package appy_generated

import (
	"github.com/nfwGytautas/appy-go"
)

func AppySetupHttp() {
	// Setup root group
	root := appy.HTTP().Root()

	// Endpoint groups
	{{ .Groups }}

	// Endpoints
}

`

const EndpointGroup = `
{{ range .Groups }}
{{ .Name }} = {{ if .Parent }} {{ .Parent }}.Group("{{ .Path }}") {{ else }} root.Group("{{ .Path }}") {{ end }}
{{ if .Middleware }} {{ .Name }}.Use({{ range .Middleware }}{{ .Name }},{{ end }}) {{ end }}
{{ end }}
`
