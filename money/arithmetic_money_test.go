package money

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoneyMulMoney tests Mul operations with Money operands
func TestMoneyMulMoney(t *testing.T) {
	t.Run("mutable MulMoney - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 2)   // $0.02

		err := m1.MulMoney(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		assert.Equal(t, "USD", m1.Currency())
		expected := NewMoneyInt("USD", 200) // $2.00
		assert.True(t, m1.Equal(expected))
	})

	t.Run("mutable MulMoney - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		err := m1.MulMoney(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable MulMoney - invalid operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 2) // invalid

		err := m1.MulMoney(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after multiplying by invalid operand")
	})

	t.Run("mutable MulMoney - zero operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		err := m1.MulMoney(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		expected := ZeroMoney("USD")
		assert.True(t, m1.Equal(expected))
	})

	t.Run("immutable MultipliedMoneyErr - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 2)

		result, err := m1.MultipliedMoneyErr(m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 200)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable MultipliedMoneyErr - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		result, err := m1.MultipliedMoneyErr(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable MultipliedMoney - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 2)

		result := m1.MultipliedMoney(m2)

		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 200)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable MultipliedMoney - different currency returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		result := m1.MultipliedMoney(m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})
}

// TestMoneyDivMoney tests Div operations with Money operands
func TestMoneyDivMoney(t *testing.T) {
	t.Run("mutable DivMoney - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 2)   // $0.02

		err := m1.DivMoney(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		assert.Equal(t, "USD", m1.Currency())
		expected := NewMoneyInt("USD", 50) // $0.50
		assert.True(t, m1.Equal(expected))
	})

	t.Run("mutable DivMoney - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		err := m1.DivMoney(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable DivMoney - division by zero", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		err := m1.DivMoney(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after division by zero")
	})

	t.Run("mutable DivMoney - invalid operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 2) // invalid

		err := m1.DivMoney(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after dividing by invalid operand")
	})

	t.Run("immutable DividedMoneyErr - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 2)

		result, err := m1.DividedMoneyErr(m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable DividedMoneyErr - division by zero", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		result, err := m1.DividedMoneyErr(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable DividedMoney - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 2)

		result := m1.DividedMoney(m2)

		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable DividedMoney - division by zero returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		result := m1.DividedMoney(m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on division by zero")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable DividedMoney - different currency returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 2)

		result := m1.DividedMoney(m2)

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

// TestMoneyMulManyFloat tests MulManyFloat varargs operations
func TestMoneyMulManyFloat(t *testing.T) {
	t.Run("mutable MulManyFloat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100) // 100 units

		err := m.MulManyFloat(2.0, 1.5) // multiply by 2.0, then 1.5

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyFloat("USD", 300.0) // 100 * 2.0 * 1.5 = 300.0
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulManyFloat - with zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulManyFloat(2.0, 0.0, 5.0) // multiply by 2.0, then 0.0, then 5.0

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := ZeroMoney("USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulManyFloat - invalid float", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulManyFloat(2.0, math.Inf(1)) // infinity

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable MulManyFloat - empty varargs", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		original := m.Amount()

		err := m.MulManyFloat()

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Amount().Equal(original), "Money should be unchanged with empty varargs")
	})

	t.Run("immutable MultipliedManyFloatErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result, err := m.MultipliedManyFloatErr(2.0, 1.5)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyFloat("USD", 300.0) // 100 * 2.0 * 1.5 = 300.0
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable MultipliedManyFloat - invalid float returns invalid", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result := m.MultipliedManyFloat(2.0, math.Inf(1)) // infinity

		assert.True(t, result.IsInvalid())
		// Original unchanged
		assert.True(t, m.IsValid())
	})
}
