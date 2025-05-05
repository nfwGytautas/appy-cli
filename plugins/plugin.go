package plugins

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type PluginMetaFields struct {
	ScriptRoot   string
	ProviderRoot string
}

type plugin struct {
	engine *PluginEngine
	t      *lua.LTable

	// Hooks
	onConfigure      lua.LValue
	onLoad           lua.LValue
	onDomainCreated  lua.LValue
	onAdapterCreated lua.LValue
}

func newPlugin(pe *PluginEngine, t *lua.LTable) *plugin {
	p := plugin{
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

func (p *plugin) loadHook(h *lua.LValue, name string) {
	hook := p.engine.luaState.GetField(p.t, name)
	if hook == lua.LNil {
		return
	}

	*h = hook
}

func (p *plugin) SetMetaFields(fields PluginMetaFields) {
	p.engine.luaState.SetField(p.t, "script_root", lua.LString(fields.ScriptRoot))
	p.engine.luaState.SetField(p.t, "provider_root", lua.LString(fields.ProviderRoot))
}

func (p *plugin) String() string {
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

func (p *plugin) OnConfigure() error {
	if p.onConfigure == nil {
		return nil
	}

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

func (p *plugin) OnLoad() error {
	if p.onLoad == nil {
		return nil
	}

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

func (p *plugin) OnDomainCreated(name string) error {
	if p.onDomainCreated == nil {
		return nil
	}

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

func (p *plugin) OnAdapterCreated(domain string, adapter string) error {
	if p.onAdapterCreated == nil {
		return nil
	}

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
