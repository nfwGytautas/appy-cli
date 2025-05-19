package variant_plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
)

type generateOpts struct {
	isRoot bool
	module string
}

func (cfg *Config) Scaffold() error {
	utils.Console.DebugLn("Scaffolding Plugin...")

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	opts := generateOpts{}

	opts.isRoot, err = promptRoot()
	if err != nil {
		return err
	}

	if opts.isRoot {
		opts.module, err = promptModuleName()
		if err != nil {
			return err
		}
	}

	cfg.Type = VariantName
	cfg.Version = shared.Version
	cfg.Name = filepath.Base(dir)
	cfg.Description = "Your plugin description here"

	err = generateFolderStructure(cfg, opts)
	if err != nil {
		return err
	}

	err = cfg.Save()
	if err != nil {
		return err
	}

	return nil
}

func promptRoot() (bool, error) {
	prompt := promptui.Select{
		Label: "Is this a single-repo plugin?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	return result == "Yes", nil
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

func generateFolderStructure(cfg *Config, opts generateOpts) error {
	utils.Console.InfoLn("Generating folder structure...")

	tree := utils.GeneratedFileTree{}

	if opts.isRoot {
		tree.AddFile("go.mod", templateGoMod, nil)
		tree.AddFile(".gitignore", templateGitignore, nil)
	}

	tree.AddDirectory("providers/")
	tree.AddDirectory("impl/")
	tree.AddDirectory("templates/")
	tree.AddFile("providers/config.go", templateConfigGo, nil)
	tree.AddFile("README.md", templateReadmeMd, nil)
	tree.AddFile("plugin.lua", templatePluginLua, nil)

	err := tree.Generate(map[string]any{
		"Config": cfg,
		"Module": opts.module,
	})
	if err != nil {
		return err
	}

	return nil
}
