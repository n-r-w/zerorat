package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/require"
)

// mustNewMoneyFloat creates Money from a float64 for tests.
func mustNewMoneyFloat(t *testing.T, currency Currency, value float64) Money {
	t.Helper()

	m, err := NewMoneyFloat(currency, value)
	require.NoError(t, err)

	return m
}

// mustNewRatFromFloat64 creates an exact Rat from float64 for Money-package tests.
func mustNewRatFromFloat64(t *testing.T, value float64) zerorat.Rat {
	t.Helper()

	r, err := zerorat.NewFromFloat64(value)
	require.NoError(t, err)

	return r
}
