package templates

import (
	"bytes"
	"text/template"
)

type TemplateParams map[string]any

func GenTemplate(templateString string, params TemplateParams) (string, error) {
	var tpl bytes.Buffer

	tmpl, err := template.New("template").Parse(templateString)
	if err != nil {
		return "", err
	}

	if err := tmpl.Execute(&tpl, params); err != nil {
		return "", err
	}

	return tpl.String(), nil
}
