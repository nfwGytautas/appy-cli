package scaffolds

import (
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
)

func scaffoldHMD(cfg *config.AppyConfig) error {
	tree := utils.GeneratedFileTree{}

	// tree.AddDirectory("tests/")
	// tree.AddDirectory("tests/mocks/")
	// tree.AddDirectory("tests/unit/")
	// tree.AddDirectory("tests/integration/")

	tree.AddDirectory(".appy/")
	tree.AddFile("go.mod", templates.GoMod, nil)

	tree.AddFile("main.go", templates.MainGo, []string{shared.ToolGoFmt})
	tree.AddFile("README.md", templates.ReadmeMd, nil)
	tree.AddFile("Dockerfile", templates.Dockerfile, nil)

	tree.AddFile(".gitignore", templates.Gitignore, nil)

	tree.AddFile(".vscode/Snippets.code-snippets", templates.VscodeSnippets, nil)
	tree.AddFile(".vscode/settings.json", templates.VscodeSettings, nil)
	tree.AddFile(".github/build.yaml", templates.GithubBuildYaml, nil)

	tree.AddDirectory("shared/")
	tree.AddFile("shared/errors.go", templates.ErrorsGo, []string{shared.ToolGoFmt})

	tree.AddDirectory("providers/")
	tree.AddFile("providers/providers.go", templates.ProvidersGo, []string{shared.ToolGoImports, shared.ToolGoFmt})

	tree.AddDirectory("domains/")
	tree.AddFile("domains/domains.go", templates.DomainsGo, []string{shared.ToolGoImports, shared.ToolGoFmt})

	tree.AddDirectory("domains/example/")
	tree.AddDirectory("domains/example/adapters/")
	tree.AddDirectory("domains/example/connectors/")
	tree.AddDirectory("domains/example/model/")
	tree.AddFile("domains/example/domain.go", templates.DomainExampleDomain, []string{shared.ToolGoFmt})
	tree.AddFile("domains/example/ping.go", templates.DomainExampleUsecase, []string{shared.ToolGoFmt})
	tree.AddFile("domains/example/model/example.go", templates.DomainExampleModel, []string{shared.ToolGoFmt})

	err := tree.Generate(map[string]any{
		"Config":      cfg,
		"DomainName":  "example",
		"UsecaseName": "ping",
	})
	if err != nil {
		return err
	}

	// Add default providers
	cfg.Repositories = append(cfg.Repositories, &config.Repository{
		Url:       "https://github.com/nfwGytautas/appy-providers",
		Branch:    "main",
		Providers: []config.Provider{},
	})

	return nil
}
