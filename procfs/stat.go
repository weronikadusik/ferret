package procfs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type SystemStat struct {
	Total   CPUStat
	PerCore []CPUStat
}

type CPUStat struct {
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
func ReadStat(procRoot string) (SystemStat, error) {
	var cpuStat []CPUStat
	var total CPUStat

	statPath := filepath.Join(procRoot, "stat")
	data, err := os.ReadFile(statPath)
	if err != nil {
		return SystemStat{}, err
	}

	statStr := string(data)

	lines := strings.Split(statStr, "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "cpu") {
			break
		}

		fields := strings.Fields(line)

		if len(fields) < 8 {
			return SystemStat{}, errors.New("incorrect stat format")
		}

		statUser, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return SystemStat{}, err
		}

		statNice, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return SystemStat{}, err
		}

		statSystem, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return SystemStat{}, err
		}

		statIdle, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return SystemStat{}, err
		}

		statIOWait, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			return SystemStat{}, err
		}

		statIRQ, err := strconv.ParseUint(fields[6], 10, 64)
		if err != nil {
			return SystemStat{}, err
		}

		statSoftIRQ, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			return SystemStat{}, err
		}

		statSteal, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return SystemStat{}, err
		}

		metrics := CPUStat{
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

	return SystemStat{
		Total:   total,
		PerCore: cpuStat,
	}, nil
}
