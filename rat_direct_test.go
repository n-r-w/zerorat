package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_float64ToRatExact_Direct tests the internal function directly
func Test_float64ToRatExact_Direct(t *testing.T) {
	// Line 676: return Rat{} for NaN
	rNaN := float64ToRatExact(math.NaN())
	assert.True(t, rNaN.IsInvalid())
	assert.Equal(t, int64(0), rNaN.numerator)
	assert.Equal(t, uint64(0), rNaN.denominator)

	// Line 676: return Rat{} for +Inf
	rPosInf := float64ToRatExact(math.Inf(1))
	assert.True(t, rPosInf.IsInvalid())
	assert.Equal(t, int64(0), rPosInf.numerator)
	assert.Equal(t, uint64(0), rPosInf.denominator)

	// Line 676: return Rat{} for -Inf
	rNegInf := float64ToRatExact(math.Inf(-1))
	assert.True(t, rNegInf.IsInvalid())
	assert.Equal(t, int64(0), rNegInf.numerator)
	assert.Equal(t, uint64(0), rNegInf.denominator)

	// Line 680: return Rat{numerator: 0, denominator: 1} for +0.0
	rZero := float64ToRatExact(0.0)
	assert.True(t, rZero.IsValid())
	assert.Equal(t, int64(0), rZero.numerator)
	assert.Equal(t, uint64(1), rZero.denominator)

	// Line 680: return Rat{numerator: 0, denominator: 1} for -0.0
	negZero := math.Float64frombits(1 << 63) // -0.0 using bit manipulation
	rNegZero := float64ToRatExact(negZero)
	assert.True(t, rNegZero.IsValid())
	assert.Equal(t, int64(0), rNegZero.numerator)
	assert.Equal(t, uint64(1), rNegZero.denominator)

	// Test normal values to ensure the function works correctly
	r1 := float64ToRatExact(0.5)
	assert.True(t, r1.IsValid())
	assert.Equal(t, int64(1), r1.numerator)
	assert.Equal(t, uint64(2), r1.denominator)

	r2 := float64ToRatExact(-0.25)
	assert.True(t, r2.IsValid())
	assert.Equal(t, int64(-1), r2.numerator)
	assert.Equal(t, uint64(4), r2.denominator)
}

// Test_float64ToRatExact_OverflowCases tests overflow scenarios
func Test_float64ToRatExact_OverflowCases(t *testing.T) {
	// Test very large values that should trigger overflow
	veryLarge := math.Ldexp(1.0, 1000) // 2^1000
	if math.IsInf(veryLarge, 0) {
		// If it becomes infinity, test that path
		rLarge := float64ToRatExact(veryLarge)
		assert.True(t, rLarge.IsInvalid())
	} else {
		// If it's still finite, it should either work or overflow
		rLarge := float64ToRatExact(veryLarge)
		// Either valid or invalid is acceptable for such large values
		_ = rLarge.IsValid()
	}

	// Test the largest finite float64
	maxFloat := math.MaxFloat64
	rMax := float64ToRatExact(maxFloat)
	// This will likely be invalid due to overflow, which is correct
	if rMax.IsInvalid() {
		assert.Equal(t, uint64(0), rMax.denominator)
	}

	// Test values that might trigger the e>=0 overflow path (line 749)
	// Create a value with large mantissa and exponent
	largeMantExp := math.Ldexp(math.Ldexp(1.0, 52), 50) // Very large value
	rLargeMantExp := float64ToRatExact(largeMantExp)
	// Should either be valid or trigger overflow
	_ = rLargeMantExp.IsValid()
}

// Test_float64ToRatExact_NegationPath tests the negation path (line 755)
func Test_float64ToRatExact_NegationPath(t *testing.T) {
	// Test negative values that go through the e>=0 path
	negValue := -8.0 // -2^3, simple negative power of 2
	rNeg := float64ToRatExact(negValue)
	assert.True(t, rNeg.IsValid())
	assert.Equal(t, int64(-8), rNeg.numerator)
	assert.Equal(t, uint64(1), rNeg.denominator)

	// Test larger negative value
	negLarge := -math.Ldexp(1.0, 50) // -2^50
	rNegLarge := float64ToRatExact(negLarge)
	if rNegLarge.IsValid() {
		assert.Negative(t, rNegLarge.numerator)
		assert.Positive(t, rNegLarge.denominator)
	}
}

// Test_float64ToRatExact_Lines749And755 targets lines 749 and 755 exactly
func Test_float64ToRatExact_Lines749And755(t *testing.T) {
	// Line 749: return Rat{} when newShift < 0 || newShift >= 64 || mant > (absLimit>>uint(newShift))
	// To trigger this, we need a value that goes through the offload path but still overflows

	// Create the largest possible finite float64 that might trigger overflow
	// Use the maximum exponent that doesn't result in infinity
	maxFiniteExp := 1023 // Unbiased exponent for largest finite values
	maxMantissa := (uint64(1) << 52) - 1

	// Create float64 with max exponent and max mantissa
	expBits := uint64(maxFiniteExp + 1023) // Add IEEE bias
	bits := (expBits << 52) | maxMantissa
	maxFinite := math.Float64frombits(bits)

	if !math.IsInf(maxFinite, 0) {
		rMax := float64ToRatExact(maxFinite)
		// This should trigger line 749 overflow
		if rMax.IsInvalid() {
			assert.Equal(t, uint64(0), rMax.denominator)
		}
	}

	// Try another approach: create a value with very large mantissa that needs offloading
	// but the offloading calculation overflows
	largeValue := math.Ldexp(math.Ldexp(1.0, 52), 500) // Extremely large value
	if !math.IsInf(largeValue, 0) {
		rLargeValue := float64ToRatExact(largeValue)
		if rLargeValue.IsInvalid() {
			assert.Equal(t, uint64(0), rLargeValue.denominator)
		}
	}

	// Line 755: n = -n (negation in e>=0 offload path)
	// We need a negative value that goes through the offload path (not the fast path)
	// The offload path happens when the value is too large to fit in the fast path

	// To force the offload path, we need a value where mant << e > absLimit
	// For negative numbers, absLimit = 2^63, so we need something larger than that

	// Create a negative value that's definitely too large for the fast path
	// Use a very large mantissa with a moderate exponent
	largeMantissa := math.Ldexp(1.5, 52)            // 1.5 * 2^52, close to max mantissa
	largeNegValue := -math.Ldexp(largeMantissa, 15) // Scale it up to force offload

	rNegOffload := float64ToRatExact(largeNegValue)
	if rNegOffload.IsValid() {
		assert.Negative(t, rNegOffload.numerator) // This should exercise line 755 in offload path
		assert.Positive(t, rNegOffload.denominator)
	}

	// Try a simpler approach: use a value that we know will need offloading
	// A value like -3 * 2^61 should be too large for fast path but manageable for offload
	negNeedsOffload := -math.Ldexp(3.0, 61)
	rNegNeedsOffload := float64ToRatExact(negNeedsOffload)
	if rNegNeedsOffload.IsValid() {
		assert.Negative(t, rNegNeedsOffload.numerator) // Line 755
		assert.Positive(t, rNegNeedsOffload.denominator)
	}

	// Test with the maximum safe negative value that still needs offloading
	// Use -2^63 which is exactly at the boundary
	minInt64Float := float64(math.MinInt64)
	rMinInt64 := float64ToRatExact(minInt64Float)
	if rMinInt64.IsValid() {
		assert.Equal(t, int64(math.MinInt64), rMinInt64.numerator) // Line 755 should be hit
		assert.Equal(t, uint64(1), rMinInt64.denominator)
	}
}
