package scaffolds

import (
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
)

func ScaffoldDomain(cfg *shared.Config, name string) error {
	tree := utils.GeneratedFileTree{}

	domainRoot := cfg.GetDomainsRoot(name)

	tree.SetPrefix(domainRoot)
	tree.AddFile("model/model.go", templates.DomainExampleModel, []string{shared.ToolGoFmt})
	tree.AddFile("ports/in/ports_in.go", templates.DomainExampleInPort, []string{shared.ToolGoFmt})
	tree.AddFile("ports/out/ports_out.go", templates.DomainExampleOutPort, []string{shared.ToolGoFmt})
	tree.AddFile("usecase/usecase.go", templates.DomainExampleUsecase, []string{shared.ToolGoFmt})
	tree.AddFile("adapters/in/adapters_in.go", templates.DomainExampleInAdapter, []string{shared.ToolGoFmt})
	tree.AddFile("adapters/out/adapters_out.go", templates.DomainExampleOutAdapter, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config":      cfg,
		"DomainName":  name,
		"DomainRoot":  cfg.Module + "/" + domainRoot,
		"UsecaseName": "Example",
	})
	if err != nil {
		return err
	}

	return nil
}
