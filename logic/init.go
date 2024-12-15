package logic

import "github.com/nfwGytautas/appy-cli/templates"

var baseStructure = []string{
	".appy",
	"endpoints",
	"jobs",
}

func Init() {
	var err error

	// Make sure we are in an empty directory
	cwd := getCurrentDirectory()
	err = ensureDirectoryIsEmpty(cwd)
	errorCheckAndPanic(err)

	// Create base structure
	chainErrorChecks(
		ensureDirectories(baseStructure),
		createFile(".gitignore", templates.Gitignore),
		createFile("go.mod", templates.GoMod),
		createFile("config.go", templates.ConfigGo),
		createFile("main.go", templates.MainGo),
		createFile(".appy/hook.go", templates.AutogenHookGo),
	)
}
