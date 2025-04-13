package scaffolds

import (
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
)

func scaffoldHSS(cfg *shared.Config) error {
	tree := utils.GeneratedFileTree{}

	cfg.Type = shared.ScaffoldHSS

	tree.AddDirectory("domain/")
	tree.AddFile("shared/errors.go", templates.ErrorsGo, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config": cfg,
	})
	if err != nil {
		return err
	}

	err = ScaffoldDomain(cfg, "domain")
	if err != nil {
		return err
	}

	return nil
}
