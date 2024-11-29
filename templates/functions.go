package templates

import (
	"strings"
	"text/template"
)

var templateFuncs = template.FuncMap{
	"fnSanitizeVar":    sanitizeVarName,
	"fnCommaSeperated": commaSeperated,
	"fnUpackMapValues": upackMapValues,
}

func sanitizeVarName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ToLower(name)
	return name
}

func commaSeperated(list []string) string {
	for i, item := range list {
		list[i] = sanitizeVarName(item)
	}
	return strings.Join(list, ", ")
}

func upackMapValues(params map[string]string) string {
	var values []string
	for _, v := range params {
		values = append(values, v)
	}
	return strings.Join(values, ", ")
}
