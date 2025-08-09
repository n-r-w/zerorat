package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew_ValidInputs tests creation of valid rational numbers
func TestNew_ValidInputs(t *testing.T) {
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

// TestNew_InvalidInputs tests creation of invalid rational numbers
func TestNew_InvalidInputs(t *testing.T) {
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

// TestNew_SignNormalization tests sign normalization
func TestNew_SignNormalization(t *testing.T) {
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

// TestNew_AutomaticReduction tests that New automatically reduces fractions
func TestNew_AutomaticReduction(t *testing.T) {
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

// TestNewFromInt tests creation of rational number from integer
func TestNewFromInt(t *testing.T) {
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

// TestNewFromFloat64 tests creation of rational number from float64 with minimum precision loss
func TestNewFromFloat64(t *testing.T) {
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

// TestNewFromFloat64_SpecialValues tests special float64 values
func TestNewFromFloat64_SpecialValues(t *testing.T) {
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

// TestNewFromFloat64_PrecisionLoss tests that the constructor minimizes precision loss
func TestNewFromFloat64_PrecisionLoss(t *testing.T) {
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

// TestNewFromFloat64_EdgeCases tests edge cases and boundary conditions
func TestNewFromFloat64_EdgeCases(t *testing.T) {
	tests := []struct {
		name            string
		value           float64
		description     string
		shouldBeInvalid bool // true if overflow to invalid state is expected
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
			name:            "largest finite",
			value:           math.MaxFloat64,
			description:     "largest finite float64 should overflow and become invalid",
			shouldBeInvalid: true, // This must overflow int64/uint64 representable bounds
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

			if tt.shouldBeInvalid {
				// For values that must overflow, the result should be invalid
				assert.True(t, r.IsInvalid(), "%s should be invalid due to overflow", tt.name)
				assert.Equal(t, uint64(0), r.denominator, "invalid rational should have denominator = 0")
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
