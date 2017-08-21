package main

import (
	"bytes"
	"html/template"
	s "strings"
)

type contents struct {
	DocsContents template.HTML
	TocContents  template.HTML
}

func documentClass(class string, template string) string {
	template = s.TrimSpace(template)
	class = s.Trim(class, ".")
	outputHTML := s.Replace(template, "{{class}}", class, -1)
	return outputHTML
}

func writeHTML() string {
	data := contents{
		DocsContents: template.HTML(readDir()),
		TocContents:  template.HTML(buildToc()),
	}

	partial := readModule("layout.html", "docs-templates")
	tmpl, err := template.New("").Parse(partial)
	checkErr(err)

	var tpl bytes.Buffer
	tmpl.Execute(&tpl, data)
	return tpl.String()
}
