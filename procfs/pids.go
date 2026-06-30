package procfs

import (
	"os"
	"strconv"
)

// ListPIDs scans procRoot and returns the PIDs of all running processes.
func ListPIDs(procRoot string) ([]int, error) {
	var PIDs []int

	processDir, err := os.ReadDir(procRoot)
	if err != nil {
		return []int{}, err
	}

	// filter out non-numeric directory names, so only process-specific ones remain
	for _, f := range processDir {
		pid, err := strconv.Atoi(f.Name())
		if err != nil {
			continue
		}

		PIDs = append(PIDs, pid)
	}

	return PIDs, nil
}
