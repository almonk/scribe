package main

import (
	"bufio"
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
			output = output + parseModule(f)
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

func parseModule(filename os.FileInfo) string {
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
			cssMap = append(cssMap, "/* "+line)
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
		moduleNameMatch := "@scribe (.*?)<template>"
		commentMatch := `\/\*(.*?)@scribe(.*?)\*\/`
		cssSelectorMatch := ".(.*?) {.*?}"

		hasTemplate, _ := regexp.MatchString(templateMatch, section)

		if hasTemplate {
			m, _ := regexp.Compile(moduleNameMatch)
			extractedModuleName := m.FindStringSubmatch(section)
			outputString = outputString + "<br><div class='ma4 f2 dark-gray'>" + extractedModuleName[1] + "</div>"

			r, _ := regexp.Compile(templateMatch)
			extractedTemplate := r.FindStringSubmatch(section)

			n, _ := regexp.Compile(commentMatch)
			cssToParse := n.ReplaceAllString(section, "")

			o, _ := regexp.Compile(cssSelectorMatch)
			cssSelectors := o.FindAllString(cssToParse, -1)

			for index := range cssSelectors {
				if s.HasPrefix(cssSelectors[index], ".") {
					class := s.Split(cssSelectors[index], "{")

					outputString = outputString + "<br><span class='code ma4'>" + class[0] + "</span><br>"
					outputString = outputString + documentClass(class[0], extractedTemplate[1])
				}
			}
		}
	}

	return outputString
}

func readModule(file string, folder string) string {
	dat, err := ioutil.ReadFile(folder + "/" + file)
	checkErr(err)
	fileBuffer := string(dat)
	return fileBuffer
}
