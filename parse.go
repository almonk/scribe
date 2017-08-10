package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

	if s.HasPrefix(filename, "_") && s.HasSuffix(filename, ".css") {
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
	fmt.Println(sanitizeModuleToHumanName(filename.Name()))

	fileBuffer := readModule(filename.Name(), workingDirectory)
	ss := css.Parse(fileBuffer)
	rules := ss.GetCSSRuleList()

	for _, rule := range rules {
		if isDocumentable(rule.Style.SelectorText) {
			documentClass(rule.Style.SelectorText, rule.Style.Styles)
		}
	}
}

func readModule(file string, folder string) string {
	dat, err := ioutil.ReadFile(folder + "/" + file)
	checkErr(err)
	fileBuffer := string(dat)
	return fileBuffer
}
