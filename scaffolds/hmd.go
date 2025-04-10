package scaffolds

import (
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
)

func scaffoldHMD(cfg *shared.Config) error {
	tree := utils.GeneratedFileTree{}

	cfg.Type = shared.ScaffoldHMD

	tree.AddDirectory("domains/")
	tree.AddDirectory("repositories/")
	tree.AddFile("shared/errors.go", templates.ErrorsGo, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config": cfg,
	})
	if err != nil {
		return err
	}
	return nil
}
