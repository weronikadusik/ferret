package procfs

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type DiskStats struct {
	Name   string
	IOTime uint64
}

func isWholeDisk(sysRoot, name string) bool {
	_, err := os.Lstat(filepath.Join(sysRoot, "sys", "class", "block", name, "device"))
	return err == nil
}

func ReadDiskStats(procRoot, sysRoot string) ([]DiskStats, error) {
	diskstatsPath := filepath.Join(procRoot, "diskstats")
	file, err := os.Open(diskstatsPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var diskStats []DiskStats
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 14 {
			return nil, errors.New("incorrect diskstats format")
		}

		name := fields[2]
		if !isWholeDisk(sysRoot, name) {
			continue // skip partitions, loop, zram, dm-*, md*, ...
		}

		ioTime, err := strconv.ParseUint(fields[12], 10, 64)
		if err != nil {
			return nil, err
		}

		diskStats = append(diskStats, DiskStats{
			Name:   name,
			IOTime: ioTime,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return diskStats, nil
}
