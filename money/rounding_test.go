package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoney_Round_InvalidState tests that invalid Money remains invalid after rounding
func TestMoney_Round_InvalidState(t *testing.T) {
	tests := []struct {
		name      string
		money     Money
		roundType zerorat.RoundType
		scale     int
	}{
		{"invalid money with RoundDown", Money{}, zerorat.RoundDown, 0},
		{"invalid money with RoundUp", Money{}, zerorat.RoundUp, 0},
		{"invalid money with RoundHalfUp", Money{}, zerorat.RoundHalfUp, 0},
		{"empty currency with positive scale", NewMoneyInt("", 100), zerorat.RoundDown, 2},
		{"invalid amount with negative scale", NewMoney("USD", zerorat.New(1, 0)), zerorat.RoundUp, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test mutable Round
			m := tt.money
			err := m.Round(tt.roundType, tt.scale)
			require.Error(t, err, "Round should return error for invalid money")
			assert.Equal(t, ErrMoneyInvalid, err, "Should return ErrMoneyInvalid")
			assert.True(t, m.IsInvalid(), "Money should remain invalid")

			// Test immutable RoundedErr
			result, err := tt.money.RoundedErr(tt.roundType, tt.scale)
			require.Error(t, err, "RoundedErr should return error for invalid money")
			assert.Equal(t, ErrMoneyInvalid, err, "Should return ErrMoneyInvalid")
			assert.True(t, result.IsInvalid(), "Result should be invalid")

			// Test immutable Rounded (no error)
			result = tt.money.Rounded(tt.roundType, tt.scale)
			assert.True(t, result.IsInvalid(), "Result should be invalid")
		})
	}
}

// TestMoney_Round_BasicRounding tests basic rounding functionality
func TestMoney_Round_BasicRounding(t *testing.T) {
	tests := []struct {
		name      string
		money     Money
		roundType zerorat.RoundType
		scale     int
		expected  Money
	}{
		// RoundDown (toward zero)
		{"positive RoundDown to integer", NewMoneyFromFraction(123, 100, "USD"), zerorat.RoundDown, 0, NewMoneyInt("USD", 1)},
		{"negative RoundDown to integer", NewMoneyFromFraction(-123, 100, "USD"), zerorat.RoundDown, 0, NewMoneyInt("USD", -1)},
		{"positive RoundDown to 1 decimal", NewMoneyFromFraction(1234, 1000, "EUR"), zerorat.RoundDown, 1, NewMoneyFromFraction(12, 10, "EUR")},
		{"negative RoundDown to 1 decimal", NewMoneyFromFraction(-1234, 1000, "EUR"), zerorat.RoundDown, 1, NewMoneyFromFraction(-12, 10, "EUR")},

		// RoundUp (away from zero)
		{"positive RoundUp to integer", NewMoneyFromFraction(123, 100, "USD"), zerorat.RoundUp, 0, NewMoneyInt("USD", 2)},
		{"negative RoundUp to integer", NewMoneyFromFraction(-123, 100, "USD"), zerorat.RoundUp, 0, NewMoneyInt("USD", -2)},
		{"positive RoundUp to 1 decimal", NewMoneyFromFraction(1234, 1000, "EUR"), zerorat.RoundUp, 1, NewMoneyFromFraction(13, 10, "EUR")},
		{"negative RoundUp to 1 decimal", NewMoneyFromFraction(-1234, 1000, "EUR"), zerorat.RoundUp, 1, NewMoneyFromFraction(-13, 10, "EUR")},

		// RoundHalfUp (financial rounding)
		{"positive half RoundHalfUp", NewMoneyFromFraction(25, 10, "JPY"), zerorat.RoundHalfUp, 0, NewMoneyInt("JPY", 3)},
		{"negative half RoundHalfUp", NewMoneyFromFraction(-25, 10, "JPY"), zerorat.RoundHalfUp, 0, NewMoneyInt("JPY", -2)},
		{"positive less than half", NewMoneyFromFraction(23, 10, "JPY"), zerorat.RoundHalfUp, 0, NewMoneyInt("JPY", 2)},
		{"negative less than half", NewMoneyFromFraction(-24, 10, "JPY"), zerorat.RoundHalfUp, 0, NewMoneyInt("JPY", -2)},

		// Zero cases
		{"zero money", ZeroMoney("GBP"), zerorat.RoundDown, 0, ZeroMoney("GBP")},
		{"zero with scale", ZeroMoney("GBP"), zerorat.RoundUp, 2, ZeroMoney("GBP")},

		// Negative scale (powers of 10)
		{"round to tens", NewMoneyInt("USD", 1234), zerorat.RoundDown, -1, NewMoneyInt("USD", 1230)},
		{"round to hundreds", NewMoneyInt("USD", 1234), zerorat.RoundUp, -2, NewMoneyInt("USD", 1300)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test mutable Round
			m := tt.money
			err := m.Round(tt.roundType, tt.scale)
			require.NoError(t, err, "Round should not return error for valid money")
			assert.True(t, m.IsValid(), "Money should remain valid")
			assert.Equal(t, tt.expected.Currency(), m.Currency(), "Currency should be preserved")
			assert.True(t, m.Equal(tt.expected), "Amount should match expected after rounding")

			// Test immutable RoundedErr
			result, err := tt.money.RoundedErr(tt.roundType, tt.scale)
			require.NoError(t, err, "RoundedErr should not return error for valid money")
			assert.True(t, result.IsValid(), "Result should be valid")
			assert.Equal(t, tt.expected.Currency(), result.Currency(), "Currency should be preserved")
			assert.True(t, result.Equal(tt.expected), "Amount should match expected")

			// Test immutable Rounded (no error)
			result = tt.money.Rounded(tt.roundType, tt.scale)
			assert.True(t, result.IsValid(), "Result should be valid")
			assert.Equal(t, tt.expected.Currency(), result.Currency(), "Currency should be preserved")
			assert.True(t, result.Equal(tt.expected), "Amount should match expected")

			// Verify original money is unchanged for immutable operations
			assert.True(t, tt.money.Equal(tt.money), "Original money should be unchanged")
		})
	}
}

// TestMoney_Ceil_BasicFunctionality tests Ceil operations
func TestMoney_Ceil_BasicFunctionality(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		scale    int
		expected Money
	}{
		// Basic ceiling operations
		{"positive to integer", NewMoneyFromFraction(123, 100, "USD"), 0, NewMoneyInt("USD", 2)},
		{"negative to integer", NewMoneyFromFraction(-123, 100, "USD"), 0, NewMoneyInt("USD", -1)},
		{"positive to 1 decimal", NewMoneyFromFraction(1234, 1000, "EUR"), 1, NewMoneyFromFraction(13, 10, "EUR")},
		{"negative to 1 decimal", NewMoneyFromFraction(-1234, 1000, "EUR"), 1, NewMoneyFromFraction(-12, 10, "EUR")},
		{"already integer", NewMoneyInt("GBP", 5), 0, NewMoneyInt("GBP", 5)},
		{"zero", ZeroMoney("JPY"), 0, ZeroMoney("JPY")},
		{"negative scale", NewMoneyInt("USD", 1234), -2, NewMoneyInt("USD", 1300)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test mutable Ceil
			m := tt.money
			err := m.Ceil(tt.scale)
			require.NoError(t, err, "Ceil should not return error for valid money")
			assert.True(t, m.IsValid(), "Money should remain valid")
			assert.Equal(t, tt.expected.Currency(), m.Currency(), "Currency should be preserved")
			assert.True(t, m.Equal(tt.expected), "Amount should match expected after ceiling")

			// Test immutable CeiledErr
			result, err := tt.money.CeiledErr(tt.scale)
			require.NoError(t, err, "CeiledErr should not return error for valid money")
			assert.True(t, result.IsValid(), "Result should be valid")
			assert.Equal(t, tt.expected.Currency(), result.Currency(), "Currency should be preserved")
			assert.True(t, result.Equal(tt.expected), "Amount should match expected")

			// Test immutable Ceiled (no error)
			result = tt.money.Ceiled(tt.scale)
			assert.True(t, result.IsValid(), "Result should be valid")
			assert.Equal(t, tt.expected.Currency(), result.Currency(), "Currency should be preserved")
			assert.True(t, result.Equal(tt.expected), "Amount should match expected")
		})
	}
}

// TestMoney_Floor_BasicFunctionality tests Floor operations
func TestMoney_Floor_BasicFunctionality(t *testing.T) {
	tests := []struct {
		name     string
		money    Money
		scale    int
		expected Money
	}{
		// Basic floor operations
		{"positive to integer", NewMoneyFromFraction(123, 100, "USD"), 0, NewMoneyInt("USD", 1)},
		{"negative to integer", NewMoneyFromFraction(-123, 100, "USD"), 0, NewMoneyInt("USD", -2)},
		{"positive to 1 decimal", NewMoneyFromFraction(1234, 1000, "EUR"), 1, NewMoneyFromFraction(12, 10, "EUR")},
		{"negative to 1 decimal", NewMoneyFromFraction(-1234, 1000, "EUR"), 1, NewMoneyFromFraction(-13, 10, "EUR")},
		{"already integer", NewMoneyInt("GBP", 5), 0, NewMoneyInt("GBP", 5)},
		{"zero", ZeroMoney("JPY"), 0, ZeroMoney("JPY")},
		{"negative scale", NewMoneyInt("USD", 1234), -2, NewMoneyInt("USD", 1200)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test mutable Floor
			m := tt.money
			err := m.Floor(tt.scale)
			require.NoError(t, err, "Floor should not return error for valid money")
			assert.True(t, m.IsValid(), "Money should remain valid")
			assert.Equal(t, tt.expected.Currency(), m.Currency(), "Currency should be preserved")
			assert.True(t, m.Equal(tt.expected), "Amount should match expected after flooring")

			// Test immutable FlooredErr
			result, err := tt.money.FlooredErr(tt.scale)
			require.NoError(t, err, "FlooredErr should not return error for valid money")
			assert.True(t, result.IsValid(), "Result should be valid")
			assert.Equal(t, tt.expected.Currency(), result.Currency(), "Currency should be preserved")
			assert.True(t, result.Equal(tt.expected), "Amount should match expected")

			// Test immutable Floored (no error)
			result = tt.money.Floored(tt.scale)
			assert.True(t, result.IsValid(), "Result should be valid")
			assert.Equal(t, tt.expected.Currency(), result.Currency(), "Currency should be preserved")
			assert.True(t, result.Equal(tt.expected), "Amount should match expected")
		})
	}
}

// TestMoney_Ceil_InvalidState tests Ceil with invalid Money
func TestMoney_Ceil_InvalidState(t *testing.T) {
	tests := []struct {
		name  string
		money Money
		scale int
	}{
		{"invalid money", Money{}, 0},
		{"empty currency", NewMoneyInt("", 100), 0},
		{"invalid amount", NewMoney("USD", zerorat.New(1, 0)), 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test mutable Ceil
			m := tt.money
			err := m.Ceil(tt.scale)
			require.Error(t, err, "Ceil should return error for invalid money")
			assert.Equal(t, ErrMoneyInvalid, err, "Should return ErrMoneyInvalid")
			assert.True(t, m.IsInvalid(), "Money should remain invalid")

			// Test immutable CeiledErr
			result, err := tt.money.CeiledErr(tt.scale)
			require.Error(t, err, "CeiledErr should return error for invalid money")
			assert.Equal(t, ErrMoneyInvalid, err, "Should return ErrMoneyInvalid")
			assert.True(t, result.IsInvalid(), "Result should be invalid")

			// Test immutable Ceiled (no error)
			result = tt.money.Ceiled(tt.scale)
			assert.True(t, result.IsInvalid(), "Result should be invalid")
		})
	}
}

// TestMoney_Floor_InvalidState tests Floor with invalid Money
func TestMoney_Floor_InvalidState(t *testing.T) {
	tests := []struct {
		name  string
		money Money
		scale int
	}{
		{"invalid money", Money{}, 0},
		{"empty currency", NewMoneyInt("", 100), 0},
		{"invalid amount", NewMoney("USD", zerorat.New(1, 0)), 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test mutable Floor
			m := tt.money
			err := m.Floor(tt.scale)
			require.Error(t, err, "Floor should return error for invalid money")
			assert.Equal(t, ErrMoneyInvalid, err, "Should return ErrMoneyInvalid")
			assert.True(t, m.IsInvalid(), "Money should remain invalid")

			// Test immutable FlooredErr
			result, err := tt.money.FlooredErr(tt.scale)
			require.Error(t, err, "FlooredErr should return error for invalid money")
			assert.Equal(t, ErrMoneyInvalid, err, "Should return ErrMoneyInvalid")
			assert.True(t, result.IsInvalid(), "Result should be invalid")

			// Test immutable Floored (no error)
			result = tt.money.Floored(tt.scale)
			assert.True(t, result.IsInvalid(), "Result should be invalid")
		})
	}
}

// TestMoney_Rounding_EdgeCases tests edge cases for all rounding operations
func TestMoney_Rounding_EdgeCases(t *testing.T) {
	t.Run("large scale values", func(t *testing.T) {
		money := NewMoneyFromFraction(1, 3, "USD") // 0.333...

		// Test with large positive scale
		result := money.Rounded(zerorat.RoundDown, 10)
		assert.True(t, result.IsValid(), "Should handle large positive scale")
		assert.Equal(t, "USD", result.Currency(), "Currency should be preserved")

		// Test with large negative scale
		bigMoney := NewMoneyInt("USD", 12345)
		result = bigMoney.Rounded(zerorat.RoundDown, -10)
		assert.True(t, result.IsValid(), "Should handle large negative scale")
		assert.Equal(t, "USD", result.Currency(), "Currency should be preserved")
	})

	t.Run("precision preservation", func(t *testing.T) {
		// Test that rounding preserves exact values when no rounding is needed
		exactMoney := NewMoneyFromFraction(5, 2, "EUR") // 2.5

		// Round to 1 decimal place - should remain exact
		result := exactMoney.Rounded(zerorat.RoundDown, 1)
		expected := NewMoneyFromFraction(25, 10, "EUR") // 2.5
		assert.True(t, result.Equal(expected), "Should preserve exact values")
	})

	t.Run("currency preservation across all operations", func(t *testing.T) {
		currencies := []string{"USD", "EUR", "JPY", "GBP", "CHF"}
		money := NewMoneyFromFraction(123, 100, "")

		for _, currency := range currencies {
			money.currency = currency
			money.amount = zerorat.New(123, 100)

			// Test Round
			result := money.Rounded(zerorat.RoundHalfUp, 0)
			assert.Equal(t, currency, result.Currency(), "Round should preserve currency: %s", currency)

			// Test Ceil
			result = money.Ceiled(0)
			assert.Equal(t, currency, result.Currency(), "Ceil should preserve currency: %s", currency)

			// Test Floor
			result = money.Floored(0)
			assert.Equal(t, currency, result.Currency(), "Floor should preserve currency: %s", currency)
		}
	})
}

// TestMoney_Rounding_ConsistencyWithRat tests consistency with underlying zerorat.Rat behavior
func TestMoney_Rounding_ConsistencyWithRat(t *testing.T) {
	tests := []struct {
		name      string
		numerator int64
		denom     uint64
		roundType zerorat.RoundType
		scale     int
	}{
		{"half up positive", 25, 10, zerorat.RoundHalfUp, 0},
		{"half up negative", -25, 10, zerorat.RoundHalfUp, 0},
		{"round down positive", 27, 10, zerorat.RoundDown, 0},
		{"round down negative", -27, 10, zerorat.RoundDown, 0},
		{"round up positive", 23, 10, zerorat.RoundUp, 0},
		{"round up negative", -23, 10, zerorat.RoundUp, 0},
		{"decimal places", 1234, 1000, zerorat.RoundHalfUp, 2},
		{"negative scale", 1234, 1, zerorat.RoundDown, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Money and equivalent Rat
			money := NewMoneyFromFraction(tt.numerator, tt.denom, "USD")
			rat := zerorat.New(tt.numerator, tt.denom)

			// Round both
			roundedMoney := money.Rounded(tt.roundType, tt.scale)
			rat.Round(tt.roundType, tt.scale)

			// They should produce the same result
			assert.True(t, roundedMoney.Amount().Equal(rat),
				"Money rounding should match Rat rounding for %s", tt.name)
		})
	}
}
