package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	s "strings"

	"github.com/gosimple/slug"
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
			output = output + parseModule(f, false)
		}
	}
	return output
}

func buildToc() string {
	output := ""
	files, _ := ioutil.ReadDir(workingDirectory)
	for _, f := range files {
		if isValidModule(f.Name()) {
			output = output + parseModule(f, true)
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

func parseModule(filename os.FileInfo, justHeaders bool) string {
	outputString := ""

	file, err := os.Open(workingDirectory + filename.Name())
	checkErr(err)
	defer file.Close()

	cssMap := make([]string, 1)
	noOfSection := 0
	noOfLines := 0
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

		noOfLines++
	}

	for i, section := range cssMap {
		templateMatch := `<template>\s(.*?)\s</template>`
		sectionNameMatch := `@scribe (.*?)\n`
		// commentMatch := `\/\*(.*?)@scribe(.*?)\*\/`
		cssSelectorMatch := ".(.*?) {.*?}"

		hasTemplate, _ := regexp.MatchString(templateMatch, section)
		hasSubSection, _ := regexp.MatchString(sectionNameMatch, section)

		matchSectionName, _ := regexp.Compile(sectionNameMatch)
		extractedSectionName := matchSectionName.FindStringSubmatch(section)

		matchTemplate, _ := regexp.Compile(templateMatch)
		extractedTemplate := matchTemplate.FindStringSubmatch(section)

		matchCSSSelector, _ := regexp.Compile(cssSelectorMatch)
		extractedCSSSelectors := matchCSSSelector.FindAllString(section, -1)

		if hasTemplate {
			if i == 1 {
				fmt.Println(humanizeModuleName(*file))
			}

			if hasSubSection {
				// Has other scribe sections in the module
				fmt.Println(extractedSectionName[1])
			}

			fmt.Println(extractedTemplate[1])

			for index := range extractedCSSSelectors {
				// Loop thru every CSS selector
				if isValidCSSClass(extractedCSSSelectors[index]) {
					cssClass := cssSelectorFromDefinition(extractedCSSSelectors[index])

					fmt.Println(cssClass)
				}
			}
		}
	}

	return outputString
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
