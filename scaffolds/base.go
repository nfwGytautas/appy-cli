package scaffolds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
)

func Base(module string) error {
	tree := utils.GeneratedFileTree{}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	projectName := filepath.Base(dir)

	tree.AddFile("go.mod", templates.GoMod, nil)
	tree.AddFile("main.go", templates.MainGo, []string{shared.ToolGoFmt})
	tree.AddFile("wiring.go", templates.WiringGo, []string{shared.ToolGoFmt})
	tree.AddFile("README.md", templates.ReadmeMd, nil)
	tree.AddFile("Dockerfile", templates.Dockerfile, nil)
	tree.AddFile(".gitignore", templates.Gitignore, nil)
	tree.AddFile(".vscode/Snippets.code-snippets", templates.VscodeSnippets, nil)
	tree.AddFile(".github/build.yaml", templates.GithubBuildYaml, nil)
	tree.AddFile(".appy/appy.yaml", templates.AppyYaml, nil)

	err = tree.Generate(map[string]string{
		"ProjectName": projectName,
		"Version":     shared.Version,
		"Module":      module,
	})
	if err != nil {
		return err
	}

	return nil
}

func Scaffold(scaffoldType string) error {
	cfg, err := shared.LoadConfig()
	if err != nil {
		return err
	}

	switch scaffoldType {
	case shared.ScaffoldHMD:
		err = scaffoldHMD(cfg)
	case shared.ScaffoldHSS:
		err = scaffoldHSS(cfg)
	default:
		return fmt.Errorf("invalid type: %s", scaffoldType)
	}

	if err != nil {
		return err
	}

	err = cfg.Save()
	if err != nil {
		return err
	}

	return nil
}
