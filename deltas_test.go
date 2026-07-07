package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weronikadusik/ferret/procfs"
)

func TestCPUStatDelta(t *testing.T) {
	before := procfs.CPUStat{
		User:   100,
		System: 50,
		Idle:   1000,
	}

	after := procfs.CPUStat{
		User:   140,
		System: 70,
		Idle:   1200,
	}

	want := procfs.CPUStat{
		User:   40,
		System: 20,
		Idle:   200,
	}

	got := CPUStatDelta(before, after)

	require.Equal(t, want, got)
}

func TestCPUUsage(t *testing.T) {
	deltas := procfs.CPUStat{
		User:    30,
		Nice:    10,
		System:  30,
		Idle:    5,
		IOWait:  5,
		IRQ:     10,
		SoftIRQ: 5,
		Steal:   5,
	}

	got := CPUUsage(deltas)

	require.Equal(t, 90.0, got)
}

func TestCPUUsageZero(t *testing.T) {
	got := CPUUsage(procfs.CPUStat{})

	require.Equal(t, 0.0, got)
}
