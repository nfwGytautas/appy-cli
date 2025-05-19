package variant_hmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/shared"
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
	tree.AddFile("go.mod", templateGoMod, nil)

	tree.AddFile("main.go", templateMainGo, []string{shared.ToolGoFmt})
	tree.AddFile("README.md", templateReadmeMd, nil)
	tree.AddFile("Dockerfile", templateDockerfile, nil)

	tree.AddFile(".gitignore", templateGitignore, nil)

	tree.AddFile(".vscode/Snippets.code-snippets", templateVscodeSnippets, nil)
	tree.AddFile(".vscode/settings.json", templateVscodeSettings, nil)
	tree.AddFile(".github/build.yaml", templateGithubBuildYaml, nil)

	tree.AddDirectory("shared/")
	tree.AddFile("shared/errors.go", templateErrorsGo, []string{shared.ToolGoFmt})

	tree.AddDirectory("providers/")
	tree.AddFile("providers/providers.go", templateProvidersGo, []string{shared.ToolGoImports, shared.ToolGoFmt})

	tree.AddDirectory("domains/")
	tree.AddFile("domains/domains.go", templateDomainsGo, []string{shared.ToolGoImports, shared.ToolGoFmt})

	tree.AddDirectory("domains/example/")
	tree.AddDirectory("domains/example/adapters/")
	tree.AddDirectory("domains/example/model/")
	tree.AddFile("domains/example/domain.go", templateDomainExampleDomain, []string{shared.ToolGoFmt})
	tree.AddFile("domains/example/ping.go", templateDomainExampleUsecase, []string{shared.ToolGoFmt})
	tree.AddFile("domains/example/model/example.go", templateDomainExampleModel, []string{shared.ToolGoFmt})

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
