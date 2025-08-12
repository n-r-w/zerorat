package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper types for arithmetic tests
type arithmeticTestCase struct {
	name      string
	receiver  Rat
	other     Rat
	wantNum   int64
	wantDenom uint64
}

type overflowTestCase struct {
	name     string
	receiver Rat
	other    Rat
	desc     string
}

type invalidStateTestCase struct {
	name     string
	receiver Rat
	other    Rat
}

// testArithmeticOperation is a helper to test arithmetic operations with common patterns
func testArithmeticOperation(t *testing.T, opName string, op func(*Rat, Rat), cases []arithmeticTestCase) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.receiver
			op(&r, tt.other)
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch for %s", opName)
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch for %s", opName)
		})
	}
}

// testImmutableOperation is a helper to test immutable operations
func testImmutableOperation(t *testing.T, opName string, op func(Rat, Rat) Rat, cases []arithmeticTestCase) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.receiver
			result := op(tt.receiver, tt.other)

			// Check the result
			assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch for %s", opName)
			assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch for %s", opName)

			// Check that the original hasn't changed
			assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified in %s", opName)
			assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified in %s", opName)
		})
	}
}

// testOverflowDetection is a helper to test overflow scenarios
func testOverflowDetection(t *testing.T, opName string, op func(*Rat, Rat), cases []overflowTestCase) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.receiver
			op(&r, tt.other)
			assert.True(t, r.IsInvalid(), "result should be invalid due to %s in %s", tt.desc, opName)
		})
	}
}

// testInvalidStatePropagation is a helper to test invalid state propagation
func testInvalidStatePropagation(t *testing.T, opName string, op func(*Rat, Rat), cases []invalidStateTestCase) {
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.receiver
			op(&r, tt.other)
			assert.True(t, r.IsInvalid(), "result should be invalid in %s", opName)
		})
	}
}

// TestRat_ArithmeticOperations tests all mutable arithmetic operations
func TestRat_ArithmeticOperations(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		cases := []arithmeticTestCase{
			{
				name:      "positive fractions",
				receiver:  New(1, 2), // 1/2
				other:     New(1, 3), // 1/3
				wantNum:   5,         // 1/2 + 1/3 = 3/6 + 2/6 = 5/6
				wantDenom: 6,
			},
			{
				name:      "negative and positive",
				receiver:  New(-1, 2), // -1/2
				other:     New(1, 4),  // 1/4
				wantNum:   -2,         // -1/2 + 1/4 = (-1*4 + 1*2)/(2*4) = (-4+2)/8 = -2/8
				wantDenom: 8,
			},
			{
				name:      "same denominator",
				receiver:  New(3, 7), // 3/7
				other:     New(2, 7), // 2/7
				wantNum:   5,         // 3/7 + 2/7 = 5/7
				wantDenom: 7,
			},
			{
				name:      "add zero",
				receiver:  New(3, 4), // 3/4
				other:     New(0, 1), // 0
				wantNum:   3,         // 3/4 + 0 = 3/4
				wantDenom: 4,
			},
			{
				name:      "add to zero",
				receiver:  New(0, 1), // 0
				other:     New(2, 5), // 2/5
				wantNum:   2,         // 0 + 2/5 = 2/5
				wantDenom: 5,
			},
			{
				name:      "result is zero",
				receiver:  New(1, 2),  // 1/2
				other:     New(-1, 2), // -1/2
				wantNum:   0,          // 1/2 + (-1/2) = 0
				wantDenom: 1,          // should normalize to 0/1
			},
		}
		testArithmeticOperation(t, "Add", (*Rat).Add, cases)
	})

	t.Run("Sub", func(t *testing.T) {
		cases := []arithmeticTestCase{
			{
				name:      "positive fractions",
				receiver:  New(3, 4), // 3/4
				other:     New(1, 4), // 1/4
				wantNum:   2,         // 3/4 - 1/4 = (3-1)/4 = 2/4 (same denominator optimization)
				wantDenom: 4,
			},
			{
				name:      "result negative",
				receiver:  New(1, 4), // 1/4
				other:     New(3, 4), // 3/4
				wantNum:   -2,        // 1/4 - 3/4 = (1-3)/4 = -2/4 (same denominator optimization)
				wantDenom: 4,
			},
			{
				name:      "subtract zero",
				receiver:  New(3, 4), // 3/4
				other:     New(0, 1), // 0
				wantNum:   3,         // 3/4 - 0 = 3/4
				wantDenom: 4,
			},
			{
				name:      "subtract from zero",
				receiver:  New(0, 1), // 0
				other:     New(2, 5), // 2/5
				wantNum:   -2,        // 0 - 2/5 = -2/5
				wantDenom: 5,
			},
		}
		testArithmeticOperation(t, "Sub", (*Rat).Sub, cases)
	})

	t.Run("Mul", func(t *testing.T) {
		cases := []arithmeticTestCase{
			{
				name:      "positive fractions",
				receiver:  New(2, 3), // 2/3
				other:     New(3, 4), // 3/4
				wantNum:   6,         // 2/3 * 3/4 = (2*3)/(3*4) = 6/12 (no auto-reduction)
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
			{
				name:      "negative result",
				receiver:  New(-2, 3), // -2/3
				other:     New(3, 5),  // 3/5
				wantNum:   -6,         // -2/3 * 3/5 = (-2*3)/(3*5) = -6/15 (no auto-reduction)
				wantDenom: 15,
			},
		}
		testArithmeticOperation(t, "Mul", (*Rat).Mul, cases)
	})

	t.Run("Div", func(t *testing.T) {
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
				wantNum:   6,         // 6 ÷ 2 = 6/1 * 1/2 = 6/2 (no auto-reduction)
				wantDenom: 2,
			},
		}
		testArithmeticOperation(t, "Div", (*Rat).Div, cases)
	})
}

// TestRat_ArithmeticInvalidStatePropagation tests invalid state propagation for all operations
func TestRat_ArithmeticInvalidStatePropagation(t *testing.T) {
	invalidStateCases := []invalidStateTestCase{
		{"invalid receiver", New(5, 0), New(1, 2)},
		{"invalid other", New(1, 2), New(3, 0)},
		{"both invalid", New(5, 0), New(3, 0)},
	}

	t.Run("Add", func(t *testing.T) {
		testInvalidStatePropagation(t, "Add", (*Rat).Add, invalidStateCases)
	})

	t.Run("Sub", func(t *testing.T) {
		testInvalidStatePropagation(t, "Sub", (*Rat).Sub, invalidStateCases)
	})

	t.Run("Mul", func(t *testing.T) {
		testInvalidStatePropagation(t, "Mul", (*Rat).Mul, invalidStateCases)
	})

	t.Run("Div", func(t *testing.T) {
		testInvalidStatePropagation(t, "Div", (*Rat).Div, invalidStateCases)
	})
}

// TestRat_ArithmeticOverflowDetection tests overflow detection for all operations
func TestRat_ArithmeticOverflowDetection(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		overflowCases := []overflowTestCase{
			{
				name:     "numerator overflow in cross multiplication",
				receiver: New(9223372036854775807, 2), // MaxInt64/2
				other:    New(9223372036854775807, 3), // MaxInt64/3
				desc:     "cross multiplication overflow",
			},
			{
				name:     "denominator overflow in multiplication",
				receiver: New(1, math.MaxUint64), // 1/MaxUint64
				other:    New(1, 2),              // 1/2 - MaxUint64 * 2 should overflow
				desc:     "denominator multiplication overflow",
			},
		}
		testOverflowDetection(t, "Add", (*Rat).Add, overflowCases)
	})

	t.Run("Sub", func(t *testing.T) {
		overflowCases := []overflowTestCase{
			{
				name:     "numerator overflow in cross multiplication",
				receiver: New(9223372036854775807, 2),  // MaxInt64/2
				other:    New(-9223372036854775807, 3), // -MaxInt64/3
				desc:     "cross multiplication overflow",
			},
		}
		testOverflowDetection(t, "Sub", (*Rat).Sub, overflowCases)
	})

	t.Run("Mul", func(t *testing.T) {
		overflowCases := []overflowTestCase{
			{
				name:     "numerator overflow",
				receiver: New(9223372036854775807, 1), // MaxInt64
				other:    New(2, 1),                   // 2
				desc:     "numerator multiplication overflow",
			},
			{
				name:     "denominator overflow",
				receiver: New(1, 9223372036854775807), // 1/MaxInt64
				other:    New(1, 3),                   // 1/3
				desc:     "denominator multiplication overflow",
			},
		}
		testOverflowDetection(t, "Mul", (*Rat).Mul, overflowCases)
	})

	t.Run("Div", func(t *testing.T) {
		overflowCases := []overflowTestCase{
			{
				name:     "division by zero",
				receiver: New(3, 4), // 3/4
				other:    New(0, 1), // 0
				desc:     "division by zero",
			},
			{
				name:     "numerator overflow in multiplication",
				receiver: New(9223372036854775807, 1), // MaxInt64
				other:    New(1, 2),                   // 1/2 -> reciprocal is 2/1
				desc:     "numerator multiplication overflow",
			},
		}
		testOverflowDetection(t, "Div", (*Rat).Div, overflowCases)
	})
}

// TestRat_ArithmeticNoAutoReduction tests that arithmetic operations don't auto-reduce
func TestRat_ArithmeticNoAutoReduction(t *testing.T) {
	t.Run("Add should not auto-reduce", func(t *testing.T) {
		// Create unreduced fractions manually to test arithmetic operations
		a := Rat{numerator: 2, denominator: 6}  // 2/6 (unreduced, equivalent to 1/3)
		b := Rat{numerator: 3, denominator: 18} // 3/18 (unreduced, equivalent to 1/6)
		a.Add(b)                                // Should be (2*18 + 3*6)/(6*18) = (36+18)/108 = 54/108 (unreduced)

		assert.Equal(t, int64(54), a.Numerator(), "Numerator should be 54 (unreduced)")
		assert.Equal(t, uint64(108), a.Denominator(), "Denominator should be 108 (unreduced)")
	})

	t.Run("Mul should not auto-reduce", func(t *testing.T) {
		// Create unreduced fractions manually to test arithmetic operations
		a := Rat{numerator: 4, denominator: 6} // 4/6 (unreduced, equivalent to 2/3)
		b := Rat{numerator: 6, denominator: 8} // 6/8 (unreduced, equivalent to 3/4)
		a.Mul(b)                               // Should be (4*6)/(6*8) = 24/48 (unreduced)

		assert.Equal(t, int64(24), a.Numerator(), "Numerator should be 24 (unreduced)")
		assert.Equal(t, uint64(48), a.Denominator(), "Denominator should be 48 (unreduced)")
	})

	t.Run("Sub should not auto-reduce", func(t *testing.T) {
		// Create unreduced fractions manually to test arithmetic operations
		a := Rat{numerator: 10, denominator: 12} // 10/12 (unreduced, equivalent to 5/6)
		b := Rat{numerator: 2, denominator: 6}   // 2/6 (unreduced, equivalent to 1/3)
		a.Sub(b)                                 // Should be (10*6 - 2*12)/(12*6) = (60-24)/72 = 36/72 (unreduced)

		assert.Equal(t, int64(36), a.Numerator(), "Numerator should be 36 (unreduced)")
		assert.Equal(t, uint64(72), a.Denominator(), "Denominator should be 72 (unreduced)")
	})

	t.Run("Div should not auto-reduce", func(t *testing.T) {
		// Create unreduced fractions manually to test arithmetic operations
		a := Rat{numerator: 4, denominator: 6}  // 4/6 (unreduced, equivalent to 2/3)
		b := Rat{numerator: 8, denominator: 12} // 8/12 (unreduced, equivalent to 2/3)
		a.Div(b)                                // Should be (4*12)/(6*8) = 48/48 (unreduced)

		assert.Equal(t, int64(48), a.Numerator(), "Numerator should be 48 (unreduced)")
		assert.Equal(t, uint64(48), a.Denominator(), "Denominator should be 48 (unreduced)")
	})
}

// TestRat_Div_NegativeDivisor tests division by negative numbers
func TestRat_Div_NegativeDivisor(t *testing.T) {
	tests := []struct {
		name      string
		receiver  Rat
		other     Rat
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "divide by negative",
			receiver:  New(6, 4),  // 6/4 = 3/2 (reduced)
			other:     New(-2, 3), // -2/3
			wantNum:   -9,         // 3/2 ÷ (-2/3) = 3/2 * (-3/2) = -9/4
			wantDenom: 4,
		},
		{
			name:      "divide negative by positive",
			receiver:  New(-6, 4), // -6/4 = -3/2 (reduced)
			other:     New(2, 3),  // 2/3
			wantNum:   -9,         // -3/2 ÷ 2/3 = -3/2 * 3/2 = -9/4
			wantDenom: 4,
		},
		{
			name:      "divide negative by negative",
			receiver:  New(-6, 4), // -6/4 = -3/2 (reduced)
			other:     New(-2, 3), // -2/3
			wantNum:   9,          // -3/2 ÷ (-2/3) = -3/2 * (-3/2) = 9/4
			wantDenom: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.receiver
			r.Div(tt.other)
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Div_MinInt64Special tests division with MinInt64 special case
func TestRat_Div_MinInt64Special(t *testing.T) {
	// Test the special case where other.numerator == math.MinInt64
	// This triggers the special handling: otherNum = uint64(math.MaxInt64) + 1
	r := New(1, 1)                 // 1/1
	other := New(math.MinInt64, 1) // MinInt64/1
	r.Div(other)                   // 1/1 ÷ MinInt64/1 = 1/1 * 1/MinInt64 = 1/MinInt64

	// The result should be 1/MinInt64 with sign change applied
	// Since MinInt64 is negative, signChange = true, so newNum = -1
	assert.Equal(t, int64(-1), r.numerator, "numerator should be -1")
	assert.Equal(t, uint64(math.MaxInt64)+1, r.denominator, "denominator should be MaxInt64+1")
}

// TestRat_Div_ZeroResult tests division resulting in zero
func TestRat_Div_ZeroResult(t *testing.T) {
	r := New(0, 5)     // 0/5
	other := New(3, 7) // 3/7
	r.Div(other)       // 0/5 ÷ 3/7 = 0

	assert.Equal(t, int64(0), r.numerator, "numerator should be 0")
	assert.Equal(t, uint64(1), r.denominator, "denominator should be normalized to 1")
}

// TestRat_Div_NegationNoOverflow ensures sign application path when negation is safe
func TestRat_Div_NegationNoOverflow(t *testing.T) {
	r := New(3, 2)      // 3/2
	other := New(-2, 5) // -2/5
	r.Div(other)        // (3/2) / (-2/5) = (3/2) * (5/2) = 15/4 then negate => -15/4

	assert.Equal(t, int64(-15), r.numerator, "numerator should be -15")
	assert.Equal(t, uint64(4), r.denominator, "denominator should be 4")
}

// TestRat_Div_MinInt64NegationOverflow tests the MinInt64 negation overflow case
func TestRat_Div_MinInt64NegationOverflow(t *testing.T) {
	// Create a case where newNum becomes MinInt64 and other.numerator < 0
	// This should trigger the MinInt64 negation overflow check (lines 284-288)

	// We need: r.numerator * other.denominator = MinInt64 and other.numerator < 0
	// Let's use: r = MinInt64/2, other = -1/2
	// Then: newNum = (MinInt64/2) * 2 = MinInt64, and other.numerator = -1 < 0
	r := New(math.MinInt64/2, 1) // MinInt64/2
	r.Div(New(-1, 2))            // divide by -1/2, which is multiply by -2

	// This should trigger the MinInt64 negation overflow and invalidate
	assert.True(t, r.IsInvalid(), "should be invalid due to MinInt64 negation overflow")
}

// TestRat_AddSub_AdditionalCases tests additional overflow and zero result cases
func TestRat_AddSub_AdditionalCases(t *testing.T) {
	t.Run("SameDenominatorOverflow", func(t *testing.T) {
		// Test overflow with same denominators
		r1 := New(math.MaxInt64, 5)
		r1.Add(New(1, 5)) // MaxInt64 + 1 should overflow
		assert.True(t, r1.IsInvalid(), "addition overflow with same denominator")

		r2 := New(math.MinInt64, 5)
		r2.Sub(New(1, 5)) // MinInt64 - 1 should overflow
		assert.True(t, r2.IsInvalid(), "subtraction overflow with same denominator")
	})

	t.Run("ZeroResults", func(t *testing.T) {
		// Test zero results in various cases
		r1 := New(3, 7)
		r1.Add(New(-3, 7)) // 3/7 + (-3/7) = 0
		assert.Equal(t, int64(0), r1.numerator, "should normalize to 0")
		assert.Equal(t, uint64(1), r1.denominator, "should normalize denominator to 1")

		r2 := New(1, 3)
		r2.Add(New(-2, 6)) // 1/3 + (-1/3) = 0
		assert.Equal(t, int64(0), r2.numerator, "should normalize to 0")
		assert.Equal(t, uint64(1), r2.denominator, "should normalize denominator to 1")
	})
}

// TestRat_AddSubCommon_EdgeCases tests specific addSubCommon paths for coverage
func TestRat_AddSubCommon_EdgeCases(t *testing.T) {
	// These tests cover specific internal paths in addSubCommon for 100% coverage

	// Test same denominator optimization paths
	cases := []struct {
		name      string
		r         Rat
		other     Rat
		isAdd     bool
		wantNum   int64
		wantDenom uint64
	}{
		{"same denom add", Rat{3, 11}, Rat{4, 11}, true, 7, 11},
		{"same denom sub", Rat{8, 13}, Rat{3, 13}, false, 5, 13},
		{"diff denom add", Rat{2, 7}, Rat{3, 11}, true, 43, 77},
		{"diff denom sub", Rat{5, 7}, Rat{2, 11}, false, 41, 77},
		{"diff denom zero result", Rat{2, 3}, Rat{4, 6}, false, 0, 1}, // Critical path for coverage
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := tc.r
			r.addSubCommon(tc.other, tc.isAdd)
			assert.Equal(t, tc.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tc.wantDenom, r.denominator, "denominator mismatch")
		})
	}

	// Additional edge cases to ensure 100% coverage
	t.Run("addSubCommon overflow edge cases", func(t *testing.T) {
		// Test overflow in first numerator term calculation
		r := Rat{math.MaxInt64, 1}
		other := Rat{1, 2}
		r.addSubCommon(other, true) // MaxInt64 * 2 should overflow
		assert.True(t, r.IsInvalid(), "should be invalid due to first term overflow")

		// Test overflow in second numerator term calculation
		r2 := Rat{1, 2}
		other2 := Rat{math.MaxInt64, 1}
		r2.addSubCommon(other2, true) // MaxInt64 * 2 should overflow
		assert.True(t, r2.IsInvalid(), "should be invalid due to second term overflow")

		// Test final addition overflow
		r3 := Rat{math.MaxInt64/2 + 1, 3}
		other3 := Rat{math.MaxInt64/2 + 1, 3}
		r3.addSubCommon(other3, true) // Large + Large should overflow
		assert.True(t, r3.IsInvalid(), "should be invalid due to final addition overflow")

		// Test final subtraction overflow
		r4 := Rat{math.MinInt64 + 1, 2}
		other4 := Rat{3, 1}
		r4.addSubCommon(other4, false) // (MinInt64+1)*1 - 3*2 should overflow
		assert.True(t, r4.IsInvalid(), "should be invalid due to final subtraction overflow")
	})
}

// TestRat_IntArithmetic tests arithmetic operations with int64 values
func TestRat_IntArithmetic(t *testing.T) {
	t.Run("AddInt", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     int64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "add positive int to fraction",
				receiver:  New(1, 2), // 1/2
				value:     3,         // 3
				wantNum:   7,         // 1/2 + 3 = 1/2 + 6/2 = 7/2
				wantDenom: 2,
			},
			{
				name:      "add negative int to fraction",
				receiver:  New(3, 4), // 3/4
				value:     -2,        // -2
				wantNum:   -5,        // 3/4 + (-2) = 3/4 - 8/4 = -5/4
				wantDenom: 4,
			},
			{
				name:      "add zero",
				receiver:  New(5, 7), // 5/7
				value:     0,         // 0
				wantNum:   5,         // 5/7 + 0 = 5/7
				wantDenom: 7,
			},
			{
				name:      "add to zero",
				receiver:  New(0, 1), // 0
				value:     42,        // 42
				wantNum:   42,        // 0 + 42 = 42/1
				wantDenom: 1,
			},
			{
				name:      "add large int",
				receiver:  New(1, 3), // 1/3
				value:     1000000,   // 1000000
				wantNum:   3000001,   // 1/3 + 1000000 = 1/3 + 3000000/3 = 3000001/3
				wantDenom: 3,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := tt.receiver
				r.AddInt(tt.value)
				assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
				assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
			})
		}
	})

	t.Run("AddedInt", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     int64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "add positive int to fraction",
				receiver:  New(1, 2), // 1/2
				value:     3,         // 3
				wantNum:   7,         // 1/2 + 3 = 7/2
				wantDenom: 2,
			},
			{
				name:      "add negative int to fraction",
				receiver:  New(3, 4), // 3/4
				value:     -1,        // -1
				wantNum:   -1,        // 3/4 + (-1) = 3/4 - 4/4 = -1/4
				wantDenom: 4,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.receiver
				result := tt.receiver.AddedInt(tt.value)

				// Check the result
				assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
				assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

				// Check that the original hasn't changed
				assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified")
				assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified")
			})
		}
	})

	t.Run("SubInt", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     int64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "subtract positive int from fraction",
				receiver:  New(7, 2), // 7/2
				value:     3,         // 3
				wantNum:   1,         // 7/2 - 3 = 7/2 - 6/2 = 1/2
				wantDenom: 2,
			},
			{
				name:      "subtract negative int from fraction",
				receiver:  New(1, 4), // 1/4
				value:     -2,        // -2
				wantNum:   9,         // 1/4 - (-2) = 1/4 + 8/4 = 9/4
				wantDenom: 4,
			},
			{
				name:      "subtract zero",
				receiver:  New(5, 7), // 5/7
				value:     0,         // 0
				wantNum:   5,         // 5/7 - 0 = 5/7
				wantDenom: 7,
			},
			{
				name:      "subtract from zero",
				receiver:  New(0, 1), // 0
				value:     42,        // 42
				wantNum:   -42,       // 0 - 42 = -42/1
				wantDenom: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := tt.receiver
				r.SubInt(tt.value)
				assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
				assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
			})
		}
	})

	t.Run("SubtractedInt", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     int64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "subtract positive int from fraction",
				receiver:  New(7, 2), // 7/2
				value:     3,         // 3
				wantNum:   1,         // 7/2 - 3 = 1/2
				wantDenom: 2,
			},
			{
				name:      "subtract negative int from fraction",
				receiver:  New(1, 4), // 1/4
				value:     -1,        // -1
				wantNum:   5,         // 1/4 - (-1) = 1/4 + 4/4 = 5/4
				wantDenom: 4,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.receiver
				result := tt.receiver.SubtractedInt(tt.value)

				// Check the result
				assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
				assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

				// Check that the original hasn't changed
				assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified")
				assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified")
			})
		}
	})

	t.Run("MulInt", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     int64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "multiply fraction by positive int",
				receiver:  New(3, 4), // 3/4
				value:     2,         // 2
				wantNum:   6,         // 3/4 * 2 = 6/4
				wantDenom: 4,
			},
			{
				name:      "multiply fraction by negative int",
				receiver:  New(2, 5), // 2/5
				value:     -3,        // -3
				wantNum:   -6,        // 2/5 * (-3) = -6/5
				wantDenom: 5,
			},
			{
				name:      "multiply by zero",
				receiver:  New(7, 3), // 7/3
				value:     0,         // 0
				wantNum:   0,         // 7/3 * 0 = 0/1
				wantDenom: 1,
			},
			{
				name:      "multiply by one",
				receiver:  New(5, 8), // 5/8
				value:     1,         // 1
				wantNum:   5,         // 5/8 * 1 = 5/8
				wantDenom: 8,
			},
			{
				name:      "multiply zero by int",
				receiver:  New(0, 1), // 0
				value:     42,        // 42
				wantNum:   0,         // 0 * 42 = 0/1
				wantDenom: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := tt.receiver
				r.MulInt(tt.value)
				assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
				assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
			})
		}
	})

	t.Run("MultipliedInt", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     int64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "multiply fraction by positive int",
				receiver:  New(3, 4), // 3/4
				value:     2,         // 2
				wantNum:   6,         // 3/4 * 2 = 6/4
				wantDenom: 4,
			},
			{
				name:      "multiply fraction by negative int",
				receiver:  New(2, 5), // 2/5
				value:     -3,        // -3
				wantNum:   -6,        // 2/5 * (-3) = -6/5
				wantDenom: 5,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.receiver
				result := tt.receiver.MultipliedInt(tt.value)

				// Check the result
				assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
				assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

				// Check that the original hasn't changed
				assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified")
				assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified")
			})
		}
	})

	t.Run("DivInt", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     int64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "divide fraction by positive int",
				receiver:  New(6, 4), // 6/4 = 3/2 (reduced by constructor)
				value:     2,         // 2
				wantNum:   3,         // 3/2 ÷ 2 = 3/2 * 1/2 = 3/4
				wantDenom: 4,
			},
			{
				name:      "divide fraction by negative int",
				receiver:  New(3, 5), // 3/5
				value:     -2,        // -2
				wantNum:   -3,        // 3/5 ÷ (-2) = 3/5 * (-1/2) = -3/10
				wantDenom: 10,
			},
			{
				name:      "divide by one",
				receiver:  New(7, 3), // 7/3
				value:     1,         // 1
				wantNum:   7,         // 7/3 ÷ 1 = 7/3
				wantDenom: 3,
			},
			{
				name:      "divide zero by int",
				receiver:  New(0, 1), // 0
				value:     42,        // 42
				wantNum:   0,         // 0 ÷ 42 = 0/1
				wantDenom: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := tt.receiver
				r.DivInt(tt.value)
				assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
				assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
			})
		}
	})

	t.Run("DividedInt", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     int64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "divide fraction by positive int",
				receiver:  New(6, 4), // 6/4 = 3/2 (reduced by constructor)
				value:     2,         // 2
				wantNum:   3,         // 3/2 ÷ 2 = 3/4
				wantDenom: 4,
			},
			{
				name:      "divide fraction by negative int",
				receiver:  New(3, 5), // 3/5
				value:     -2,        // -2
				wantNum:   -3,        // 3/5 ÷ (-2) = -3/10
				wantDenom: 10,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.receiver
				result := tt.receiver.DividedInt(tt.value)

				// Check the result
				assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
				assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

				// Check that the original hasn't changed
				assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified")
				assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified")
			})
		}
	})
}

// TestRat_FloatArithmetic tests arithmetic operations with float64 values
func TestRat_FloatArithmetic(t *testing.T) {
	t.Run("AddFloat", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     float64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "add positive float to fraction",
				receiver:  New(1, 2), // 1/2 = 0.5
				value:     0.25,      // 0.25 = 1/4
				wantNum:   6,         // 1/2 + 1/4 = (1*4 + 1*2)/(2*4) = 6/8 (not reduced)
				wantDenom: 8,
			},
			{
				name:      "add negative float to fraction",
				receiver:  New(3, 4), // 3/4 = 0.75
				value:     -0.5,      // -0.5 = -1/2
				wantNum:   2,         // 3/4 + (-1/2) = (3*2 + (-1)*4)/(4*2) = (6-4)/8 = 2/8 (not reduced)
				wantDenom: 8,
			},
			{
				name:      "add zero float",
				receiver:  New(5, 7), // 5/7
				value:     0.0,       // 0
				wantNum:   5,         // 5/7 + 0 = 5/7
				wantDenom: 7,
			},
			{
				name:      "add float to zero",
				receiver:  New(0, 1), // 0
				value:     0.125,     // 0.125 = 1/8
				wantNum:   1,         // 0 + 1/8 = 1/8
				wantDenom: 8,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := tt.receiver
				r.AddFloat(tt.value)
				assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
				assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
			})
		}
	})

	t.Run("AddedFloat", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     float64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "add positive float to fraction",
				receiver:  New(1, 2), // 1/2 = 0.5
				value:     0.25,      // 0.25 = 1/4
				wantNum:   6,         // 1/2 + 1/4 = 6/8 (not reduced)
				wantDenom: 8,
			},
			{
				name:      "add negative float to fraction",
				receiver:  New(3, 4), // 3/4 = 0.75
				value:     -0.25,     // -0.25 = -1/4
				wantNum:   2,         // 3/4 + (-1/4) = (3*4 + (-1)*4)/(4*4) = (12-4)/16 = 8/16 but reduced by constructor to 2/4
				wantDenom: 4,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.receiver
				result := tt.receiver.AddedFloat(tt.value)

				// Check the result
				assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
				assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

				// Check that the original hasn't changed
				assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified")
				assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified")
			})
		}
	})

	t.Run("SubFloat", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     float64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "subtract positive float from fraction",
				receiver:  New(3, 4), // 3/4 = 0.75
				value:     0.25,      // 0.25 = 1/4
				wantNum:   2,         // 3/4 - 1/4 = (3*4 - 1*4)/(4*4) = 8/16 = 2/4 (reduced)
				wantDenom: 4,
			},
			{
				name:      "subtract negative float from fraction",
				receiver:  New(1, 4), // 1/4 = 0.25
				value:     -0.5,      // -0.5 = -1/2
				wantNum:   6,         // 1/4 - (-1/2) = (1*2 - (-1)*4)/(4*2) = (2+4)/8 = 6/8 (not reduced)
				wantDenom: 8,
			},
			{
				name:      "subtract zero float",
				receiver:  New(5, 7), // 5/7
				value:     0.0,       // 0
				wantNum:   5,         // 5/7 - 0 = 5/7
				wantDenom: 7,
			},
			{
				name:      "subtract float from zero",
				receiver:  New(0, 1), // 0
				value:     0.125,     // 0.125 = 1/8
				wantNum:   -1,        // 0 - 1/8 = -1/8
				wantDenom: 8,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := tt.receiver
				r.SubFloat(tt.value)
				assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
				assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
			})
		}
	})

	t.Run("SubtractedFloat", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     float64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "subtract positive float from fraction",
				receiver:  New(3, 4), // 3/4 = 0.75
				value:     0.25,      // 0.25 = 1/4
				wantNum:   2,         // 3/4 - 1/4 = 2/4 (reduced)
				wantDenom: 4,
			},
			{
				name:      "subtract negative float from fraction",
				receiver:  New(1, 4), // 1/4 = 0.25
				value:     -0.25,     // -0.25 = -1/4
				wantNum:   2,         // 1/4 - (-1/4) = 1/4 + 1/4 = 2/4 (reduced)
				wantDenom: 4,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.receiver
				result := tt.receiver.SubtractedFloat(tt.value)

				// Check the result
				assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
				assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

				// Check that the original hasn't changed
				assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified")
				assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified")
			})
		}
	})

	t.Run("MulFloat", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     float64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "multiply fraction by positive float",
				receiver:  New(3, 4), // 3/4 = 0.75
				value:     0.5,       // 0.5 = 1/2
				wantNum:   3,         // 3/4 * 1/2 = 3/8
				wantDenom: 8,
			},
			{
				name:      "multiply fraction by negative float",
				receiver:  New(2, 5), // 2/5 = 0.4
				value:     -0.25,     // -0.25 = -1/4
				wantNum:   -2,        // 2/5 * (-1/4) = -2/20 (not reduced)
				wantDenom: 20,
			},
			{
				name:      "multiply by zero float",
				receiver:  New(7, 3), // 7/3
				value:     0.0,       // 0
				wantNum:   0,         // 7/3 * 0 = 0/1
				wantDenom: 1,
			},
			{
				name:      "multiply by one float",
				receiver:  New(5, 8), // 5/8
				value:     1.0,       // 1
				wantNum:   5,         // 5/8 * 1 = 5/8
				wantDenom: 8,
			},
			{
				name:      "multiply zero by float",
				receiver:  New(0, 1), // 0
				value:     0.125,     // 0.125 = 1/8
				wantNum:   0,         // 0 * 1/8 = 0/1
				wantDenom: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := tt.receiver
				r.MulFloat(tt.value)
				assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
				assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
			})
		}
	})

	t.Run("MultipliedFloat", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     float64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "multiply fraction by positive float",
				receiver:  New(3, 4), // 3/4 = 0.75
				value:     0.5,       // 0.5 = 1/2
				wantNum:   3,         // 3/4 * 1/2 = 3/8
				wantDenom: 8,
			},
			{
				name:      "multiply fraction by negative float",
				receiver:  New(2, 5), // 2/5 = 0.4
				value:     -0.25,     // -0.25 = -1/4
				wantNum:   -2,        // 2/5 * (-1/4) = -2/20 (not reduced)
				wantDenom: 20,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.receiver
				result := tt.receiver.MultipliedFloat(tt.value)

				// Check the result
				assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
				assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

				// Check that the original hasn't changed
				assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified")
				assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified")
			})
		}
	})

	t.Run("DivFloat", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     float64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "divide fraction by positive float",
				receiver:  New(3, 4), // 3/4 = 0.75
				value:     0.5,       // 0.5 = 1/2
				wantNum:   6,         // 3/4 ÷ 1/2 = 3/4 * 2/1 = 6/4 (not reduced)
				wantDenom: 4,
			},
			{
				name:      "divide fraction by negative float",
				receiver:  New(1, 2), // 1/2 = 0.5
				value:     -0.25,     // -0.25 = -1/4
				wantNum:   -4,        // 1/2 ÷ (-1/4) = 1/2 * (-4/1) = -4/2 (not reduced)
				wantDenom: 2,
			},
			{
				name:      "divide by one float",
				receiver:  New(7, 3), // 7/3
				value:     1.0,       // 1
				wantNum:   7,         // 7/3 ÷ 1 = 7/3
				wantDenom: 3,
			},
			{
				name:      "divide zero by float",
				receiver:  New(0, 1), // 0
				value:     0.125,     // 0.125 = 1/8
				wantNum:   0,         // 0 ÷ 1/8 = 0/1
				wantDenom: 1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r := tt.receiver
				r.DivFloat(tt.value)
				assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
				assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
			})
		}
	})

	t.Run("DividedFloat", func(t *testing.T) {
		tests := []struct {
			name      string
			receiver  Rat
			value     float64
			wantNum   int64
			wantDenom uint64
		}{
			{
				name:      "divide fraction by positive float",
				receiver:  New(3, 4), // 3/4 = 0.75
				value:     0.5,       // 0.5 = 1/2
				wantNum:   6,         // 3/4 ÷ 1/2 = 6/4 (not reduced)
				wantDenom: 4,
			},
			{
				name:      "divide fraction by negative float",
				receiver:  New(1, 2), // 1/2 = 0.5
				value:     -0.25,     // -0.25 = -1/4
				wantNum:   -4,        // 1/2 ÷ (-1/4) = -4/2 (not reduced)
				wantDenom: 2,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				original := tt.receiver
				result := tt.receiver.DividedFloat(tt.value)

				// Check the result
				assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
				assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

				// Check that the original hasn't changed
				assert.Equal(t, original.numerator, tt.receiver.numerator, "receiver should not be modified")
				assert.Equal(t, original.denominator, tt.receiver.denominator, "receiver should not be modified")
			})
		}
	})
}

// TestRat_IntArithmetic_EdgeCases tests edge cases for int64 arithmetic operations
func TestRat_IntArithmetic_EdgeCases(t *testing.T) {
	t.Run("invalid state propagation", func(t *testing.T) {
		// Test that invalid state is propagated through operations
		invalid := New(1, 0) // invalid rational

		// Mutable operations should remain invalid
		invalid.AddInt(5)
		assert.True(t, invalid.IsInvalid(), "AddInt should propagate invalid state")

		invalid = New(1, 0)
		invalid.SubInt(3)
		assert.True(t, invalid.IsInvalid(), "SubInt should propagate invalid state")

		invalid = New(1, 0)
		invalid.MulInt(2)
		assert.True(t, invalid.IsInvalid(), "MulInt should propagate invalid state")

		invalid = New(1, 0)
		invalid.DivInt(4)
		assert.True(t, invalid.IsInvalid(), "DivInt should propagate invalid state")
	})

	t.Run("immutable operations with invalid state", func(t *testing.T) {
		// Test that immutable operations return invalid results for invalid inputs
		invalid := New(1, 0) // invalid rational

		result := invalid.AddedInt(5)
		assert.True(t, result.IsInvalid(), "AddedInt should return invalid result for invalid input")

		result = invalid.SubtractedInt(3)
		assert.True(t, result.IsInvalid(), "SubtractedInt should return invalid result for invalid input")

		result = invalid.MultipliedInt(2)
		assert.True(t, result.IsInvalid(), "MultipliedInt should return invalid result for invalid input")

		result = invalid.DividedInt(4)
		assert.True(t, result.IsInvalid(), "DividedInt should return invalid result for invalid input")
	})

	t.Run("division by zero", func(t *testing.T) {
		// Test division by zero
		r := New(5, 3)
		r.DivInt(0)
		assert.True(t, r.IsInvalid(), "DivInt by zero should result in invalid state")

		r2 := New(7, 2)
		result := r2.DividedInt(0)
		assert.True(t, result.IsInvalid(), "DividedInt by zero should return invalid result")
	})
}

// TestRat_FloatArithmetic_EdgeCases tests edge cases for float64 arithmetic operations
func TestRat_FloatArithmetic_EdgeCases(t *testing.T) {
	t.Run("invalid state propagation", func(t *testing.T) {
		// Test that invalid state is propagated through operations
		invalid := New(1, 0) // invalid rational

		// Mutable operations should remain invalid
		invalid.AddFloat(0.5)
		assert.True(t, invalid.IsInvalid(), "AddFloat should propagate invalid state")

		invalid = New(1, 0)
		invalid.SubFloat(0.25)
		assert.True(t, invalid.IsInvalid(), "SubFloat should propagate invalid state")

		invalid = New(1, 0)
		invalid.MulFloat(2.0)
		assert.True(t, invalid.IsInvalid(), "MulFloat should propagate invalid state")

		invalid = New(1, 0)
		invalid.DivFloat(0.5)
		assert.True(t, invalid.IsInvalid(), "DivFloat should propagate invalid state")
	})

	t.Run("immutable operations with invalid state", func(t *testing.T) {
		// Test that immutable operations return invalid results for invalid inputs
		invalid := New(1, 0) // invalid rational

		result := invalid.AddedFloat(0.5)
		assert.True(t, result.IsInvalid(), "AddedFloat should return invalid result for invalid input")

		result = invalid.SubtractedFloat(0.25)
		assert.True(t, result.IsInvalid(), "SubtractedFloat should return invalid result for invalid input")

		result = invalid.MultipliedFloat(2.0)
		assert.True(t, result.IsInvalid(), "MultipliedFloat should return invalid result for invalid input")

		result = invalid.DividedFloat(0.5)
		assert.True(t, result.IsInvalid(), "DividedFloat should return invalid result for invalid input")
	})

	t.Run("special float values", func(t *testing.T) {
		// Test operations with special float values (NaN, Inf)
		r := New(1, 2)

		// NaN should result in invalid state
		r.AddFloat(math.NaN())
		assert.True(t, r.IsInvalid(), "AddFloat with NaN should result in invalid state")

		r = New(1, 2)
		r.AddFloat(math.Inf(1))
		assert.True(t, r.IsInvalid(), "AddFloat with +Inf should result in invalid state")

		r = New(1, 2)
		r.AddFloat(math.Inf(-1))
		assert.True(t, r.IsInvalid(), "AddFloat with -Inf should result in invalid state")

		// Test division by zero float
		r = New(3, 4)
		r.DivFloat(0.0)
		assert.True(t, r.IsInvalid(), "DivFloat by zero should result in invalid state")
	})
}
