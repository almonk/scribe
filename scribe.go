package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Get current working path
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Working from " + pwd + "...")

	inputDirPtr := flag.String("input", "myfile", "Folder containing css modules")
	outputDoc := pwd + "/public_html/documentation.html"
	glossaryFilePtr := flag.String("glossary", "", "File to compile glossary from")
	flag.Parse()

	// Let's get this show on the road
	workingDirectory = *inputDirPtr

	outputFile, err := os.Create(outputDoc)
	checkErr(err)

	outputFile.WriteString(writeHTML())
	outputFile.Sync()

	// Now build the table of styles
	if *glossaryFilePtr != "" {
		buildToS(*glossaryFilePtr)
		fmt.Println("Built glossary file")
	}

	// Now build the static pages
	buildStaticSite()

	fmt.Println("Docs generated at ./" + outputDoc)
}
