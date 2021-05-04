package vic_ii

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestColormap(t *testing.T) {
	require.Equal(t, 16, len(C64Colors))
}
