package templates

const VariablesAutogen = `
package appy_generated

import (
	"github.com/nfwGytautas/appy-go"
)

{{ range .Variables }}
var appyVar_{{ .Name }} {{ .Type }}  = {{ .Value }}
{{ end }}

`
