package procfs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListPIDs(t *testing.T) {
	t.Parallel()

	want := []int{1, 42}

	got, err := ListPIDs("./testdata/proc_valid")

	require.NoError(t, err)
	require.ElementsMatch(t, want, got)
}
