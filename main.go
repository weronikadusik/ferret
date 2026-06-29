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
	PID      int
	Comm     string
	State    byte
	Priority int
	Nice     int

	vszBytes uint64
	rssBytes uint64
}

// read /proc/[pid]/stat; return filled Process struct
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

	processStats := strings.Fields(statStr[rParen+1:])

	processPriority, err := strconv.Atoi(processStats[15])
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processNice, err := strconv.Atoi(processStats[16])
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processVSZ, err := strconv.Atoi(processStats[20])
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processRSS, err := strconv.Atoi(processStats[21])
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	return Process{
		PID:      pid,
		Comm:     statStr[lParen+1 : rParen],
		State:    processStats[0][0],
		Priority: processPriority,
		vszBytes: uint64(processVSZ),
		rssBytes: uint64(processRSS) * uint64(os.Getpagesize()),
		Nice:     processNice,
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

		fmt.Printf("Process %d: %s;\t\t State:%q Priority:%d Nice:%d; \t\t Virtual memory size (B):%d Resident Memory size (B):%d \n", process.PID, process.Comm, process.State, process.Priority, process.Nice, process.vszBytes, process.rssBytes)
		processList = append(processList, process)
	}

	fmt.Printf("%d processes found\n", len(processList))
}
