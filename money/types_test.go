package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/assert"
)

// TestNewMoney tests the NewMoney constructor
func TestNewMoney(t *testing.T) {
	t.Run("valid currency and amount", func(t *testing.T) {
		amount := zerorat.New(123, 100) // 1.23
		m := NewMoney("USD", amount)

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "USD", m.Currency(), "Currency should match")
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(amount), "Amount should match")
	})

	t.Run("empty currency makes invalid", func(t *testing.T) {
		amount := zerorat.New(123, 100)
		m := NewMoney("", amount)

		assert.True(t, m.IsInvalid(), "Money with empty currency should be invalid")
		assert.False(t, m.IsValid(), "Money with empty currency should not be valid")
	})

	t.Run("invalid amount propagates", func(t *testing.T) {
		invalidAmount := zerorat.New(1, 0) // invalid Rat
		m := NewMoney("USD", invalidAmount)

		assert.True(t, m.IsInvalid(), "Money with invalid amount should be invalid")
		assert.False(t, m.IsValid(), "Money with invalid amount should not be valid")
	})

	t.Run("amount is reduced automatically", func(t *testing.T) {
		amount := zerorat.New(6, 4) // 6/4 should reduce to 3/2
		m := NewMoney("USD", amount)

		assert.True(t, m.IsValid(), "Money should be valid")
		// The amount should be reduced by the zerorat constructor
		expected := zerorat.New(3, 2)
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should be reduced")
	})
}

// TestNewMoneyInt tests the NewMoneyInt constructor
func TestNewMoneyInt(t *testing.T) {
	t.Run("positive integer", func(t *testing.T) {
		m := NewMoneyInt("EUR", 42)

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "EUR", m.Currency(), "Currency should match")
		expected := zerorat.NewFromInt(42)
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should match")
	})

	t.Run("negative integer", func(t *testing.T) {
		m := NewMoneyInt("JPY", -100)

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "JPY", m.Currency(), "Currency should match")
		expected := zerorat.NewFromInt(-100)
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should match")
	})

	t.Run("zero integer", func(t *testing.T) {
		m := NewMoneyInt("GBP", 0)

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "GBP", m.Currency(), "Currency should match")
		expected := zerorat.Zero()
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should be zero")
	})

	t.Run("empty currency makes invalid", func(t *testing.T) {
		m := NewMoneyInt("", 42)

		assert.True(t, m.IsInvalid(), "Money with empty currency should be invalid")
	})
}

// TestNewMoneyFloat tests the NewMoneyFloat constructor
func TestNewMoneyFloat(t *testing.T) {
	t.Run("valid float", func(t *testing.T) {
		m := NewMoneyFloat("USD", 1.23)

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "USD", m.Currency(), "Currency should match")
		expected := zerorat.NewFromFloat64(1.23)
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should match")
	})

	t.Run("negative float", func(t *testing.T) {
		m := NewMoneyFloat("EUR", -5.67)

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "EUR", m.Currency(), "Currency should match")
		expected := zerorat.NewFromFloat64(-5.67)
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should match")
	})

	t.Run("zero float", func(t *testing.T) {
		m := NewMoneyFloat("JPY", 0.0)

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "JPY", m.Currency(), "Currency should match")
		expected := zerorat.Zero()
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should be zero")
	})

	t.Run("invalid float propagates", func(t *testing.T) {
		// Test with a float that would create invalid Rat
		m := NewMoneyFloat("USD", 1e100) // Very large number that might overflow

		// If the float conversion fails, Money should be invalid
		actualAmount := m.Amount()
		if actualAmount.IsInvalid() {
			assert.True(t, m.IsInvalid(), "Money with invalid float should be invalid")
		}
	})

	t.Run("empty currency makes invalid", func(t *testing.T) {
		m := NewMoneyFloat("", 1.23)

		assert.True(t, m.IsInvalid(), "Money with empty currency should be invalid")
	})
}

// TestNewMoneyFromFraction tests the NewMoneyFromFraction constructor
func TestNewMoneyFromFraction(t *testing.T) {
	t.Run("valid fraction", func(t *testing.T) {
		m := NewMoneyFromFraction(123, 100, "USD") // 1.23

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "USD", m.Currency(), "Currency should match")
		expected := zerorat.New(123, 100)
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should match")
	})

	t.Run("zero numerator", func(t *testing.T) {
		m := NewMoneyFromFraction(0, 100, "EUR")

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "EUR", m.Currency(), "Currency should match")
		expected := zerorat.Zero()
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should be zero")
	})

	t.Run("negative numerator", func(t *testing.T) {
		m := NewMoneyFromFraction(-50, 100, "JPY")

		assert.True(t, m.IsValid(), "Money should be valid")
		assert.Equal(t, "JPY", m.Currency(), "Currency should match")
		expected := zerorat.New(-50, 100)
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should match")
	})

	t.Run("zero denominator makes invalid", func(t *testing.T) {
		m := NewMoneyFromFraction(123, 0, "USD")

		assert.True(t, m.IsInvalid(), "Money with zero denominator should be invalid")
	})

	t.Run("empty currency makes invalid", func(t *testing.T) {
		m := NewMoneyFromFraction(123, 100, "")

		assert.True(t, m.IsInvalid(), "Money with empty currency should be invalid")
	})

	t.Run("fraction is reduced automatically", func(t *testing.T) {
		m := NewMoneyFromFraction(6, 4, "USD") // 6/4 should reduce to 3/2

		assert.True(t, m.IsValid(), "Money should be valid")
		expected := zerorat.New(3, 2)
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should be reduced")
	})
}

// TestZeroMoney tests the ZeroMoney constructor
func TestZeroMoney(t *testing.T) {
	t.Run("valid currency", func(t *testing.T) {
		m := ZeroMoney("USD")

		assert.True(t, m.IsValid(), "Zero money should be valid")
		assert.Equal(t, "USD", m.Currency(), "Currency should match")
		expected := zerorat.Zero()
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(expected), "Amount should be zero")
	})

	t.Run("empty currency makes invalid", func(t *testing.T) {
		m := ZeroMoney("")

		assert.True(t, m.IsInvalid(), "Zero money with empty currency should be invalid")
	})
}

// TestMoneyValidityMethods tests IsValid, IsInvalid, and Invalidate methods
func TestMoneyValidityMethods(t *testing.T) {
	t.Run("valid money", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		assert.True(t, m.IsValid(), "Valid money should return true for IsValid")
		assert.False(t, m.IsInvalid(), "Valid money should return false for IsInvalid")
	})

	t.Run("invalid money - empty currency", func(t *testing.T) {
		m := NewMoneyInt("", 100)

		assert.False(t, m.IsValid(), "Invalid money should return false for IsValid")
		assert.True(t, m.IsInvalid(), "Invalid money should return true for IsInvalid")
	})

	t.Run("invalid money - invalid amount", func(t *testing.T) {
		invalidAmount := zerorat.New(1, 0)
		m := NewMoney("USD", invalidAmount)

		assert.False(t, m.IsValid(), "Invalid money should return false for IsValid")
		assert.True(t, m.IsInvalid(), "Invalid money should return true for IsInvalid")
	})

	t.Run("invalidate method", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		assert.True(t, m.IsValid(), "Money should start valid")

		m.Invalidate()

		assert.False(t, m.IsValid(), "Money should be invalid after Invalidate")
		assert.True(t, m.IsInvalid(), "Money should return true for IsInvalid after Invalidate")
		assert.Empty(t, m.Currency(), "Currency should be empty after Invalidate")
		actualAmount := m.Amount()
		assert.True(t, actualAmount.IsInvalid(), "Amount should be invalid after Invalidate")
	})
}

// TestMoneyAccessorMethods tests Currency and Amount accessor methods
func TestMoneyAccessorMethods(t *testing.T) {
	t.Run("valid money accessors", func(t *testing.T) {
		amount := zerorat.New(123, 100)
		m := NewMoney("USD", amount)

		assert.Equal(t, "USD", m.Currency(), "Currency accessor should return correct value")
		actualAmount := m.Amount()
		assert.True(t, actualAmount.Equal(amount), "Amount accessor should return correct value")
	})

	t.Run("invalid money accessors", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid due to empty currency

		assert.Empty(t, m.Currency(), "Invalid money should return empty currency")
		actualAmount := m.Amount()
		assert.True(t, actualAmount.IsInvalid(), "Invalid money should return invalid amount")
	})
}

// TestSameCurrency tests the SameCurrency methods
func TestSameCurrency(t *testing.T) {
	t.Run("same currency - both valid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 200)

		assert.True(t, m1.SameCurrency(m2), "Same currency should return true")
		assert.True(t, SameCurrency(m1, m2), "SameCurrency function should return true")
	})

	t.Run("different currency - both valid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 200)

		assert.False(t, m1.SameCurrency(m2), "Different currency should return false")
		assert.False(t, SameCurrency(m1, m2), "SameCurrency function should return false")
	})

	t.Run("one invalid money", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 200) // invalid

		assert.False(t, m1.SameCurrency(m2), "Invalid money should return false")
		assert.False(t, SameCurrency(m1, m2), "SameCurrency function should return false")
	})

	t.Run("both invalid money", func(t *testing.T) {
		m1 := NewMoneyInt("", 100) // invalid
		m2 := NewMoneyInt("", 200) // invalid

		assert.False(t, m1.SameCurrency(m2), "Both invalid should return false")
		assert.False(t, SameCurrency(m1, m2), "SameCurrency function should return false")
	})
}

// TestSameCurrencies tests the SameCurrencies variadic function
func TestSameCurrencies(t *testing.T) {
	t.Run("no arguments", func(t *testing.T) {
		result := SameCurrencies()
		assert.True(t, result, "SameCurrencies with no arguments should return true")
	})

	t.Run("single money - valid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		result := SameCurrencies(m1)
		assert.True(t, result, "SameCurrencies with single valid money should return true")
	})

	t.Run("single money - invalid", func(t *testing.T) {
		m1 := NewMoneyInt("", 100) // invalid
		result := SameCurrencies(m1)
		assert.False(t, result, "SameCurrencies with single invalid money should return false")
	})

	t.Run("two money - same currency both valid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 200)
		result := SameCurrencies(m1, m2)
		assert.True(t, result, "SameCurrencies with same currency should return true")
	})

	t.Run("two money - different currency both valid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 200)
		result := SameCurrencies(m1, m2)
		assert.False(t, result, "SameCurrencies with different currencies should return false")
	})

	t.Run("two money - one invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 200) // invalid
		result := SameCurrencies(m1, m2)
		assert.False(t, result, "SameCurrencies with one invalid money should return false")
	})

	t.Run("two money - both invalid", func(t *testing.T) {
		m1 := NewMoneyInt("", 100) // invalid
		m2 := NewMoneyInt("", 200) // invalid
		result := SameCurrencies(m1, m2)
		assert.False(t, result, "SameCurrencies with both invalid money should return false")
	})

	t.Run("three money - all same currency", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 200)
		m3 := NewMoneyInt("USD", 300)
		result := SameCurrencies(m1, m2, m3)
		assert.True(t, result, "SameCurrencies with three same currency should return true")
	})

	t.Run("three money - first two same, third different", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 200)
		m3 := NewMoneyInt("EUR", 300)
		result := SameCurrencies(m1, m2, m3)
		assert.False(t, result, "SameCurrencies with mixed currencies should return false")
	})

	t.Run("three money - first different from second and third", func(t *testing.T) {
		m1 := NewMoneyInt("GBP", 100)
		m2 := NewMoneyInt("USD", 200)
		m3 := NewMoneyInt("USD", 300)
		result := SameCurrencies(m1, m2, m3)
		assert.False(t, result, "SameCurrencies with first different should return false")
	})

	t.Run("three money - all different currencies", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 200)
		m3 := NewMoneyInt("GBP", 300)
		result := SameCurrencies(m1, m2, m3)
		assert.False(t, result, "SameCurrencies with all different currencies should return false")
	})

	t.Run("three money - one invalid in middle", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 200) // invalid
		m3 := NewMoneyInt("USD", 300)
		result := SameCurrencies(m1, m2, m3)
		assert.False(t, result, "SameCurrencies with invalid money in middle should return false")
	})

	t.Run("five money - all same currency", func(t *testing.T) {
		m1 := NewMoneyInt("EUR", 100)
		m2 := NewMoneyInt("EUR", 200)
		m3 := NewMoneyInt("EUR", 300)
		m4 := NewMoneyInt("EUR", 400)
		m5 := NewMoneyInt("EUR", 500)
		result := SameCurrencies(m1, m2, m3, m4, m5)
		assert.True(t, result, "SameCurrencies with five same currency should return true")
	})

	t.Run("five money - last one different", func(t *testing.T) {
		m1 := NewMoneyInt("EUR", 100)
		m2 := NewMoneyInt("EUR", 200)
		m3 := NewMoneyInt("EUR", 300)
		m4 := NewMoneyInt("EUR", 400)
		m5 := NewMoneyInt("JPY", 500)
		result := SameCurrencies(m1, m2, m3, m4, m5)
		assert.False(t, result, "SameCurrencies with last different should return false")
	})

	t.Run("mixed valid and invalid money", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 200) // invalid
		m3 := NewMoneyInt("USD", 300)
		m4 := NewMoneyInt("", 400) // invalid
		result := SameCurrencies(m1, m2, m3, m4)
		assert.False(t, result, "SameCurrencies with mixed valid/invalid should return false")
	})

	t.Run("zero amounts with same currency", func(t *testing.T) {
		m1 := ZeroMoney("USD")
		m2 := NewMoneyInt("USD", 0)
		m3 := NewMoneyFromFraction(0, 1, "USD")
		result := SameCurrencies(m1, m2, m3)
		assert.True(t, result, "SameCurrencies with zero amounts same currency should return true")
	})

	t.Run("negative amounts with same currency", func(t *testing.T) {
		m1 := NewMoneyInt("USD", -100)
		m2 := NewMoneyInt("USD", -200)
		m3 := NewMoneyInt("USD", 300)
		result := SameCurrencies(m1, m2, m3)
		assert.True(t, result, "SameCurrencies with negative amounts same currency should return true")
	})
}

// TestMoney_StatusMethods tests the status predicate methods
func TestMoney_StatusMethods(t *testing.T) {
	t.Run("IsNegative", func(t *testing.T) {
		t.Run("valid negative money", func(t *testing.T) {
			m := NewMoneyInt("USD", -100)

			result := m.IsNegative()

			assert.True(t, result, "Negative money should return true for IsNegative")
		})

		t.Run("valid positive money", func(t *testing.T) {
			m := NewMoneyInt("USD", 100)

			result := m.IsNegative()

			assert.False(t, result, "Positive money should return false for IsNegative")
		})

		t.Run("valid zero money", func(t *testing.T) {
			m := ZeroMoney("USD")

			result := m.IsNegative()

			assert.False(t, result, "Zero money should return false for IsNegative")
		})

		t.Run("invalid money", func(t *testing.T) {
			m := NewMoneyInt("", 100) // invalid

			result := m.IsNegative()

			assert.False(t, result, "Invalid money should return false for IsNegative")
		})
	})

	t.Run("IsPositive", func(t *testing.T) {
		t.Run("valid positive money", func(t *testing.T) {
			m := NewMoneyInt("USD", 100)

			result := m.IsPositive()

			assert.True(t, result, "Positive money should return true for IsPositive")
		})

		t.Run("valid negative money", func(t *testing.T) {
			m := NewMoneyInt("USD", -100)

			result := m.IsPositive()

			assert.False(t, result, "Negative money should return false for IsPositive")
		})

		t.Run("valid zero money", func(t *testing.T) {
			m := ZeroMoney("USD")

			result := m.IsPositive()

			assert.False(t, result, "Zero money should return false for IsPositive")
		})

		t.Run("invalid money", func(t *testing.T) {
			m := NewMoneyInt("", 100) // invalid

			result := m.IsPositive()

			assert.False(t, result, "Invalid money should return false for IsPositive")
		})
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("valid money", func(t *testing.T) {
			m := NewMoneyInt("USD", 100)

			result := m.IsEmpty()

			assert.False(t, result, "Valid money should return false for IsEmpty")
		})

		t.Run("invalid money with empty currency", func(t *testing.T) {
			m := NewMoneyInt("", 100) // invalid

			result := m.IsEmpty()

			assert.True(t, result, "Invalid money should return true for IsEmpty")
		})

		t.Run("invalid money with invalid amount", func(t *testing.T) {
			invalidAmount := zerorat.New(1, 0) // invalid Rat
			m := NewMoney("USD", invalidAmount)

			result := m.IsEmpty()

			assert.True(t, result, "Invalid money should return true for IsEmpty")
		})

		t.Run("consistency with IsInvalid", func(t *testing.T) {
			// Test various Money states to ensure IsEmpty == IsInvalid
			testCases := []Money{
				NewMoneyInt("USD", 100),            // valid
				NewMoneyInt("", 100),               // invalid currency
				NewMoney("USD", zerorat.New(1, 0)), // invalid amount
				ZeroMoney("EUR"),                   // valid zero
			}

			for i, m := range testCases {
				assert.Equal(t, m.IsInvalid(), m.IsEmpty(),
					"IsEmpty should equal IsInvalid for test case %d", i)
			}
		})
	})
}
