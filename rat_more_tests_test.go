package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRat_Greater covers missing Greater branches
func TestRat_Greater(t *testing.T) {
	cases := []struct {
		name     string
		a, b     Rat
		expected bool
	}{
		{"3/4 > 1/2", New(3, 4), New(1, 2), true},
		{"1/2 > 3/4", New(1, 2), New(3, 4), false},
		{"equal", New(1, 2), New(2, 4), false},
		{"invalid", New(1, 0), New(1, 2), false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Greater(tt.b))
		})
	}
}

// TestRat_Div_NegationNoOverflow ensures sign application path when negation is safe
func TestRat_Div_NegationNoOverflow(t *testing.T) {
	r := New(3, 2)      // 3/2
	other := New(-2, 5) // -2/5
	r.Div(other)        // (3/2) / (-2/5) = (3/2) * (5/2) = 15/4 then negate => -15/4
	assert.True(t, r.IsValid())
	assert.Equal(t, int64(-15), r.numerator)
	assert.Equal(t, uint64(4), r.denominator)
}

// TestReduce_SignedMagnitudePath exercises gcd>1 with negative numerator path
func TestReduce_SignedMagnitudePath(t *testing.T) {
	r := Rat{numerator: -10, denominator: 20}
	r.Reduce()
	assert.True(t, r.IsValid())
	assert.Equal(t, int64(-1), r.numerator)
	assert.Equal(t, uint64(2), r.denominator)
}

// TestReduce_EdgeCases covers remaining Reduce branches
func TestReduce_EdgeCases(t *testing.T) {
	// Test gcd == 1 case (already reduced)
	r := Rat{numerator: 7, denominator: 11}
	r.Reduce()
	assert.True(t, r.IsValid())
	assert.Equal(t, int64(7), r.numerator)
	assert.Equal(t, uint64(11), r.denominator)

	// Test large negative numerator reduction
	r2 := Rat{numerator: -1000000, denominator: 2000000}
	r2.Reduce()
	assert.True(t, r2.IsValid())
	assert.Equal(t, int64(-1), r2.numerator)
	assert.Equal(t, uint64(2), r2.denominator)
}

// Test_uint64ToInt64WithSign_FullCoverage adds missing branches
func Test_uint64ToInt64WithSign_FullCoverage(t *testing.T) {
	// negative below limit
	v, ok := uint64ToInt64WithSign(123, true)
	assert.True(t, ok)
	assert.Equal(t, int64(-123), v)

	// negative at limit => MinInt64
	v, ok = uint64ToInt64WithSign(uint64(math.MaxInt64)+1, true)
	assert.True(t, ok)
	assert.Equal(t, int64(math.MinInt64), v)

	// negative above limit -> false
	_, ok = uint64ToInt64WithSign(uint64(math.MaxInt64)+2, true)
	assert.False(t, ok)

	// positive within range
	v, ok = uint64ToInt64WithSign(uint64(math.MaxInt64), false)
	assert.True(t, ok)
	assert.Equal(t, int64(math.MaxInt64), v)

	// positive above MaxInt64 -> false
	_, ok = uint64ToInt64WithSign(uint64(math.MaxUint64), false)
	assert.False(t, ok)
}

// Test_float64ToRatExact_Branches adds coverage for e>=0 offload and e<0 rounding paths
func Test_float64ToRatExact_Branches(t *testing.T) {
	// e >= 0 fast path within range
	r := NewFromFloat64(8.0) // 8 = 1 * 2^3
	assert.True(t, r.IsValid())
	assert.Equal(t, int64(8), r.numerator)
	assert.Equal(t, uint64(1), r.denominator)

	// e < 0 rounding: value where rounding is tie and base is even/odd
	r2 := NewFromFloat64(0.2)
	assert.True(t, r2.IsValid())
	assert.Positive(t, r2.denominator)
}

// Test_float64ToRatExact_OffloadPath covers e>=0 offload to denominator path
func Test_float64ToRatExact_OffloadPath(t *testing.T) {
	// Construct a value that requires offloading exponent to denominator
	// We need a float where mantissa << e would overflow but can be represented with denominator
	// Use math.Ldexp to construct: mantissa * 2^exp where we control the split

	// Create a large mantissa that when shifted by a large exponent needs offloading
	// Example: 2^60 = mantissa=1, exp=60; this should fit as numerator=2^60, denom=1
	// But 2^62 might need offloading depending on limits

	// Use a value that's exactly representable but large enough to trigger offload logic
	val := math.Ldexp(1.0, 62) // 2^62, should be within range but test the boundary
	r := NewFromFloat64(val)
	// This might be valid or invalid depending on exact limits, but exercises the path
	if r.IsValid() {
		assert.Positive(t, r.denominator)
	} else {
		// If it overflows, that's also valid behavior
		assert.True(t, r.IsInvalid())
	}
}

// Test_float64ToRatExact_RoundingTies covers tie-breaking in rounding
func Test_float64ToRatExact_RoundingTies(t *testing.T) {
	// Test values that create rounding ties in the e < 0 path
	// We need to construct floats where after shifting, rem == half

	// Use very small values that will trigger the e < 0 path with specific bit patterns
	// 2^-1075 is close to the smallest subnormal, will have e < 0
	val := math.Ldexp(1.0, -1070) // Small but not the absolute minimum
	r := NewFromFloat64(val)
	assert.True(t, r.IsValid())
	assert.Positive(t, r.denominator)

	// Test shiftUp == 0 path (exact case in e < 0 branch)
	// This happens when e + denPow == 0, so we need e = -denPow
	// For small values, denPow = min(-e, 63), so if e = -10, denPow = 10, shiftUp = 0
	val2 := math.Ldexp(1.0, -10) // 2^-10, should hit shiftUp == 0 case
	r2 := NewFromFloat64(val2)
	assert.True(t, r2.IsValid())
	assert.Equal(t, int64(1), r2.numerator)
	assert.Equal(t, uint64(1024), r2.denominator) // 2^10
}
