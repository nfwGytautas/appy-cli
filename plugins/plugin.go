package plugins

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type PluginMetaFields struct {
	Root string
}

type plugin struct {
	engine *PluginEngine
	t      *lua.LTable

	// Hooks
	onLoad             lua.LValue
	onDomainCreated    lua.LValue
	onConnectorCreated lua.LValue
	onAdapterCreated   lua.LValue
}

func newPlugin(pe *PluginEngine, t *lua.LTable) *plugin {
	p := plugin{
		engine: pe,
		t:      t,

		// Hooks
		onLoad:             nil,
		onDomainCreated:    nil,
		onConnectorCreated: nil,
		onAdapterCreated:   nil,
	}

	// Resolve hooks
	p.loadHook(&p.onLoad, "on_load")
	p.loadHook(&p.onDomainCreated, "on_domain_created")
	p.loadHook(&p.onConnectorCreated, "on_connector_created")
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
	p.engine.luaState.SetField(p.t, "root", lua.LString(fields.Root))
}

func (p *plugin) String() string {
	return fmt.Sprintf(`
  Hooks:
    + onLoad: %v
    + onDomainCreated: %v
    + onConnectorCreated: %v
    + onAdapterCreated: %v
	`,
		p.onLoad != nil,
		p.onDomainCreated != nil,
		p.onConnectorCreated != nil,
		p.onAdapterCreated != nil,
	)
}

func (p *plugin) OnLoad() error {
	if p.onLoad == nil {
		return nil
	}

	err := p.engine.luaState.CallByParam(lua.P{
		Fn:      p.onLoad,
		NRet:    0,
		Protect: true,
	})

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
	}, lua.LString(name))

	if err != nil {
		return err
	}

	return nil
}

func (p *plugin) OnConnectorCreated(domain string, connector string) error {
	if p.onConnectorCreated == nil {
		return nil
	}

	err := p.engine.luaState.CallByParam(lua.P{
		Fn:      p.onConnectorCreated,
		NRet:    0,
		Protect: true,
	}, lua.LString(domain), lua.LString(connector))

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
	}, lua.LString(domain), lua.LString(adapter))

	if err != nil {
		return err
	}

	return nil
}
