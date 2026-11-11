package procfs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDiskStats(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		procRoot string
		sysRoot  string
		want     []DiskStats
		wantErr  bool
	}{
		{
			name:     "valid diskstats",
			procRoot: "./testdata/proc_valid",
			sysRoot:  "./testdata/sysroot",
			want: []DiskStats{
				{Name: "sda", IOTime: 113517},
				{Name: "sdb", IOTime: 48},
			},
			wantErr: false,
		},
		{
			name:     "malformed diskstats",
			procRoot: "./testdata/proc_malformed",
			sysRoot:  "./testdata/sysroot",
			wantErr:  true,
		},
		{
			name:     "missing diskstats",
			procRoot: "./testdata/proc_missing",
			sysRoot:  "./testdata/sysroot",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDiskStats(tt.procRoot, tt.sysRoot)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}
