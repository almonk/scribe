package main

import (
	"bufio"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	s "strings"
)

var workingDirectory = ""

func checkErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}

func readDir() string {
	output := ""
	files, _ := ioutil.ReadDir(workingDirectory)
	for _, f := range files {
		if isValidModule(f.Name()) {
			output = output + parseModuleForDocs(f)
		}
	}
	return output
}

func buildToc() string {
	output := ""
	files, _ := ioutil.ReadDir(workingDirectory)
	for _, f := range files {
		if isValidModule(f.Name()) {
			output = output + parseModuleForDocs(f)
		}
	}
	return output
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

func isDocumentable(class string) bool {
	var redflags = []string{}

	for _, term := range redflags {
		if s.Contains(class, term) {
			return false
		}
	}

	return true
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

func parseModuleForDocs(filename os.FileInfo) string {
	outputString := ""

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
			isInScribeSection = true
			cssMap = append(cssMap, "/* "+line+"\n")
			noOfSection++
		}

		if isInScribeSection {
			cssMap[noOfSection] = cssMap[noOfSection] + line + "\n"
		}

		if s.Contains(line, "@scribe nodoc") {
			isInScribeSection = false
		}
	}

	for i, section := range cssMap {
		sectionNameMatch := `@scribe (.*?)\n`
		cssSelectorMatch := ".(.*?) {"

		hasTemplate := s.Contains(section, "<template>")
		hasSubSection, _ := regexp.MatchString(sectionNameMatch, section)

		matchSectionName, _ := regexp.Compile(sectionNameMatch)
		extractedSectionName := matchSectionName.FindStringSubmatch(section)

		matchCSSSelector, _ := regexp.Compile(cssSelectorMatch)
		extractedCSSSelectors := matchCSSSelector.FindAllString(section, -1)

		if hasTemplate {
			if i == 1 {
				fmt.Println("Found documentation module " + humanizeModuleName(*file))
				outputString = outputString + heading(humanizeModuleName(*file), slugifyModuleName(*file))
			}

			if hasSubSection {
				// Has other scribe sections in the module
				outputString = outputString + subheading(extractedSectionName[1])
			}

			for index := range extractedCSSSelectors {
				// Loop thru every CSS selector
				if isValidCSSClass(extractedCSSSelectors[index]) {
					cssClass := cssSelectorFromDefinition(extractedCSSSelectors[index])
					template := getInnerSubstring(section, "<template>", "</template>")

					outputString = outputString + documentClass(cssClass)
					outputString = outputString +
						"<div class='flex'><div class='w-50'>" +
						templateForClass(cssClass, template) +
						"</div><div class='w-50'><pre class='h3 ma0 overflow-auto gray'>" +
						html.EscapeString(templateForClass(cssClass, template)) +
						"</pre></div></div>"
				}
			}
		}
	}

	return outputString
}

func makeTOC() string {
	output := ""
	files, _ := ioutil.ReadDir(workingDirectory)
	for _, f := range files {
		if isValidModule(f.Name()) {
			file, err := os.Open(workingDirectory + f.Name())
			checkErr(err)
			defer file.Close()

			cssMap := make([]string, 1)
			noOfSection := 0
			isInScribeSection := false

			fileScanner := bufio.NewScanner(file)

			for fileScanner.Scan() {
				line := fileScanner.Text()

				if s.Contains(line, "/*") {
					isInScribeSection = true
					cssMap = append(cssMap, "/* "+line+"\n")
					noOfSection++
				}

				if isInScribeSection {
					cssMap[noOfSection] = cssMap[noOfSection] + line + "\n"
				}

				if s.Contains(line, "@scribe nodoc") {
					isInScribeSection = false
				}
			}

			for i, section := range cssMap {
				hasTemplate := s.Contains(section, "<template>")

				if hasTemplate {
					if i == 1 {
						output = output + tocItem(humanizeModuleName(*file), slugifyModuleName(*file))
					}
				}
			}
		}
	}

	return output
}

func isValidCSSClass(class string) bool {
	if s.HasPrefix(class, ".") {
		return true
	}
	return false
}

func cssSelectorFromDefinition(rule string) string {
	class := s.Split(rule, " {")
	return class[0]
}

func readModule(file string, folder string) string {
	dat, err := ioutil.ReadFile(folder + "/" + file)
	checkErr(err)
	fileBuffer := string(dat)
	return fileBuffer
}

func getInnerSubstring(str string, prefix string, suffix string) string {
	var beginIndex, endIndex int
	beginIndex = s.Index(str, prefix)
	if beginIndex == -1 {
		beginIndex = 0
		endIndex = 0
	} else if len(prefix) == 0 {
		beginIndex = 0
		endIndex = s.Index(str, suffix)
		if endIndex == -1 || len(suffix) == 0 {
			endIndex = len(str)
		}
	} else {
		beginIndex += len(prefix)
		endIndex = s.Index(str[beginIndex:], suffix)
		if endIndex == -1 {
			if s.Index(str, suffix) < beginIndex {
				endIndex = beginIndex
			} else {
				endIndex = len(str)
			}
		} else {
			if len(suffix) == 0 {
				endIndex = len(str)
			} else {
				endIndex += beginIndex
			}
		}
	}

	return str[beginIndex:endIndex]
}
