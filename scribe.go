package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	inputDirPtr := flag.String("input", "myfile", "Folder containing css modules")
	outputDoc := "public_html/documentation.html"
	flag.Parse()

	// Let's get this show on the road
	workingDirectory = *inputDirPtr

	outputFile, err := os.Create(outputDoc)
	checkErr(err)

	outputFile.WriteString(writeHTML())
	outputFile.Sync()

	// Now build the static pages
	buildStaticSite()

	fmt.Println("Docs generated at ./" + outputDoc)
}
