package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	inputDirPtr := flag.String("input", "myfile", "Folder containing css modules")
	outputDoc := "test.html"
	flag.Parse()

	// Let's get this show on the road
	workingDirectory = *inputDirPtr

	outputFile, err := os.Create(outputDoc)
	checkErr(err)

	outputFile.WriteString(writeHTMLHeader())
	outputFile.Sync()

	outputFile.WriteString(readDir())
	outputFile.Sync()

	outputFile.WriteString(writeHTMLFooter())
	outputFile.Sync()

	fmt.Println("Docs generated at ./" + outputDoc)
}
