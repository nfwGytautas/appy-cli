package variant_plugin

const templateGoMod = //
`module {{.Module}}

go 1.23

require github.com/nfwGytautas/appy-go v0.2.0
`

const templateGitignore = //
`# Go
bin/
build/
dist/

# Appy
.appy/

# General
.env
.DS_Store
`

const templateReadmeMd = //
`# {{.Config.Name}}

## Changelog

### 0.1.0 [xxxx/xx/xxx]

#### Changes

- TODO

#### Bugfixes

- N/A
`

const templatePluginLua = //
`local appy = require("appy")

-- This is a plugin for the appy framework
-- It provides a basic structure for a plugin
-- You can remove hooks that you don't need

local plugin = {}

function plugin:on_configure()
	-- This function is called when the plugin is configured
end

function plugin:on_load()
    -- This function is called when the plugin is loaded
end

function plugin:on_domain_created(domain)
    -- This function is called when a domain is created
end

function plugin:on_adapter_created(domain, adapter)
    -- This function is called when an adapter is created inside a domain
end

return plugin
`

const templateConfigGo = //
`package providers_{{.Config.Name}}

//
// This file is an auto-generated template for appy plugin, do not remove Init or Start functions
//

type Provider struct {
}

func Init() (*Provider, error) {
	provider := &Provider{}

	// TODO: Implement the initialization logic of plugin in code here

	return provider, nil
}

func (p *Provider) Start() error {
	// TODO: Implement the start logic of plugin in code here
	return nil
}
`
