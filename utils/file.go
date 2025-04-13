package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GeneratedFileTree struct {
	files  []generateFileEntry
	prefix string
}

type generateFileEntry struct {
	path     string
	template string
	tools    []string
}

func (g *GeneratedFileTree) Generate(data any) error {
	for _, file := range g.files {
		err := g.generateFile(&file, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *GeneratedFileTree) AddFile(path string, template string, tools []string) {
	g.files = append(g.files, generateFileEntry{
		path:     path,
		template: template,
		tools:    tools,
	})
}

func (g *GeneratedFileTree) AddDirectory(path string) {
	g.files = append(g.files, generateFileEntry{
		path: path,
	})
}

func (g *GeneratedFileTree) SetPrefix(prefix string) {
	g.prefix = prefix
}

func (g *GeneratedFileTree) generateFile(file *generateFileEntry, data any) error {
	var err error

	path := filepath.Join(g.prefix, file.path)

	// Directory
	if strings.HasSuffix(file.path, "/") {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}

		return nil
	}

	// File

	// Create parent directory if not exists
	dir := filepath.Dir(path)
	if dir != "." {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	// Create file
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	// Write template
	tmpl := template.Must(template.New(path).Parse(file.template))

	err = tmpl.Execute(f, data)
	if err != nil {
		f.Close()
		return err
	}

	// Run tools
	f.Close()
	for _, tool := range file.tools {
		toolArgs := strings.Split(tool, " ")
		toolArgs = append(toolArgs, path)
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

func CalculateFileHash(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
