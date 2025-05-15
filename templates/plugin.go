package templates

const PluginLua = //
`local appy = require("appy")

-- This is a plugin for the appy framework
-- It provides a basic structure for a plugin

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
