package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoneyMulManyInt tests MulManyInt varargs operations
func TestMoneyMulManyInt(t *testing.T) {
	t.Run("mutable MulManyInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 10) // $0.10

		err := m.MulManyInt(2, 3, 5) // multiply by 2, then 3, then 5

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 300) // 10 * 2 * 3 * 5 = 300
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulManyInt - with zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulManyInt(2, 0, 5) // multiply by 2, then 0, then 5

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := ZeroMoney("USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulManyInt - empty varargs", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		original := m.Amount()

		err := m.MulManyInt()

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Amount().Equal(original), "Money should be unchanged with empty varargs")
	})

	t.Run("mutable MulManyInt - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid

		err := m.MulManyInt(2, 3)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("immutable MultipliedManyIntErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 10)

		result, err := m.MultipliedManyIntErr(2, 3)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 60) // 10 * 2 * 3 = 60
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 10)))
	})

	t.Run("immutable MultipliedManyInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 10)

		result := m.MultipliedManyInt(2, 3)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 60)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 10)))
	})
}

// TestMoneyMulManyRat tests MulManyRat varargs operations.
func TestMoneyMulManyRat(t *testing.T) {
	t.Run("mutable MulManyRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulManyRat(zerorat.NewFromInt64(2), mustNewRatFromFloat64(t, 1.5))

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := mustNewMoneyFloat(t, "USD", 300.0)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulManyRat - with zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulManyRat(zerorat.NewFromInt64(2), zerorat.Zero(), zerorat.NewFromInt64(5))

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Equal(ZeroMoney("USD")))
	})

	t.Run("mutable MulManyRat - invalid Rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulManyRat(zerorat.NewFromInt64(2), zerorat.Rat{})

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable MulManyRat - empty varargs", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		original := m.Amount()

		err := m.MulManyRat()

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Amount().Equal(original), "Money should be unchanged with empty varargs")
	})

	t.Run("immutable MultipliedManyRatErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result, err := m.MultipliedManyRatErr(zerorat.NewFromInt64(2), mustNewRatFromFloat64(t, 1.5))

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := mustNewMoneyFloat(t, "USD", 300.0)
		assert.True(t, result.Equal(expected))
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable MultipliedManyRat - invalid Rat returns invalid", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result := m.MultipliedManyRat(zerorat.NewFromInt64(2), zerorat.Rat{})

		assert.True(t, result.IsInvalid())
		assert.True(t, m.IsValid())
	})
}

// TestMoneyReduce tests Reduce and Reduced methods.
func TestMoneyReduce(t *testing.T) {
	t.Run("mutable Reduce - reduces fraction to lowest terms", func(t *testing.T) {
		m := NewMoneyFromFraction(6, 4, "USD") // 6/4 should reduce to 3/2

		m.Reduce()

		assert.True(t, m.IsValid())
		expected := NewMoneyFromFraction(3, 2, "USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable Reduce - already reduced fraction stays unchanged", func(t *testing.T) {
		m := NewMoneyFromFraction(3, 2, "USD") // already in lowest terms
		original := m.Amount()

		m.Reduce()

		assert.True(t, m.IsValid())
		assert.True(t, m.Amount().Equal(original), "Already reduced fraction should remain unchanged")
	})

	t.Run("mutable Reduce - zero value stays zero", func(t *testing.T) {
		m := ZeroMoney("USD")

		m.Reduce()

		assert.True(t, m.IsValid())
		assert.True(t, m.Equal(ZeroMoney("USD")))
	})

	t.Run("mutable Reduce - large fraction reduction", func(t *testing.T) {
		m := NewMoneyFromFraction(100, 25, "USD") // 100/25 should reduce to 4/1

		m.Reduce()

		assert.True(t, m.IsValid())
		expected := NewMoneyFromFraction(4, 1, "USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable Reduce - fraction with common factor", func(t *testing.T) {
		m := NewMoneyFromFraction(15, 10, "USD") // 15/10 should reduce to 3/2

		m.Reduce()

		assert.True(t, m.IsValid())
		expected := NewMoneyFromFraction(3, 2, "USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("immutable Reduced - returns new reduced Money", func(t *testing.T) {
		m := NewMoneyFromFraction(6, 4, "USD") // 6/4
		original := m.Amount()

		result := m.Reduced()

		assert.True(t, result.IsValid())
		expected := NewMoneyFromFraction(3, 2, "USD")
		assert.True(t, result.Equal(expected))
		assert.True(t, m.Amount().Equal(original), "Original Money should be unchanged")
	})

	t.Run("immutable Reduced - already reduced returns equivalent", func(t *testing.T) {
		m := NewMoneyFromFraction(3, 2, "USD")

		result := m.Reduced()

		assert.True(t, result.IsValid())
		assert.True(t, result.Equal(m))
	})

	t.Run("immutable Reduced - zero value", func(t *testing.T) {
		m := ZeroMoney("USD")

		result := m.Reduced()

		assert.True(t, result.IsValid())
		assert.True(t, result.Equal(ZeroMoney("USD")))
	})

	t.Run("mutable Reduce - invalid Money remains invalid", func(t *testing.T) {
		m := Money{} // invalid

		m.Reduce()

		assert.True(t, m.IsInvalid())
	})

	t.Run("immutable Reduced - invalid Money returns invalid", func(t *testing.T) {
		m := Money{} // invalid

		result := m.Reduced()

		assert.True(t, result.IsInvalid())
	})

	t.Run("mutable Reduce - whole number remains whole", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		m.Reduce()

		assert.True(t, m.IsValid())
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})
}
