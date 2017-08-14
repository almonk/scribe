package main

import (
	"flag"
	"os"
)

func main() {
	inputDirPtr := flag.String("input", "myfile", "Folder containing css modules")
	flag.Parse()

	// Let's get this show on the road
	workingDirectory = *inputDirPtr

	outputFile, err := os.Create("test.html")
	checkErr(err)

	outputFile.WriteString(writeHTMLHeader())
	outputFile.Sync()

	outputFile.WriteString(readDir())
	outputFile.Sync()

	outputFile.WriteString(writeHTMLFooter())
	outputFile.Sync()
}
