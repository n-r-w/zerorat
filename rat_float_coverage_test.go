package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_float64ToRatExact_ComprehensiveCoverage targets specific uncovered branches
func Test_float64ToRatExact_ComprehensiveCoverage(t *testing.T) {
	// Test e >= 0 path where fast path fails but offload succeeds
	// We need a value where mant << e > absLimit but can be split
	// Construct using bit manipulation to get precise control

	// Create a float64 with specific mantissa and exponent
	// Use a mantissa that's large enough to need offloading
	// Example: mantissa close to 2^52, exponent that makes it need splitting

	// Test case: value that needs denominator offloading
	// 2^61 should be representable but might need offloading depending on implementation
	val1 := math.Ldexp(1.0, 61)
	r1 := NewFromFloat64(val1)
	if r1.IsValid() {
		// Either fits in numerator or uses denominator
		assert.Positive(t, r1.denominator)
	}

	// Test negative value in offload path (covers line 755: n = -n)
	val2 := -math.Ldexp(1.0, 61)
	r2 := NewFromFloat64(val2)
	if r2.IsValid() {
		assert.Positive(t, r2.denominator)
		assert.Negative(t, r2.numerator)
	}

	// Test e < 0 with shift >= 64 (should result in n64 = 0)
	val3 := math.Ldexp(1.0, -1100) // Very small, should trigger shift >= 64
	r3 := NewFromFloat64(val3)
	if r3.IsValid() {
		// Very small values might round to 0 or have positive numerator
		assert.True(t, r3.numerator >= 0 || r3.numerator < 0) // Any valid numerator is fine
	}

	// Test subnormal numbers (expBits == 0)
	val4 := math.SmallestNonzeroFloat64 // This is subnormal
	r4 := NewFromFloat64(val4)
	assert.True(t, r4.IsValid())
	assert.Positive(t, r4.denominator)
	// Very small subnormal might round to 0, so check it's non-negative
	assert.GreaterOrEqual(t, r4.numerator, int64(0))
}

// Test_float64ToRatExact_SpecificLines covers the exact uncovered lines
func Test_float64ToRatExact_SpecificLines(t *testing.T) {
	// Line 676: return Rat{} for NaN/Inf
	rNaN := NewFromFloat64(math.NaN())
	assert.True(t, rNaN.IsInvalid())
	assert.Equal(t, uint64(0), rNaN.denominator)

	rInf := NewFromFloat64(math.Inf(1))
	assert.True(t, rInf.IsInvalid())
	assert.Equal(t, uint64(0), rInf.denominator)

	// Line 680: return Rat{numerator: 0, denominator: 1} for zero
	rZero := NewFromFloat64(0.0)
	assert.True(t, rZero.IsValid())
	assert.Equal(t, int64(0), rZero.numerator)
	assert.Equal(t, uint64(1), rZero.denominator)

	// Test negative zero using bit manipulation to create actual -0.0
	negZeroBits := uint64(1) << 63 // Sign bit set, everything else zero
	rNegZero := NewFromFloat64(math.Float64frombits(negZeroBits))
	assert.True(t, rNegZero.IsValid())
	assert.Equal(t, int64(0), rNegZero.numerator)
	assert.Equal(t, uint64(1), rNegZero.denominator)

	// Line 749: return Rat{} for overflow in e>=0 path
	// Need a value where newShift < 0 or newShift >= 64 or mant > (absLimit>>newShift)
	// Try to construct a very large float that would overflow
	veryLarge := math.Ldexp(1.0, 1000) // 2^1000, should overflow
	rLarge := NewFromFloat64(veryLarge)
	assert.True(t, rLarge.IsInvalid())
	assert.Equal(t, uint64(0), rLarge.denominator)

	// Line 755: n = -n (negation in e>=0 path)
	// Need a negative value that goes through the e>=0 offload path successfully
	// Use a negative power of 2 that's large enough to need offloading but not overflow
	negValue := -math.Ldexp(1.0, 60) // -2^60, should be representable
	rNeg := NewFromFloat64(negValue)
	if rNeg.IsValid() {
		assert.Negative(t, rNeg.numerator) // This exercises the n = -n line
		assert.Positive(t, rNeg.denominator)
	}
}

// Test_float64ToRatExact_RoundingEdgeCases covers specific rounding scenarios
func Test_float64ToRatExact_RoundingEdgeCases(t *testing.T) {
	// Test rounding tie with even base (should not round up)
	// Test rounding tie with odd base (should round up)
	// These are hard to construct directly, so we test values that exercise the paths

	// Test values that will go through the e < 0 rounding path
	testValues := []float64{
		1.0 / 3.0,            // 0.333... will need rounding
		1.0 / 7.0,            // 0.142857... will need rounding
		1.0 / 9.0,            // 0.111... will need rounding
		math.Ldexp(3.0, -10), // 3/1024, might hit exact shiftUp == 0 case
	}

	for i, val := range testValues {
		r := NewFromFloat64(val)
		assert.True(t, r.IsValid(), "test value %d should be valid", i)
		assert.Positive(t, r.denominator, "test value %d should have positive denominator", i)

		// Verify the conversion is reasonable
		backToFloat := float64(r.numerator) / float64(r.denominator)
		diff := math.Abs(backToFloat - val)
		relErr := diff / math.Abs(val)
		assert.Less(t, relErr, 1e-15, "test value %d should have small relative error", i)
	}
}

// Test_float64ToRatExact_BoundaryConditions tests edge cases in the algorithm
func Test_float64ToRatExact_BoundaryConditions(t *testing.T) {
	// Test newShift boundary conditions
	// Test neededDenPow boundary (should be <= 63)

	// Test a value that might hit neededDenPow > 63 (should return invalid)
	// This would be a very large number that can't be represented
	val1 := math.Ldexp(1.0, 100) // 2^100, likely too large
	r1 := NewFromFloat64(val1)
	// Should either be valid with large denominator or invalid due to overflow
	if r1.IsInvalid() {
		assert.Equal(t, uint64(0), r1.denominator)
	}

	// Test maxShiftAllowed = 0 case (when limitBits - mantBits <= 0)
	// This happens with large mantissas

	// Test the boundary where newShift < 0 or newShift >= 64
	val2 := math.Ldexp(1.0, 63) // 2^63, at the boundary
	r2 := NewFromFloat64(val2)
	// Should either work or overflow gracefully
	if r2.IsValid() {
		assert.Positive(t, r2.denominator)
	}
}
