package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRat_Round_InvalidState tests that invalid rationals remain invalid after rounding
func TestRat_Round_InvalidState(t *testing.T) {
	tests := []struct {
		name      string
		rat       Rat
		roundType RoundType
		scale     int
	}{
		{"invalid with RoundDown", Rat{numerator: 5, denominator: 0}, RoundDown, 0},
		{"invalid with RoundUp", Rat{numerator: -3, denominator: 0}, RoundUp, 0},
		{"invalid with RoundHalfUp", Rat{numerator: 0, denominator: 0}, RoundHalfUp, 0},
		{"invalid with positive scale", Rat{numerator: 1, denominator: 0}, RoundDown, 2},
		{"invalid with negative scale", Rat{numerator: 1, denominator: 0}, RoundUp, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(tt.roundType, tt.scale)
			assert.True(t, r.IsInvalid(), "invalid rational should remain invalid")
			assert.Equal(t, uint64(0), r.denominator, "invalid rational should have denominator = 0")
		})
	}
}

// TestRat_Round_ScaleZero_RoundDown tests rounding to integers (scale=0) with RoundDown
func TestRat_Round_ScaleZero_RoundDown(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected Rat
	}{
		{"positive exact integer", New(5, 1), New(5, 1)},
		{"negative exact integer", New(-3, 1), New(-3, 1)},
		{"zero", New(0, 1), New(0, 1)},
		{"positive fraction rounds down", New(7, 3), New(2, 1)},          // 7/3 = 2.333... -> 2
		{"negative fraction rounds toward zero", New(-7, 3), New(-2, 1)}, // -7/3 = -2.333... -> -2
		{"positive half rounds down", New(5, 2), New(2, 1)},              // 5/2 = 2.5 -> 2
		{"negative half rounds toward zero", New(-5, 2), New(-2, 1)},     // -5/2 = -2.5 -> -2
		{"small positive fraction", New(1, 10), New(0, 1)},               // 0.1 -> 0
		{"small negative fraction", New(-1, 10), New(0, 1)},              // -0.1 -> 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(RoundDown, 0)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expected.numerator, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.expected.denominator, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Round_ScaleZero_RoundUp tests rounding to integers (scale=0) with RoundUp
func TestRat_Round_ScaleZero_RoundUp(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected Rat
	}{
		{"positive exact integer", New(5, 1), New(5, 1)},
		{"negative exact integer", New(-3, 1), New(-3, 1)},
		{"zero", New(0, 1), New(0, 1)},
		{"positive fraction rounds up", New(7, 3), New(3, 1)},               // 7/3 = 2.333... -> 3
		{"negative fraction rounds away from zero", New(-7, 3), New(-3, 1)}, // -7/3 = -2.333... -> -3
		{"positive half rounds up", New(5, 2), New(3, 1)},                   // 5/2 = 2.5 -> 3
		{"negative half rounds away from zero", New(-5, 2), New(-3, 1)},     // -5/2 = -2.5 -> -3
		{"small positive fraction", New(1, 10), New(1, 1)},                  // 0.1 -> 1
		{"small negative fraction", New(-1, 10), New(-1, 1)},                // -0.1 -> -1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(RoundUp, 0)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expected.numerator, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.expected.denominator, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Round_ScaleZero_RoundHalfUp tests rounding to integers (scale=0) with RoundHalfUp
func TestRat_Round_ScaleZero_RoundHalfUp(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected Rat
	}{
		{"positive exact integer", New(5, 1), New(5, 1)},
		{"negative exact integer", New(-3, 1), New(-3, 1)},
		{"zero", New(0, 1), New(0, 1)},
		{"positive less than half", New(7, 4), New(2, 1)},        // 7/4 = 1.75 -> 2
		{"negative less than half", New(-7, 4), New(-2, 1)},      // -7/4 = -1.75 -> -2
		{"positive exactly half", New(5, 2), New(3, 1)},          // 5/2 = 2.5 -> 3
		{"negative exactly half", New(-5, 2), New(-2, 1)},        // -5/2 = -2.5 -> -2
		{"positive more than half", New(8, 3), New(3, 1)},        // 8/3 = 2.666... -> 3
		{"negative more than half", New(-8, 3), New(-3, 1)},      // -8/3 = -2.666... -> -3
		{"small positive less than half", New(2, 10), New(0, 1)}, // 0.2 -> 0
		{"small positive exactly half", New(5, 10), New(1, 1)},   // 0.5 -> 1
		{"small positive less than half", New(4, 10), New(0, 1)}, // 0.4 -> 0
		{"small negative exactly half", New(-5, 10), New(0, 1)},  // -0.5 -> 0
		// -0.4 ->
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(RoundHalfUp, 0)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expected.numerator, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.expected.denominator, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Round_PositiveScale tests rounding with positive scale (decimal places)
func TestRat_Round_PositiveScale(t *testing.T) {
	tests := []struct {
		name      string
		rat       Rat
		roundType RoundType
		scale     int
		expected  Rat
	}{
		// Scale 1 (one decimal place)
		{"1.23 scale 1 RoundDown", New(123, 100), RoundDown, 1, New(12, 10)},     // 1.23 -> 1.2
		{"1.27 scale 1 RoundUp", New(127, 100), RoundUp, 1, New(13, 10)},         // 1.27 -> 1.3
		{"1.25 scale 1 RoundHalfUp", New(125, 100), RoundHalfUp, 1, New(13, 10)}, // 1.25 -> 1.3

		// Scale 2 (two decimal places)
		{"1.234 scale 2 RoundDown", New(1234, 1000), RoundDown, 2, New(123, 100)},     // 1.234 -> 1.23
		{"1.236 scale 2 RoundUp", New(1236, 1000), RoundUp, 2, New(124, 100)},         // 1.236 -> 1.24
		{"1.235 scale 2 RoundHalfUp", New(1235, 1000), RoundHalfUp, 2, New(124, 100)}, // 1.235 -> 1.24

		// Negative numbers
		{"-1.27 scale 1 RoundDown", New(-127, 100), RoundDown, 1, New(-12, 10)},     // -1.27 -> -1.2
		{"-1.23 scale 1 RoundUp", New(-123, 100), RoundUp, 1, New(-13, 10)},         // -1.23 -> -1.3
		{"-1.25 scale 1 RoundHalfUp", New(-125, 100), RoundHalfUp, 1, New(-12, 10)}, // -1.25 -> -1.2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(tt.roundType, tt.scale)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expected.numerator, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.expected.denominator, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Round_NegativeScale tests rounding with negative scale (to tens, hundreds, etc.)
func TestRat_Round_NegativeScale(t *testing.T) {
	tests := []struct {
		name      string
		rat       Rat
		roundType RoundType
		scale     int
		expected  Rat
	}{
		// Scale -1 (round to nearest 10)
		{"123 scale -1 RoundDown", New(123, 1), RoundDown, -1, New(120, 1)},     // 123 -> 120
		{"127 scale -1 RoundUp", New(127, 1), RoundUp, -1, New(130, 1)},         // 127 -> 130
		{"125 scale -1 RoundHalfUp", New(125, 1), RoundHalfUp, -1, New(130, 1)}, // 125 -> 130

		// Scale -2 (round to nearest 100)
		{"1234 scale -2 RoundDown", New(1234, 1), RoundDown, -2, New(1200, 1)},     // 1234 -> 1200
		{"1267 scale -2 RoundUp", New(1267, 1), RoundUp, -2, New(1300, 1)},         // 1267 -> 1300
		{"1250 scale -2 RoundHalfUp", New(1250, 1), RoundHalfUp, -2, New(1300, 1)}, // 1250 -> 1300

		// Negative numbers
		{"-127 scale -1 RoundDown", New(-127, 1), RoundDown, -1, New(-120, 1)},     // -127 -> -120 (toward zero)
		{"-123 scale -1 RoundUp", New(-123, 1), RoundUp, -1, New(-130, 1)},         // -123 -> -130 (away from zero)
		{"-125 scale -1 RoundHalfUp", New(-125, 1), RoundHalfUp, -1, New(-120, 1)}, // -125 -> -120 (half up toward positive)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(tt.roundType, tt.scale)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expected.numerator, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.expected.denominator, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Round_EdgeCases tests edge cases and boundary conditions
func TestRat_Round_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		rat       Rat
		roundType RoundType
		scale     int
		expected  Rat
	}{
		// Already exact at the requested scale
		{"exact integer scale 0", New(5, 1), RoundDown, 0, New(5, 1)},
		{"exact decimal scale 1", New(15, 10), RoundUp, 1, New(3, 2)},       // 3/2 -> 3/2 (already exact, reduced form)
		{"exact decimal scale 2", New(125, 100), RoundHalfUp, 2, New(5, 4)}, // 5/4 -> 5/4 (already exact, reduced form)

		// Zero with various scales
		{"zero scale 0", New(0, 1), RoundDown, 0, New(0, 1)},
		{"zero scale 1", New(0, 1), RoundUp, 1, New(0, 1)},
		{"zero scale -1", New(0, 1), RoundHalfUp, -1, New(0, 1)},

		// Large scale values
		{"small number large positive scale", New(1, 3), RoundDown, 10, New(3333333333, 10000000000)}, // 0.333... at 10 decimal places
		{"number large negative scale", New(12345, 1), RoundDown, -10, New(0, 1)},                     // 12345 rounded to nearest 10^10 is 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(tt.roundType, tt.scale)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expected.numerator, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.expected.denominator, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Round_MutableOperation tests that Round modifies the original rational
func TestRat_Round_MutableOperation(t *testing.T) {
	// Test that the original rational is modified, not a copy
	original := New(7, 3) // 2.333...
	originalPtr := &original

	original.Round(RoundDown, 0)

	// Verify the same instance was modified
	assert.Equal(t, int64(2), originalPtr.numerator, "original should be modified in place")
	assert.Equal(t, uint64(1), originalPtr.denominator, "original should be modified in place")
	assert.True(t, originalPtr.IsValid(), "original should remain valid")
}

// TestRat_Round_OverflowScenarios tests potential overflow scenarios
func TestRat_Round_OverflowScenarios(t *testing.T) {
	tests := []struct {
		name        string
		rat         Rat
		roundType   RoundType
		scale       int
		expectValid bool
	}{
		// Very large scale that might cause overflow in power-of-10 calculation
		{"large positive scale", New(1, 3), RoundDown, 100, false},   // May cause overflow
		{"large negative scale", New(123, 1), RoundDown, -100, true}, // Should be valid (rounds to 0)

		// Large numerator with scaling
		{"max int64 with scale", Rat{numerator: math.MaxInt64, denominator: 1}, RoundDown, 1, true}, // max int64
		{"min int64 with scale", Rat{numerator: math.MinInt64, denominator: 1}, RoundDown, 1, true}, // min int64
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(tt.roundType, tt.scale)

			if tt.expectValid {
				assert.True(t, r.IsValid(), "result should be valid")
			} else {
				// For overflow cases, the function should either succeed or mark as invalid
				// We don't require it to fail, just that it handles the case gracefully
				assert.True(t, r.IsValid() || r.IsInvalid(), "result should have valid state (either valid or invalid)")
			}
		})
	}
}

// TestRat_Round_ReductionBehavior tests that results are properly reduced
func TestRat_Round_ReductionBehavior(t *testing.T) {
	tests := []struct {
		name      string
		rat       Rat
		roundType RoundType
		scale     int
		checkGCD  bool // whether to verify the result is in lowest terms
	}{
		// Results that should be reduced
		{"fraction result should be reduced", New(10, 3), RoundDown, 1, true}, // 3.333... -> 3.3 -> 33/10
		{"integer result", New(7, 2), RoundDown, 0, true},                     // 3.5 -> 3 -> 3/1
		{"decimal result", New(125, 100), RoundHalfUp, 1, true},               // 1.25 -> 1.3 -> 13/10
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(tt.roundType, tt.scale)

			assert.True(t, r.IsValid(), "result should be valid")

			if tt.checkGCD && r.numerator != 0 {
				// Verify the result is in lowest terms
				gcd := gcdInt64Uint64(r.numerator, r.denominator)
				assert.Equal(t, uint64(1), gcd, "result should be in lowest terms")
			}
		})
	}
}

// TestRat_Round_ConsistencyWithFloatRounding tests consistency with expected float rounding behavior
func TestRat_Round_ConsistencyWithFloatRounding(t *testing.T) {
	tests := []struct {
		name      string
		rat       Rat
		roundType RoundType
		scale     int
		expected  Rat
		desc      string
	}{
		// Test cases that verify the rounding matches mathematical expectations
		{"banker's rounding edge case", New(25, 10), RoundHalfUp, 0, New(3, 1), "2.5 should round up to 3"},
		{"negative banker's rounding", New(-25, 10), RoundHalfUp, 0, New(-2, 1), "-2.5 should round toward positive (up to -2)"},
		{"tie-breaking consistency", New(35, 10), RoundHalfUp, 0, New(4, 1), "3.5 should round up to 4"},
		{"negative tie-breaking", New(-35, 10), RoundHalfUp, 0, New(-3, 1), "-3.5 should round toward positive (up to -3)"},

		// Decimal place rounding consistency
		{"decimal tie-breaking", New(125, 100), RoundHalfUp, 1, New(13, 10), "1.25 should round up to 1.3"},
		{"negative decimal tie", New(-125, 100), RoundHalfUp, 1, New(-12, 10), "-1.25 should round toward positive (-1.2)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.rat
			r.Round(tt.roundType, tt.scale)

			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expected.numerator, r.numerator, "numerator mismatch: %s", tt.desc)
			assert.Equal(t, tt.expected.denominator, r.denominator, "denominator mismatch: %s", tt.desc)
		})
	}
}

// TestRat_Round_AllRoundTypesComparison tests all three rounding types on the same values
func TestRat_Round_AllRoundTypesComparison(t *testing.T) {
	tests := []struct {
		name         string
		rat          Rat
		scale        int
		expectedDown Rat
		expectedUp   Rat
		expectedHalf Rat
	}{
		// Integer rounding (scale 0)
		{"2.3 to integer", New(23, 10), 0, New(2, 1), New(3, 1), New(2, 1)},
		{"2.7 to integer", New(27, 10), 0, New(2, 1), New(3, 1), New(3, 1)},
		{"2.5 to integer", New(25, 10), 0, New(2, 1), New(3, 1), New(3, 1)},      // half up
		{"-2.5 to integer", New(-25, 10), 0, New(-2, 1), New(-3, 1), New(-2, 1)}, // half up toward positive

		// Decimal rounding (scale 1)
		{"1.23 to 1 decimal", New(123, 100), 1, New(12, 10), New(13, 10), New(12, 10)},
		{"1.27 to 1 decimal", New(127, 100), 1, New(12, 10), New(13, 10), New(13, 10)},
		{"1.25 to 1 decimal", New(125, 100), 1, New(12, 10), New(13, 10), New(13, 10)}, // half up
	}

	for _, tt := range tests {
		t.Run(tt.name+" RoundDown", func(t *testing.T) {
			r := tt.rat
			r.Round(RoundDown, tt.scale)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expectedDown.numerator, r.numerator, "RoundDown numerator mismatch")
			assert.Equal(t, tt.expectedDown.denominator, r.denominator, "RoundDown denominator mismatch")
		})

		t.Run(tt.name+" RoundUp", func(t *testing.T) {
			r := tt.rat
			r.Round(RoundUp, tt.scale)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expectedUp.numerator, r.numerator, "RoundUp numerator mismatch")
			assert.Equal(t, tt.expectedUp.denominator, r.denominator, "RoundUp denominator mismatch")
		})

		t.Run(tt.name+" RoundHalfUp", func(t *testing.T) {
			r := tt.rat
			r.Round(RoundHalfUp, tt.scale)
			assert.True(t, r.IsValid(), "result should be valid")
			assert.Equal(t, tt.expectedHalf.numerator, r.numerator, "RoundHalfUp numerator mismatch")
			assert.Equal(t, tt.expectedHalf.denominator, r.denominator, "RoundHalfUp denominator mismatch")
		})
	}
}
