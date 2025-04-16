package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/config"
	"github.com/nfwGytautas/appy-cli/scaffolds"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
	"github.com/nfwGytautas/appy-cli/watchers"
)

func main() {
	// Check if verbose flag is set
	flag.BoolVar(&utils.Verbose, "debug", false, "Debug output")
	flag.Parse()

	utils.Console.ClearEntireConsole()

	// Check if empty directory
	entries, err := os.ReadDir(".")
	if err != nil {
		utils.Console.Fatal(err)
	}

	if len(entries) == 0 || (len(entries) == 1 && entries[0].Name() == ".git") {
		utils.Console.InfoLn("Empty project. Scaffolding...")
		scaffold()
	}

	// Check if config exists
	if _, err := os.Stat("appy.yaml"); os.IsNotExist(err) {
		utils.Console.InfoLn("appy.yaml not found. Either misconfigured or incorrect directory.")
		return
	}

	utils.Console.ClearEntireConsole()

	cfg, err := config.LoadConfig()
	if err != nil {
		utils.Console.Fatal(err)
	}

	// TODO: Terminal ui

	err = cfg.StartProviders()
	if err != nil {
		utils.Console.Fatal(err)
	}

	watchers.Watch()
	if err != nil {
		utils.Console.Fatal(err)
	}
}

func scaffold() {
	{
		prompt := promptui.Select{
			Label: "Select an option",
			Items: []string{
				shared.ScaffoldHMD,
				shared.ScaffoldHSS,
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			utils.Console.Fatal(err)
		}

		utils.Console.ClearLines(2)

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
			utils.Console.Fatal(err)
		}

		utils.Console.InfoLn("Scaffolding: %s\n", result)

		err = scaffolds.Base(module)
		if err != nil {
			utils.Console.Fatal(err)
		}

		err = scaffolds.Scaffold(result)
		if err != nil {
			utils.Console.Fatal(err)
		}

		utils.Console.InfoLn("Done!")
	}

	{
		prompt := promptui.Prompt{
			Label: "Press Enter to continue",
			Validate: func(input string) error {
				return nil
			},
		}

		_, err := prompt.Run()
		if err != nil {
			utils.Console.Fatal(err)
		}
	}
}
