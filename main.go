package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/scaffolds"
	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/utils"
	"github.com/nfwGytautas/appy-cli/watchers"
)

func main() {
	utils.ClearEntireConsole()
	printHeader()

	// Check if empty directory
	entries, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	if len(entries) == 0 || (len(entries) == 1 && entries[0].Name() == ".git") {
		fmt.Println("Empty project. Scaffolding...")
		scaffold()
	}

	// Check if config exists
	if _, err := os.Stat(".appy/appy.yaml"); os.IsNotExist(err) {
		fmt.Println(".appy/appy.yaml not found. Either misconfigured or incorrect directory.")
		return
	}

	utils.ClearEntireConsole()
	printHeader()

	watchers.Watch()
	if err != nil {
		log.Fatal(err)
	}
}

func printHeader() {
	fmt.Println("------------------------------------------------------------------------------------------------")
	fmt.Printf("                                   Appy CLI: %s\n", shared.Version)
	fmt.Println("------------------------------------------------------------------------------------------------")
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
			log.Fatal(err)
		}

		utils.ClearLines(2)

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
			log.Fatal(err)
		}

		fmt.Printf("Scaffolding: %s\n", result)

		err = scaffolds.Base(module)
		if err != nil {
			log.Fatal(err)
		}

		err = scaffolds.Scaffold(result)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Done!")
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
			log.Fatal(err)
		}
	}

	utils.ClearLines(4)
}
