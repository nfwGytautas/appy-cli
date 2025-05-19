package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/nfwGytautas/appy-cli/utils"
	"github.com/nfwGytautas/appy-cli/variants"
)

func main() {
	// Check if verbose flag is set
	flag.BoolVar(&utils.Verbose, "debug", false, "Debug output")
	flag.Parse()

	utils.Console.ClearEntireConsole()

	// Get variant
	variant, err := variants.Load()
	if err != nil {
		utils.Console.Fatal(err)
	}

	if variant == nil {
		// No variant scaffold one
		variantType, err := promptVariantScaffold()
		if err != nil {
			utils.Console.Fatal(err)
		}

		variant, err = variants.CreateEmptyVariant(variantType)
		if err != nil {
			utils.Console.Fatal(err)
		}

		err = variant.Scaffold()
		if err != nil {
			utils.Console.Fatal(err)
		}

		err = pauseUntilInput()
		if err != nil {
			utils.Console.Fatal(err)
		}

		return
	}

	// Check if config exists
	if _, err := os.Stat("appy.yaml"); os.IsNotExist(err) {
		utils.Console.ErrorLn("appy.yaml not found. Either misconfigured or incorrect directory.")
		return
	}

	utils.Console.ClearEntireConsole()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: Terminal ui
	variant.Start(ctx)

	<-stop
	utils.Console.ClearLines(1)
	utils.Console.DebugLn("Received signal, shutting down...")
}

func promptVariantScaffold() (string, error) {
	prompt := promptui.Select{
		Label: "Select an option",
		Items: variants.GetVariantTypes(),
	}

	_, scaffoldType, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return scaffoldType, nil
}

func pauseUntilInput() error {
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

	return nil
}
