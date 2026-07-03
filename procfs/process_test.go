package procfs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadProcessStat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pid     int
		want    Process
		wantErr bool
	}{
		{
			name: "valid pid",
			pid:  1,
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
			name: "special characters in comm",
			pid:  2,
			want: Process{
				PID:      2,
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
			name:    "insufficient fields",
			pid:     3,
			wantErr: true,
		},
		{
			name:    "invalid numeric format",
			pid:     4,
			wantErr: true,
		},
		{
			name:    "missing pid directory",
			pid:     5,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadProcessStat("./testdata/proc", tt.pid)

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

	want := []int{1, 2, 3, 4}

	got, err := ListPIDs("./testdata/proc")

	require.NoError(t, err)
	require.ElementsMatch(t, want, got)
}
