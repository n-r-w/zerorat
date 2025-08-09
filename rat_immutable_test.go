package zerorat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRat_ImmutableOperations tests all immutable arithmetic operations
func TestRat_ImmutableOperations(t *testing.T) {
	t.Run("Added", func(t *testing.T) {
		cases := []arithmeticTestCase{
			{
				name:      "positive fractions",
				receiver:  New(1, 2), // 1/2
				other:     New(1, 3), // 1/3
				wantNum:   5,         // 1/2 + 1/3 = 5/6
				wantDenom: 6,
			},
			{
				name:      "add zero",
				receiver:  New(3, 4), // 3/4
				other:     New(0, 1), // 0
				wantNum:   3,         // 3/4 + 0 = 3/4
				wantDenom: 4,
			},
			{
				name:      "negative and positive",
				receiver:  New(-1, 2), // -1/2
				other:     New(1, 4),  // 1/4
				wantNum:   -2,         // -1/2 + 1/4 = -2/8
				wantDenom: 8,
			},
		}
		testImmutableOperation(t, "Added", func(r Rat, other Rat) Rat { return r.Added(other) }, cases)
	})

	t.Run("Subtracted", func(t *testing.T) {
		cases := []arithmeticTestCase{
			{
				name:      "positive fractions",
				receiver:  New(3, 4), // 3/4
				other:     New(1, 4), // 1/4
				wantNum:   2,         // 3/4 - 1/4 = 2/4
				wantDenom: 4,
			},
			{
				name:      "result negative",
				receiver:  New(1, 4), // 1/4
				other:     New(3, 4), // 3/4
				wantNum:   -2,        // 1/4 - 3/4 = -2/4
				wantDenom: 4,
			},
			{
				name:      "subtract zero",
				receiver:  New(3, 4), // 3/4
				other:     New(0, 1), // 0
				wantNum:   3,         // 3/4 - 0 = 3/4
				wantDenom: 4,
			},
		}
		testImmutableOperation(t, "Subtracted", func(r Rat, other Rat) Rat { return r.Subtracted(other) }, cases)
	})

	t.Run("Multiplied", func(t *testing.T) {
		cases := []arithmeticTestCase{
			{
				name:      "positive fractions",
				receiver:  New(2, 3), // 2/3
				other:     New(3, 4), // 3/4
				wantNum:   6,         // 2/3 * 3/4 = 6/12
				wantDenom: 12,
			},
			{
				name:      "multiply by zero",
				receiver:  New(5, 7), // 5/7
				other:     New(0, 1), // 0
				wantNum:   0,         // 5/7 * 0 = 0
				wantDenom: 1,
			},
			{
				name:      "multiply by one",
				receiver:  New(3, 4), // 3/4
				other:     New(1, 1), // 1
				wantNum:   3,         // 3/4 * 1 = 3/4
				wantDenom: 4,
			},
		}
		testImmutableOperation(t, "Multiplied", func(r Rat, other Rat) Rat { return r.Multiplied(other) }, cases)
	})

	t.Run("Divided", func(t *testing.T) {
		cases := []arithmeticTestCase{
			{
				name:      "positive fractions",
				receiver:  New(2, 3), // 2/3
				other:     New(3, 4), // 3/4
				wantNum:   8,         // 2/3 ÷ 3/4 = 2/3 * 4/3 = 8/9
				wantDenom: 9,
			},
			{
				name:      "divide by one",
				receiver:  New(3, 4), // 3/4
				other:     New(1, 1), // 1
				wantNum:   3,         // 3/4 ÷ 1 = 3/4
				wantDenom: 4,
			},
			{
				name:      "divide integer",
				receiver:  New(6, 1), // 6
				other:     New(2, 1), // 2
				wantNum:   6,         // 6 ÷ 2 = 6/2
				wantDenom: 2,
			},
		}
		testImmutableOperation(t, "Divided", func(r Rat, other Rat) Rat { return r.Divided(other) }, cases)
	})
}

// TestRat_ImmutableInvalidStatePropagation tests invalid state propagation for immutable operations
func TestRat_ImmutableInvalidStatePropagation(t *testing.T) {
	invalidReceiver := New(5, 0)
	invalidOther := New(3, 0)
	validOther := New(1, 2)

	t.Run("Added", func(t *testing.T) {
		// Invalid receiver
		result1 := invalidReceiver.Added(validOther)
		assert.True(t, result1.IsInvalid(), "Added should propagate invalid receiver")

		// Invalid other
		validReceiver := New(1, 2)
		result2 := validReceiver.Added(invalidOther)
		assert.True(t, result2.IsInvalid(), "Added should propagate invalid other")

		// Both invalid
		result3 := invalidReceiver.Added(invalidOther)
		assert.True(t, result3.IsInvalid(), "Added should handle both invalid")
	})

	t.Run("Subtracted", func(t *testing.T) {
		// Invalid receiver
		result1 := invalidReceiver.Subtracted(validOther)
		assert.True(t, result1.IsInvalid(), "Subtracted should propagate invalid receiver")

		// Invalid other
		validReceiver := New(1, 2)
		result2 := validReceiver.Subtracted(invalidOther)
		assert.True(t, result2.IsInvalid(), "Subtracted should propagate invalid other")

		// Both invalid
		result3 := invalidReceiver.Subtracted(invalidOther)
		assert.True(t, result3.IsInvalid(), "Subtracted should handle both invalid")
	})

	t.Run("Multiplied", func(t *testing.T) {
		// Invalid receiver
		result1 := invalidReceiver.Multiplied(validOther)
		assert.True(t, result1.IsInvalid(), "Multiplied should propagate invalid receiver")

		// Invalid other
		validReceiver := New(1, 2)
		result2 := validReceiver.Multiplied(invalidOther)
		assert.True(t, result2.IsInvalid(), "Multiplied should propagate invalid other")

		// Both invalid
		result3 := invalidReceiver.Multiplied(invalidOther)
		assert.True(t, result3.IsInvalid(), "Multiplied should handle both invalid")
	})

	t.Run("Divided", func(t *testing.T) {
		// Invalid receiver
		result1 := invalidReceiver.Divided(validOther)
		assert.True(t, result1.IsInvalid(), "Divided should propagate invalid receiver")

		// Invalid other
		validReceiver := New(1, 2)
		result2 := validReceiver.Divided(invalidOther)
		assert.True(t, result2.IsInvalid(), "Divided should propagate invalid other")

		// Both invalid
		result3 := invalidReceiver.Divided(invalidOther)
		assert.True(t, result3.IsInvalid(), "Divided should handle both invalid")

		// Division by zero
		result4 := validReceiver.Divided(New(0, 1))
		assert.True(t, result4.IsInvalid(), "Divided should handle division by zero")
	})
}

// TestRat_ImmutableOverflowDetection tests overflow detection for immutable operations
func TestRat_ImmutableOverflowDetection(t *testing.T) {
	t.Run("Added", func(t *testing.T) {
		// Test overflow in cross multiplication
		receiver := New(9223372036854775807, 2) // MaxInt64/2
		other := New(9223372036854775807, 3)    // MaxInt64/3
		result := receiver.Added(other)
		assert.True(t, result.IsInvalid(), "Added should detect overflow")
	})

	t.Run("Subtracted", func(t *testing.T) {
		// Test overflow in cross multiplication
		receiver := New(9223372036854775807, 2) // MaxInt64/2
		other := New(-9223372036854775807, 3)   // -MaxInt64/3
		result := receiver.Subtracted(other)
		assert.True(t, result.IsInvalid(), "Subtracted should detect overflow")
	})

	t.Run("Multiplied", func(t *testing.T) {
		// Test numerator overflow
		receiver := New(9223372036854775807, 1) // MaxInt64
		other := New(2, 1)                      // 2
		result := receiver.Multiplied(other)
		assert.True(t, result.IsInvalid(), "Multiplied should detect numerator overflow")
	})

	t.Run("Divided", func(t *testing.T) {
		// Test overflow in cross multiplication
		receiver := New(9223372036854775807, 1) // MaxInt64
		other := New(1, 2)                      // 1/2
		result := receiver.Divided(other)
		assert.True(t, result.IsInvalid(), "Divided should detect overflow")
	})
}

// TestRat_String tests string representation
func TestRat_String(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected string
	}{
		{
			name:     "positive fraction",
			rat:      New(3, 4),
			expected: "3/4",
		},
		{
			name:     "negative fraction",
			rat:      New(-5, 7),
			expected: "-5/7",
		},
		{
			name:     "positive integer",
			rat:      New(42, 1),
			expected: "42",
		},
		{
			name:     "negative integer",
			rat:      New(-17, 1),
			expected: "-17",
		},
		{
			name:     "zero",
			rat:      New(0, 1),
			expected: "0",
		},
		{
			name:     "zero with different denominator",
			rat:      New(0, 5),
			expected: "0",
		},
		{
			name:     "invalid rational",
			rat:      New(5, 0),
			expected: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.String()
			assert.Equal(t, tt.expected, result, "String representation mismatch")
		})
	}
}

// TestRat_String_EdgeCases tests edge cases for string representation
func TestRat_String_EdgeCases(t *testing.T) {
	// Test MinInt64 special case
	r1 := New(-9223372036854775808, 1) // MinInt64
	result1 := r1.String()
	assert.Equal(t, "-9223372036854775808", result1, "MinInt64 string representation")

	// Test large denominator
	r2 := New(1, 18446744073709551615) // MaxUint64
	result2 := r2.String()
	assert.Equal(t, "1/18446744073709551615", result2, "MaxUint64 denominator string representation")

	// Test reduced fraction string
	r3 := New(6, 8) // Should be reduced to 3/4
	result3 := r3.String()
	assert.Equal(t, "3/4", result3, "Reduced fraction string representation")
}
