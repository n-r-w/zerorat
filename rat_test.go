package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper types and functions for consolidated tests
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

// TestNewRat_ValidInputs tests creation of valid rational numbers
func TestNewRat_ValidInputs(t *testing.T) {
	tests := []struct {
		name        string
		numerator   int64
		denominator uint64
		wantNum     int64
		wantDenom   uint64
	}{
		{
			name:        "positive fraction",
			numerator:   3,
			denominator: 4,
			wantNum:     3,
			wantDenom:   4,
		},
		{
			name:        "negative fraction",
			numerator:   -5,
			denominator: 7,
			wantNum:     -5,
			wantDenom:   7,
		},
		{
			name:        "positive integer",
			numerator:   42,
			denominator: 1,
			wantNum:     42,
			wantDenom:   1,
		},
		{
			name:        "negative integer",
			numerator:   -42,
			denominator: 1,
			wantNum:     -42,
			wantDenom:   1,
		},
		{
			name:        "zero numerator",
			numerator:   0,
			denominator: 5,
			wantNum:     0,
			wantDenom:   1, // should normalize to 0/1
		},
		{
			name:        "large values",
			numerator:   9223372036854775807,  // MaxInt64
			denominator: 18446744073709551615, // MaxUint64
			wantNum:     9223372036854775807,
			wantDenom:   18446744073709551615,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.numerator, tt.denominator)
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestNewRat_InvalidInputs tests creation of invalid rational numbers
func TestNewRat_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		numerator   int64
		denominator uint64
	}{
		{
			name:        "zero denominator positive numerator",
			numerator:   5,
			denominator: 0,
		},
		{
			name:        "zero denominator negative numerator",
			numerator:   -3,
			denominator: 0,
		},
		{
			name:        "zero denominator zero numerator",
			numerator:   0,
			denominator: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.numerator, tt.denominator)
			assert.Equal(t, uint64(0), r.denominator, "should be invalid (denominator = 0)")
		})
	}
}

// TestNewRat_SignNormalization tests sign normalization
func TestNewRat_SignNormalization(t *testing.T) {
	// Note: in Go uint64 is always positive, so sign is always in numerator
	// This test verifies sign is handled correctly
	tests := []struct {
		name        string
		numerator   int64
		denominator uint64
		wantNum     int64
		wantDenom   uint64
	}{
		{
			name:        "positive/positive",
			numerator:   3,
			denominator: 4,
			wantNum:     3,
			wantDenom:   4,
		},
		{
			name:        "negative/positive",
			numerator:   -3,
			denominator: 4,
			wantNum:     -3,
			wantDenom:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.numerator, tt.denominator)
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestNewRatFromInt tests creation of rational number from integer
func TestNewRatFromInt(t *testing.T) {
	tests := []struct {
		name      string
		value     int64
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "positive integer",
			value:     42,
			wantNum:   42,
			wantDenom: 1,
		},
		{
			name:      "negative integer",
			value:     -17,
			wantNum:   -17,
			wantDenom: 1,
		},
		{
			name:      "zero",
			value:     0,
			wantNum:   0,
			wantDenom: 1,
		},
		{
			name:      "max int64",
			value:     9223372036854775807,
			wantNum:   9223372036854775807,
			wantDenom: 1,
		},
		{
			name:      "min int64",
			value:     -9223372036854775808,
			wantNum:   -9223372036854775808,
			wantDenom: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFromInt(tt.value)
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestNewRatFromFloat64 tests creation of rational number from float64 with minimum precision loss
func TestNewRatFromFloat64(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		wantNum   int64
		wantDenom uint64
	}{
		// Simple cases
		{
			name:      "zero",
			value:     0.0,
			wantNum:   0,
			wantDenom: 1,
		},
		{
			name:      "positive integer",
			value:     42.0,
			wantNum:   42,
			wantDenom: 1,
		},
		{
			name:      "negative integer",
			value:     -17.0,
			wantNum:   -17,
			wantDenom: 1,
		},
		// Simple fractions
		{
			name:      "one half",
			value:     0.5,
			wantNum:   1,
			wantDenom: 2,
		},
		{
			name:      "negative one half",
			value:     -0.5,
			wantNum:   -1,
			wantDenom: 2,
		},
		{
			name:      "one quarter",
			value:     0.25,
			wantNum:   1,
			wantDenom: 4,
		},
		{
			name:      "three quarters",
			value:     0.75,
			wantNum:   3,
			wantDenom: 4,
		},
		// Decimal fractions - these will be exact binary representations, not simplified decimals
		// 0.1 in binary is exactly 3602879701896397/36028797018963968 after reduction
		{
			name:      "one tenth",
			value:     0.1,
			wantNum:   3602879701896397,
			wantDenom: 36028797018963968,
		},
		// 0.01 in binary is exactly 5764607523034235/576460752303423488 after reduction
		{
			name:      "one hundredth",
			value:     0.01,
			wantNum:   5764607523034235,
			wantDenom: 576460752303423488,
		},
		{
			name:      "decimal 0.125",
			value:     0.125,
			wantNum:   1,
			wantDenom: 8,
		},
		{
			name:      "decimal 0.375",
			value:     0.375,
			wantNum:   3,
			wantDenom: 8,
		},
		// Mixed numbers
		{
			name:      "one and half",
			value:     1.5,
			wantNum:   3,
			wantDenom: 2,
		},
		{
			name:      "two and quarter",
			value:     2.25,
			wantNum:   9,
			wantDenom: 4,
		},
		{
			name:      "negative mixed",
			value:     -3.75,
			wantNum:   -15,
			wantDenom: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFromFloat64(tt.value)
			assert.True(t, r.IsValid(), "rational should be valid")
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestNewRatFromFloat64_SpecialValues tests special float64 values
func TestNewRatFromFloat64_SpecialValues(t *testing.T) {
	tests := []struct {
		name            string
		value           float64
		shouldBeInvalid bool
	}{
		{
			name:            "positive infinity",
			value:           math.Inf(1),
			shouldBeInvalid: true,
		},
		{
			name:            "negative infinity",
			value:           math.Inf(-1),
			shouldBeInvalid: true,
		},
		{
			name:            "NaN",
			value:           math.NaN(),
			shouldBeInvalid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFromFloat64(tt.value)
			if tt.shouldBeInvalid {
				assert.True(t, r.IsInvalid(), "should be invalid for %s (denominator should be 0)", tt.name)
				assert.Equal(t, uint64(0), r.denominator, "invalid rational should have denominator = 0")
			} else {
				assert.True(t, r.IsValid(), "should be valid for %s", tt.name)
				assert.Positive(t, r.denominator, "valid rational should have denominator > 0")
			}
		})
	}
}

// TestNewRatFromFloat64_PrecisionLoss tests that the constructor minimizes precision loss
func TestNewRatFromFloat64_PrecisionLoss(t *testing.T) {
	tests := []struct {
		name        string
		value       float64
		description string
	}{
		{
			name:        "repeating decimal 1/3",
			value:       1.0 / 3.0,
			description: "should find a good rational approximation for 1/3",
		},
		{
			name:        "repeating decimal 2/3",
			value:       2.0 / 3.0,
			description: "should find a good rational approximation for 2/3",
		},
		{
			name:        "pi approximation",
			value:       3.141592653589793,
			description: "should find a good rational approximation for pi",
		},
		{
			name:        "e approximation",
			value:       2.718281828459045,
			description: "should find a good rational approximation for e",
		},
		{
			name:        "very small positive",
			value:       1e-10,
			description: "should handle very small positive numbers",
		},
		{
			name:        "very small negative",
			value:       -1e-10,
			description: "should handle very small negative numbers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFromFloat64(tt.value)
			assert.True(t, r.IsValid(), "should be valid: %s", tt.description)

			// Convert back to float64 and check that we're reasonably close
			backToFloat := float64(r.numerator) / float64(r.denominator)
			diff := math.Abs(backToFloat - tt.value)
			relativeError := diff / math.Abs(tt.value)

			// Allow for some reasonable tolerance (e.g., 1e-15 for most cases)
			tolerance := 1e-15
			if math.Abs(tt.value) < 1e-10 {
				// For very small numbers, use absolute tolerance
				tolerance = 1e-20
			}

			assert.True(t, relativeError < tolerance || diff < tolerance,
				"precision loss too high: value=%g, rational=%d/%d, back=%g, diff=%g, rel_err=%g",
				tt.value, r.numerator, r.denominator, backToFloat, diff, relativeError)
		})
	}
}

// TestNewRatFromFloat64_EdgeCases tests edge cases and boundary conditions
func TestNewRatFromFloat64_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		description  string
		mayBeInvalid bool // true if overflow to invalid state is acceptable
	}{
		{
			name:        "max safe integer",
			value:       9007199254740992.0, // 2^53, largest integer exactly representable in float64
			description: "should handle max safe integer",
		},
		{
			name:        "min safe integer",
			value:       -9007199254740992.0, // -2^53
			description: "should handle min safe integer",
		},
		{
			name:        "smallest positive normal",
			value:       math.SmallestNonzeroFloat64,
			description: "should handle smallest positive normal float64",
		},
		{
			name:         "largest finite",
			value:        math.MaxFloat64,
			description:  "largest finite float64 may overflow to invalid state",
			mayBeInvalid: true, // This is likely to overflow int64/uint64 limits
		},
		{
			name:        "negative zero",
			value:       math.Copysign(0.0, -1),
			description: "should handle negative zero same as positive zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFromFloat64(tt.value)

			if tt.mayBeInvalid {
				// For values that may overflow, either valid or invalid is acceptable
				t.Logf("%s result: valid=%v, num=%d, denom=%d", tt.name, r.IsValid(), r.numerator, r.denominator)
				if r.IsInvalid() {
					assert.Equal(t, uint64(0), r.denominator, "invalid rational should have denominator = 0")
				}
				return
			}

			assert.True(t, r.IsValid(), "should be valid: %s", tt.description)
			assert.Positive(t, r.denominator, "valid rational should have denominator > 0")

			// For negative zero, should be same as positive zero
			if tt.value == 0.0 || tt.value == math.Copysign(0.0, -1) {
				assert.Equal(t, int64(0), r.numerator, "zero should have numerator 0")
				assert.Equal(t, uint64(1), r.denominator, "zero should have denominator 1")
			}
		})
	}
}

// TestZero tests creation of zero rational number
func TestZero(t *testing.T) {
	r := Zero()
	assert.Equal(t, int64(0), r.numerator, "zero numerator expected")
	assert.Equal(t, uint64(1), r.denominator, "denominator should be 1")
}

// TestOne tests creation of one rational number
func TestOne(t *testing.T) {
	r := One()
	assert.Equal(t, int64(1), r.numerator, "numerator should be 1")
	assert.Equal(t, uint64(1), r.denominator, "denominator should be 1")
}

// TestNewRat_AutomaticReduction tests that NewRat automatically reduces fractions
func TestNewRat_AutomaticReduction(t *testing.T) {
	tests := []struct {
		name      string
		num       int64
		denom     uint64
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "reduce 6/9 to 2/3",
			num:       6,
			denom:     9,
			wantNum:   2,
			wantDenom: 3,
		},
		{
			name:      "reduce 10/15 to 2/3",
			num:       10,
			denom:     15,
			wantNum:   2,
			wantDenom: 3,
		},
		{
			name:      "reduce -12/18 to -2/3",
			num:       -12,
			denom:     18,
			wantNum:   -2,
			wantDenom: 3,
		},
		{
			name:      "already reduced 3/4 stays 3/4",
			num:       3,
			denom:     4,
			wantNum:   3,
			wantDenom: 4,
		},
		{
			name:      "reduce 100/200 to 1/2",
			num:       100,
			denom:     200,
			wantNum:   1,
			wantDenom: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rat := New(tt.num, tt.denom)
			assert.Equal(t, tt.wantNum, rat.Numerator(), "Numerator should be reduced")
			assert.Equal(t, tt.wantDenom, rat.Denominator(), "Denominator should be reduced")
		})
	}
}

// TestRat_FieldAccess tests struct field access
func TestRat_FieldAccess(t *testing.T) {
	r := New(3, 4)

	// Verify fields are accessible (this compiles)
	assert.Equal(t, int64(3), r.numerator)
	assert.Equal(t, uint64(4), r.denominator)
}

// TestRat_IsValid tests validity check of rational number
func TestRat_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected bool
	}{
		{
			name:     "valid positive fraction",
			rat:      New(3, 4),
			expected: true,
		},
		{
			name:     "valid negative fraction",
			rat:      New(-5, 7),
			expected: true,
		},
		{
			name:     "valid zero",
			rat:      New(0, 1),
			expected: true,
		},
		{
			name:     "valid integer",
			rat:      New(42, 1),
			expected: true,
		},
		{
			name:     "invalid zero denominator",
			rat:      New(5, 0),
			expected: false,
		},
		{
			name:     "invalid zero denominator with zero numerator",
			rat:      New(0, 0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.IsValid()
			assert.Equal(t, tt.expected, result, "IsValid() result mismatch")
		})
	}
}

// TestRat_IsInvalid tests invalidity check of rational number
func TestRat_IsInvalid(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected bool
	}{
		{
			name:     "valid positive fraction",
			rat:      New(3, 4),
			expected: false,
		},
		{
			name:     "valid negative fraction",
			rat:      New(-5, 7),
			expected: false,
		},
		{
			name:     "valid zero",
			rat:      New(0, 1),
			expected: false,
		},
		{
			name:     "invalid zero denominator",
			rat:      New(5, 0),
			expected: true,
		},
		{
			name:     "invalid zero denominator with zero numerator",
			rat:      New(0, 0),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.IsInvalid()
			assert.Equal(t, tt.expected, result, "IsInvalid() result mismatch")
		})
	}
}

// TestRat_Invalidate tests forced invalidation
func TestRat_Invalidate(t *testing.T) {
	r := New(3, 4)
	assert.True(t, r.IsValid(), "should be valid initially")

	r.Invalidate()
	assert.True(t, r.IsInvalid(), "should be invalid after Invalidate()")
	assert.False(t, r.IsValid(), "should not be valid after Invalidate()")
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
				receiver: New(1, 18446744073709551615), // MaxUint64
				other:    New(1, 2),                    // Multiplying MaxUint64 * 2 should overflow
				desc:     "denominator overflow",
			},
		}
		testOverflowDetection(t, "Add", (*Rat).Add, overflowCases)
	})

	t.Run("Mul", func(t *testing.T) {
		overflowCases := []overflowTestCase{
			{
				name:     "numerator overflow",
				receiver: New(math.MaxInt64, 1),
				other:    New(2, 1), // MaxInt64 * 2 should overflow
				desc:     "numerator multiplication overflow",
			},
			{
				name:     "denominator overflow",
				receiver: New(1, math.MaxUint64),
				other:    New(1, 2), // MaxUint64 * 2 should overflow
				desc:     "denominator multiplication overflow",
			},
		}
		testOverflowDetection(t, "Mul", (*Rat).Mul, overflowCases)
	})

	t.Run("Div", func(t *testing.T) {
		overflowCases := []overflowTestCase{
			{
				name:     "numerator overflow in cross multiplication",
				receiver: New(math.MaxInt64, 1),
				other:    New(1, 2), // MaxInt64 * 2 should overflow
				desc:     "numerator cross multiplication overflow",
			},
			{
				name:     "denominator overflow in cross multiplication",
				receiver: New(1, math.MaxUint64),
				other:    New(2, 1), // MaxUint64 * 2 should overflow
				desc:     "denominator cross multiplication overflow",
			},
		}
		testOverflowDetection(t, "Div", (*Rat).Div, overflowCases)
	})
}

// TestRat_Div_DivisionByZero tests division by zero
func TestRat_Div_DivisionByZero(t *testing.T) {
	r := New(3, 4)
	r.Div(New(0, 1)) // division by zero
	assert.True(t, r.IsInvalid(), "division by zero should result in invalid state")
}

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
		}
		testImmutableOperation(t, "Added", func(r Rat, other Rat) Rat { return r.Added(other) }, cases)
	})

	t.Run("Subtracted", func(t *testing.T) {
		cases := []arithmeticTestCase{
			{
				name:      "positive fractions",
				receiver:  New(3, 4), // 3/4
				other:     New(1, 4), // 1/4
				wantNum:   2,         // 3/4 - 1/4 = 2/4 (no auto-reduction)
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
				wantNum:   6,         // 2/3 * 3/4 = 6/12 (no auto-reduction)
				wantDenom: 12,
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
				wantNum:   8,         // 2/3 ÷ 3/4 = 8/9
				wantDenom: 9,
			},
		}
		testImmutableOperation(t, "Divided", func(r Rat, other Rat) Rat { return r.Divided(other) }, cases)
	})
}

// TestRat_Equal tests equality comparison
func TestRat_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected bool
	}{
		{"equal fractions", New(1, 2), New(1, 2), true},
		{"equal reduced", New(2, 4), New(1, 2), true},
		{"different fractions", New(1, 2), New(1, 3), false},
		{"zero equal", New(0, 1), New(0, 5), true},
		{"invalid equal", New(1, 0), New(2, 0), false}, // invalid numbers are not equal
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Equal(tt.b))
		})
	}
}

// TestRat_Less tests "less than" comparison
func TestRat_Less(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected bool
	}{
		{"1/2 < 3/4", New(1, 2), New(3, 4), true},
		{"3/4 < 1/2", New(3, 4), New(1, 2), false},
		{"equal", New(1, 2), New(2, 4), false},
		{"negative", New(-1, 2), New(1, 2), true},
		{"invalid", New(1, 0), New(1, 2), false}, // invalid numbers always return false
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Less(tt.b))
		})
	}
}

// TestRat_Less_OverflowBug demonstrates the critical bug where Less() returns false on overflow
func TestRat_Less_OverflowBug(t *testing.T) {
	// This test demonstrates the bug: when cross-multiplication would overflow,
	// the current implementation incorrectly returns false instead of the correct comparison

	// Demonstrate the bug: (MaxInt64-1)/MaxUint64 vs MaxInt64/MaxUint64
	// (MaxInt64-1)/MaxUint64 < MaxInt64/MaxUint64, so Less should return true
	// Cross multiplication: (MaxInt64-1) * MaxUint64 vs MaxInt64 * MaxUint64
	// Both will overflow, current implementation returns false (WRONG!)

	a := New(math.MaxInt64-1, math.MaxUint64)
	b := New(math.MaxInt64, math.MaxUint64)

	result := a.Less(b)
	// Current implementation returns false due to overflow, but correct answer is true
	// This demonstrates the bug!

	t.Logf("(MaxInt64-1)/MaxUint64 < MaxInt64/MaxUint64 should be true, but got %v", result)

	// The bug: current implementation returns false when overflow occurs
	// But (MaxInt64-1)/MaxUint64 < MaxInt64/MaxUint64 should return true
	// This assertion will PASS with buggy implementation, but should FAIL after we fix it
	// Then we'll change it to assert.True to make it pass with the correct implementation
	assert.True(t, result, "This should FAIL with current buggy implementation - (MaxInt64-1)/MaxUint64 < MaxInt64/MaxUint64 should be true")
}

// TestRat_Less_OverflowSafe tests Less() method with overflow scenarios using 128-bit arithmetic
func TestRat_Less_OverflowSafe(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected bool
	}{
		{
			name:     "overflow case that should return true",
			a:        New(math.MaxInt64-1, math.MaxUint64),
			b:        New(math.MaxInt64, math.MaxUint64),
			expected: true, // (MaxInt64-1)/MaxUint64 < MaxInt64/MaxUint64
		},
		{
			name:     "overflow case that should return false",
			a:        New(math.MaxInt64, math.MaxUint64),
			b:        New(math.MaxInt64-1, math.MaxUint64),
			expected: false, // MaxInt64/MaxUint64 > (MaxInt64-1)/MaxUint64
		},
		{
			name:     "large positive vs small positive with overflow",
			a:        New(math.MaxInt64, 1),
			b:        New(1, math.MaxUint64),
			expected: false, // MaxInt64/1 > 1/MaxUint64
		},
		{
			name:     "small positive vs large positive with overflow",
			a:        New(1, math.MaxUint64),
			b:        New(math.MaxInt64, 1),
			expected: true, // 1/MaxUint64 < MaxInt64/1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Now that we have overflow-safe implementation, verify correct results
			result := tt.a.Less(tt.b)
			assert.Equal(t, tt.expected, result, "Less() should return correct result even with potential overflow")

			// Also verify the method doesn't panic
			assert.NotPanics(t, func() {
				tt.a.Less(tt.b)
			})
		})
	}
}

// TestRat_Equal_OverflowBug demonstrates the issue with Equal() method overflow handling
func TestRat_Equal_OverflowBug(t *testing.T) {
	// Test case where cross-multiplication would overflow but rationals are equal
	// Create two equivalent fractions that would cause overflow in cross-multiplication
	// Better example: MaxInt64/MaxUint64 == (MaxInt64/2)/(MaxUint64/2)
	// Wait, that's not right either. Let me use: 2*MaxInt64/2*MaxUint64 == MaxInt64/MaxUint64
	// But we can't have 2*MaxInt64. Let me try a different approach:

	// Use a case that definitely causes overflow but should be equal
	// Let's try: 1000000000000000000 / 2000000000000000000 == 500000000000000000 / 1000000000000000000
	// Both equal 0.5, but cross-multiplication will be very large

	a := New(1000000000000000000, 2000000000000000000)
	b := New(500000000000000000, 1000000000000000000)

	// These fractions are mathematically equal (both = 0.5):
	// 1000000000000000000 / 2000000000000000000 = 500000000000000000 / 1000000000000000000
	// Cross-multiplication: 1000000000000000000 * 1000000000000000000 vs 500000000000000000 * 2000000000000000000
	// Both equal 1000000000000000000000000000000000000, but this will overflow int64

	result := a.Equal(b)
	t.Logf("Equal comparison result: %v", result)
	t.Logf("a = %d/%d", a.numerator, a.denominator)
	t.Logf("b = %d/%d", b.numerator, b.denominator)

	// Let's verify manually: both fractions equal 0.5
	// Cross multiply: 1000000000000000000 * 1000000000000000000 vs 500000000000000000 * 2000000000000000000
	// Both equal 1000000000000000000000000000000000000
	// These should be equal

	// Let's test with simpler numbers first
	c := New(1, 2) // 1/2
	d := New(2, 4) // 2/4 = 1/2
	simpleResult := c.Equal(d)
	t.Logf("Simple test: 1/2 == 2/4 = %v", simpleResult)

	// Current implementation uses equalWithGCD fallback, which should work correctly
	// But we want to replace this with 128-bit arithmetic for consistency
	// This test documents the current behavior before we change it
	assert.True(t, result, "These fractions should be equal: both equal 0.5")
}

// TestRat_Equal_OverflowSafe tests Equal() method with 128-bit arithmetic
func TestRat_Equal_OverflowSafe(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected bool
	}{
		{
			name:     "equal fractions with potential overflow",
			a:        New(1000000000000000000, 2000000000000000000),
			b:        New(500000000000000000, 1000000000000000000),
			expected: true, // Both equal 0.5, cross-multiplication overflows
		},
		{
			name:     "unequal fractions with potential overflow",
			a:        New(1000000000000000000, 2000000000000000000),
			b:        New(math.MaxInt64-1, math.MaxUint64),
			expected: false, // 0.5 != (MaxInt64-1)/MaxUint64
		},
		{
			name:     "large equal integers",
			a:        New(math.MaxInt64, 1),
			b:        New(math.MaxInt64, 1),
			expected: true,
		},
		{
			name:     "large unequal integers",
			a:        New(math.MaxInt64, 1),
			b:        New(math.MaxInt64-1, 1),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Equal(tt.b)
			assert.Equal(t, tt.expected, result, "Equal() should return correct result with 128-bit arithmetic")
		})
	}
}

// TestRat_Compare_OverflowSafe tests Compare() method with overflow scenarios
func TestRat_Compare_OverflowSafe(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected int
	}{
		{
			name:     "equal with overflow",
			a:        New(1000000000000000000, 2000000000000000000),
			b:        New(500000000000000000, 1000000000000000000),
			expected: 0, // Both equal 0.5
		},
		{
			name:     "less than with overflow",
			a:        New(math.MaxInt64-1, math.MaxUint64),
			b:        New(math.MaxInt64, math.MaxUint64),
			expected: -1, // (MaxInt64-1)/MaxUint64 < MaxInt64/MaxUint64
		},
		{
			name:     "greater than with overflow",
			a:        New(math.MaxInt64, math.MaxUint64),
			b:        New(math.MaxInt64-1, math.MaxUint64),
			expected: 1, // MaxInt64/MaxUint64 > (MaxInt64-1)/MaxUint64
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			assert.Equal(t, tt.expected, result, "Compare() should return correct result with 128-bit arithmetic")
		})
	}
}

// TestRat_ArithmeticNoAutoReduction tests that arithmetic operations don't automatically reduce
func TestRat_ArithmeticNoAutoReduction(t *testing.T) {
	t.Run("Add should not auto-reduce", func(t *testing.T) {
		// Create unreduced fractions manually to test arithmetic operations
		// We need to bypass NewRat's automatic reduction for this test
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

// TestRat_OverflowDetectionMathBits tests overflow detection using math/bits
func TestRat_OverflowDetectionMathBits(t *testing.T) {
	// These tests verify that our new math/bits-based overflow detection works correctly
	// We'll test the edge cases that the current willOverflow functions handle

	t.Run("int64 multiplication overflow detection", func(t *testing.T) {
		// Test cases that should overflow
		overflowCases := []struct {
			a, b int64
			name string
		}{
			{math.MaxInt64, 2, "MaxInt64 * 2"},
			{math.MaxInt64 / 2, 3, "(MaxInt64/2) * 3"},
			{math.MinInt64, -1, "MinInt64 * -1"},
			{math.MinInt64 / 2, 3, "(MinInt64/2) * 3"},
		}

		for _, tc := range overflowCases {
			t.Run(tc.name, func(t *testing.T) {
				// Current implementation should detect overflow
				shouldOverflow := willOverflowInt64Mul(tc.a, tc.b)
				assert.True(t, shouldOverflow, "Should detect overflow for %s", tc.name)

				// After refactoring with math/bits, this should still work
				// We'll verify the math/bits implementation gives same result
			})
		}
	})

	t.Run("uint64 multiplication overflow detection", func(t *testing.T) {
		// Test cases that should overflow
		overflowCases := []struct {
			a, b uint64
			name string
		}{
			{math.MaxUint64, 2, "MaxUint64 * 2"},
			{math.MaxUint64 / 2, 3, "(MaxUint64/2) * 3"},
		}

		for _, tc := range overflowCases {
			t.Run(tc.name, func(t *testing.T) {
				// Current implementation should detect overflow
				shouldOverflow := willOverflowUint64Mul(tc.a, tc.b)
				assert.True(t, shouldOverflow, "Should detect overflow for %s", tc.name)
			})
		}
	})
}

// TestRat_Compare tests three-way comparison
func TestRat_Compare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected int
	}{
		{"1/2 vs 3/4", New(1, 2), New(3, 4), -1},
		{"3/4 vs 1/2", New(3, 4), New(1, 2), 1},
		{"equal", New(1, 2), New(2, 4), 0},
		{"invalid", New(1, 0), New(1, 2), 0}, // invalid numbers return 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Compare(tt.b))
		})
	}
}

// TestRat_String tests string representation
func TestRat_String(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected string
	}{
		{"positive fraction", New(3, 4), "3/4"},
		{"negative fraction", New(-5, 7), "-5/7"},
		{"integer", New(42, 1), "42"},
		{"zero", New(0, 1), "0"},
		{"invalid", New(1, 0), "invalid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.rat.String())
		})
	}
}

// TestRat_UtilityMethods tests utility methods
func TestRat_UtilityMethods(t *testing.T) {
	// Numerator, Denominator
	r := New(3, 4)
	assert.Equal(t, int64(3), r.Numerator())
	assert.Equal(t, uint64(4), r.Denominator())

	// Sign
	assert.Equal(t, 1, New(3, 4).Sign())
	assert.Equal(t, -1, New(-3, 4).Sign())
	assert.Equal(t, 0, New(0, 1).Sign())
	assert.Equal(t, 0, New(1, 0).Sign()) // invalid

	// IsZero, IsOne
	assert.True(t, New(0, 1).IsZero())
	assert.False(t, New(1, 2).IsZero())
	assert.True(t, New(1, 1).IsOne())
	assert.False(t, New(2, 1).IsOne())
}

// TestRat_Reduce tests mutable reduction to lowest terms
func TestRat_Reduce(t *testing.T) {
	tests := []struct {
		name      string
		input     Rat
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "already reduced",
			input:     New(3, 4),
			wantNum:   3,
			wantDenom: 4,
		},
		{
			name:      "reduce simple fraction",
			input:     New(6, 8),
			wantNum:   3,
			wantDenom: 4,
		},
		{
			name:      "reduce to integer",
			input:     New(10, 5),
			wantNum:   2,
			wantDenom: 1,
		},
		{
			name:      "reduce negative fraction",
			input:     New(-12, 18),
			wantNum:   -2,
			wantDenom: 3,
		},
		{
			name:      "reduce zero",
			input:     New(0, 15),
			wantNum:   0,
			wantDenom: 1,
		},
		{
			name:      "large numbers",
			input:     New(1000000, 2000000),
			wantNum:   1,
			wantDenom: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.input
			r.Reduce()
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Reduce_Invalid tests reduction with invalid state
func TestRat_Reduce_Invalid(t *testing.T) {
	r := New(6, 0) // invalid
	r.Reduce()
	assert.True(t, r.IsInvalid(), "invalid state should be preserved")
}

// TestRat_Reduced tests immutable reduction to lowest terms
func TestRat_Reduced(t *testing.T) {
	tests := []struct {
		name      string
		input     Rat
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "reduce fraction",
			input:     New(6, 8),
			wantNum:   3,
			wantDenom: 4,
		},
		{
			name:      "already reduced",
			input:     New(5, 7),
			wantNum:   5,
			wantDenom: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.input
			result := tt.input.Reduced()

			// Check result
			assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
			assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

			// Check original is unchanged
			assert.Equal(t, original.numerator, tt.input.numerator, "original should not be modified")
			assert.Equal(t, original.denominator, tt.input.denominator, "original should not be modified")
		})
	}
}

// TestRat_Reduced_Invalid tests immutable reduction with invalid state
func TestRat_Reduced_Invalid(t *testing.T) {
	original := New(6, 0) // invalid
	result := original.Reduced()
	assert.True(t, result.IsInvalid(), "result should be invalid")
	assert.True(t, original.IsInvalid(), "original should remain invalid")
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

// TestRat_Compare_ZeroComparison tests comparison when both numerators are zero
func TestRat_Compare_ZeroComparison(t *testing.T) {
	// This tests the specific line: if r.numerator == 0 && other.numerator == 0
	a := New(0, 3) // 0/3
	b := New(0, 7) // 0/7

	result := a.Compare(b)
	assert.Equal(t, 0, result, "0/3 should equal 0/7")
}

// TestWillOverflowUint64Mul_ZeroInputs tests overflow detection with zero inputs
func TestWillOverflowUint64Mul_ZeroInputs(t *testing.T) {
	// Test the specific line: if a == 0 || b == 0 { return false }
	assert.False(t, willOverflowUint64Mul(0, 100), "0 * 100 should not overflow")
	assert.False(t, willOverflowUint64Mul(100, 0), "100 * 0 should not overflow")
	assert.False(t, willOverflowUint64Mul(0, 0), "0 * 0 should not overflow")
}

// TestWillOverflowInt64Mul_MinInt64Cases tests overflow detection with MinInt64
func TestWillOverflowInt64Mul_MinInt64Cases(t *testing.T) {
	// Test the specific lines for MinInt64 special case handling
	// Lines 468-469: if b == math.MinInt64 { bAbs = uint64(math.MaxInt64) + 1 }

	// Test case where b == MinInt64
	result := willOverflowInt64Mul(2, math.MinInt64)
	assert.True(t, result, "2 * MinInt64 should overflow")

	// Test case where a == MinInt64
	result2 := willOverflowInt64Mul(math.MinInt64, 2)
	assert.True(t, result2, "MinInt64 * 2 should overflow")
}

// TestGcdInt64Uint64_MinInt64Special tests GCD calculation with MinInt64
func TestGcdInt64Uint64_MinInt64Special(t *testing.T) {
	// Test the specific lines 514-515: if a == math.MinInt64 { absA = uint64(math.MaxInt64) + 1 }
	result := gcdInt64Uint64(math.MinInt64, 2)
	expected := gcdUint64(uint64(math.MaxInt64)+1, 2)
	assert.Equal(t, expected, result, "GCD with MinInt64 should handle special case")
}

// TestCompareRationalsCrossMul_SignHandling tests cross multiplication with different sign combinations
func TestCompareRationalsCrossMul_SignHandling(t *testing.T) {
	tests := []struct {
		name     string
		aNum     int64
		aDenom   uint64
		cNum     int64
		cDenom   uint64
		expected int
	}{
		{
			name:     "negative vs positive",
			aNum:     -3,
			aDenom:   4,
			cNum:     2,
			cDenom:   5,
			expected: -1, // -3/4 < 2/5
		},
		{
			name:     "positive vs negative",
			aNum:     3,
			aDenom:   4,
			cNum:     -2,
			cDenom:   5,
			expected: 1, // 3/4 > -2/5
		},
		{
			name:     "both negative with MinInt64",
			aNum:     math.MinInt64,
			aDenom:   1,
			cNum:     -1,
			cDenom:   1,
			expected: -1, // MinInt64/1 < -1/1
		},
		{
			name:     "MinInt64 special case in cNum",
			aNum:     -1,
			aDenom:   1,
			cNum:     math.MinInt64,
			cDenom:   1,
			expected: 1, // -1/1 > MinInt64/1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareRationalsCrossMul(tt.aNum, tt.aDenom, tt.cNum, tt.cDenom)
			assert.Equal(t, tt.expected, result, "comparison result mismatch")
		})
	}
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
