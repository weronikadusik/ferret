// ferret v. 0.0.2
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var procRoot = "/proc"

type Process struct {
	PID  int
	Comm string
}

// read /proc/[pid]/stat to obtain comm (the filename of the executable)
func getProcessInfo(pid int) (Process, error) {
	statPath := filepath.Join(procRoot, strconv.Itoa(pid), "stat")
	data, err := os.ReadFile(statPath)
	if err != nil {
		return Process{}, err
	}

	statStr := string(data)

	lParen := strings.IndexByte(statStr, '(')
	rParen := strings.LastIndexByte(statStr, ')')

	if lParen == -1 || rParen == -1 {
		return Process{}, errors.New("incorrect stat format")
	}

	return Process{
		PID:  pid,
		Comm: statStr[lParen+1 : rParen],
	}, nil
}

func main() {
	var processList []Process

	fmt.Println("Hi! I'm ferret 🦦")

	// read proc directory, to obtain a list of processes
	processDir, err := os.ReadDir(procRoot)
	if err != nil {
		log.Fatalf("could not read process list: %v", err)
	}

	// filter out non-numeric directory names, so only process-specific ones remain
	for _, f := range processDir {
		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}

		process, err := getProcessInfo(pid)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				continue // process exited between reading proc directory and reading process-specific stat, or access denied
			}
			log.Fatalf("could not read process stat: %v", err)
		}

		fmt.Printf("Process %d: %s\n", process.PID, process.Comm)
		processList = append(processList, process)
	}

	fmt.Printf("%d processes found\n", len(processList))
}
