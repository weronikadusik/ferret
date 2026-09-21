package procfs

import (
	"path/filepath"
	"strconv"
)

// ReadProcessDiskUsage reads /proc/[pid]/io and returns the total number of bytes read from and written to the storage layer
func ReadProcessDiskUsage(procRoot string, pid int) (uint64, error) {
	ioPath := filepath.Join(procRoot, strconv.Itoa(pid), "io")
	io, err := readNamedFields(ioPath, []string{
		"read_bytes",
		"write_bytes",
	})
	if err != nil {
		return 0, err
	}

	return io["read_bytes"] + io["write_bytes"], nil
}
