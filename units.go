package main

import (
	"math"

	"github.com/tklauser/go-sysconf"
)

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

func KBtoGB(amountKB uint64) float64 {
	gb := float64(amountKB) / 1024 / 1024
	return math.Round(gb*10) / 10
}
