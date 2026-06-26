// ferret v. 0.0.1
package main

import (
	"fmt"
	"log"
	"os"
	"unicode"
)

func main() {
	var procRoot = "/proc"
	var processList []os.DirEntry

	fmt.Println("Hi! I'm ferret 🦦")

	// read proc directory, to obtain a list of processes
	processDir, err := os.ReadDir(procRoot)
	if err != nil {
		log.Fatal("could not read process list.")
	}

	// filter out non-numeric directory names, so only process-specific ones remain
	for _, f := range processDir {
		isProcess := true

		for _, c := range f.Name() {
			if !unicode.IsDigit(c) {
				isProcess = false
				break
			}
		}

		if isProcess {
			processList = append(processList, f)
		}
	}

	fmt.Printf("%d processes found\n", len(processList))
}
