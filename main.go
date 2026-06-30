package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

const procRoot = "/proc"

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

		process, err := getProcessInfo(procRoot, pid)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				continue // process exited between reading proc directory and reading process-specific stat, or access denied
			}
			log.Fatalf("could not read process stat: %v", err)
		}

		fmt.Printf("Process %d: %s:\n", process.PID, process.Comm)
		fmt.Printf("\t├─ State:%q  Priority:%d  Nice:%d\n", process.State, process.Priority, process.Nice)
		fmt.Printf("\t├─ Virtual memory size (B):%d  Resident Memory size (B):%d\n", process.VSZBytes, process.RSSBytes)
		fmt.Printf("\t└─ User mode time (s):%.2f  Kernel mode time (s):%.2f\n", ticksToSeconds(float64(process.UTimeTicks)), ticksToSeconds(float64(process.STimeTicks)))
		processList = append(processList, process)
	}

	fmt.Printf("%d processes found\n", len(processList))
}
