package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoney_String tests the String method
func TestMoney_String(t *testing.T) {
	t.Run("valid positive money", func(t *testing.T) {
		// USD/123/100 represents $1.23
		money := NewMoneyFromFraction(123, 100, "USD")
		expected := "USD/123/100"

		assert.Equal(t, expected, money.String(), "Positive money should format correctly")
	})

	t.Run("valid negative money", func(t *testing.T) {
		// EUR/-5/2 represents -2.5 EUR
		money := NewMoneyFromFraction(-5, 2, "EUR")
		expected := "EUR/-5/2"

		assert.Equal(t, expected, money.String(), "Negative money should format correctly")
	})

	t.Run("zero money", func(t *testing.T) {
		// JPY/0 represents 0 JPY (zerorat.Rat.String() returns "0" for zero)
		money := ZeroMoney("JPY")
		expected := "JPY/0"

		assert.Equal(t, expected, money.String(), "Zero money should format correctly")
	})

	t.Run("integer money", func(t *testing.T) {
		// When denominator is 1, zerorat.Rat.String() returns just the numerator
		money := NewMoneyInt("GBP", 42)
		expected := "GBP/42"

		assert.Equal(t, expected, money.String(), "Integer money should format correctly")
	})

	t.Run("invalid money", func(t *testing.T) {
		invalid := Money{}
		expected := "invalid"

		assert.Equal(t, expected, invalid.String(), "Invalid money should return 'invalid'")
	})

	t.Run("money with invalid amount", func(t *testing.T) {
		// Create money with invalid amount (denominator = 0)
		money := NewMoney("USD", zerorat.New(1, 0))
		expected := "invalid"

		assert.Equal(t, expected, money.String(), "Money with invalid amount should return 'invalid'")
	})

	t.Run("money with empty currency", func(t *testing.T) {
		money := NewMoneyInt("", 100)
		expected := "invalid"

		assert.Equal(t, expected, money.String(), "Money with empty currency should return 'invalid'")
	})

	t.Run("reduced fractions", func(t *testing.T) {
		// 6/4 should reduce to 3/2
		money := NewMoneyFromFraction(6, 4, "CHF")
		expected := "CHF/3/2"

		assert.Equal(t, expected, money.String(), "Reduced fractions should format correctly")
	})
}

// TestParseMoney tests the ParseMoney function
func TestParseMoney(t *testing.T) {
	t.Run("valid positive money", func(t *testing.T) {
		input := "USD/123/100"
		money, err := ParseMoney(input)

		require.NoError(t, err)
		assert.True(t, money.IsValid(), "Parsed money should be valid")
		assert.Equal(t, "USD", money.Currency(), "Currency should match")

		expected := NewMoneyFromFraction(123, 100, "USD")
		assert.True(t, money.Equal(expected), "Parsed money should equal expected")
	})

	t.Run("valid negative money", func(t *testing.T) {
		input := "EUR/-5/2"
		money, err := ParseMoney(input)

		require.NoError(t, err)
		assert.True(t, money.IsValid(), "Parsed money should be valid")
		assert.Equal(t, "EUR", money.Currency(), "Currency should match")

		expected := NewMoneyFromFraction(-5, 2, "EUR")
		assert.True(t, money.Equal(expected), "Parsed money should equal expected")
	})

	t.Run("zero money", func(t *testing.T) {
		input := "JPY/0/1"
		money, err := ParseMoney(input)

		require.NoError(t, err)
		assert.True(t, money.IsValid(), "Parsed money should be valid")
		assert.Equal(t, "JPY", money.Currency(), "Currency should match")

		expected := ZeroMoney("JPY")
		assert.True(t, money.Equal(expected), "Parsed money should equal expected")
	})

	t.Run("integer money", func(t *testing.T) {
		input := "GBP/42"
		money, err := ParseMoney(input)

		require.NoError(t, err)
		assert.True(t, money.IsValid(), "Parsed money should be valid")
		assert.Equal(t, "GBP", money.Currency(), "Currency should match")

		expected := NewMoneyInt("GBP", 42)
		assert.True(t, money.Equal(expected), "Parsed money should equal expected")
	})

	t.Run("invalid format - no slashes", func(t *testing.T) {
		input := "USD100"
		money, err := ParseMoney(input)

		require.Error(t, err)
		assert.True(t, money.IsInvalid(), "Parsed money should be invalid")
	})

	t.Run("invalid format - empty currency", func(t *testing.T) {
		input := "/123/100"
		money, err := ParseMoney(input)

		require.Error(t, err)
		assert.True(t, money.IsInvalid(), "Parsed money should be invalid")
	})

	t.Run("invalid format - invalid fraction", func(t *testing.T) {
		input := "USD/123/0"
		money, err := ParseMoney(input)

		require.Error(t, err)
		assert.True(t, money.IsInvalid(), "Parsed money should be invalid")
	})

	t.Run("invalid format - non-numeric numerator", func(t *testing.T) {
		input := "USD/abc/100"
		money, err := ParseMoney(input)

		require.Error(t, err)
		assert.True(t, money.IsInvalid(), "Parsed money should be invalid")
	})

	t.Run("invalid format - non-numeric denominator", func(t *testing.T) {
		input := "USD/123/abc"
		money, err := ParseMoney(input)

		require.Error(t, err)
		assert.True(t, money.IsInvalid(), "Parsed money should be invalid")
	})

	t.Run("invalid format - too many parts", func(t *testing.T) {
		input := "USD/123/100/extra"
		money, err := ParseMoney(input)

		require.Error(t, err)
		assert.True(t, money.IsInvalid(), "Parsed money should be invalid")
	})

	t.Run("integer format - currency/numerator", func(t *testing.T) {
		input := "USD/123"
		money, err := ParseMoney(input)

		require.NoError(t, err)
		assert.True(t, money.IsValid(), "Parsed money should be valid")

		expected := NewMoneyInt("USD", 123)
		assert.True(t, money.Equal(expected), "Parsed money should equal expected")
	})

	t.Run("empty string", func(t *testing.T) {
		input := ""
		money, err := ParseMoney(input)

		require.Error(t, err)
		assert.True(t, money.IsInvalid(), "Parsed money should be invalid")
	})

	t.Run("invalid string", func(t *testing.T) {
		input := "invalid"
		money, err := ParseMoney(input)

		require.Error(t, err)
		assert.True(t, money.IsInvalid(), "Parsed money should be invalid")
	})

	t.Run("round trip consistency", func(t *testing.T) {
		testCases := []Money{
			NewMoneyFromFraction(123, 100, "USD"),
			NewMoneyFromFraction(-5, 2, "EUR"),
			ZeroMoney("JPY"),
			NewMoneyInt("GBP", 42),
			NewMoneyFromFraction(3, 2, "CHF"),
		}

		for _, original := range testCases {
			str := original.String()
			parsed, err := ParseMoney(str)

			require.NoError(t, err, "Parsing should succeed for: %s", str)
			assert.True(t, original.Equal(parsed), "Round trip should preserve value for: %s", str)
		}
	})
}
