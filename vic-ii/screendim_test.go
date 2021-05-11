package vic_ii

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScreenDim(t *testing.T) {
	require.Equal(t, 403, PalRightBorderWidth38Cols+PalContentWidth40Cols+PalRightBorderWidth40Cols, "Visible width doesn't add up")
}
