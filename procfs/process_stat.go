package procfs

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReadProcessStat reads /proc/[pid]/stat and returns a Process
func ReadProcessStat(procRoot string, pid int) (Process, error) {
	statPath := filepath.Join(procRoot, strconv.Itoa(pid), "stat")
	data, err := os.ReadFile(statPath)
	if err != nil {
		return Process{}, err
	}

	statStr := string(data)

	lParen := strings.IndexByte(statStr, '(')
	rParen := strings.LastIndexByte(statStr, ')')

	if lParen == -1 || rParen == -1 {
		return Process{}, errors.New("incorrect stat format")
	}

	processStats := strings.Fields(statStr[rParen+1:])

	if len(processStats) < 22 {
		return Process{}, errors.New("incorrect stat format")
	}

	processPriority, err := strconv.Atoi(processStats[15])
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processNice, err := strconv.Atoi(processStats[16])
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processVSZ, err := strconv.ParseUint(processStats[20], 10, 64)
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processRSS, err := strconv.ParseUint(processStats[21], 10, 64)
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processUTime, err := strconv.ParseUint(processStats[11], 10, 64)
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processSTime, err := strconv.ParseUint(processStats[12], 10, 64)
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	return Process{
		PID:        pid,
		Comm:       statStr[lParen+1 : rParen],
		State:      processStats[0][0],
		Priority:   processPriority,
		Nice:       processNice,
		VSZBytes:   processVSZ,
		RSSBytes:   processRSS * uint64(os.Getpagesize()),
		UTimeTicks: processUTime,
		STimeTicks: processSTime,
	}, nil
}
