package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"syscall"
	"time"

	"github.com/weronikadusik/ferret/procfs"
)

const procRoot = "/proc"

type ProcessMetrics struct {
	Process         procfs.Process
	CPUUsagePercent float64
	PrivateKB       uint64
}

type Snapshot struct {
	Processes map[int]procfs.Process
	CPUStats  procfs.CPUStats
}

type CPUUsage struct {
	Total  float64
	PerCPU []float64
}

func isSkippable(err error) bool {
	return errors.Is(err, os.ErrNotExist) ||
		errors.Is(err, os.ErrPermission) ||
		errors.Is(err, syscall.ESRCH)
}

func getStablePIDs(before, after Snapshot) []int {
	var pids []int

	for pid := range before.Processes {
		if _, exists := after.Processes[pid]; !exists {
			continue
		}

		pids = append(pids, pid)
	}

	return pids
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
			if isSkippable(err) {
				continue // process exited between reading proc directory and reading process-specific stat, or access denied
			}
			return nil, fmt.Errorf("reading stat for pid %d: %w", pid, err)
		}

		processes[pid] = proc
	}
	return processes, nil
}

func getProcessMetrics(before, after Snapshot, stablePIDs []int, cpuDelta procfs.CPUTimes) ([]ProcessMetrics, error) {
	systemTicksDelta := TotalTicks(cpuDelta)

	processesMetrics := make([]ProcessMetrics, 0, len(stablePIDs))
	for _, pid := range stablePIDs {
		privateKB, err := memoryUsageByPID(pid)
		if err != nil {
			return nil, err
		}

		processesMetrics = append(processesMetrics, ProcessMetrics{
			Process:         after.Processes[pid],
			CPUUsagePercent: cpuUsageByPID(pid, before, after, systemTicksDelta),
			PrivateKB:       privateKB,
		})
	}

	return processesMetrics, nil
}

func memoryUsageByPID(pid int) (uint64, error) {
	privateMemoryUsage, err := procfs.ReadProcessPrivateMemoryUsage(procRoot, pid)
	if err != nil {
		if isSkippable(err) {
			return 0, nil // process exited between reading proc directory and reading process-specific stat, or access denied
		}
		return 0, fmt.Errorf("reading smaps_rollup for pid %d: %w", pid, err)
	}

	return privateMemoryUsage, nil
}

func cpuUsageByPID(pid int, before Snapshot, after Snapshot, systemTicksDelta uint64) float64 {
	startTicks := before.Processes[pid].UTimeTicks + before.Processes[pid].STimeTicks
	stopTicks := after.Processes[pid].UTimeTicks + after.Processes[pid].STimeTicks
	if stopTicks < startTicks || systemTicksDelta == 0 {
		return 0
	}

	procTicksDelta := stopTicks - startTicks
	return (float64(procTicksDelta) / float64(systemTicksDelta)) * 100.0
}

func systemCPUUsage(before Snapshot, after Snapshot, cpuDelta procfs.CPUTimes) CPUUsage {
	cpuUsagePerCPU := make([]float64, len(after.CPUStats.PerCPU))
	for i, cpu := range after.CPUStats.PerCPU {
		cpuUsagePerCPU[i] = CPUUtilisation(
			CPUStatDelta(before.CPUStats.PerCPU[i], cpu),
		)
	}

	return CPUUsage{
		Total:  CPUUtilisation(cpuDelta),
		PerCPU: cpuUsagePerCPU,
	}
}

func printProcesses(metrics []ProcessMetrics) {
	for _, p := range metrics {
		fmt.Printf("Process %d: %s:\n", p.Process.PID, p.Process.Comm)
		fmt.Printf("\t├─ State:%q  Priority:%d  Nice:%d\n", p.Process.State, p.Process.Priority, p.Process.Nice)
		fmt.Printf("\t├─ Virtual memory size: %.1f MB  Resident Memory size: %.1f MB\n", BtoMB(p.Process.VSZBytes), BtoMB(p.Process.RSSBytes))
		fmt.Printf("\t├─ Private memory usage: %.1f MB\n", KBtoMB(p.PrivateKB))
		fmt.Printf("\t└─ CPU Usage: %.2f%%\n", p.CPUUsagePercent)
	}

	fmt.Printf("%d Processes found\n\n", len(metrics))
}

func printSystemCPUUsage(cpuUsage CPUUsage) {
	fmt.Print("System CPU Usage:\n")

	for i, cpuUsage := range cpuUsage.PerCPU {
		fmt.Printf("\t├─ CPU %d: %.1f%%\n", i, cpuUsage)
	}
	fmt.Printf("\t└─ Total: %.1f%%\n\n", cpuUsage.Total)
}

func printSystemMemoryUsage(memoryUsage procfs.MemInfo) {
	fmt.Printf("System Memory Usage: %.1f/%.1f GB\n", KBtoGB(memoryUsage.InUseKB), KBtoGB(memoryUsage.TotalKB))
	fmt.Printf("Memory Available: %.1f GB\n", KBtoGB(memoryUsage.AvailableKB))
}

func main() {
	sortBy := flag.String("sort", "pid", "Sort process list by `pid`, `cpu`, or `memory`")
	flag.Parse()

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

	cpuDelta := CPUStatDelta(before.CPUStats.Total, after.CPUStats.Total)
	cpuUsage := systemCPUUsage(before, after, cpuDelta)

	memoryUsage, err := procfs.ReadMemInfo(procRoot)
	if err != nil {
		log.Fatalf("could not get memory usage info: %v", err)
	}

	stablePIDs := getStablePIDs(before, after)
	processesMetrics, err := getProcessMetrics(before, after, stablePIDs, cpuDelta)
	if err != nil {
		log.Fatalf("could not get process metrics: %v", err)
	}

	switch *sortBy {
	case "pid":
		sort.Slice(processesMetrics, func(i, j int) bool {
			return processesMetrics[i].Process.PID < processesMetrics[j].Process.PID
		})
	case "cpu":
		sort.Slice(processesMetrics, func(i, j int) bool {
			return processesMetrics[i].CPUUsagePercent > processesMetrics[j].CPUUsagePercent
		})
	case "memory":
		sort.Slice(processesMetrics, func(i, j int) bool {
			return processesMetrics[i].PrivateKB > processesMetrics[j].PrivateKB
		})
	default:
		log.Fatalf("invalid sort option: %q (valid: pid, cpu, memory)", *sortBy)
	}

	printProcesses(processesMetrics)
	printSystemCPUUsage(cpuUsage)
	printSystemMemoryUsage(memoryUsage)
}
