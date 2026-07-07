package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/weronikadusik/ferret/procfs"
)

const procRoot = "/proc"

func main() {
	var PIDs []int
	var processList []procfs.Process

	fmt.Println("Hi! I'm ferret 🦦")

	PIDs, err := procfs.ListPIDs(procRoot)
	if err != nil {
		log.Fatalf("could not read process list: %v", err)
	}

	for _, pid := range PIDs {
		process, err := procfs.ReadProcessStat(procRoot, pid)
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

	fmt.Printf("%d processes found\n\n", len(processList))

	cpuTimesStart, err := procfs.ReadStat(procRoot)
	if err != nil {
		log.Fatalf("could not parse CPU statistics: %v", err)
	}

	time.Sleep(time.Second)

	cpuTimesStop, err := procfs.ReadStat(procRoot)
	if err != nil {
		log.Fatalf("could not parse CPU statistics: %v", err)
	}

	deltas := CPUStatDelta(cpuTimesStart.Total, cpuTimesStop.Total)
	cpuUsage := CPUUsage(deltas)

	fmt.Printf("System CPU Usage: %.1f%%\n", cpuUsage)
}
