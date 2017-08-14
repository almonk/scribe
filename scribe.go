package main

import (
	"flag"
)

func main() {
	inputDirPtr := flag.String("input", "myfile", "Folder containing css modules")
	flag.Parse()

	// Let's get this show on the road
	workingDirectory = *inputDirPtr

	writeHTMLHeader()
	readDir()
	writeHTMLFooter()
}
