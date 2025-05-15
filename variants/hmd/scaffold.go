package variant_hmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/nfwGytautas/appy-cli/utils"
)

func (cfg *Config) Scaffold() error {
	utils.Console.DebugLn("Scaffolding HMD...")

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	moduleName, err := promptModuleName()
	if err != nil {
		return err
	}

	cfg.Type = VariantName
	cfg.Module = moduleName
	cfg.Project = filepath.Base(dir)
	cfg.Version = shared.Version

	err = generateFolderStructure(cfg)
	if err != nil {
		return err
	}

	err = cfg.Reconfigure()
	if err != nil {
		return err
	}

	return nil
}

func promptModuleName() (string, error) {
	promptModule := promptui.Prompt{
		Label: "Enter module name",
		Validate: func(input string) error {
			if input == "" {
				return fmt.Errorf("module name is required")
			}

			matched, err := regexp.MatchString(`^[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)?\/[a-zA-Z0-9]+\/[a-zA-Z0-9\-]+$`, input)
			if err != nil {
				return err
			}
			if !matched {
				return fmt.Errorf("module name must follow the pattern: domain.extension/module/repository")
			}

			return nil
		},
	}

	module, err := promptModule.Run()
	if err != nil {
		return "", err
	}

	return module, nil
}

func generateFolderStructure(cfg *Config) error {
	utils.Console.InfoLn("Generating folder structure...")

	tree := utils.GeneratedFileTree{}

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
	cfg.Repositories = append(cfg.Repositories, &Repository{
		Url:       "https://github.com/nfwGytautas/appy-providers",
		Branch:    "main",
		Providers: []Provider{},
	})

	return nil
}
