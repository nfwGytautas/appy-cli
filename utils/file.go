package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
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
		path:     filepath.Join(g.prefix, path),
		template: template,
		tools:    tools,
	})
}

func (g *GeneratedFileTree) AddDirectory(path string) {
	g.files = append(g.files, generateFileEntry{
		path: filepath.Join(g.prefix, path) + "/",
	})
}

func (g *GeneratedFileTree) SetPrefix(prefix string) {
	g.prefix = prefix
}

func (g *GeneratedFileTree) generateFile(file *generateFileEntry, data any) error {
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
	tmpl := NewTemplate(file.template)

	err = tmpl.Execute(f, data)
	if err != nil {
		f.Close()
		return err
	}

	// Run tools
	f.Close()
	err = RunTools(file.path, file.tools)
	if err != nil {
		return err
	}

	return nil
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
