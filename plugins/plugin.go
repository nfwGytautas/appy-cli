package plugins

import (
	"fmt"

	"github.com/nfwGytautas/appy-cli/utils"
	lua "github.com/yuin/gopher-lua"
)

type PluginMetaFields struct {
	ScriptRoot   string
	ProviderRoot string
}

type Plugin struct {
	engine *PluginEngine
	t      *lua.LTable

	// Hooks
	onConfigure      lua.LValue
	onLoad           lua.LValue
	onDomainCreated  lua.LValue
	onAdapterCreated lua.LValue
}

func newPlugin(pe *PluginEngine, t *lua.LTable) *Plugin {
	p := Plugin{
		engine: pe,
		t:      t,

		// Hooks
		onConfigure:      nil,
		onLoad:           nil,
		onDomainCreated:  nil,
		onAdapterCreated: nil,
	}

	// Resolve hooks
	p.loadHook(&p.onConfigure, "on_configure")
	p.loadHook(&p.onLoad, "on_load")
	p.loadHook(&p.onDomainCreated, "on_domain_created")
	p.loadHook(&p.onAdapterCreated, "on_adapter_created")

	return &p
}

func (p *Plugin) loadHook(h *lua.LValue, name string) {
	hook := p.engine.luaState.GetField(p.t, name)
	if hook == lua.LNil {
		return
	}

	*h = hook
}

func (p *Plugin) SetMetaFields(fields PluginMetaFields) {
	p.engine.luaState.SetField(p.t, "script_root", lua.LString(fields.ScriptRoot))
	p.engine.luaState.SetField(p.t, "provider_root", lua.LString(fields.ProviderRoot))
}

func (p *Plugin) String() string {
	return fmt.Sprintf(`
  Hooks:
    + onConfigure: %v
    + onLoad: %v
    + onDomainCreated: %v
    + onAdapterCreated: %v
	`,
		p.onConfigure != nil,
		p.onLoad != nil,
		p.onDomainCreated != nil,
		p.onAdapterCreated != nil,
	)
}

func (p *Plugin) OnConfigure() error {
	if p.onConfigure == nil {
		return nil
	}

	utils.Console.DebugLn("Plugin %p: onConfigure", p)

	err := p.engine.luaState.CallByParam(lua.P{
		Fn:      p.onConfigure,
		NRet:    0,
		Protect: true,
	}, p.t)

	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) OnLoad() error {
	if p.onLoad == nil {
		return nil
	}

	utils.Console.DebugLn("Plugin %p: onLoad", p)

	err := p.engine.luaState.CallByParam(lua.P{
		Fn:      p.onLoad,
		NRet:    0,
		Protect: true,
	}, p.t)

	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) OnDomainCreated(name string) error {
	if p.onDomainCreated == nil {
		return nil
	}

	utils.Console.DebugLn("Plugin %p: onDomainCreated", p)

	err := p.engine.luaState.CallByParam(lua.P{
		Fn:      p.onDomainCreated,
		NRet:    0,
		Protect: true,
	}, p.t, lua.LString(name))

	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) OnAdapterCreated(domain string, adapter string) error {
	if p.onAdapterCreated == nil {
		return nil
	}

	utils.Console.DebugLn("Plugin %p: onAdapterCreated", p)

	err := p.engine.luaState.CallByParam(lua.P{
		Fn:      p.onAdapterCreated,
		NRet:    0,
		Protect: true,
		Handler: p.engine.errorHandler,
	}, p.t, lua.LString(domain), lua.LString(adapter))

	if err != nil {
		return err
	}

	return nil
}
