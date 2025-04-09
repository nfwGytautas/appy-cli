package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nfwGytautas/appy-cli/shared"
	"github.com/nfwGytautas/appy-cli/templates"
	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold the base structure of your app",
	Run: func(cmd *cobra.Command, args []string) {
		base(cmd, args)
	},
}

type scaffoldFileEntry struct {
	path     string
	template string
	tools    []string
}

type scaffoldInput struct {
	ProjectName string
}

func init() {
	rootCmd.AddCommand(scaffoldCmd)
}

func base(cmd *cobra.Command, args []string) {
	in := scaffoldInput{}

	fmt.Println("Scaffolding...")

	files := []scaffoldFileEntry{
		{
			path:     "main.go",
			template: templates.MainGo,
			tools:    []string{shared.ToolGoFmt},
		},
		{
			path:     "go.mod",
			template: templates.GoMod,
		},
		{
			path:     "README.md",
			template: templates.ReadmeMd,
		},
		{
			path:     "Dockerfile",
			template: templates.Dockerfile,
		},
		{
			path:     ".gitignore",
			template: templates.Gitignore,
		},
		{
			path:     ".vscode/Snippets.code-snippets",
			template: templates.VscodeSnippets,
		},
		{
			path:     ".github/build.yaml",
			template: templates.GithubBuildYaml,
		},
		{
			path: "repositories/",
		},
		{
			path: "hooks/",
		},
		{
			path:     "shared/errors.go",
			template: templates.ErrorsGo,
			tools:    []string{"gofmt -w"},
		},
		{
			path: "domains/",
		},
	}

	err := gatherInput(&in)
	if err != nil {
		panic(err)
	}

	err = generateFiles(files, in)
	if err != nil {
		panic(err)
	}

	fmt.Println("Done!")
}

func gatherInput(in *scaffoldInput) error {
	var err error

	// Get project name from parent directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	in.ProjectName = filepath.Base(dir)

	return err
}

func generateFiles(files []scaffoldFileEntry, in any) error {
	var err error

	for _, file := range files {
		err := generateFile(&file, in)
		if err != nil {
			return err
		}
	}

	return err
}

func generateFile(file *scaffoldFileEntry, in any) error {
	var err error

	// Directory
	if strings.HasSuffix(file.path, "/") {
		err = os.MkdirAll(file.path, 0755)
		if err != nil {
			return err
		}

		return nil
	}

	// File

	// Create parent directory if not exists
	dir := filepath.Dir(file.path)
	if dir != "." {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// Create file
	f, err := os.Create(file.path)
	if err != nil {
		return err
	}

	// Write template
	tmpl := template.Must(template.New(file.path).Parse(file.template))

	err = tmpl.Execute(f, in)
	if err != nil {
		f.Close()
		return err
	}

	// Run tools
	f.Close()
	for _, tool := range file.tools {
		toolArgs := strings.Split(tool, " ")
		toolArgs = append(toolArgs, file.path)
		cmd := exec.Command(toolArgs[0], toolArgs[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return err
}
