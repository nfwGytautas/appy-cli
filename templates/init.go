package templates

const Gitignore = `
.appy/
`

const GoMod = `
module github.com/username/your-app-name

go 1.20

replace github.com/nfwGytautas/appy-go => ../appy-go
`

const ConfigGo = `
package main

import (
	"github.com/nfwGytautas/appy-go/config"
)

func getConfig() appy_config.AppyConfig {
	// TODO: Specify your appy configuration here
	return appy_config.AppyConfig{}
}

`

const MainGo = `
package main

import (
	appy_runtime "github.com/nfwGytautas/appy-go/runtime"
	appy_autogen "github.com/username/your-app-name/.appy"
)

func main() {
	ctx := appy_runtime.Initialize(getConfig(), appy_autogen.Hook)

	//
	// Add your custom initialization logic here
	//

	ctx.Takeover()
}

`

const AutogenHookGo = `
package appy_autogen

import (
	appy_logger "github.com/nfwGytautas/appy-go/logger"
	appy_runtime "github.com/nfwGytautas/appy-go/runtime"
)

func Hook(ctx *appy_runtime.AppyContext) {
	appy_logger.Logger().Debug("Hooking into appy_autogen")

	registerEndpoints(ctx)
}

`
