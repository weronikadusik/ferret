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
