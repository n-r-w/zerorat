package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_float64ToRatExact tests the internal exact conversion helper.
func Test_float64ToRatExact(t *testing.T) {
	t.Run("rejects non-finite values", func(t *testing.T) {
		_, err := float64ToRatExact(math.NaN())
		require.ErrorIs(t, err, ErrNonFiniteFloat)

		_, err = float64ToRatExact(math.Inf(1))
		require.ErrorIs(t, err, ErrNonFiniteFloat)

		_, err = float64ToRatExact(math.Inf(-1))
		require.ErrorIs(t, err, ErrNonFiniteFloat)
	})

	t.Run("keeps exact representable values", func(t *testing.T) {
		r, err := float64ToRatExact(0.5)
		require.NoError(t, err)
		assert.Equal(t, Rat{numerator: 1, denominator: 2}, r)

		r, err = float64ToRatExact(-0.25)
		require.NoError(t, err)
		assert.Equal(t, Rat{numerator: -1, denominator: 4}, r)

		r, err = float64ToRatExact(float64(math.MinInt64))
		require.NoError(t, err)
		assert.Equal(t, Rat{numerator: math.MinInt64, denominator: 1}, r)
	})

	t.Run("rejects exact values outside Rat model", func(t *testing.T) {
		_, err := float64ToRatExact(math.Ldexp(1, 63))
		require.ErrorIs(t, err, ErrNotRepresentable)

		_, err = float64ToRatExact(math.MaxFloat64)
		require.ErrorIs(t, err, ErrNotRepresentable)

		_, err = float64ToRatExact(math.SmallestNonzeroFloat64)
		require.ErrorIs(t, err, ErrNotRepresentable)

		_, err = float64ToRatExact(math.Ldexp(1, -64))
		require.ErrorIs(t, err, ErrNotRepresentable)

		_, err = float64ToRatExact(3 * math.Ldexp(1, -64))
		require.ErrorIs(t, err, ErrNotRepresentable)
	})
}

// Test_float64ToRatApprox tests the internal approximation helper.
func Test_float64ToRatApprox(t *testing.T) {
	t.Run("rejects non-finite values", func(t *testing.T) {
		_, err := float64ToRatApprox(math.NaN())
		require.ErrorIs(t, err, ErrNonFiniteFloat)

		_, err = float64ToRatApprox(math.Inf(1))
		require.ErrorIs(t, err, ErrNonFiniteFloat)
	})

	t.Run("returns zero unchanged", func(t *testing.T) {
		r, err := float64ToRatApprox(0)
		require.NoError(t, err)
		assert.True(t, r.Equal(Zero()))
	})

	t.Run("collapses deep subnormals to zero", func(t *testing.T) {
		r, err := float64ToRatApprox(math.SmallestNonzeroFloat64)
		require.NoError(t, err)
		assert.True(t, r.Equal(Zero()))

		r, err = float64ToRatApprox(-math.SmallestNonzeroFloat64)
		require.NoError(t, err)
		assert.True(t, r.Equal(Zero()))
	})

	t.Run("keeps exact values when they already fit", func(t *testing.T) {
		r, err := float64ToRatApprox(0.5)
		require.NoError(t, err)
		assert.Equal(t, Rat{numerator: 1, denominator: 2}, r)
	})

	t.Run("rounds values that need denominator above 2^63", func(t *testing.T) {
		r, err := float64ToRatApprox(math.Ldexp(1, -64))
		require.NoError(t, err)
		assert.True(t, r.Equal(Zero()))

		r, err = float64ToRatApprox(3 * math.Ldexp(1, -64))
		require.NoError(t, err)
		assert.Equal(t, Rat{numerator: 1, denominator: 1 << 62}, r)

		r, err = float64ToRatApprox(-3 * math.Ldexp(1, -64))
		require.NoError(t, err)
		assert.Equal(t, Rat{numerator: -1, denominator: 1 << 62}, r)
	})

	t.Run("still rejects out-of-range finite values", func(t *testing.T) {
		_, err := float64ToRatApprox(math.MaxFloat64)
		require.ErrorIs(t, err, ErrNotRepresentable)
	})
}

// TestNewApproxFromFloat64 tests the exported approximation constructor.
func TestNewApproxFromFloat64(t *testing.T) {
	t.Run("returns reduced approximate value", func(t *testing.T) {
		r, err := NewApproxFromFloat64(3 * math.Ldexp(1, -64))
		require.NoError(t, err)
		assert.Equal(t, mustNewFromFloat64(t, math.Ldexp(1, -62)), r)
	})

	t.Run("returns exact-specific errors unchanged", func(t *testing.T) {
		_, err := NewApproxFromFloat64(math.Inf(1))
		require.ErrorIs(t, err, ErrNonFiniteFloat)

		_, err = NewApproxFromFloat64(math.MaxFloat64)
		require.ErrorIs(t, err, ErrNotRepresentable)
	})
}
