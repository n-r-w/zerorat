package money

import (
	"math"
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoneyAddMoney tests Add operation with Money operands
func TestMoneyAddMoney(t *testing.T) {
	t.Run("mutable Add - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 50)  // $0.50

		err := m1.Add(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		assert.Equal(t, "USD", m1.Currency())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, m1.Equal(expected))
	})

	t.Run("mutable Add - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 50)

		err := m1.Add(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable Add - invalid operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 50) // invalid

		err := m1.Add(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after adding invalid operand")
	})

	t.Run("immutable AddedErr - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 50)

		result, err := m1.AddedErr(m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable AddedErr - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 50)

		result, err := m1.AddedErr(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable Added - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 50)

		result := m1.Added(m2)

		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable Added - different currency returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 50)

		result := m1.Added(m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})
}

// TestMoneySubMoney tests Sub operation with Money operands
func TestMoneySubMoney(t *testing.T) {
	t.Run("mutable Sub - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 30)  // $0.30

		err := m1.Sub(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		assert.Equal(t, "USD", m1.Currency())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, m1.Equal(expected))
	})

	t.Run("mutable Sub - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 30)

		err := m1.Sub(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable Sub - invalid operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 30) // invalid

		err := m1.Sub(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after subtracting invalid operand")
	})

	t.Run("immutable SubtractedErr - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 30)

		result, err := m1.SubtractedErr(m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable SubtractedErr - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 30)

		result, err := m1.SubtractedErr(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable Subtracted - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 30)

		result := m1.Subtracted(m2)

		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m1.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable Subtracted - different currency returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 30)

		result := m1.Subtracted(m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("negative result", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 30)
		m2 := NewMoneyInt("USD", 100)

		result := m1.Subtracted(m2)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", -70)
		assert.True(t, result.Equal(expected))
	})
}

// TestMoneyArithmeticOverflow tests overflow scenarios
func TestMoneyArithmeticOverflow(t *testing.T) {
	t.Run("Add overflow - mutable", func(t *testing.T) {
		// Create Money values that would cause overflow when added
		m1 := NewMoneyFromFraction(9223372036854775807, 1, "USD") // max int64
		m2 := NewMoneyInt("USD", 1)

		err := m1.Add(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after overflow")
	})

	t.Run("Sub overflow - immutable", func(t *testing.T) {
		// Create Money values that would cause overflow when subtracted
		m1 := NewMoneyFromFraction(-9223372036854775808, 1, "USD") // min int64
		m2 := NewMoneyInt("USD", 1)

		result, err := m1.SubtractedErr(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid after overflow")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})
}

// TestMoneyArithmeticEdgeCases tests edge cases
func TestMoneyArithmeticEdgeCases(t *testing.T) {
	t.Run("Add zero", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		result := m1.Added(m2)

		assert.True(t, result.IsValid())
		assert.True(t, result.Equal(m1))
	})

	t.Run("Sub zero", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")

		result := m1.Subtracted(m2)

		assert.True(t, result.IsValid())
		assert.True(t, result.Equal(m1))
	})

	t.Run("Add with fractions", func(t *testing.T) {
		m1 := NewMoneyFromFraction(1, 3, "USD") // 1/3
		m2 := NewMoneyFromFraction(1, 6, "USD") // 1/6

		result := m1.Added(m2)

		assert.True(t, result.IsValid())
		// 1/3 + 1/6 = 2/6 + 1/6 = 3/6 = 1/2
		expected := NewMoneyFromFraction(1, 2, "USD")
		assert.True(t, result.Equal(expected))
	})

	t.Run("both operands invalid", func(t *testing.T) {
		m1 := NewMoneyInt("", 100) // invalid
		m2 := NewMoneyInt("", 50)  // invalid

		err := m1.Add(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid())
	})
}

// TestMoneyAddScalar tests Add operations with scalar operands
func TestMoneyAddScalar(t *testing.T) {
	t.Run("mutable AddInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100) // $1.00

		err := m.AddInt(50) // add $0.50

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable AddInt - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid

		err := m.AddInt(50)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("immutable AddedIntErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result, err := m.AddedIntErr(50)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable AddedInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result := m.AddedInt(50)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("explicit float conversion before Add - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		other := mustNewMoneyFloat(t, "USD", 0.5)

		err := m.Add(other)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Equal(mustNewMoneyFloat(t, "USD", 100.5)))
	})

	t.Run("float constructor error is visible before Add", func(t *testing.T) {
		_, err := NewMoneyFloat("USD", math.Inf(1))
		require.ErrorIs(t, err, zerorat.ErrNonFiniteFloat)

		_, err = NewMoneyFloat("USD", 1e100)
		require.ErrorIs(t, err, zerorat.ErrNotRepresentable)
	})

	t.Run("explicit float conversion before AddedErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		other := mustNewMoneyFloat(t, "USD", 0.5)

		result, err := m.AddedErr(other)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.True(t, result.Equal(mustNewMoneyFloat(t, "USD", 100.5)))
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})
}

// TestMoneySubScalar tests Sub operations with scalar operands
func TestMoneySubScalar(t *testing.T) {
	t.Run("mutable SubInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100) // $1.00

		err := m.SubInt(30) // subtract $0.30

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable SubInt - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid

		err := m.SubInt(30)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("immutable SubtractedIntErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result, err := m.SubtractedIntErr(30)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable SubtractedInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result := m.SubtractedInt(30)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("explicit float conversion before Sub - success", func(t *testing.T) {
		m := mustNewMoneyFloat(t, "USD", 2.0)
		other := mustNewMoneyFloat(t, "USD", 0.5)

		err := m.Sub(other)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Equal(mustNewMoneyFloat(t, "USD", 1.5)))
	})

	t.Run("float constructor error is visible before Sub", func(t *testing.T) {
		_, err := NewMoneyFloat("USD", math.Inf(1))
		require.ErrorIs(t, err, zerorat.ErrNonFiniteFloat)

		_, err = NewMoneyFloat("USD", 1e100)
		require.ErrorIs(t, err, zerorat.ErrNotRepresentable)
	})

	t.Run("explicit float conversion before SubtractedErr - success", func(t *testing.T) {
		m := mustNewMoneyFloat(t, "USD", 2.0)
		other := mustNewMoneyFloat(t, "USD", 0.5)

		result, err := m.SubtractedErr(other)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.True(t, result.Equal(mustNewMoneyFloat(t, "USD", 1.5)))
		assert.True(t, m.Equal(mustNewMoneyFloat(t, "USD", 2.0)))
	})

	t.Run("negative result", func(t *testing.T) {
		m := NewMoneyInt("USD", 30)

		result := m.SubtractedInt(100)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", -70)
		assert.True(t, result.Equal(expected))
	})
}

// TestMoneyAddMany tests AddMany varargs operations
func TestMoneyAddMany(t *testing.T) {
	t.Run("mutable AddMany - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100) // $1.00
		m1 := NewMoneyInt("USD", 50) // $0.50
		m2 := NewMoneyInt("USD", 25) // $0.25

		err := m.AddMany(m1, m2)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 175) // $1.75
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable AddMany - currency mismatch", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		m1 := NewMoneyInt("USD", 50)
		m2 := NewMoneyInt("EUR", 25) // different currency

		err := m.AddMany(m1, m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable AddMany - invalid operand", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		m1 := NewMoneyInt("USD", 50)
		m2 := NewMoneyInt("", 25) // invalid

		err := m.AddMany(m1, m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid(), "Money should be invalid after adding invalid operand")
	})

	t.Run("mutable AddMany - empty varargs", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		original := m.Amount()

		err := m.AddMany()

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Amount().Equal(original), "Money should be unchanged with empty varargs")
	})

	t.Run("immutable AddedManyErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		m1 := NewMoneyInt("USD", 50)
		m2 := NewMoneyInt("USD", 25)

		result, err := m.AddedManyErr(m1, m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 175)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable AddedMany - currency mismatch returns invalid", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		m1 := NewMoneyInt("USD", 50)
		m2 := NewMoneyInt("EUR", 25)

		result := m.AddedMany(m1, m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original unchanged
		assert.True(t, m.IsValid())
	})
}

// TestMoneySubMany tests SubMany varargs operations
func TestMoneySubMany(t *testing.T) {
	t.Run("mutable SubMany - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 200) // $2.00
		m1 := NewMoneyInt("USD", 50) // $0.50
		m2 := NewMoneyInt("USD", 25) // $0.25

		err := m.SubMany(m1, m2)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 125) // $1.25
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable SubMany - currency mismatch", func(t *testing.T) {
		m := NewMoneyInt("USD", 200)
		m1 := NewMoneyInt("USD", 50)
		m2 := NewMoneyInt("EUR", 25) // different currency

		err := m.SubMany(m1, m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable SubMany - invalid operand", func(t *testing.T) {
		m := NewMoneyInt("USD", 200)
		m1 := NewMoneyInt("USD", 50)
		m2 := NewMoneyInt("", 25) // invalid

		err := m.SubMany(m1, m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid(), "Money should be invalid after subtracting invalid operand")
	})

	t.Run("mutable SubMany - empty varargs", func(t *testing.T) {
		m := NewMoneyInt("USD", 200)
		original := m.Amount()

		err := m.SubMany()

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Amount().Equal(original), "Money should be unchanged with empty varargs")
	})

	t.Run("immutable SubtractedManyErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 200)
		m1 := NewMoneyInt("USD", 50)
		m2 := NewMoneyInt("USD", 25)

		result, err := m.SubtractedManyErr(m1, m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 125)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 200)))
	})

	t.Run("immutable SubtractedMany - currency mismatch returns invalid", func(t *testing.T) {
		m := NewMoneyInt("USD", 200)
		m1 := NewMoneyInt("USD", 50)
		m2 := NewMoneyInt("EUR", 25)

		result := m.SubtractedMany(m1, m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original unchanged
		assert.True(t, m.IsValid())
	})

	t.Run("negative result", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		m1 := NewMoneyInt("USD", 150)

		result := m.SubtractedMany(m1)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", -50)
		assert.True(t, result.Equal(expected))
	})
}

// TestMoneyProfitMoney tests Profit operations (aliases for Sub operations)
func TestMoneyProfitMoney(t *testing.T) {
	t.Run("mutable Profit - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 30)  // $0.30

		err := m1.Profit(m2)

		require.NoError(t, err)
		assert.True(t, m1.IsValid())
		assert.Equal(t, "USD", m1.Currency())
		expected := NewMoneyInt("USD", 70) // $1.00 - $0.30 = $0.70
		assert.True(t, m1.Equal(expected))
	})

	t.Run("mutable Profit - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 30)

		err := m1.Profit(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after currency mismatch")
	})

	t.Run("mutable Profit - invalid operand", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 30) // invalid

		err := m1.Profit(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m1.IsInvalid(), "Money should be invalid after profit with invalid operand")
	})

	t.Run("immutable ProfitedErr - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 30)

		result, err := m1.ProfitedErr(m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		original := NewMoneyInt("USD", 100)
		assert.True(t, m1.Equal(original))
	})

	t.Run("immutable ProfitedErr - different currency failure", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 30)

		result, err := m1.ProfitedErr(m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("immutable Profited - same currency success", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("USD", 30)

		result := m1.Profited(m2)

		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		original := NewMoneyInt("USD", 100)
		assert.True(t, m1.Equal(original))
	})

	t.Run("immutable Profited - different currency returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 30)

		result := m1.Profited(m2)

		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
		// Original should be unchanged
		assert.True(t, m1.IsValid())
	})

	t.Run("negative profit result", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 50)  // $0.50
		m2 := NewMoneyInt("USD", 100) // $1.00

		result := m1.Profited(m2) // $0.50 - $1.00 = -$0.50

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", -50)
		assert.True(t, result.Equal(expected))
	})
}

// TestMoneyPercentOperations tests Percent operations (percentage calculations)
func TestMoneyPercentOperations(t *testing.T) {
	t.Run("PercentInt operations", func(t *testing.T) {
		t.Run("mutable PercentInt - success", func(t *testing.T) {
			m := NewMoneyInt("USD", 100) // $1.00

			err := m.PercentInt(50) // 50% of $1.00 = $0.50

			require.NoError(t, err)
			assert.True(t, m.IsValid())
			expected := NewMoneyInt("USD", 50)
			assert.True(t, m.Equal(expected))
		})

		t.Run("mutable PercentInt - invalid money", func(t *testing.T) {
			m := NewMoneyInt("", 100) // invalid

			err := m.PercentInt(50)

			require.Error(t, err)
			assert.Equal(t, ErrMoneyInvalid, err)
			assert.True(t, m.IsInvalid())
		})

		t.Run("immutable PercentIntErr - success", func(t *testing.T) {
			m := NewMoneyInt("USD", 200)

			result, err := m.PercentIntErr(25) // 25% of $2.00 = $0.50

			require.NoError(t, err)
			assert.True(t, result.IsValid())
			expected := NewMoneyInt("USD", 50)
			assert.True(t, result.Equal(expected))
			// Original unchanged
			original := NewMoneyInt("USD", 200)
			assert.True(t, m.Equal(original))
		})

		t.Run("immutable PercentedInt - success", func(t *testing.T) {
			m := NewMoneyInt("USD", 300)

			result := m.PercentedInt(10) // 10% of $3.00 = $0.30

			assert.True(t, result.IsValid())
			expected := NewMoneyInt("USD", 30)
			assert.True(t, result.Equal(expected))
			// Original unchanged
			original := NewMoneyInt("USD", 300)
			assert.True(t, m.Equal(original))
		})

		t.Run("immutable PercentedInt - invalid money returns invalid", func(t *testing.T) {
			m := NewMoneyInt("", 100) // invalid

			result := m.PercentedInt(50)

			assert.True(t, result.IsInvalid())
		})
	})

	t.Run("Percent with Rat operations", func(t *testing.T) {
		t.Run("mutable Percent - success", func(t *testing.T) {
			m := NewMoneyInt("USD", 100)

			err := m.Percent(zerorat.NewFromInt64(50))

			require.NoError(t, err)
			assert.True(t, m.IsValid())
			assert.True(t, m.Equal(NewMoneyInt("USD", 50)))
		})

		t.Run("mutable Percent - invalid Rat", func(t *testing.T) {
			m := NewMoneyInt("USD", 100)

			err := m.Percent(zerorat.Rat{})

			require.Error(t, err)
			assert.Equal(t, ErrMoneyInvalid, err)
		})

		t.Run("immutable PercentedErr - success", func(t *testing.T) {
			m := NewMoneyInt("USD", 200)

			result, err := m.PercentedErr(zerorat.NewFromInt64(25))

			require.NoError(t, err)
			assert.True(t, result.IsValid())
			assert.True(t, result.Equal(NewMoneyInt("USD", 50)))
			assert.True(t, m.Equal(NewMoneyInt("USD", 200)))
		})
	})

	t.Run("PercentMoney operations", func(t *testing.T) {
		t.Run("mutable PercentMoney - same currency success", func(t *testing.T) {
			m1 := NewMoneyInt("USD", 200) // $2.00
			m2 := NewMoneyInt("USD", 50)  // 50 (as percentage rate)

			err := m1.PercentMoney(m2) // $2.00 * (50 / 100) = $1.00

			require.NoError(t, err)
			assert.True(t, m1.IsValid())
			expected := NewMoneyInt("USD", 100)
			assert.True(t, m1.Equal(expected))
		})

		t.Run("mutable PercentMoney - different currency failure", func(t *testing.T) {
			m1 := NewMoneyInt("USD", 200)
			m2 := NewMoneyInt("EUR", 50)

			err := m1.PercentMoney(m2)

			require.Error(t, err)
			assert.Equal(t, ErrMoneyCurrencyMismatch, err)
			assert.True(t, m1.IsInvalid())
		})

		t.Run("immutable PercentMoneyErr - same currency success", func(t *testing.T) {
			m1 := NewMoneyInt("USD", 400) // $4.00
			m2 := NewMoneyInt("USD", 25)  // 25 (as percentage rate)

			result, err := m1.PercentMoneyErr(m2) // $4.00 * (25 / 100) = $1.00

			require.NoError(t, err)
			assert.True(t, result.IsValid())
			expected := NewMoneyInt("USD", 100)
			assert.True(t, result.Equal(expected))
			// Original unchanged
			original := NewMoneyInt("USD", 400)
			assert.True(t, m1.Equal(original))
		})

		t.Run("immutable PercentedMoney - same currency success", func(t *testing.T) {
			m1 := NewMoneyInt("USD", 100) // $1.00
			m2 := NewMoneyInt("USD", 50)  // 50 (as percentage rate)

			result := m1.PercentedMoney(m2) // $1.00 * (50 / 100) = $0.50

			assert.True(t, result.IsValid())
			expected := NewMoneyInt("USD", 50) // $1.00 * 0.5 = 50 cents
			assert.True(t, result.Equal(expected))
			// Original unchanged
			original := NewMoneyInt("USD", 100)
			assert.True(t, m1.Equal(original))
		})

		t.Run("immutable PercentedMoney - different currency returns invalid", func(t *testing.T) {
			m1 := NewMoneyInt("USD", 100)
			m2 := NewMoneyInt("EUR", 50)

			result := m1.PercentedMoney(m2)

			assert.True(t, result.IsInvalid())
			// Original unchanged
			assert.True(t, m1.IsValid())
		})
	})
}

// TestMoneyMulScalar tests Mul operations with scalar operands
func TestMoneyMulScalar(t *testing.T) {
	t.Run("mutable MulInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100) // $1.00

		err := m.MulInt(3) // multiply by 3

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 300)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulInt - zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulInt(0)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := ZeroMoney("USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulInt - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid

		err := m.MulInt(3)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("immutable MultipliedIntErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result, err := m.MultipliedIntErr(3)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 300)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable MultipliedInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result := m.MultipliedInt(3)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 300)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("mutable MulRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulRat(mustNewRatFromFloat64(t, 2.5))

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Equal(mustNewMoneyFloat(t, "USD", 250.0)))
	})

	t.Run("mutable MulRat - invalid Rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.MulRat(zerorat.Rat{})

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("immutable MultipliedRatErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result, err := m.MultipliedRatErr(mustNewRatFromFloat64(t, 2.5))

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.True(t, result.Equal(mustNewMoneyFloat(t, "USD", 250.0)))
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})
}

// TestMoneyDivScalar tests Div operations with scalar operands
func TestMoneyDivScalar(t *testing.T) {
	t.Run("mutable DivInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100) // $1.00

		err := m.DivInt(2) // divide by 2

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable DivInt - division by zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.DivInt(0)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid(), "Money should be invalid after division by zero")
	})

	t.Run("mutable DivInt - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid

		err := m.DivInt(2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("immutable DividedIntErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result, err := m.DividedIntErr(2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable DividedInt - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result := m.DividedInt(2)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, result.Equal(expected))
		// Original unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("mutable DivRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.DivRat(mustNewRatFromFloat64(t, 2.0))

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		assert.True(t, m.Equal(mustNewMoneyFloat(t, "USD", 50.0)))
	})

	t.Run("mutable DivRat - division by zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		err := m.DivRat(zerorat.Zero())

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("immutable DividedRatErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result, err := m.DividedRatErr(mustNewRatFromFloat64(t, 2.0))

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.True(t, result.Equal(mustNewMoneyFloat(t, "USD", 50.0)))
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable DividedRat - division by zero returns invalid", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)

		result := m.DividedRat(zerorat.Zero())

		assert.True(t, result.IsInvalid())
		assert.True(t, m.IsValid())
	})
}

// TestSum tests Sum and SumErr varargs operations
func TestSum(t *testing.T) {
	t.Run("SumErr - empty slice", func(t *testing.T) {
		result, err := SumErr()

		require.NoError(t, err)
		assert.True(t, result.IsInvalid(), "Sum of empty slice should return invalid Money")
		assert.Empty(t, result.Currency(), "Empty sum should have empty currency")
	})

	t.Run("SumErr - single money", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)

		result, err := SumErr(m1)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		assert.True(t, result.Equal(m1), "Sum of single Money should equal itself")
	})

	t.Run("SumErr - two moneys same currency", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 50)  // $0.50

		result, err := SumErr(m1, m2)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 150) // $1.50
		assert.True(t, result.Equal(expected))
	})

	t.Run("SumErr - multiple moneys same currency", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", 50)  // $0.50
		m3 := NewMoneyInt("USD", 25)  // $0.25
		m4 := NewMoneyInt("USD", 75)  // $0.75

		result, err := SumErr(m1, m2, m3, m4)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		assert.Equal(t, "USD", result.Currency())
		expected := NewMoneyInt("USD", 250) // $2.50
		assert.True(t, result.Equal(expected))
	})

	t.Run("SumErr - currency mismatch", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("EUR", 50) // different currency

		result, err := SumErr(m1, m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyCurrencyMismatch, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid on currency mismatch")
	})

	t.Run("SumErr - invalid money in first position", func(t *testing.T) {
		m1 := NewMoneyInt("", 100) // invalid
		m2 := NewMoneyInt("USD", 50)

		result, err := SumErr(m1, m2)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid when first operand is invalid")
	})

	t.Run("SumErr - invalid money in middle", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 50) // invalid
		m3 := NewMoneyInt("USD", 25)

		result, err := SumErr(m1, m2, m3)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid(), "Result should be invalid when any operand is invalid")
	})

	t.Run("SumErr - with fractions", func(t *testing.T) {
		m1 := NewMoneyFromFraction(1, 3, "USD") // 1/3
		m2 := NewMoneyFromFraction(1, 6, "USD") // 1/6
		m3 := NewMoneyFromFraction(1, 2, "USD") // 1/2

		result, err := SumErr(m1, m2, m3)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		// 1/3 + 1/6 + 1/2 = 2/6 + 1/6 + 3/6 = 6/6 = 1
		expected := NewMoneyInt("USD", 1)
		assert.True(t, result.Equal(expected))
	})

	t.Run("SumErr - with negative values", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100) // $1.00
		m2 := NewMoneyInt("USD", -30) // -$0.30
		m3 := NewMoneyInt("USD", -20) // -$0.20

		result, err := SumErr(m1, m2, m3)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 50) // $0.50
		assert.True(t, result.Equal(expected))
	})

	t.Run("SumErr - with zero values", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := ZeroMoney("USD")
		m3 := NewMoneyInt("USD", 50)

		result, err := SumErr(m1, m2, m3)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, result.Equal(expected))
	})

	t.Run("Sum - empty slice", func(t *testing.T) {
		result := Sum()

		assert.True(t, result.IsInvalid(), "Sum of empty slice should return invalid Money")
	})

	t.Run("Sum - single money", func(t *testing.T) {
		m1 := NewMoneyInt("EUR", 200)

		result := Sum(m1)

		assert.True(t, result.IsValid())
		assert.Equal(t, "EUR", result.Currency())
		assert.True(t, result.Equal(m1))
	})

	t.Run("Sum - multiple moneys success", func(t *testing.T) {
		m1 := NewMoneyInt("GBP", 100)
		m2 := NewMoneyInt("GBP", 200)
		m3 := NewMoneyInt("GBP", 300)

		result := Sum(m1, m2, m3)

		assert.True(t, result.IsValid())
		assert.Equal(t, "GBP", result.Currency())
		expected := NewMoneyInt("GBP", 600)
		assert.True(t, result.Equal(expected))
	})

	t.Run("Sum - currency mismatch returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("CAD", 50) // different currency

		result := Sum(m1, m2)

		assert.True(t, result.IsInvalid(), "Sum should return invalid Money on currency mismatch")
	})

	t.Run("Sum - invalid money returns invalid", func(t *testing.T) {
		m1 := NewMoneyInt("USD", 100)
		m2 := NewMoneyInt("", 50) // invalid

		result := Sum(m1, m2)

		assert.True(t, result.IsInvalid(), "Sum should return invalid Money when any operand is invalid")
	})
}
