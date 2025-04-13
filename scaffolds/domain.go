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
	tree.AddDirectory("adapters/")
	tree.AddFile("model/model.go", templates.DomainExampleModel, []string{shared.ToolGoFmt})
	tree.AddFile("ports/in_example.go", templates.DomainExampleInPort, []string{shared.ToolGoFmt})
	tree.AddFile("ports/out_example.go", templates.DomainExampleOutPort, []string{shared.ToolGoFmt})
	tree.AddFile("usecase/usecase.go", templates.DomainExampleUsecase, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config":      cfg,
		"DomainName":  name,
		"DomainRoot":  cfg.Module + "/" + domainRoot,
		"UsecaseName": "Example",
	})
	if err != nil {
		return err
	}

	err = cfg.RunHook(shared.HookOnDomainCreated, map[string]any{
		"DomainName": name,
		"Module":     cfg.Module,
	})
	if err != nil {
		return err
	}

	return nil
}
