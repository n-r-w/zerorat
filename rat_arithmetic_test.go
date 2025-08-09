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
