package procfs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadProcessPrivateMemoryUsage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		procRoot string
		pid      int
		want     uint64
		wantErr  bool
	}{
		{
			name:     "valid smaps_rollup",
			procRoot: "./testdata/proc_valid",
			pid:      1,
			want:     37264 + 237144,
			wantErr:  false,
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
			name:     "missing smaps_rollup",
			procRoot: "./testdata/proc_missing",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadProcessPrivateMemoryUsage(tt.procRoot, tt.pid)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
