package main

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
	s "strings"

	"github.com/gosimple/slug"
)

type contents struct {
	DocsContents template.HTML
	TocContents  template.HTML
}

func templateForClass(class string, template string) string {
	template = s.TrimSpace(template)
	class = s.Trim(class, ".")
	outputHTML := s.Replace(template, "{{class}}", class, -1)
	return outputHTML
}

func documentClass(class string) string {
	output := s.Trim(class, ".")
	return "<pre class='bg-light-silver pa2 lh-copy'>" + output + "</pre>"
}

func heading(text string, slug string) string {
	return "<div class='mv5 bb b--light-gray' id='" + slug + "'></div><div class='mv4 f3 dark-gray'>" + text + "</div>"
}

func subheading(text string) string {
	return "<h4 class='dark-gray ttu f5 mt6 db'>" + text + "</h4>"
}

func tocItem(text string, slug string) string {
	return "<li><a class='db ph2 pv1 dark-gray link' href='#" + slug + "'>" + text + "</a></li>"
}

func writeHTML() string {
	data := contents{
		DocsContents: template.HTML(readDir()),
		TocContents:  template.HTML(makeTOC()),
	}

	partial := readModule("layout.html", "docs-templates")
	tmpl, err := template.New("").Parse(partial)
	checkErr(err)

	var tpl bytes.Buffer
	tmpl.Execute(&tpl, data)
	return tpl.String()
}

func humanizeModuleName(file os.File) string {
	cleanFileName := filepath.Base(file.Name())
	replacer := s.NewReplacer("_", "", ".css", "", "-", " ")
	output := replacer.Replace(cleanFileName)
	output = s.Title(output)
	return output
}

func slugifyModuleName(file os.File) string {
	return slug.Make(humanizeModuleName(file))
}
