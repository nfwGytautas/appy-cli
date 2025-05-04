package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/plugins"
	"github.com/nfwGytautas/appy-cli/project"
	"github.com/nfwGytautas/appy-cli/scaffolds"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
)

func main() {
	// Check if verbose flag is set
	flag.BoolVar(&utils.Verbose, "debug", false, "Debug output")
	flag.Parse()

	pe := plugins.NewPluginEngine()

	p, err := pe.LoadPlugin("plugin.lua")
	if err != nil {
		utils.Console.Fatal(err)
	}

	log.Println(p.String())

	cwd, err := os.Getwd()
	if err != nil {
		utils.Console.Fatal(err)
	}

	p.SetMetaFields(plugins.PluginMetaFields{
		ScriptRoot:   cwd + "/",
		ProviderRoot: cwd + "/",
	})

	err = p.OnLoad()
	if err != nil {
		utils.Console.Fatal(err)
	}

	err = p.OnDomainCreated("test")
	if err != nil {
		utils.Console.Fatal(err)
	}

	err = p.OnAdapterCreated("test", "adapter")
	if err != nil {
		utils.Console.Fatal(err)
	}

	err = p.OnConnectorCreated("test", "connector")
	if err != nil {
		utils.Console.Fatal(err)
	}

	defer pe.Shutdown()

	// Block
	<-make(chan struct{})

	return

	utils.Console.ClearEntireConsole()

	// Check if empty directory
	entries, err := os.ReadDir(".")
	if err != nil {
		utils.Console.Fatal(err)
	}

	if len(entries) == 0 || (len(entries) == 1 && entries[0].Name() == ".git") {
		utils.Console.InfoLn("Empty project. Scaffolding...")
		err = scaffold()
		if err != nil {
			utils.Console.Fatal(err)
		}
	}

	// Check if config exists
	if _, err := os.Stat("appy.yaml"); os.IsNotExist(err) {
		utils.Console.InfoLn("appy.yaml not found. Either misconfigured or incorrect directory.")
		return
	}

	utils.Console.ClearEntireConsole()

	// TODO: Terminal ui

	project.Watch()
	if err != nil {
		utils.Console.Fatal(err)
	}
}

func scaffold() error {
	{
		prompt := promptui.Select{
			Label: "Select an option",
			Items: []string{
				shared.ScaffoldHMD,
			},
		}

		_, scaffoldType, err := prompt.Run()
		if err != nil {
			return err
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
			return err
		}

		err = scaffolds.Scaffold(module, scaffoldType)
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
			return err
		}
	}

	return nil
}
