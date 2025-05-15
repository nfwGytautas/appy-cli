package utils

import (
	"os"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func NewTemplate(content string) *template.Template {
	name := GenerateRandomString(10)
	tmpl := template.Must(template.New(name).Funcs(template.FuncMap{
		"TitleString":        titleString,
		"HyphenToUnderscore": hyphenToUnderscore,
	}).Parse(content))
	return tmpl
}

func TemplateAString(content string, data any) (string, error) {
	// Create template
	tmpl := NewTemplate(content)

	// Execute template
	var result strings.Builder
	err := tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func TemplateAFile(pathIn string, pathOut string, data any) error {
	// Copy file
	templateData, err := os.ReadFile(pathIn)
	if err != nil {
		return err
	}

	file, err := os.Create(pathOut)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write template
	tmpl := NewTemplate(string(templateData))

	err = tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}

func TemplateAStringToFile(pathOut string, content string, data any) error {
	var file *os.File

	// Create file if not exists
	if _, err := os.Stat(pathOut); os.IsNotExist(err) {
		file, err = os.Create(pathOut)
		if err != nil {
			return err
		}
	} else {
		// Open file
		file, err = os.OpenFile(pathOut, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	}
	defer file.Close()

	// Write template
	tmpl := NewTemplate(content)

	err := tmpl.Execute(file, data)
	if err != nil {
		return err
	}

	return nil
}

func titleString(s string) string {
	return cases.Title(language.English).String(strings.ToLower(s))
}

func hyphenToUnderscore(s string) string {
	return strings.ReplaceAll(s, "-", "_")
}
