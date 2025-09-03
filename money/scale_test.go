package money

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoney_ScaleDown tests the mutable ScaleDown method
func TestMoney_ScaleDown(t *testing.T) {
	t.Run("basic scaling down", func(t *testing.T) {
		tests := []struct {
			name     string
			money    Money
			scale    int
			expected Money
		}{
			{"USD 1/2 scale down 1", NewMoneyFromFraction(1, 2, "USD"), 1, NewMoneyFromFraction(1, 20, "USD")},      // $0.5 -> $0.05
			{"USD 1/2 scale down 2", NewMoneyFromFraction(1, 2, "USD"), 2, NewMoneyFromFraction(1, 200, "USD")},     // $0.5 -> $0.005
			{"EUR 3/4 scale down 1", NewMoneyFromFraction(3, 4, "EUR"), 1, NewMoneyFromFraction(3, 40, "EUR")},      // €0.75 -> €0.075
			{"USD 5 scale down 2", NewMoneyInt("USD", 5), 2, NewMoneyFromFraction(5, 100, "USD")},                   // $5 -> $0.05
			{"USD negative scale down", NewMoneyFromFraction(-7, 3, "USD"), 1, NewMoneyFromFraction(-7, 30, "USD")}, // -$2.33... -> -$0.233...
			{"USD zero scale down", ZeroMoney("USD"), 3, NewMoneyFromFraction(0, 1000, "USD")},                      // $0 -> $0 (but denom changes)
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				m := tt.money
				err := m.ScaleDown(tt.scale)
				require.NoError(t, err, "ScaleDown should not return error for valid money")
				assert.True(t, m.IsValid(), "Money should remain valid")
				assert.Equal(t, tt.expected.Currency(), m.Currency(), "Currency should be preserved")
				assert.True(t, m.Equal(tt.expected), "Amount should match expected after scaling down")
			})
		}
	})

	t.Run("zero scale does nothing", func(t *testing.T) {
		m := NewMoneyFromFraction(3, 7, "USD")
		original := m
		err := m.ScaleDown(0)
		require.NoError(t, err)
		assert.Equal(t, original.Currency(), m.Currency(), "Currency should be unchanged")
		assert.True(t, m.Equal(original), "Money should be unchanged with zero scale")
	})

	t.Run("negative scale calls ScaleUp", func(t *testing.T) {
		m1 := NewMoneyFromFraction(1, 200, "USD")
		m2 := NewMoneyFromFraction(1, 200, "USD")

		err1 := m1.ScaleDown(-2) // Should be equivalent to ScaleUp(2)
		err2 := m2.ScaleUp(2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.True(t, m1.Equal(m2), "ScaleDown(-2) should equal ScaleUp(2)")
	})

	t.Run("invalid state handling", func(t *testing.T) {
		m := NewMoneyInt("", 100) // Invalid currency
		err := m.ScaleDown(2)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid(), "invalid input should remain invalid")
	})

	t.Run("overflow handling", func(t *testing.T) {
		// Test denominator overflow - create Money with large denominator
		m := NewMoneyFromFraction(1, math.MaxUint64/5, "USD") // Large denominator
		err := m.ScaleDown(10)                                // This should cause overflow
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid(), "overflow should result in invalid state")
	})
}

// TestMoney_ScaledDownErr tests the immutable ScaledDownErr method
func TestMoney_ScaledDownErr(t *testing.T) {
	t.Run("basic scaling down immutable with error", func(t *testing.T) {
		tests := []struct {
			name     string
			money    Money
			scale    int
			expected Money
		}{
			{"USD 1/2 scaled down 1", NewMoneyFromFraction(1, 2, "USD"), 1, NewMoneyFromFraction(1, 20, "USD")},
			{"USD 1/2 scaled down 2", NewMoneyFromFraction(1, 2, "USD"), 2, NewMoneyFromFraction(1, 200, "USD")},
			{"EUR 3/4 scaled down 1", NewMoneyFromFraction(3, 4, "EUR"), 1, NewMoneyFromFraction(3, 40, "EUR")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.money
				result, err := tt.money.ScaledDownErr(tt.scale)

				require.NoError(t, err, "ScaledDownErr should not return error for valid money")

				// Check original is unchanged
				assert.Equal(t, original.Currency(), tt.money.Currency(), "original currency should be unchanged")
				assert.True(t, tt.money.Equal(original), "original money should be unchanged")

				// Check result
				assert.True(t, result.IsValid(), "result should be valid")
				assert.Equal(t, tt.expected.Currency(), result.Currency(), "result currency should match expected")
				assert.True(t, result.Equal(tt.expected), "result should match expected after scaling down")
			})
		}
	})

	t.Run("invalid money returns error", func(t *testing.T) {
		m := NewMoneyInt("", 100) // Invalid currency
		result, err := m.ScaledDownErr(2)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid(), "result should be invalid")
	})
}

// TestMoney_ScaledDown tests the immutable ScaledDown method (without error)
func TestMoney_ScaledDown(t *testing.T) {
	t.Run("basic scaling down immutable without error", func(t *testing.T) {
		m := NewMoneyFromFraction(1, 2, "USD")
		original := m
		result := m.ScaledDown(2)

		// Check original is unchanged
		assert.True(t, m.Equal(original), "original money should be unchanged")

		// Check result
		expected := NewMoneyFromFraction(1, 200, "USD")
		assert.True(t, result.Equal(expected), "result should match expected")
	})

	t.Run("invalid money returns invalid", func(t *testing.T) {
		m := NewMoneyInt("", 100) // Invalid currency
		result := m.ScaledDown(2)
		assert.True(t, result.IsInvalid(), "result should be invalid for invalid input")
	})
}

// TestMoney_ScaleUp tests the mutable ScaleUp method
func TestMoney_ScaleUp(t *testing.T) {
	t.Run("basic scaling up", func(t *testing.T) {
		tests := []struct {
			name     string
			money    Money
			scale    int
			expected Money
		}{
			{"USD 1/200 scale up 2", NewMoneyFromFraction(1, 200, "USD"), 2, NewMoneyFromFraction(100, 200, "USD")}, // $0.005 -> $0.5
			{"USD 1/20 scale up 1", NewMoneyFromFraction(1, 20, "USD"), 1, NewMoneyFromFraction(10, 20, "USD")},     // $0.05 -> $0.5
			{"EUR 3/40 scale up 1", NewMoneyFromFraction(3, 40, "EUR"), 1, NewMoneyFromFraction(30, 40, "EUR")},     // €0.075 -> €0.75
			{"USD 5/100 scale up 2", NewMoneyFromFraction(5, 100, "USD"), 2, NewMoneyFromFraction(100, 20, "USD")},  // $0.05 -> $5 (5/100 reduces to 1/20, then 1*100/20 = 100/20 = 5)
			{"USD negative scale up", NewMoneyFromFraction(-7, 30, "USD"), 1, NewMoneyFromFraction(-70, 30, "USD")}, // -$0.233... -> -$2.33...
			{"USD zero scale up", ZeroMoney("USD"), 3, ZeroMoney("USD")},                                            // $0 -> $0 (num stays 0)
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				m := tt.money
				err := m.ScaleUp(tt.scale)
				require.NoError(t, err, "ScaleUp should not return error for valid money")
				assert.True(t, m.IsValid(), "Money should remain valid")
				assert.Equal(t, tt.expected.Currency(), m.Currency(), "Currency should be preserved")
				assert.True(t, m.Equal(tt.expected), "Amount should match expected after scaling up")
			})
		}
	})

	t.Run("zero scale does nothing", func(t *testing.T) {
		m := NewMoneyFromFraction(3, 7, "USD")
		original := m
		err := m.ScaleUp(0)
		require.NoError(t, err)
		assert.Equal(t, original.Currency(), m.Currency(), "Currency should be unchanged")
		assert.True(t, m.Equal(original), "Money should be unchanged with zero scale")
	})

	t.Run("negative scale calls ScaleDown", func(t *testing.T) {
		m1 := NewMoneyFromFraction(1, 2, "USD")
		m2 := NewMoneyFromFraction(1, 2, "USD")

		err1 := m1.ScaleUp(-2) // Should be equivalent to ScaleDown(2)
		err2 := m2.ScaleDown(2)

		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.True(t, m1.Equal(m2), "ScaleUp(-2) should equal ScaleDown(2)")
	})

	t.Run("invalid state handling", func(t *testing.T) {
		m := NewMoneyInt("", 100) // Invalid currency
		err := m.ScaleUp(2)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid(), "invalid input should remain invalid")
	})

	t.Run("overflow handling", func(t *testing.T) {
		// Test numerator overflow - create Money with large numerator
		m := NewMoneyInt("USD", math.MaxInt64/5) // Large numerator
		err := m.ScaleUp(10)                     // This should cause overflow
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid(), "overflow should result in invalid state")
	})
}

// TestMoney_ScaledUpErr tests the immutable ScaledUpErr method
func TestMoney_ScaledUpErr(t *testing.T) {
	t.Run("basic scaling up immutable with error", func(t *testing.T) {
		tests := []struct {
			name     string
			money    Money
			scale    int
			expected Money
		}{
			{"USD 1/200 scaled up 2", NewMoneyFromFraction(1, 200, "USD"), 2, NewMoneyFromFraction(100, 200, "USD")},
			{"USD 1/20 scaled up 1", NewMoneyFromFraction(1, 20, "USD"), 1, NewMoneyFromFraction(10, 20, "USD")},
			{"EUR 3/40 scaled up 1", NewMoneyFromFraction(3, 40, "EUR"), 1, NewMoneyFromFraction(30, 40, "EUR")},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.money
				result, err := tt.money.ScaledUpErr(tt.scale)

				require.NoError(t, err, "ScaledUpErr should not return error for valid money")

				// Check original is unchanged
				assert.Equal(t, original.Currency(), tt.money.Currency(), "original currency should be unchanged")
				assert.True(t, tt.money.Equal(original), "original money should be unchanged")

				// Check result
				assert.True(t, result.IsValid(), "result should be valid")
				assert.Equal(t, tt.expected.Currency(), result.Currency(), "result currency should match expected")
				assert.True(t, result.Equal(tt.expected), "result should match expected after scaling up")
			})
		}
	})

	t.Run("invalid money returns error", func(t *testing.T) {
		m := NewMoneyInt("", 100) // Invalid currency
		result, err := m.ScaledUpErr(2)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid(), "result should be invalid")
	})
}

// TestMoney_ScaledUp tests the immutable ScaledUp method (without error)
func TestMoney_ScaledUp(t *testing.T) {
	t.Run("basic scaling up immutable without error", func(t *testing.T) {
		m := NewMoneyFromFraction(1, 200, "USD")
		original := m
		result := m.ScaledUp(2)

		// Check original is unchanged
		assert.True(t, m.Equal(original), "original money should be unchanged")

		// Check result
		expected := NewMoneyFromFraction(100, 200, "USD")
		assert.True(t, result.Equal(expected), "result should match expected")
	})

	t.Run("invalid money returns invalid", func(t *testing.T) {
		m := NewMoneyInt("", 100) // Invalid currency
		result := m.ScaledUp(2)
		assert.True(t, result.IsInvalid(), "result should be invalid for invalid input")
	})
}
