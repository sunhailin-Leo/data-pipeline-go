package sink

import (
	"bytes"
	_ "embed"
	"errors"
	"strings"
	"text/template"

	"github.com/sunhailin-Leo/data-pipeline-go/pkg/config"
)

//go:embed create_table.tmpl
var createTableTmpl string

// RenderClickhouseCreateTableTemplate render clickhouse create table template
func RenderClickhouseCreateTableTemplate(templateName string, data config.ClickhouseSinkConfig) (string, error) {
	var result bytes.Buffer
	funcMap := template.FuncMap{
		"sub":      func(a, b int) int { return a - b },
		"joinWith": strings.Join,
	}

	tmpl, err := template.New(templateName).Funcs(funcMap).Parse(createTableTmpl)
	if err != nil {
		return "", errors.New("Error parsing template: " + err.Error())
	}
	executeErr := tmpl.Execute(&result, data)
	if executeErr != nil {
		return "", errors.New("Error Execute template: " + executeErr.Error())
	}
	return result.String(), nil
}
