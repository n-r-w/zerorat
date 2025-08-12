package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoney_Compare tests the Compare method
func TestMoney_Compare(t *testing.T) {
	t.Run("same currency comparisons", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		usd200 := NewMoneyInt("USD", 200)
		usd100_copy := NewMoneyInt("USD", 100)

		// Less than
		assert.Equal(t, -1, usd100.Compare(usd200), "100 < 200")

		// Greater than
		assert.Equal(t, 1, usd200.Compare(usd100), "200 > 100")

		// Equal
		assert.Equal(t, 0, usd100.Compare(usd100_copy), "100 == 100")
	})

	t.Run("different currency returns 0", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		eur100 := NewMoneyInt("EUR", 100)

		assert.Equal(t, 0, usd100.Compare(eur100), "Different currencies should return 0")
	})

	t.Run("invalid money returns 0", func(t *testing.T) {
		valid := NewMoneyInt("USD", 100)
		invalid := Money{}

		assert.Equal(t, 0, valid.Compare(invalid), "Valid vs invalid should return 0")
		assert.Equal(t, 0, invalid.Compare(valid), "Invalid vs valid should return 0")
		assert.Equal(t, 0, invalid.Compare(invalid), "Invalid vs invalid should return 0")
	})

	t.Run("zero and negative values", func(t *testing.T) {
		usdZero := ZeroMoney("USD")
		usdPositive := NewMoneyInt("USD", 100)
		usdNegative := NewMoneyInt("USD", -100)

		assert.Equal(t, -1, usdZero.Compare(usdPositive), "0 < 100")
		assert.Equal(t, 1, usdZero.Compare(usdNegative), "0 > -100")
		assert.Equal(t, -1, usdNegative.Compare(usdPositive), "-100 < 100")
	})

	t.Run("fractional comparisons", func(t *testing.T) {
		usd_half := NewMoneyFromFraction(1, 2, "USD")    // 0.5
		usd_quarter := NewMoneyFromFraction(1, 4, "USD") // 0.25

		assert.Equal(t, 1, usd_half.Compare(usd_quarter), "1/2 > 1/4")
		assert.Equal(t, -1, usd_quarter.Compare(usd_half), "1/4 < 1/2")
	})
}

// TestMoney_CompareErr tests the CompareErr method
func TestMoney_CompareErr(t *testing.T) {
	t.Run("same currency comparisons", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		usd200 := NewMoneyInt("USD", 200)

		result, err := usd100.CompareErr(usd200)
		require.NoError(t, err)
		assert.Equal(t, -1, result)
	})

	t.Run("different currency returns error", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		eur100 := NewMoneyInt("EUR", 100)

		result, err := usd100.CompareErr(eur100)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.Equal(t, 0, result)
	})

	t.Run("invalid money returns error", func(t *testing.T) {
		valid := NewMoneyInt("USD", 100)
		invalid := Money{}

		result, err := valid.CompareErr(invalid)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.Equal(t, 0, result)
	})
}

// TestMoney_Equal tests the Equal method
func TestMoney_Equal(t *testing.T) {
	t.Run("same currency equal values", func(t *testing.T) {
		usd100_1 := NewMoneyInt("USD", 100)
		usd100_2 := NewMoneyInt("USD", 100)

		assert.True(t, usd100_1.Equal(usd100_2), "Equal values should return true")
	})

	t.Run("same currency different values", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		usd200 := NewMoneyInt("USD", 200)

		assert.False(t, usd100.Equal(usd200), "Different values should return false")
	})

	t.Run("different currency returns false", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		eur100 := NewMoneyInt("EUR", 100)

		assert.False(t, usd100.Equal(eur100), "Different currencies should return false")
	})

	t.Run("invalid money returns false", func(t *testing.T) {
		valid := NewMoneyInt("USD", 100)
		invalid := Money{}

		assert.False(t, valid.Equal(invalid), "Valid vs invalid should return false")
		assert.False(t, invalid.Equal(valid), "Invalid vs valid should return false")
		assert.False(t, invalid.Equal(invalid), "Invalid vs invalid should return false")
	})

	t.Run("fractional equality", func(t *testing.T) {
		usd_half_1 := NewMoneyFromFraction(1, 2, "USD")
		usd_half_2 := NewMoneyFromFraction(2, 4, "USD") // Should reduce to 1/2

		assert.True(t, usd_half_1.Equal(usd_half_2), "Equivalent fractions should be equal")
	})
}

// TestMoney_EqualErr tests the EqualErr method
func TestMoney_EqualErr(t *testing.T) {
	t.Run("same currency equal values", func(t *testing.T) {
		usd100_1 := NewMoneyInt("USD", 100)
		usd100_2 := NewMoneyInt("USD", 100)

		result, err := usd100_1.EqualErr(usd100_2)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("different currency returns error", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		eur100 := NewMoneyInt("EUR", 100)

		result, err := usd100.EqualErr(eur100)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.False(t, result)
	})

	t.Run("invalid money returns error", func(t *testing.T) {
		valid := NewMoneyInt("USD", 100)
		invalid := Money{}

		result, err := valid.EqualErr(invalid)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.False(t, result)
	})
}

// TestMoney_Less tests the Less method
func TestMoney_Less(t *testing.T) {
	t.Run("same currency less than", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		usd200 := NewMoneyInt("USD", 200)

		assert.True(t, usd100.Less(usd200), "100 < 200")
		assert.False(t, usd200.Less(usd100), "200 not < 100")
		assert.False(t, usd100.Less(usd100), "100 not < 100")
	})

	t.Run("different currency returns false", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		eur200 := NewMoneyInt("EUR", 200)

		assert.False(t, usd100.Less(eur200), "Different currencies should return false")
	})

	t.Run("invalid money returns false", func(t *testing.T) {
		valid := NewMoneyInt("USD", 100)
		invalid := Money{}

		assert.False(t, valid.Less(invalid), "Valid vs invalid should return false")
		assert.False(t, invalid.Less(valid), "Invalid vs valid should return false")
	})

	t.Run("negative values", func(t *testing.T) {
		usdNeg100 := NewMoneyInt("USD", -100)
		usdPos100 := NewMoneyInt("USD", 100)
		usdZero := ZeroMoney("USD")

		assert.True(t, usdNeg100.Less(usdZero), "-100 < 0")
		assert.True(t, usdNeg100.Less(usdPos100), "-100 < 100")
		assert.True(t, usdZero.Less(usdPos100), "0 < 100")
	})
}

// TestMoney_LessErr tests the LessErr method
func TestMoney_LessErr(t *testing.T) {
	t.Run("same currency less than", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		usd200 := NewMoneyInt("USD", 200)

		result, err := usd100.LessErr(usd200)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("different currency returns error", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		eur200 := NewMoneyInt("EUR", 200)

		result, err := usd100.LessErr(eur200)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.False(t, result)
	})

	t.Run("invalid money returns error", func(t *testing.T) {
		valid := NewMoneyInt("USD", 100)
		invalid := Money{}

		result, err := valid.LessErr(invalid)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.False(t, result)
	})
}

// TestMoney_Greater tests the Greater method
func TestMoney_Greater(t *testing.T) {
	t.Run("same currency greater than", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		usd200 := NewMoneyInt("USD", 200)

		assert.True(t, usd200.Greater(usd100), "200 > 100")
		assert.False(t, usd100.Greater(usd200), "100 not > 200")
		assert.False(t, usd100.Greater(usd100), "100 not > 100")
	})

	t.Run("different currency returns false", func(t *testing.T) {
		usd200 := NewMoneyInt("USD", 200)
		eur100 := NewMoneyInt("EUR", 100)

		assert.False(t, usd200.Greater(eur100), "Different currencies should return false")
	})

	t.Run("invalid money returns false", func(t *testing.T) {
		valid := NewMoneyInt("USD", 100)
		invalid := Money{}

		assert.False(t, valid.Greater(invalid), "Valid vs invalid should return false")
		assert.False(t, invalid.Greater(valid), "Invalid vs valid should return false")
	})

	t.Run("negative values", func(t *testing.T) {
		usdNeg100 := NewMoneyInt("USD", -100)
		usdPos100 := NewMoneyInt("USD", 100)
		usdZero := ZeroMoney("USD")

		assert.True(t, usdZero.Greater(usdNeg100), "0 > -100")
		assert.True(t, usdPos100.Greater(usdNeg100), "100 > -100")
		assert.True(t, usdPos100.Greater(usdZero), "100 > 0")
	})
}

// TestMoney_GreaterErr tests the GreaterErr method
func TestMoney_GreaterErr(t *testing.T) {
	t.Run("same currency greater than", func(t *testing.T) {
		usd100 := NewMoneyInt("USD", 100)
		usd200 := NewMoneyInt("USD", 200)

		result, err := usd200.GreaterErr(usd100)
		require.NoError(t, err)
		assert.True(t, result)
	})

	t.Run("different currency returns error", func(t *testing.T) {
		usd200 := NewMoneyInt("USD", 200)
		eur100 := NewMoneyInt("EUR", 100)

		result, err := usd200.GreaterErr(eur100)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.False(t, result)
	})

	t.Run("invalid money returns error", func(t *testing.T) {
		valid := NewMoneyInt("USD", 100)
		invalid := Money{}

		result, err := valid.GreaterErr(invalid)
		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.False(t, result)
	})
}

// TestMoney_ComparisonEdgeCases tests edge cases for all comparison methods
func TestMoney_ComparisonEdgeCases(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		usdZero1 := ZeroMoney("USD")
		usdZero2 := ZeroMoney("USD")

		assert.Equal(t, 0, usdZero1.Compare(usdZero2), "Zero should equal zero")
		assert.True(t, usdZero1.Equal(usdZero2), "Zero should equal zero")
		assert.False(t, usdZero1.Less(usdZero2), "Zero not less than zero")
		assert.False(t, usdZero1.Greater(usdZero2), "Zero not greater than zero")
	})

	t.Run("fractional values", func(t *testing.T) {
		usd_half := NewMoneyFromFraction(1, 2, "USD")           // 0.5
		usd_quarter := NewMoneyFromFraction(1, 4, "USD")        // 0.25
		usd_three_quarters := NewMoneyFromFraction(3, 4, "USD") // 0.75

		assert.Equal(t, 1, usd_half.Compare(usd_quarter), "1/2 > 1/4")
		assert.Equal(t, -1, usd_half.Compare(usd_three_quarters), "1/2 < 3/4")
		assert.True(t, usd_half.Greater(usd_quarter), "1/2 > 1/4")
		assert.True(t, usd_quarter.Less(usd_half), "1/4 < 1/2")
	})

	t.Run("equivalent fractions", func(t *testing.T) {
		usd_half_1 := NewMoneyFromFraction(1, 2, "USD")
		usd_half_2 := NewMoneyFromFraction(2, 4, "USD") // Should reduce to 1/2

		assert.Equal(t, 0, usd_half_1.Compare(usd_half_2), "Equivalent fractions should be equal")
		assert.True(t, usd_half_1.Equal(usd_half_2), "Equivalent fractions should be equal")
		assert.False(t, usd_half_1.Less(usd_half_2), "Equivalent fractions not less")
		assert.False(t, usd_half_1.Greater(usd_half_2), "Equivalent fractions not greater")
	})

	t.Run("currency preservation", func(t *testing.T) {
		currencies := []string{"USD", "EUR", "JPY", "GBP", "CHF"}

		for _, currency := range currencies {
			money1 := NewMoneyInt(currency, 100)
			money2 := NewMoneyInt(currency, 200)

			// All comparison methods should work with any valid currency
			assert.Equal(t, -1, money1.Compare(money2), "Comparison should work for currency: %s", currency)
			assert.True(t, money1.Less(money2), "Less should work for currency: %s", currency)
			assert.True(t, money2.Greater(money1), "Greater should work for currency: %s", currency)
			assert.False(t, money1.Equal(money2), "Equal should work for currency: %s", currency)
		}
	})
}

// TestMoney_IsZero tests the IsZero method
func TestMoney_IsZero(t *testing.T) {
	t.Run("valid zero money", func(t *testing.T) {
		t.Run("zero from ZeroMoney constructor", func(t *testing.T) {
			m := ZeroMoney("USD")

			result := m.IsZero()

			assert.True(t, result, "ZeroMoney should return true for IsZero")
		})

		t.Run("zero from NewMoneyInt with 0", func(t *testing.T) {
			m := NewMoneyInt("EUR", 0)

			result := m.IsZero()

			assert.True(t, result, "NewMoneyInt with 0 should return true for IsZero")
		})

		t.Run("zero from NewMoneyFromFraction with 0 numerator", func(t *testing.T) {
			m := NewMoneyFromFraction(0, 100, "JPY")

			result := m.IsZero()

			assert.True(t, result, "NewMoneyFromFraction with 0 numerator should return true for IsZero")
		})

		t.Run("zero from NewMoneyFloat with 0.0", func(t *testing.T) {
			m := NewMoneyFloat("GBP", 0.0)

			result := m.IsZero()

			assert.True(t, result, "NewMoneyFloat with 0.0 should return true for IsZero")
		})

		t.Run("zero from NewMoney with zerorat.Zero", func(t *testing.T) {
			m := NewMoney("CHF", zerorat.Zero())

			result := m.IsZero()

			assert.True(t, result, "NewMoney with zerorat.Zero should return true for IsZero")
		})
	})

	t.Run("valid non-zero money", func(t *testing.T) {
		t.Run("positive integer", func(t *testing.T) {
			m := NewMoneyInt("USD", 100)

			result := m.IsZero()

			assert.False(t, result, "Positive money should return false for IsZero")
		})

		t.Run("negative integer", func(t *testing.T) {
			m := NewMoneyInt("USD", -100)

			result := m.IsZero()

			assert.False(t, result, "Negative money should return false for IsZero")
		})

		t.Run("positive fraction", func(t *testing.T) {
			m := NewMoneyFromFraction(1, 2, "EUR") // 0.5

			result := m.IsZero()

			assert.False(t, result, "Positive fractional money should return false for IsZero")
		})

		t.Run("negative fraction", func(t *testing.T) {
			m := NewMoneyFromFraction(-3, 4, "JPY") // -0.75

			result := m.IsZero()

			assert.False(t, result, "Negative fractional money should return false for IsZero")
		})

		t.Run("positive float", func(t *testing.T) {
			m := NewMoneyFloat("GBP", 1.23)

			result := m.IsZero()

			assert.False(t, result, "Positive float money should return false for IsZero")
		})

		t.Run("negative float", func(t *testing.T) {
			m := NewMoneyFloat("CHF", -2.50)

			result := m.IsZero()

			assert.False(t, result, "Negative float money should return false for IsZero")
		})
	})
	t.Run("invalid money", func(t *testing.T) {
		t.Run("default zero-value Money", func(t *testing.T) {
			m := Money{}

			result := m.IsZero()

			assert.False(t, result, "Invalid Money (zero-value) should return false for IsZero")
		})

		t.Run("money with empty currency", func(t *testing.T) {
			m := NewMoneyInt("", 0) // invalid due to empty currency

			result := m.IsZero()

			assert.False(t, result, "Money with empty currency should return false for IsZero")
		})

		t.Run("money with invalid amount", func(t *testing.T) {
			invalidAmount := zerorat.New(1, 0) // invalid Rat (denominator = 0)
			m := NewMoney("USD", invalidAmount)

			result := m.IsZero()

			assert.False(t, result, "Money with invalid amount should return false for IsZero")
		})

		t.Run("invalidated money", func(t *testing.T) {
			m := NewMoneyInt("USD", 0) // start with valid zero money
			assert.True(t, m.IsZero(), "Should be zero before invalidation")

			m.Invalidate() // make it invalid

			result := m.IsZero()

			assert.False(t, result, "Invalidated money should return false for IsZero")
		})
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Run("very small positive fraction", func(t *testing.T) {
			m := NewMoneyFromFraction(1, 1000000, "USD") // 0.000001

			result := m.IsZero()

			assert.False(t, result, "Very small positive fraction should return false for IsZero")
		})

		t.Run("very small negative fraction", func(t *testing.T) {
			m := NewMoneyFromFraction(-1, 1000000, "USD") // -0.000001

			result := m.IsZero()

			assert.False(t, result, "Very small negative fraction should return false for IsZero")
		})

		t.Run("different currencies with zero", func(t *testing.T) {
			currencies := []string{"USD", "EUR", "JPY", "GBP", "CHF", "CAD", "AUD"}

			for _, currency := range currencies {
				m := ZeroMoney(currency)

				result := m.IsZero()

				assert.True(t, result, "Zero money should return true for IsZero regardless of currency: %s", currency)
			}
		})

		t.Run("reduced fraction that equals zero", func(t *testing.T) {
			// This creates 0/4 which should reduce to 0/1 (zero)
			m := NewMoneyFromFraction(0, 4, "USD")

			result := m.IsZero()

			assert.True(t, result, "Reduced fraction equal to zero should return true for IsZero")
		})
	})

	t.Run("consistency with underlying zerorat", func(t *testing.T) {
		t.Run("zero money amount should be zero", func(t *testing.T) {
			m := ZeroMoney("USD")

			assert.True(t, m.IsZero(), "Money should be zero")
			assert.True(t, m.Amount().IsZero(), "Underlying amount should also be zero")
		})

		t.Run("non-zero money amount should not be zero", func(t *testing.T) {
			m := NewMoneyInt("USD", 100)

			assert.False(t, m.IsZero(), "Money should not be zero")
			assert.False(t, m.Amount().IsZero(), "Underlying amount should also not be zero")
		})
	})
}
