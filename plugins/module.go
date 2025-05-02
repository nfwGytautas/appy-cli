package plugins

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
	lua "github.com/yuin/gopher-lua"
)

var appyModuleExports = map[string]lua.LGFunction{
	"get_domain_root":    getDomainRoot,
	"get_adapter_root":   getAdapterRoot,
	"get_connector_root": getConnectorRoot,
	"apply_template":     applyTemplate,
}

func appyModuleLoader(l *lua.LState) int {
	// Create a new table for the appy module
	mod := l.SetFuncs(l.NewTable(), appyModuleExports)

	// Register the functions in the appy module
	l.SetField(mod, "version", lua.LString(shared.Version))

	// Set the appy module in the global table
	l.Push(mod)
	return 1
}

func getDomainRoot(l *lua.LState) int {
	domainRoot := l.Get(-1)
	l.Pop(1)

	if domainRoot == lua.LNil {
		utils.Console.ErrorLn("'get_domain_root: incorrect number of arguments'")
		l.Push(lua.LNil)
		return 1
	}

	domainName, ok := domainRoot.(lua.LString)
	if !ok {
		utils.Console.Error("'get_domain_root: expected 1st arg to be string got: %v'\n", domainRoot.Type())
		l.Push(lua.LNil)
		return 1
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.Push(lua.LNil)
		return 1
	}

	l.Push(lua.LString(fmt.Sprintf("%s/domains/%s", cwd, domainName.String())))
	return 1
}

func getAdapterRoot(l *lua.LState) int {
	domainRoot := l.Get(-2)
	adapter := l.Get(-1)
	l.Pop(2)

	if domainRoot == lua.LNil {
		utils.Console.ErrorLn("'get_adapter_root: incorrect number of arguments'")
		l.Push(lua.LNil)
		return 1
	}

	if adapter == lua.LNil {
		utils.Console.ErrorLn("'get_adapter_root: incorrect number of arguments'")
		l.Push(lua.LNil)
		return 1
	}

	domainName, ok := domainRoot.(lua.LString)
	if !ok {
		utils.Console.Error("'get_adapter_root: expected 1st arg to be string got: %v'\n", domainRoot.Type())
		l.Push(lua.LNil)
		return 1
	}

	adapterName, ok := adapter.(lua.LString)
	if !ok {
		utils.Console.Error("'get_adapter_root: expected 2nd arg to be string got: %v'\n", adapter.Type())
		l.Push(lua.LNil)
		return 1
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.Push(lua.LNil)
		return 1
	}

	l.Push(lua.LString(fmt.Sprintf("%s/domains/%s/adapters/%s/", cwd, domainName.String(), adapterName.String())))
	return 1
}

func getConnectorRoot(l *lua.LState) int {
	domainRoot := l.Get(-2)
	connector := l.Get(-1)

	l.Pop(2)

	if domainRoot == lua.LNil {
		utils.Console.ErrorLn("'get_connector_root: incorrect number of arguments'")
		l.Push(lua.LNil)
		return 1
	}

	if connector == lua.LNil {
		utils.Console.ErrorLn("'get_connector_root: incorrect number of arguments'")
		l.Push(lua.LNil)
		return 1
	}

	domainName, ok := domainRoot.(lua.LString)
	if !ok {
		utils.Console.Error("'get_connector_root: expected 1st arg to be string got: %v'\n", domainRoot.Type())
		l.Push(lua.LNil)
		return 1
	}

	connectorName, ok := connector.(lua.LString)
	if !ok {
		utils.Console.Error("'get_connector_root: expected 2nd arg to be string got: %v'\n", connector.Type())
		l.Push(lua.LNil)
		return 1
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.Push(lua.LNil)
		return 1
	}

	l.Push(lua.LString(fmt.Sprintf("%s/domains/%s/connectors/%s/", cwd, domainName.String(), connectorName.String())))
	return 1
}

func applyTemplate(l *lua.LState) int {
	utils.Console.DebugLn("apply_template called")
	argTable := l.Get(-1)
	l.Pop(1)

	if argTable == lua.LNil {
		utils.Console.ErrorLn("'apply_template': incorrect number of arguments'")
		l.Push(lua.LNil)
		return 1
	}

	table, ok := argTable.(*lua.LTable)
	if !ok {
		utils.Console.Error("'apply_template': expected 1st arg to be table got: %v'\n", argTable.Type())
		l.Push(lua.LNil)
		return 1
	}

	arguments := make(map[string]any)

	template := l.GetField(table, "template")
	target := l.GetField(table, "target")
	args := l.GetField(table, "args")

	if template == lua.LNil || target == lua.LNil {
		utils.Console.ErrorLn("'apply_template': arguments, 'template' and 'target' keys are required'")
		l.Push(lua.LNil)
		return 1
	}

	templateName, ok := template.(lua.LString)
	if !ok {
		utils.Console.Error("'apply_template': expected 'template' to be string got: %v'\n", template.Type())
		l.Push(lua.LNil)
		return 1
	}

	targetName, ok := target.(lua.LString)
	if !ok {
		utils.Console.Error("'apply_template': expected 'target' to be string got: %v'\n", target.Type())
		l.Push(lua.LNil)
		return 1
	}

	if args != lua.LNil {
		argsTable, ok := args.(*lua.LTable)
		if !ok {
			utils.Console.Error("'apply_template': expected 'args' to be table got: %v'\n", args.Type())
			l.Push(lua.LNil)
			return 1
		}

		argsTable.ForEach(func(k lua.LValue, v lua.LValue) {
			key, ok := k.(lua.LString)
			if !ok {
				utils.Console.Error("'apply_template': expected key to be string got: %v'\n", k.Type())
				l.Push(lua.LNil)
				return
			}

			arguments[key.String()] = v
		})
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetName.String()), 0755); err != nil {
		utils.Console.Error("'apply_template': failed to create provider directory: %v'\n", err)
		l.Push(lua.LNil)
		return 1
	}

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Console.Error("'apply_template': failed to load config: %v\n", err)
		l.Push(lua.LNil)
		return 1
	}

	arguments["Config"] = cfg

	err = utils.TemplateAFile(templateName.String(), targetName.String(), arguments)
	if err != nil {
		utils.Console.Error("'apply_template': error applying template: %v'\n", err)
		l.Push(lua.LNil)
		return 1
	}

	return 0
}
