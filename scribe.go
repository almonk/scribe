package main

import (
	"flag"
	"fmt"
)

func main() {
	inputDirPtr := flag.String("input", "myfile", "Folder containing css modules")
	flag.Parse()

	fmt.Println("The scribe is sitting down to write...")

	// Let's get this show on the road
	workingDirectory = *inputDirPtr
	readDir()
}
