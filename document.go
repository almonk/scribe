package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	s "strings"

	"github.com/almonk/css"
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
	return "<pre class='lh-copy'>" + output + "</pre>"
}

func heading(text string, slug string) string {
	return "<div class='mv5 bb b--light-gray' id='" + slug + "'></div><div class='purple lh-copy measure-wide f3 fw4'>" + text + "</div>"
}

func subheading(text string) string {
	return "<h4 class='dark-gray ttu f5 mt6 db'>" + text + "</h4>"
}

func tocItem(text string, slug string) string {
	return "<li><a class='scroll db ph2 pv1 dark-gray link lh-copy hk-focus-ring hover-bg-silver' href='#" + slug + "'>" + text + "</a></li>"
}

func writeHTML() string {
	data := contents{
		DocsContents: template.HTML(readDir()),
		TocContents:  template.HTML(makeTOC()),
	}

	partial := readModule("docs-layout.html", "scribe/templates")
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

func buildStaticSite() {
	files, _ := ioutil.ReadDir("./scribe/pages/")
	for _, f := range files {
		dat, err := ioutil.ReadFile("./scribe/pages/" + f.Name())
		outputFile, err := os.Create("./public_html/" + filepath.Base(f.Name()))
		fileBuffer := string(dat)

		checkErr(err)

		compiledPage := wrapStaticPage(fileBuffer)
		outputFile.WriteString(compiledPage)
		outputFile.Sync()
	}
}

func wrapStaticPage(pageHTML string) string {
	data := contents{
		DocsContents: template.HTML(pageHTML),
	}

	partial := readModule("layout.html", "scribe/templates")
	tmpl, err := template.New("").Parse(partial)
	checkErr(err)

	var tpl bytes.Buffer
	tmpl.Execute(&tpl, data)
	return tpl.String()
}

func buildToS(glossaryFile string) {
	distFile := readFile(glossaryFile)
	outputFile, err := os.Create("./scribe/pages/glossary.html")
	checkErr(err)

	ss := css.Parse(distFile)
	rules := ss.GetCSSRuleList()

	outputFile.WriteString(`
	<div class="ml2">
		<h4 class="purple lh-copy measure-wide f3 fw4">Glossary of styles</h4>
		<div class="dark-gray lh-copy measure-wide">
			<p>A list of all classes and their properties</p>
		</div>
	</div>
	`)

	outputFile.WriteString("<table class='w-100 f4 lh-copy'><tbody>")

	for _, rule := range rules {
		outputFile.WriteString("<tr><td class='ph2 w-50 bb b--silver v-top'><pre class='measure dark-gray'>" + rule.Style.SelectorText + "</pre></td><td class='ph2 w-50 bb b--silver v-top'><pre class='measure truncate'>")

		for _, style := range rule.Style.Styles {
			outputFile.WriteString("<span class='purple'>" + style.Property + "</span>: <span class='blue'>" + style.Value + "</span><br/>")
		}

		outputFile.WriteString("</pre></div></td></tr></tbody>")
	}

	outputFile.WriteString("</table>")

	outputFile.Sync()
}
