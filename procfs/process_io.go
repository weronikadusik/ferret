package procfs

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadProcessDiskUsage reads /proc/[pid]/io and returns the total number of bytes read from and written to the storage layer
func ReadProcessDiskUsage(procRoot string, pid int) (uint64, error) {
	ioPath := filepath.Join(procRoot, strconv.Itoa(pid), "io")
	file, err := os.Open(ioPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var readB, writeB uint64
	var haveReadB, haveWriteB bool

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "read_bytes:":
			readB, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, err
			}
			haveReadB = true
		case "write_bytes:":
			writeB, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, err
			}
			haveWriteB = true
		}

		if haveReadB && haveWriteB {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	if !haveReadB {
		return 0, errors.New("read_bytes value not found in io file")
	}
	if !haveWriteB {
		return 0, errors.New("write_bytes value not found in io file")
	}

	return readB + writeB, nil
}
