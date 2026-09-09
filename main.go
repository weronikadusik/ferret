package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/weronikadusik/ferret/procfs"
)

const procRoot = "/proc"

type Snapshot struct {
	Processes map[int]procfs.Process
	CPUStats  procfs.CPUStats
}

type CPUUsage struct {
	Total  float64
	PerCPU []float64
}

func readSnapshot() (Snapshot, error) {
	Processes, err := readProcesses()
	if err != nil {
		return Snapshot{}, err
	}

	CPUStats, err := procfs.ReadStat(procRoot)
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		Processes: Processes,
		CPUStats:  CPUStats,
	}, nil
}

func readProcesses() (map[int]procfs.Process, error) {
	pids, err := procfs.ListPIDs(procRoot)
	if err != nil {
		return nil, fmt.Errorf("listing PIDs: %w", err)
	}

	processes := make(map[int]procfs.Process, len(pids))
	for _, pid := range pids {
		proc, err := procfs.ReadProcessStat(procRoot, pid)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) {
				continue // process exited between reading proc directory and reading process-specific stat, or access denied
			}
			return nil, fmt.Errorf("reading stat for pid %d: %w", pid, err)
		}

		privateMemoryUsage, err := procfs.ReadProcessPrivateMemoryUsage(procRoot, pid)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrPermission) || errors.Is(err, syscall.ESRCH) {
				continue // process exited between reading proc directory and reading process-specific stat, or access denied
			}
			return nil, fmt.Errorf("reading smaps_rollup for pid %d: %w", pid, err)
		}

		proc.PrivateKB = privateMemoryUsage
		processes[pid] = proc
	}
	return processes, nil
}

func calculateSystemCPUUsage(before Snapshot, after Snapshot) CPUUsage {
	deltas := CPUStatDelta(before.CPUStats.Total, after.CPUStats.Total)
	cpuUsageTotal := CPUUtilisation(deltas)

	cpuUsagePerCPU := make([]float64, len(after.CPUStats.PerCPU))
	for i, cpu := range after.CPUStats.PerCPU {
		cpuUsagePerCPU[i] = CPUUtilisation(
			CPUStatDelta(before.CPUStats.PerCPU[i], cpu),
		)
	}

	return CPUUsage{
		Total:  cpuUsageTotal,
		PerCPU: cpuUsagePerCPU,
	}
}

func main() {
	fmt.Println("Hi! I'm ferret 🦦")

	before, err := readSnapshot()
	if err != nil {
		log.Fatalf("could not get initial snapshot: %v", err)
	}

	time.Sleep(time.Second)

	after, err := readSnapshot()
	if err != nil {
		log.Fatalf("could not get final snapshot: %v", err)
	}

	cpuUsage := calculateSystemCPUUsage(before, after)

	deltas := CPUStatDelta(before.CPUStats.Total, after.CPUStats.Total)
	systemTicksDelta := TotalTicks(deltas)
	cpuUsageByPID := make(map[int]float64, len(after.Processes))
	for pid, process := range after.Processes {
		start, exists := before.Processes[pid]
		if !exists {
			continue // process was created during the sleep window
		}

		startTicks := start.UTimeTicks + start.STimeTicks
		stopTicks := process.UTimeTicks + process.STimeTicks
		if stopTicks < startTicks || systemTicksDelta == 0 {
			continue
		}

		procTicksDelta := stopTicks - startTicks
		cpuUsageByPID[pid] = (float64(procTicksDelta) / float64(systemTicksDelta)) * 100.0
	}

	pids := make([]int, 0, len(cpuUsageByPID))
	for pid := range cpuUsageByPID {
		pids = append(pids, pid)
	}
	sort.Slice(pids, func(i, j int) bool { return cpuUsageByPID[pids[i]] > cpuUsageByPID[pids[j]] })

	memoryUsage, err := procfs.ReadMemInfo(procRoot)
	if err != nil {
		log.Fatalf("could not get memory usage info: %v", err)
	}

	for _, pid := range pids {
		process := after.Processes[pid]
		fmt.Printf("Process %d: %s:\n", process.PID, process.Comm)
		fmt.Printf("\t├─ State:%q  Priority:%d  Nice:%d\n", process.State, process.Priority, process.Nice)
		fmt.Printf("\t├─ Virtual memory size (B):%d  Resident Memory size (B):%d\n", process.VSZBytes, process.RSSBytes)
		fmt.Printf("\t├─ Private memory usage: %d KB\n", process.PrivateKB)
		fmt.Printf("\t└─ CPU Usage: %.2f%%\n", cpuUsageByPID[pid])
	}

	fmt.Printf("%d Processes found\n\n", len(after.Processes))
	fmt.Print("System CPU Usage:\n")

	for i, cpuUsage := range cpuUsage.PerCPU {
		fmt.Printf("\t├─ CPU %d: %.1f%%\n", i, cpuUsage)
	}
	fmt.Printf("\t└─ Total: %.1f%%\n\n", cpuUsage.Total)

	fmt.Printf("System Memory Usage: %.1f/%.1f GB\n", KBtoGB(memoryUsage.InUseKB), KBtoGB(memoryUsage.TotalKB))
	fmt.Printf("Memory Available: %.1f GB\n", KBtoGB(memoryUsage.AvailableKB))
}
