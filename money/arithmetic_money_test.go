package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoneyMul tests Mul operations with Money operands
func TestMoneyMul(t *testing.T) {
	t.Run("mutable Mul - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 2)   // $0.02

		err := m1.Mul(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		assert.Equal(t, "USD", m1.Currency())
		expected := NewMoneyInt("USD", 200) // $2.00
		assert.True(t, m1.Equal(expected))
	})

	t.Run("mutable Mul - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		err := m1.Mul(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable Mul - invalid operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 2) // invalid

		err := m1.Mul(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after multiplying by invalid operand")
	})

	t.Run("mutable Mul - zero operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		err := m1.Mul(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		expected := ZeroMoney("USD")
		assert.True(t, m1.Equal(expected))
	})

	t.Run("immutable MultipliedErr - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 2)

		result, err := m1.MultipliedErr(m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 200)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable MultipliedErr - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		result, err := m1.MultipliedErr(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable Multiplied - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 2)

		result := m1.Multiplied(m2)

		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 200)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable Multiplied - different currency returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		result := m1.Multiplied(m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})
}

// TestMoneyDiv tests Div operations with Money operands
func TestMoneyDiv(t *testing.T) {
	t.Run("mutable Div - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 2)   // $0.02

		err := m1.Div(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		assert.Equal(t, "USD", m1.Currency())
		expected := NewMoneyInt("USD", 50) // $0.50
		assert.True(t, m1.Equal(expected))
	})

	t.Run("mutable Div - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		err := m1.Div(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable Div - division by zero", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		err := m1.Div(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after division by zero")
	})

	t.Run("mutable Div - invalid operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 2) // invalid

		err := m1.Div(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after dividing by invalid operand")
	})

	t.Run("immutable DividedErr - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 2)

		result, err := m1.DividedErr(m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable DividedErr - division by zero", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		result, err := m1.DividedErr(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable Divided - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 2)

		result := m1.Divided(m2)

		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable Divided - division by zero returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		result := m1.Divided(m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on division by zero")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable Divided - different currency returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		result := m1.Divided(m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})
}

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
