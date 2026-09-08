package procfs

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadProcessPrivateMemoryUsage(procRoot string, pid int) (uint64, error) {
	smapsPath := filepath.Join(procRoot, strconv.Itoa(pid), "smaps_rollup")
	file, err := os.Open(smapsPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var clean, dirty uint64
	var havePrivateClean, havePrivateDirty bool

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		switch fields[0] {
		case "Private_Clean:":
			clean, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, err
			}
			havePrivateClean = true
		case "Private_Dirty:":
			dirty, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return 0, err
			}
			havePrivateDirty = true
		}

		if havePrivateClean && havePrivateDirty {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	if !havePrivateClean {
		return 0, errors.New("Private_Clean value not found in smaps_rollup")
	}
	if !havePrivateDirty {
		return 0, errors.New("Private_Dirty value not found in smaps_rollup")
	}

	return clean + dirty, nil
}
