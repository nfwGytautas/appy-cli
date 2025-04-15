package utils

import (
	"html/template"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func NewTemplate(content string) *template.Template {
	name := GenerateRandomString(10)
	tmpl := template.Must(template.New(name).Funcs(template.FuncMap{
		"TitleString": titleString,
	}).Parse(content))
	return tmpl
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

func titleString(s string) string {
	return cases.Title(language.English).String(strings.ToLower(s))
}
