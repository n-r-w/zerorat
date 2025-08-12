package money

import (
	"testing"

	"github.com/n-r-w/zerorat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMoneyAddRat tests AddRat operations with zerorat.Rat operands
func TestMoneyAddRat(t *testing.T) {
	t.Run("mutable AddRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)   // $1.00
		ratValue := zerorat.New(50, 1) // 50/1 = 50

		err := m.AddRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable AddRat - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(50, 1)

		err := m.AddRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable AddRat - invalid rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(1, 0) // invalid: division by zero

		err := m.AddRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable AddRat - zero value", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.Zero()

		err := m.AddRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 100) // unchanged
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable AddRat - fractional value", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)  // 100
		ratValue := zerorat.New(1, 2) // 1/2 = 0.5

		err := m.AddRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		// 100 + 1/2 = 200/2 + 1/2 = 201/2
		expected := NewMoneyFromFraction(201, 2, "USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("immutable AddedRatErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(50, 1)

		result, err := m.AddedRatErr(ratValue)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable AddedRatErr - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(50, 1)

		result, err := m.AddedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable AddedRatErr - invalid rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(1, 0) // invalid

		result, err := m.AddedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable AddedRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(50, 1)

		result := m.AddedRat(ratValue)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 150)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable AddedRat - invalid returns invalid", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(50, 1)

		result := m.AddedRat(ratValue)

		assert.True(t, result.IsInvalid(), "Result should be invalid on invalid operand")
	})

	t.Run("negative rat value", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(-30, 1) // -30

		result := m.AddedRat(ratValue)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 70) // 100 + (-30) = 70
		assert.True(t, result.Equal(expected))
	})

	t.Run("complex fraction addition", func(t *testing.T) {
		m := NewMoneyFromFraction(1, 3, "USD") // 1/3
		ratValue := zerorat.New(1, 6)          // 1/6

		result := m.AddedRat(ratValue)

		assert.True(t, result.IsValid())
		// 1/3 + 1/6 = 2/6 + 1/6 = 3/6 = 1/2
		expected := NewMoneyFromFraction(1, 2, "USD")
		assert.True(t, result.Equal(expected))
	})
}

// TestMoneySubRat tests SubRat operations with zerorat.Rat operands
func TestMoneySubRat(t *testing.T) {
	t.Run("mutable SubRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)   // $1.00
		ratValue := zerorat.New(30, 1) // 30/1 = 30

		err := m.SubRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable SubRat - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(30, 1)

		err := m.SubRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable SubRat - invalid rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(1, 0) // invalid: division by zero

		err := m.SubRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable SubRat - zero value", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.Zero()

		err := m.SubRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 100) // unchanged
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable SubRat - fractional value", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)  // 100
		ratValue := zerorat.New(1, 2) // 1/2 = 0.5

		err := m.SubRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		// 100 - 1/2 = 200/2 - 1/2 = 199/2
		expected := NewMoneyFromFraction(199, 2, "USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("immutable SubtractedRatErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(30, 1)

		result, err := m.SubtractedRatErr(ratValue)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable SubtractedRatErr - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(30, 1)

		result, err := m.SubtractedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable SubtractedRatErr - invalid rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(1, 0) // invalid

		result, err := m.SubtractedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable SubtractedRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(30, 1)

		result := m.SubtractedRat(ratValue)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 70)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable SubtractedRat - invalid returns invalid", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(30, 1)

		result := m.SubtractedRat(ratValue)

		assert.True(t, result.IsInvalid(), "Result should be invalid on invalid operand")
	})

	t.Run("negative result", func(t *testing.T) {
		m := NewMoneyInt("USD", 30)
		ratValue := zerorat.New(100, 1) // 100

		result := m.SubtractedRat(ratValue)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", -70) // 30 - 100 = -70
		assert.True(t, result.Equal(expected))
	})

	t.Run("complex fraction subtraction", func(t *testing.T) {
		m := NewMoneyFromFraction(1, 2, "USD") // 1/2
		ratValue := zerorat.New(1, 6)          // 1/6

		result := m.SubtractedRat(ratValue)

		assert.True(t, result.IsValid())
		// 1/2 - 1/6 = 3/6 - 1/6 = 2/6 = 1/3
		expected := NewMoneyFromFraction(1, 3, "USD")
		assert.True(t, result.Equal(expected))
	})
}

// TestMoneyMulRat tests MulRat operations with zerorat.Rat operands
func TestMoneyMulRat(t *testing.T) {
	t.Run("mutable MulRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)  // $1.00
		ratValue := zerorat.New(3, 1) // 3/1 = 3

		err := m.MulRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 300)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulRat - zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.Zero()

		err := m.MulRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := ZeroMoney("USD")
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable MulRat - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(3, 1)

		err := m.MulRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable MulRat - invalid rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(1, 0) // invalid: division by zero

		err := m.MulRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable MulRat - fractional value", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)  // 100
		ratValue := zerorat.New(1, 2) // 1/2 = 0.5

		err := m.MulRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		// 100 * 1/2 = 100/2 = 50
		expected := NewMoneyInt("USD", 50)
		assert.True(t, m.Equal(expected))
	})

	t.Run("immutable MultipliedRatErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(3, 1)

		result, err := m.MultipliedRatErr(ratValue)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 300)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable MultipliedRatErr - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(3, 1)

		result, err := m.MultipliedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable MultipliedRatErr - invalid rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(1, 0) // invalid

		result, err := m.MultipliedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable MultipliedRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(3, 1)

		result := m.MultipliedRat(ratValue)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 300)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable MultipliedRat - invalid returns invalid", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(3, 1)

		result := m.MultipliedRat(ratValue)

		assert.True(t, result.IsInvalid(), "Result should be invalid on invalid operand")
	})

	t.Run("negative rat value", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(-2, 1) // -2

		result := m.MultipliedRat(ratValue)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", -200) // 100 * (-2) = -200
		assert.True(t, result.Equal(expected))
	})

	t.Run("complex fraction multiplication", func(t *testing.T) {
		m := NewMoneyFromFraction(2, 3, "USD") // 2/3
		ratValue := zerorat.New(3, 4)          // 3/4

		result := m.MultipliedRat(ratValue)

		assert.True(t, result.IsValid())
		// 2/3 * 3/4 = 6/12 = 1/2
		expected := NewMoneyFromFraction(1, 2, "USD")
		assert.True(t, result.Equal(expected))
	})
}

// TestMoneyDivRat tests DivRat operations with zerorat.Rat operands
func TestMoneyDivRat(t *testing.T) {
	t.Run("mutable DivRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)  // $1.00
		ratValue := zerorat.New(2, 1) // 2/1 = 2

		err := m.DivRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, m.Equal(expected))
	})

	t.Run("mutable DivRat - division by zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.Zero() // 0

		err := m.DivRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid(), "Money should be invalid after division by zero")
	})

	t.Run("mutable DivRat - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(2, 1)

		err := m.DivRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable DivRat - invalid rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(1, 0) // invalid: division by zero

		err := m.DivRat(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, m.IsInvalid())
	})

	t.Run("mutable DivRat - fractional value", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)  // 100
		ratValue := zerorat.New(1, 2) // 1/2 = 0.5

		err := m.DivRat(ratValue)

		require.NoError(t, err)
		assert.True(t, m.IsValid())
		// 100 / (1/2) = 100 * (2/1) = 200
		expected := NewMoneyInt("USD", 200)
		assert.True(t, m.Equal(expected))
	})

	t.Run("immutable DividedRatErr - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(2, 1)

		result, err := m.DividedRatErr(ratValue)

		require.NoError(t, err)
		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable DividedRatErr - division by zero", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.Zero()

		result, err := m.DividedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable DividedRatErr - invalid money", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(2, 1)

		result, err := m.DividedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable DividedRatErr - invalid rat", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(1, 0) // invalid

		result, err := m.DividedRatErr(ratValue)

		require.Error(t, err)
		assert.Equal(t, ErrMoneyInvalid, err)
		assert.True(t, result.IsInvalid())
	})

	t.Run("immutable DividedRat - success", func(t *testing.T) {
		m := NewMoneyInt("USD", 100)
		ratValue := zerorat.New(2, 1)

		result := m.DividedRat(ratValue)

		assert.True(t, result.IsValid())
		expected := NewMoneyInt("USD", 50)
		assert.True(t, result.Equal(expected))
		// Original should be unchanged
		assert.True(t, m.Equal(NewMoneyInt("USD", 100)))
	})

	t.Run("immutable DividedRat - invalid returns invalid", func(t *testing.T) {
		m := NewMoneyInt("", 100) // invalid
		ratValue := zerorat.New(2, 1)

		result := m.DividedRat(ratValue)

		assert.True(t, result.IsInvalid(), "Result should be invalid on invalid operand")
	})

	t.Run("complex fraction division", func(t *testing.T) {
		m := NewMoneyFromFraction(2, 3, "USD") // 2/3
		ratValue := zerorat.New(3, 4)          // 3/4

		result := m.DividedRat(ratValue)

		assert.True(t, result.IsValid())
		// 2/3 ÷ 3/4 = 2/3 * 4/3 = 8/9
		expected := NewMoneyFromFraction(8, 9, "USD")
		assert.True(t, result.Equal(expected))
	})
}
