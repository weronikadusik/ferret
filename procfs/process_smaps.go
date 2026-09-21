package procfs

import (
	"path/filepath"
	"strconv"
)

// ReadProcessPrivateMemoryUsage reads /proc/[pid]/smaps_rollup and returns the sum of Private_Clean and Private_Dirty memory.
func ReadProcessPrivateMemoryUsage(procRoot string, pid int) (uint64, error) {
	smapsPath := filepath.Join(procRoot, strconv.Itoa(pid), "smaps_rollup")
	smaps, err := readNamedFields(smapsPath, []string{
		"Private_Clean",
		"Private_Dirty",
	})
	if err != nil {
		return 0, err
	}

	return smaps["Private_Clean"] + smaps["Private_Dirty"], nil
}
