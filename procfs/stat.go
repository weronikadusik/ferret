package procfs

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CPUStats struct {
	Total  CPUTimes
	PerCPU []CPUTimes
}

type CPUTimes struct {
	User    uint64
	Nice    uint64
	System  uint64
	Idle    uint64
	IOWait  uint64
	IRQ     uint64
	SoftIRQ uint64
	Steal   uint64
}

// ReadStat reads /proc/stat and returns filled SystemStat struct
func ReadStat(procRoot string) (CPUStats, error) {
	var cpuStat []CPUTimes
	var total CPUTimes

	statPath := filepath.Join(procRoot, "stat")
	file, err := os.Open(statPath)
	if err != nil {
		return CPUStats{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if !strings.HasPrefix(scanner.Text(), "cpu") {
			break
		}

		fields := strings.Fields(scanner.Text())

		if len(fields) < 8 {
			return CPUStats{}, errors.New("incorrect stat format")
		}

		statUser, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return CPUStats{}, err
		}

		statNice, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return CPUStats{}, err
		}

		statSystem, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return CPUStats{}, err
		}

		statIdle, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return CPUStats{}, err
		}

		statIOWait, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			return CPUStats{}, err
		}

		statIRQ, err := strconv.ParseUint(fields[6], 10, 64)
		if err != nil {
			return CPUStats{}, err
		}

		statSoftIRQ, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			return CPUStats{}, err
		}

		statSteal, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return CPUStats{}, err
		}

		metrics := CPUTimes{
			User:    statUser,
			Nice:    statNice,
			System:  statSystem,
			Idle:    statIdle,
			IOWait:  statIOWait,
			IRQ:     statIRQ,
			SoftIRQ: statSoftIRQ,
			Steal:   statSteal,
		}

		if fields[0] == "cpu" {
			total = metrics
		} else {
			cpuStat = append(cpuStat, metrics)
		}
	}

	if err := scanner.Err(); err != nil {
		return CPUStats{}, err
	}

	return CPUStats{
		Total:  total,
		PerCPU: cpuStat,
	}, nil
}
