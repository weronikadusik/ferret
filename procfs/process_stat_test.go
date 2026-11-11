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
