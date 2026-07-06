package main

import "github.com/weronikadusik/ferret/procfs"

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
