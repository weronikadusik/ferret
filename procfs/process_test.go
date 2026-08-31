package procfs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadProcessStat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		procRoot string
		pid      int
		want     Process
		wantErr  bool
	}{
		{
			name:     "valid pid",
			procRoot: "./testdata/proc_valid",
			pid:      1,
			want: Process{
				PID:      1,
				Comm:     "zsh",
				State:    'S',
				Priority: 20,
				Nice:     0,

				VSZBytes: 242819072,
				RSSBytes: 2179 * uint64(os.Getpagesize()),

				UTimeTicks: 20,
				STimeTicks: 7,
			},
			wantErr: false,
		},
		{
			name:     "special characters in comm",
			procRoot: "./testdata/proc_valid",
			pid:      42,
			want: Process{
				PID:      42,
				Comm:     "abc/XYZ (Preview)",
				State:    'S',
				Priority: 20,
				Nice:     0,

				VSZBytes: 242819072,
				RSSBytes: 2179 * uint64(os.Getpagesize()),

				UTimeTicks: 20,
				STimeTicks: 7,
			},
			wantErr: false,
		},
		{
			name:     "insufficient fields",
			procRoot: "./testdata/proc_malformed",
			pid:      123,
			wantErr:  true,
		},
		{
			name:     "invalid numeric format",
			procRoot: "./testdata/proc_malformed",
			pid:      999,
			wantErr:  true,
		},
		{
			name:     "missing pid directory",
			procRoot: "./testdata/proc_missing",
			pid:      50,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadProcessStat(tt.procRoot, tt.pid)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestReadStat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		procRoot string
		want     SystemStat
		wantErr  bool
	}{
		{
			name:     "valid stat",
			procRoot: "./testdata/proc_valid",
			want: SystemStat{
				Total: CPUStat{
					User:    434667,
					Nice:    2646,
					System:  124574,
					Idle:    3630324,
					IOWait:  6822,
					IRQ:     29795,
					SoftIRQ: 9612,
					Steal:   0,
				},
				PerCPU: []CPUStat{
					{
						User:    72398,
						Nice:    440,
						System:  21302,
						Idle:    604918,
						IOWait:  1034,
						IRQ:     4483,
						SoftIRQ: 1496,
						Steal:   0,
					},
					{
						User:    362269,
						Nice:    2206,
						System:  103272,
						Idle:    3025406,
						IOWait:  5788,
						IRQ:     25312,
						SoftIRQ: 8116,
						Steal:   0,
					},
				},
			},
		},
		{
			name:     "invalid numeric format",
			procRoot: "./testdata/proc_malformed",
			wantErr:  true,
		},
		{
			name:     "missing stat",
			procRoot: "./testdata/proc_missing",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadStat(tt.procRoot)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestReadMemInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		procRoot string
		want     MemInfo
		wantErr  bool
	}{
		{
			name:     "valid meminfo",
			procRoot: "./testdata/proc_valid",
			want: MemInfo{
				TotalKB:     16282144,
				AvailableKB: 9363644,
				InUseKB:     16282144 - 9363644,
			},
		},
		{
			name:     "invalid numeric format",
			procRoot: "./testdata/proc_malformed",
			wantErr:  true,
		},
		{
			name:     "missing meminfo",
			procRoot: "./testdata/proc_missing",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadMemInfo(tt.procRoot)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestListPIDs(t *testing.T) {
	t.Parallel()

	want := []int{1, 42}

	got, err := ListPIDs("./testdata/proc_valid")

	require.NoError(t, err)
	require.ElementsMatch(t, want, got)
}
