package utils

import (
	"html/template"
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

func titleString(s string) string {
	return cases.Title(language.English).String(strings.ToLower(s))
}
