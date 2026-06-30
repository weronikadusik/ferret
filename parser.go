package main

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tklauser/go-sysconf"
)

type Process struct {
	PID      int
	Comm     string
	State    byte
	Priority int
	Nice     int

	VSZBytes uint64
	RSSBytes uint64

	UTimeTicks uint64
	STimeTicks uint64
}

// clkTck is the number of clock ticks per second used by the kernel to track CPU time
var clkTck = func() int64 {
	clkTck, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
	if err != nil {
		return 100
	}
	return clkTck
}()

// ticksToSeconds converts CPU ticks to seconds
func ticksToSeconds(ticks float64) float64 {
	return ticks / float64(clkTck)
}

// read /proc/[pid]/stat; return filled Process struct
func getProcessInfo(procRoot string, pid int) (Process, error) {
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

	processVSZ, err := strconv.Atoi(processStats[20])
	if err != nil {
		return Process{}, errors.New("incorrect stat format")
	}

	processRSS, err := strconv.Atoi(processStats[21])
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
		VSZBytes:   uint64(processVSZ),
		RSSBytes:   uint64(processRSS) * uint64(os.Getpagesize()),
		UTimeTicks: processUTime,
		STimeTicks: processSTime,
	}, nil
}
