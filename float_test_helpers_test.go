package zerorat

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// mustNewFromFloat64 creates an exact Rat from float64 for tests.
func mustNewFromFloat64(t *testing.T, value float64) Rat {
	t.Helper()

	r, err := NewFromFloat64(value)
	require.NoError(t, err)

	return r
}
