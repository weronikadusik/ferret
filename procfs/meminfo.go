package procfs

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type MemInfo struct {
	TotalKB     uint64
	InUseKB     uint64
	AvailableKB uint64
}

func ReadMemInfo(procRoot string) (MemInfo, error) {
	meminfoPath := filepath.Join(procRoot, "meminfo")
	file, err := os.Open(meminfoPath)
	if err != nil {
		return MemInfo{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var total, available uint64
	var haveTotal, haveAvailable bool

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}

		switch fields[0] {
		case "MemTotal:":
			total, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return MemInfo{}, err
			}
			haveTotal = true
		case "MemAvailable:":
			available, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return MemInfo{}, err
			}
			haveAvailable = true
		}

		if haveTotal && haveAvailable {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return MemInfo{}, err
	}

	if !haveTotal {
		return MemInfo{}, errors.New("MemTotal not found in meminfo")
	}
	if !haveAvailable {
		return MemInfo{}, errors.New("MemAvailable not found in meminfo")
	}

	return MemInfo{
		TotalKB:     total,
		InUseKB:     total - available,
		AvailableKB: available,
	}, nil
}
