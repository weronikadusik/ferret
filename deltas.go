package main

import "github.com/weronikadusik/ferret/procfs"

// TotalTicks calculates the total system ticks across all CPU states for a given stat struct
func TotalTicks(stat procfs.CPUStat) uint64 {
	return stat.User + stat.Nice + stat.System + stat.Idle +
		stat.IOWait + stat.IRQ + stat.SoftIRQ + stat.Steal
}

// CPUStatDelta returns the difference between two CPU samples
func CPUStatDelta(before, after procfs.CPUStat) procfs.CPUStat {
	return procfs.CPUStat{
		User:    after.User - before.User,
		Nice:    after.Nice - before.Nice,
		System:  after.System - before.System,
		Idle:    after.Idle - before.Idle,
		IOWait:  after.IOWait - before.IOWait,
		IRQ:     after.IRQ - before.IRQ,
		SoftIRQ: after.SoftIRQ - before.SoftIRQ,
		Steal:   after.Steal - before.Steal,
	}
}

// CPUUsage returns the average CPU utilisation percentage from a procfs.CPUStat containing cpu time deltas
func CPUUsage(deltas procfs.CPUStat) float64 {
	total := deltas.User + deltas.Nice + deltas.System + deltas.Idle + deltas.IOWait + deltas.IRQ + deltas.SoftIRQ + deltas.Steal
	busy := deltas.User + deltas.Nice + deltas.System + deltas.IRQ + deltas.SoftIRQ + deltas.Steal

	if total == 0 {
		return 0.0
	}

	return (float64(busy) / float64(total)) * 100.0
}
