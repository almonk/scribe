package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	s "strings"

	"github.com/almonk/css"
)

var workingDirectory = ""

func checkErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

func readDir() {
	files, _ := ioutil.ReadDir(workingDirectory)
	for _, f := range files {
		if isValidModule(f.Name()) {
			parseModule(f)
		}
	}
}

func isValidModule(filename string) bool {
	// Check the file has the hallmarks of a css module

	if isOnBlocklist(filename) {
		return false
	}

	if s.HasSuffix(filename, ".css") {
		return true
	}

	return false
}

func isOnBlocklist(filenameToCheck string) bool {
	// Css modules we don't want to document
	var blocklist = []string{}

	for _, term := range blocklist {
		if s.Contains(filenameToCheck, term) {
			return true
		}
	}

	return false
}

func parseModule(filename os.FileInfo) {
	file, err := os.Open(workingDirectory + filename.Name())
	checkErr(err)
	defer file.Close()

	cssMap := make([]string, 1)
	noOfSection := 0
	isInScribeSection := false

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if s.Contains(line, "/*") {
			cssMap = append(cssMap, line)
			isInScribeSection = true
			noOfSection++
		}

		if isInScribeSection {
			cssMap[noOfSection] = cssMap[noOfSection] + line
		}

		if s.Contains(line, "@scribe nodoc") {
			isInScribeSection = false
		}
	}

	for _, section := range cssMap {
		templateMatch := "<template>(.*?)</template>"
		moduleNameMatch := "@scribe(.*)<template>"
		commentMatch := "/*(.*?)*/"

		hasTemplate, _ := regexp.MatchString(templateMatch, section)

		if hasTemplate {
			m, _ := regexp.Compile(moduleNameMatch)
			extractedModuleName := m.FindStringSubmatch(section)
			fmt.Println("<br><div class='ma4 f2 dark-gray'>" + extractedModuleName[1] + "</div>")

			r, _ := regexp.Compile(templateMatch)
			extractedTemplate := r.FindStringSubmatch(section)

			n, _ := regexp.Compile(commentMatch)
			cssToParse := n.ReplaceAllString(section, "")

			ss := css.Parse(cssToParse)
			// fmt.Println("\nCSS:" + cssToParse)
			rules := ss.GetCSSRuleList()

			for _, rule := range rules {
				if isDocumentable(rule.Style.SelectorText) {
					fmt.Println("<br><span class='code ma4'>" + rule.Style.SelectorText + "</span>")
					fmt.Println(documentClass(rule.Style.SelectorText, rule.Style.Styles, extractedTemplate[1]))
				}
			}
		}
	}
}

func readModule(file string, folder string) string {
	dat, err := ioutil.ReadFile(folder + "/" + file)
	checkErr(err)
	fileBuffer := string(dat)
	return fileBuffer
}
