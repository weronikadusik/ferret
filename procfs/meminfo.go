package procfs

import (
	"path/filepath"
)

type MemInfo struct {
	TotalKB     uint64
	InUseKB     uint64
	AvailableKB uint64
}

func ReadMemInfo(procRoot string) (MemInfo, error) {
	meminfoPath := filepath.Join(procRoot, "meminfo")
	meminfo, err := readNamedFields(meminfoPath, []string{
		"MemTotal",
		"MemAvailable",
	})
	if err != nil {
		return MemInfo{}, err
	}

	return MemInfo{
		TotalKB:     meminfo["MemTotal"],
		InUseKB:     meminfo["MemTotal"] - meminfo["MemAvailable"],
		AvailableKB: meminfo["MemAvailable"],
	}, nil
}
