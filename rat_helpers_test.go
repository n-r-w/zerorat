package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test_willOverflowUint64Mul tests overflow detection for uint64 multiplication
func Test_willOverflowUint64Mul(t *testing.T) {
	tests := []struct {
		name     string
		a, b     uint64
		expected bool
	}{
		// Zero cases - should never overflow
		{"zero * anything", 0, 100, false},
		{"anything * zero", 100, 0, false},
		{"zero * zero", 0, 0, false},

		// Small values - should not overflow
		{"small values", 100, 200, false},
		{"medium values", 1000000, 1000000, false},

		// Edge cases that should overflow
		{"MaxUint64 * 2", math.MaxUint64, 2, true},
		{"2 * MaxUint64", 2, math.MaxUint64, true},
		{"MaxUint64 * MaxUint64", math.MaxUint64, math.MaxUint64, true},

		// Boundary cases
		{"sqrt(MaxUint64) * sqrt(MaxUint64)", 4294967296, 4294967296, true}, // 2^32 * 2^32 = 2^64 > MaxUint64
		{"large but safe", 1000000000, 18, false},                           // 18 billion, within uint64 range
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := willOverflowUint64Mul(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "overflow detection mismatch")
		})
	}
}

// Test_willOverflowInt64Mul tests overflow detection for int64 multiplication
func Test_willOverflowInt64Mul(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int64
		expected bool
	}{
		// Zero cases - should never overflow
		{"zero * anything", 0, 100, false},
		{"anything * zero", 100, 0, false},
		{"zero * zero", 0, 0, false},

		// Small values - should not overflow
		{"small positive", 100, 200, false},
		{"small negative", -100, 200, false},
		{"both negative", -100, -200, false},

		// Edge cases that should overflow
		{"MaxInt64 * 2", math.MaxInt64, 2, true},
		{"2 * MaxInt64", 2, math.MaxInt64, true},
		{"MinInt64 * 2", math.MinInt64, 2, true},
		{"MinInt64 * -1", math.MinInt64, -1, true}, // Special case: -MinInt64 overflows

		// MinInt64 special cases
		{"positive * MinInt64", 2, math.MinInt64, true},
		{"negative * MinInt64", -1, math.MinInt64, true},

		// Boundary cases
		{"large but safe positive", 1000000000, 9, false}, // 9 billion, within int64 range
		{"large but safe negative", -1000000000, 9, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := willOverflowInt64Mul(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "overflow detection mismatch")
		})
	}
}

// Test_willOverflowInt64Add tests overflow detection for int64 addition
func Test_willOverflowInt64Add(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int64
		expected bool
	}{
		// Normal cases - should not overflow
		{"small positive", 100, 200, false},
		{"small negative", -100, -200, false},
		{"mixed signs", 100, -50, false},
		{"zero cases", 0, math.MaxInt64, false},

		// Edge cases that should overflow
		{"MaxInt64 + 1", math.MaxInt64, 1, true},
		{"MaxInt64 + MaxInt64", math.MaxInt64, math.MaxInt64, true},
		{"MinInt64 + (-1)", math.MinInt64, -1, true},
		{"MinInt64 + MinInt64", math.MinInt64, math.MinInt64, true},

		// Boundary cases
		{"MaxInt64 + 0", math.MaxInt64, 0, false},
		{"MinInt64 + 0", math.MinInt64, 0, false},
		{"MaxInt64 + (-1)", math.MaxInt64, -1, false},
		{"MinInt64 + 1", math.MinInt64, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := willOverflowInt64Add(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "overflow detection mismatch")
		})
	}
}

// Test_willOverflowInt64Sub tests overflow detection for int64 subtraction
func Test_willOverflowInt64Sub(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int64
		expected bool
	}{
		// Normal cases - should not overflow
		{"small positive", 200, 100, false},
		{"small negative", -100, -200, false},
		{"mixed signs", 100, -50, false},

		// Edge cases that should overflow
		{"MinInt64 - 1", math.MinInt64, 1, true},
		{"MaxInt64 - (-1)", math.MaxInt64, -1, true},
		{"MinInt64 - MaxInt64", math.MinInt64, math.MaxInt64, true},

		// Boundary cases
		{"MaxInt64 - 0", math.MaxInt64, 0, false},
		{"MinInt64 - 0", math.MinInt64, 0, false},
		{"MaxInt64 - 1", math.MaxInt64, 1, false},
		{"MinInt64 - (-1)", math.MinInt64, -1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := willOverflowInt64Sub(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "overflow detection mismatch")
		})
	}
}

// Test_mulInt64ByUint64ToInt64 tests safe multiplication with overflow detection
func Test_mulInt64ByUint64ToInt64(t *testing.T) {
	tests := []struct {
		name     string
		a        int64
		b        uint64
		expected int64
		shouldOK bool
	}{
		// Zero cases
		{"zero int64", 0, 123, 0, true},
		{"zero uint64", 123, 0, 0, true},
		{"both zero", 0, 0, 0, true},

		// Normal cases
		{"positive * positive", 7, 9, 63, true},
		{"negative * positive", -7, 9, -63, true},

		// Edge cases
		{"MaxInt64 * 1", math.MaxInt64, 1, math.MaxInt64, true},
		{"MinInt64 * 1", math.MinInt64, 1, math.MinInt64, true},

		// Special MinInt64 case
		{"negative exact MinInt64", -1, uint64(math.MaxInt64) + 1, math.MinInt64, true},

		// Overflow cases
		{"MinInt64 * 2", math.MinInt64, 2, 0, false},
		{"MaxInt64 * 2", math.MaxInt64, 2, 0, false},
		{"negative overflow", -3, 1 << 62, 0, false},
		{"positive overflow", 3, 1 << 62, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := mulInt64ByUint64ToInt64(tt.a, tt.b)
			assert.Equal(t, tt.shouldOK, ok, "overflow detection mismatch")
			if tt.shouldOK {
				assert.Equal(t, tt.expected, result, "result value mismatch")
			}
		})
	}
}

// Test_gcdUint64 tests GCD calculation for uint64 values
func Test_gcdUint64(t *testing.T) {
	tests := []struct {
		name     string
		a, b     uint64
		expected uint64
	}{
		{"gcd(0, 5)", 0, 5, 5},
		{"gcd(5, 0)", 5, 0, 5},
		{"gcd(0, 0)", 0, 0, 0},
		{"gcd(1, 1)", 1, 1, 1},
		{"gcd(12, 8)", 12, 8, 4},
		{"gcd(17, 13)", 17, 13, 1}, // coprime
		{"gcd(48, 18)", 48, 18, 6},
		{"gcd(100, 75)", 100, 75, 25},
		{"large values", 123456789, 987654321, 9},
		{"powers of 2", 64, 128, 64},
		{"MaxUint64 and 1", math.MaxUint64, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gcdUint64(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "GCD calculation mismatch")
		})
	}
}

// Test_gcdInt64Uint64 tests GCD calculation between int64 and uint64
func Test_gcdInt64Uint64(t *testing.T) {
	tests := []struct {
		name     string
		a        int64
		b        uint64
		expected uint64
	}{
		{"positive int64", 12, 8, 4},
		{"negative int64", -12, 8, 4},
		{"zero int64", 0, 5, 5},
		{"int64 with zero uint64", 5, 0, 5},
		{"MinInt64 special case", math.MinInt64, 2, 2},
		{"large negative", -123456789, 987654321, 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gcdInt64Uint64(tt.a, tt.b)
			assert.Equal(t, tt.expected, result, "GCD calculation mismatch")
		})
	}
}

// Test_absInt64ToUint64 tests absolute value conversion from int64 to uint64
func Test_absInt64ToUint64(t *testing.T) {
	tests := []struct {
		name     string
		input    int64
		expected uint64
	}{
		{"positive value", 123, 123},
		{"negative value", -456, 456},
		{"zero", 0, 0},
		{"MaxInt64", math.MaxInt64, uint64(math.MaxInt64)},
		{"MinInt64 special case", math.MinInt64, uint64(math.MaxInt64) + 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := absInt64ToUint64(tt.input)
			assert.Equal(t, tt.expected, result, "absolute value conversion mismatch")
		})
	}
}

// Test_uint64ToInt64WithSign tests conversion from uint64 to int64 with sign
func Test_uint64ToInt64WithSign(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
		negative bool
		expected int64
		shouldOK bool
	}{
		// Positive cases
		{"positive small", 123, false, 123, true},
		{"positive MaxInt64", uint64(math.MaxInt64), false, math.MaxInt64, true},
		{"positive overflow", uint64(math.MaxInt64) + 1, false, 0, false},
		{"positive MaxUint64", math.MaxUint64, false, 0, false},

		// Negative cases
		{"negative small", 123, true, -123, true},
		{"negative at limit", uint64(math.MaxInt64) + 1, true, math.MinInt64, true},
		{"negative above limit", uint64(math.MaxInt64) + 2, true, 0, false},
		{"negative MaxUint64", math.MaxUint64, true, 0, false},

		// Zero case
		{"zero positive", 0, false, 0, true},
		{"zero negative", 0, true, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := uint64ToInt64WithSign(tt.value, tt.negative)
			assert.Equal(t, tt.shouldOK, ok, "conversion success mismatch")
			if tt.shouldOK {
				assert.Equal(t, tt.expected, result, "converted value mismatch")
			}
		})
	}
}

// Test_compareRationalsCrossMul tests overflow-safe cross multiplication comparison
func Test_compareRationalsCrossMul(t *testing.T) {
	tests := []struct {
		name     string
		aNum     int64
		aDenom   uint64
		cNum     int64
		cDenom   uint64
		expected int
	}{
		{"equal fractions", 1, 2, 2, 4, 0},
		{"first smaller", 1, 3, 1, 2, -1},
		{"first larger", 2, 3, 1, 2, 1},
		{"negative vs positive", -1, 2, 1, 2, -1},
		{"both negative", -2, 3, -1, 2, -1},
		{"zero vs positive", 0, 1, 1, 2, -1},
		{"positive vs zero", 1, 2, 0, 1, 1},
		{"both zero", 0, 1, 0, 2, 0},

		// Overflow cases that require 128-bit arithmetic
		{"overflow case 1", math.MaxInt64 - 1, math.MaxUint64, math.MaxInt64, math.MaxUint64, -1},
		{"overflow case 2", math.MaxInt64, math.MaxUint64, math.MaxInt64 - 1, math.MaxUint64, 1},
		{"large equal", 1000000000000000000, 2000000000000000000, 500000000000000000, 1000000000000000000, 0},

		// MinInt64 special cases
		{"MinInt64 vs negative", math.MinInt64, 1, -1, 1, -1},
		{"negative vs MinInt64", -1, 1, math.MinInt64, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareRationalsCrossMul(tt.aNum, tt.aDenom, tt.cNum, tt.cDenom)
			assert.Equal(t, tt.expected, result, "comparison result mismatch")
		})
	}
}
