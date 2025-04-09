package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new component to your project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var addDomainCmd = &cobra.Command{
	Use:   "domain",
	Short: "Add a new domain to your project",
	Run: func(cmd *cobra.Command, args []string) {
		addDomain(cmd, args)
	},
}

func init() {
	addCmd.AddCommand(addDomainCmd)
	rootCmd.AddCommand(addCmd)
}

func addDomain(cmd *cobra.Command, args []string) {
	// Get domain name from user
	fmt.Print("Name: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	domain := scanner.Text()

	// Create domain files
	prefix := "domains/" + domain
	files := []scaffoldFileEntry{
		{
			path:     prefix + "/model/model.go",
			template: templates.DomainExampleModel,
			tools:    []string{shared.ToolGoFmt},
		},
		{
			path:     prefix + "/usecase/usecase.go",
			template: templates.DomainExampleUsecase,
			tools:    []string{shared.ToolGoFmt},
		},
		{
			path:     prefix + "/ports/in/in.go",
			template: templates.DomainExampleInPort,
			tools:    []string{shared.ToolGoFmt},
		},
		{
			path:     prefix + "/ports/out/out.go",
			template: templates.DomainExampleOutPort,
			tools:    []string{shared.ToolGoFmt},
		},
		{
			path:     prefix + "/adapter/in/in.go",
			template: templates.DomainExampleInAdapter,
			tools:    []string{shared.ToolGoFmt},
		},
		{
			path:     prefix + "/adapter/out/out.go",
			template: templates.DomainExampleOutAdapter,
			tools:    []string{shared.ToolGoFmt},
		},
		{
			path:     prefix + "/wiring.go",
			template: templates.DomainExampleWiring,
		},
	}

	err := generateFiles(files, map[string]any{
		"DomainName": domain,
	})
	if err != nil {
		panic(err)
	}
}
