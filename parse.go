package main

import (
	"bufio"
	"html"
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

			if noOfLines > 1 {
				noOfSection++
			}
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
		moduleNameMatch := `@scribe(.*?)\n`
		commentMatch := `\/\*(.*?)@scribe(.*?)\*\/`
		cssSelectorMatch := ".(.*?) {.*?}"

		hasTemplate, _ := regexp.MatchString(templateMatch, section)

		if hasTemplate {
			if justHeaders {
				outputString = outputString + "<li><a href='#" + slugifyModuleName(*file) + "' class='link dark-gray pa2 db hover-bg-light-silver'>" + humanizeModuleName(*file) + "</a></li>"
			} else {
				outputString = outputString + "<div class='mv5 bb b--light-gray' id='" + slugifyModuleName(*file) + "'></div><div class='mv4 f3 dark-gray'>" + humanizeModuleName(*file) + "</div>"
			}

			m, _ := regexp.Compile(moduleNameMatch)
			extractedModuleName := m.FindStringSubmatch(section)

			if justHeaders && len(extractedModuleName[1]) > 0 {
				outputString = outputString + "<li><a href='#' class='link dark-gray pa2 db hover-bg-light-silver'>" + extractedModuleName[1] + "</a></li>"
			}

			if !justHeaders && i > 0 {
				outputString = outputString + "<h4 class='gray ttu f6 mt2 dib'>" + extractedModuleName[1] + "</h4>"
			}

			r, _ := regexp.Compile(templateMatch)
			extractedTemplate := r.FindStringSubmatch(section)

			n, _ := regexp.Compile(commentMatch)
			cssToParse := n.ReplaceAllString(section, "")

			o, _ := regexp.Compile(cssSelectorMatch)
			cssSelectors := o.FindAllString(cssToParse, -1)

			if !justHeaders {
				for index := range cssSelectors {
					if s.HasPrefix(cssSelectors[index], ".") {
						class := s.Split(cssSelectors[index], " {")

						outputString = outputString + "<pre class='bg-light-silver pa2 mt2'>" + html.EscapeString(class[0]) + "</pre>"
						outputString = outputString + documentClass(class[0], extractedTemplate[1])
					}
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

func readModule(file string, folder string) string {
	dat, err := ioutil.ReadFile(folder + "/" + file)
	checkErr(err)
	fileBuffer := string(dat)
	return fileBuffer
}
