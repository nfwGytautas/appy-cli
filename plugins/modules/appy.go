package plugins_modules

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
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
	"mkdir":              mkdir,
	"copy_file":          copyFile,
	"execute_shell":      executeShell,
	"watch_directory":    watchDirectory,
	"get_project_name":   getProjectName,
}

var watchers = []*utils.Watcher{}

func AppyModuleLoader(l *lua.LState) int {
	// Create a new table for the appy module
	mod := l.SetFuncs(l.NewTable(), appyModuleExports)

	// Register the functions in the appy module
	l.SetField(mod, "version", lua.LString(shared.Version))

	l.SetField(mod, "FS_OP_CREATE", lua.LNumber(fsnotify.Create))
	l.SetField(mod, "FS_OP_WRITE", lua.LNumber(fsnotify.Write))
	l.SetField(mod, "FS_OP_REMOVE", lua.LNumber(fsnotify.Remove))
	l.SetField(mod, "FS_OP_RENAME", lua.LNumber(fsnotify.Rename))

	// Set the appy module in the global table
	l.Push(mod)
	return 1
}

func StopModuleWatchers() {
	for _, watcher := range watchers {
		watcher.Stop()
	}

	watchers = []*utils.Watcher{}
}

func getDomainRoot(l *lua.LState) int {
	domainRoot := l.Get(-1)
	l.Pop(1)

	if domainRoot == lua.LNil {
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	domainName, ok := domainRoot.(lua.LString)
	if !ok {
		l.RaiseError("expected 1st arg to be string got: %v'", domainRoot.Type())
		return 1
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.RaiseError("failed to get working dir: %v'", err)
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
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	if adapter == lua.LNil {
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	domainName, ok := domainRoot.(lua.LString)
	if !ok {
		l.RaiseError("expected 1st arg to be string got: %v'", domainRoot.Type())
		return 1
	}

	adapterName, ok := adapter.(lua.LString)
	if !ok {
		l.RaiseError("expected 2nd arg to be string got: %v'", adapter.Type())
		return 1
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.RaiseError("failed to get working dir: %v'", err)
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
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	if connector == lua.LNil {
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	domainName, ok := domainRoot.(lua.LString)
	if !ok {
		l.RaiseError("expected 1st arg to be string got: %v'", domainRoot.Type())
		return 1
	}

	connectorName, ok := connector.(lua.LString)
	if !ok {
		l.RaiseError("expected 2nd arg to be string got: %v'", connector.Type())
		return 1
	}

	cwd, err := os.Getwd()
	if err != nil {
		l.RaiseError("failed to get working dir: %v'", err)
		return 1
	}

	l.Push(lua.LString(fmt.Sprintf("%s/domains/%s/connectors/%s/", cwd, domainName.String(), connectorName.String())))
	return 1
}

func applyTemplate(l *lua.LState) int {
	argTable := l.Get(-1)
	l.Pop(1)

	if argTable == lua.LNil {
		utils.Console.ErrorLn("'apply_template': incorrect number of arguments'")
		l.Push(lua.LNil)
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	table, ok := argTable.(*lua.LTable)
	if !ok {
		l.RaiseError("expected 1st arg to be table got: %v'", argTable.Type())
		return 1
	}

	arguments := make(map[string]any)

	template := l.GetField(table, "template")
	target := l.GetField(table, "target")
	args := l.GetField(table, "args")

	if template == lua.LNil || target == lua.LNil {
		l.RaiseError("arguments, 'template' and 'target' keys are required'")
		return 1
	}

	templateName, ok := template.(lua.LString)
	if !ok {
		l.RaiseError("expected 'template' to be string got: %v'\n", template.Type())
		return 1
	}

	targetName, ok := target.(lua.LString)
	if !ok {
		l.RaiseError("expected 'target' to be string got: %v'\n", target.Type())
		return 1
	}

	if args != lua.LNil {
		argsTable, ok := args.(*lua.LTable)
		if !ok {
			l.RaiseError("expected 'args' to be table got: %v'\n", args.Type())
			return 1
		}

		argsTable.ForEach(func(k lua.LValue, v lua.LValue) {
			key, ok := k.(lua.LString)
			if !ok {
				l.RaiseError("expected key to be string got: %v'\n", k.Type())
				return
			}

			value, ok := v.(lua.LString)
			if !ok {
				l.RaiseError("expected value to be string got: %v'\n", v.Type())
				return
			}

			arguments[key.String()] = string(value)
		})
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(targetName.String()), 0755); err != nil {
		l.RaiseError("failed to create provider directory: %v'\n", err)
		return 1
	}

	// Load config
	arguments["Config"] = config.GetConfig()

	err := utils.TemplateAFile(templateName.String(), targetName.String(), arguments)
	if err != nil {
		l.RaiseError("error applying template: %v'\n", err)
		return 1
	}

	return 0
}

func mkdir(l *lua.LState) int {
	path := l.Get(-1)
	l.Pop(1)

	if path == lua.LNil {
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	pathName, ok := path.(lua.LString)
	if !ok {
		l.RaiseError("expected 1st arg to be string got: %v'\n", path.Type())
		return 1
	}

	err := os.MkdirAll(pathName.String(), 0755)
	if err != nil {
		l.RaiseError("failed to create directory: %v'\n", err)
		return 1
	}

	return 0
}

func copyFile(l *lua.LState) int {
	src := l.Get(-2)
	dst := l.Get(-1)
	l.Pop(2)

	if src == lua.LNil || dst == lua.LNil {
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	srcName, ok := src.(lua.LString)
	if !ok {
		l.RaiseError("expected 1st arg to be string got: %v'\n", src.Type())
		return 1
	}

	dstName, ok := dst.(lua.LString)
	if !ok {
		l.RaiseError("expected 2nd arg to be string got: %v'\n", dst.Type())
		return 1
	}

	err := os.MkdirAll(filepath.Dir(dstName.String()), 0755)
	if err != nil {
		l.RaiseError("failed to create directory: %v'\n", err)
		return 1
	}

	// Copy the file
	err = utils.CopyFile(srcName.String(), dstName.String())
	if err != nil {
		l.RaiseError("failed to copy file: %v'\n", err)
		return 1
	}

	return 0
}

func executeShell(l *lua.LState) int {
	commandToRun := ""

	shell := l.Get(-2)
	args := l.Get(-1)
	l.Pop(1)

	if shell == lua.LNil {
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	if args != lua.LNil {
		l.Pop(1)

		argsTable, ok := args.(*lua.LTable)
		if !ok {
			l.RaiseError("expected 2nd arg to be table got: %v'\n", args.Type())
			return 1
		}

		argsTable.ForEach(func(k lua.LValue, v lua.LValue) {
			value, ok := v.(lua.LString)
			if !ok {
				l.RaiseError("expected value to be string got: %v", k.Type())
				return
			}

			commandToRun += " " + string(value)
		})
	}

	shellCommand, ok := shell.(lua.LString)
	if !ok {
		l.RaiseError("expected 1st arg to be string got: %v'\n", shell.Type())
		return 1
	}

	commandToRun = string(shellCommand) + commandToRun

	cwd, err := os.Getwd()
	if err != nil {
		l.RaiseError("failed to get working dir: %v'\n", err)
		return 1
	}

	err = utils.RunCommand(cwd, commandToRun)
	if err != nil {
		l.RaiseError("failed to execute command: %v'\n", err)
		return 1
	}

	return 0
}

func watchDirectory(l *lua.LState) int {
	path := l.Get(-2)
	callback := l.Get(-1)
	l.Pop(2)

	if path == lua.LNil || callback == lua.LNil {
		l.RaiseError("incorrect number of arguments'")
		return 1
	}

	pathName, ok := path.(lua.LString)
	if !ok {
		l.RaiseError("expected 1st arg to be string got: %v'\n", path.Type())
		return 1
	}

	callbackFunc, ok := callback.(*lua.LFunction)
	if !ok {
		l.RaiseError("expected 2nd arg to be function got: %v'\n", callback.Type())
		return 1
	}

	watcher, err := utils.NewWatcher(string(pathName), func(event fsnotify.Event) {
		if event.Op&fsnotify.Chmod == fsnotify.Chmod {
			return
		}

		err := l.CallByParam(lua.P{
			Fn:      callbackFunc,
			NRet:    0,
			Protect: true,
		}, lua.LString(event.Name), lua.LNumber(event.Op))
		if err != nil {
			utils.Console.ErrorLn("failed to call callback function: %v", err)
			return
		}
	})
	if err != nil {
		l.RaiseError("failed to watch directory: %v'\n", err)
		return 1
	}

	watcher.Start()

	watchers = append(watchers, watcher)

	return 0
}

func getProjectName(l *lua.LState) int {
	l.Push(lua.LString(config.GetConfig().Project))
	return 1
}
