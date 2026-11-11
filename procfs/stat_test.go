package procfs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadStat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		procRoot string
		want     CPUStats
		wantErr  bool
	}{
		{
			name:     "valid stat",
			procRoot: "./testdata/proc_valid",
			want: CPUStats{
				Total: CPUTimes{
					User:    434667,
					Nice:    2646,
					System:  124574,
					Idle:    3630324,
					IOWait:  6822,
					IRQ:     29795,
					SoftIRQ: 9612,
					Steal:   0,
				},
				PerCPU: []CPUTimes{
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
