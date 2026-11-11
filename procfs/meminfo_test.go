package procfs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

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
