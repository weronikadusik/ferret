package procfs

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadProcessStat(t *testing.T) {
	want := Process{
		PID:      1,
		Comm:     "zsh",
		State:    'S',
		Priority: 20,
		Nice:     0,

		VSZBytes: 242819072,
		RSSBytes: 2179 * uint64(os.Getpagesize()),

		UTimeTicks: 20,
		STimeTicks: 7,
	}

	got, err := ReadProcessStat("./testdata/proc", 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("Mismatch (-x, +y):\n%s", diff)
	}
}

func TestReadProcessStatMalformed(t *testing.T) {
	_, err := ReadProcessStat("./testdata/proc", 99)
	if err == nil {
		t.Errorf("Expected to parse with errors but was success")
	}
}

func TestReadProcessStatNonExistent(t *testing.T) {
	_, err := ReadProcessStat("./testdata/proc", 50)
	if err == nil {
		t.Errorf("Expected to parse with errors but was success")
	}
}

func TestListPIDs(t *testing.T) {
	want := []int{1, 99}

	got, err := ListPIDs("./testdata/proc")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("Mismatch (-x, +y):\n%s", diff)
	}
}
