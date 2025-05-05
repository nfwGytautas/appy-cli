package plugins

import (
	"fmt"

	plugins_modules "github.com/nfwGytautas/appy-cli/plugins/modules"
	lua "github.com/yuin/gopher-lua"
)

type PluginEngine struct {
	luaState *lua.LState

	errorHandler *lua.LFunction
}

func NewPluginEngine() *PluginEngine {
	pe := PluginEngine{
		luaState: lua.NewState(),
	}

	pe.luaState.PreloadModule("appy", plugins_modules.AppyModuleLoader)
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
	plugins_modules.StopModuleWatchers()
	pe.luaState.Close()
}

func (pe *PluginEngine) LoadPlugin(path string) (*plugin, error) {
	err := pe.luaState.DoFile(path)
	if err != nil {
		return nil, err
	}

	p, ok := pe.luaState.Get(-1).(*lua.LTable)
	if !ok {
		return nil, fmt.Errorf("'%v' is not a valid plugin", path)
	}

	return newPlugin(pe, p), nil
}
