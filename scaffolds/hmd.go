package scaffolds

import (
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
)

func scaffoldHMD(cfg *config.AppyConfig) error {
	tree := utils.GeneratedFileTree{}

	cfg.Type = shared.ScaffoldHMD

	tree.AddDirectory("domains/")
	tree.AddFile("shared/errors.go", templates.ErrorsGo, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config": cfg,
	})
	if err != nil {
		return err
	}
	return nil
}
