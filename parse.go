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
	var blocklist = []string{
		"_background-position",
		"_background-size",
		"_clears",
		"_debug",
		"_debug-grid",
		"_debug-children",
		"_display",
		"_floats",
		"_hk",
		"_hk-base",
		"_hk--compiled",
		"_hk-spinner",
		"_malibu",
	}

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

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		line := fileScanner.Text()

		if s.Contains(line, "@scribe") {
			cssMap = append(cssMap, line)
			noOfSection++
		} else {
			cssMap[noOfSection] = cssMap[noOfSection] + line
		}
	}

	for _, section := range cssMap {
		templateMatch := "<template>(.*?)</template>"
		moduleNameMatch := "@scribe \"(.*?)\""
		iterativeClassMatch := "{{class}}"
		commentMatch := "/*(.*?)*/"

		hasTemplate, _ := regexp.MatchString(templateMatch, section)

		if hasTemplate {
			m, _ := regexp.Compile(moduleNameMatch)
			extractedModuleName := m.FindStringSubmatch(section)
			fmt.Println(extractedModuleName[1])

			r, _ := regexp.Compile(templateMatch)
			extractedTemplate := r.FindStringSubmatch(section)
			// fmt.Println("Template:" + extractedTemplate[1])

			n, _ := regexp.Compile(commentMatch)
			cssToParse := n.ReplaceAllString(section, "")

			shouldLoop, _ := regexp.MatchString(iterativeClassMatch, section)
			fmt.Println(shouldLoop)

			ss := css.Parse(cssToParse)
			rules := ss.GetCSSRuleList()

			for _, rule := range rules {
				if isDocumentable(rule.Style.SelectorText) {
					fmt.Println("<br><span class='code'>" + rule.Style.SelectorText + "</span>")
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
