package templates

const EndpointsScaffold = `
package appy_autogen

import (
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_http "github.com/nfwGytautas/appy-go/http"
	appy_runtime "github.com/nfwGytautas/appy-go/runtime"
	{{ range .Imports }}
	{{ . }}
	{{ end }}
)

func registerEndpoints(ctx *appy_runtime.AppyContext) {
	{{ if .Endpoints }}
	appy_logger.Logger().Debug("Registering endpoints")
	{{ range .Endpoints }}
	ctx.HttpEngine.RegisterEndpoint("{{ .Method }}", "{{ .Path}}", func() appy_http.EndpointIf { return &{{ .FullType }}{} })
	{{ end }}
	{{ end }}
}

`
