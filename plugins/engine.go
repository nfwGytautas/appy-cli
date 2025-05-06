package plugins

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

type PluginEngine struct {
	luaState      *lua.LState
	errorHandler  *lua.LFunction
	loadedPlugins []*Plugin
}

func NewPluginEngine(config map[string]any) *PluginEngine {
	pe := PluginEngine{
		luaState:      lua.NewState(),
		loadedPlugins: []*Plugin{},
	}

	pe.luaState.PreloadModule("appy", appyModuleLoader(config))
	pe.errorHandler = pe.luaState.NewFunction(func(l *lua.LState) int {
		err := l.ToString(1)
		if err != "" {
			l.RaiseError(err)
		}
		return 0
	})

	return &pe
}

func (pe *PluginEngine) Shutdown() {
	stopModuleWatchers()
	pe.luaState.Close()
}

func (pe *PluginEngine) LoadPlugin(path string) (*Plugin, error) {
	err := pe.luaState.DoFile(path)
	if err != nil {
		return nil, err
	}

	p, ok := pe.luaState.Get(-1).(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("'%v' is not a valid plugin", path)
	}

	pInstance := newPlugin(pe, p)

	pe.loadedPlugins = append(pe.loadedPlugins, pInstance)

	return pInstance, nil
}

func (pe *PluginEngine) GetLoadedPlugins() []*Plugin {
	return pe.loadedPlugins
}
