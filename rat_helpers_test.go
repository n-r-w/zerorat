package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// assertSigned128Equal verifies the sign and magnitude of an internal signed128 value.
func assertSigned128Equal(t *testing.T, expected signed128, actual signed128) {
	t.Helper()
	assert.Equal(t, expected.negative, actual.negative, "sign mismatch")
	assert.Equal(t, expected.hi, actual.hi, "high word mismatch")
	assert.Equal(t, expected.lo, actual.lo, "low word mismatch")
}

// Test_signed128BasicHelpers covers sign normalization and negation invariants.
func Test_signed128BasicHelpers(t *testing.T) {
	t.Run("negative zero normalizes to positive zero", func(t *testing.T) {
		value := signed128{negative: true}.normalize()

		assert.False(t, value.negative, "zero must not keep a negative sign")
		assert.True(t, value.isZero(), "normalized zero should stay zero")
	})

	t.Run("negating zero keeps positive zero", func(t *testing.T) {
		value := negateSigned128(signed128{negative: true})

		assertSigned128Equal(t, signed128{}, value)
	})

	t.Run("negating non-zero flips sign", func(t *testing.T) {
		value := negateSigned128(signed128{lo: 7})

		assertSigned128Equal(t, signed128{negative: true, lo: 7}, value)
	})
}

// Test_mulInt64ByUint64ToSigned128 verifies exact 128-bit products used by Add and Sub recovery.
func Test_mulInt64ByUint64ToSigned128(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		factor   uint64
		expected signed128
	}{
		{"positive product", 6, 7, signed128{lo: 42}},
		{"negative product", -6, 7, signed128{negative: true, lo: 42}},
		{"zero product clears sign", -6, 0, signed128{}},
		{"MinInt64 product crosses low word", math.MinInt64, 2, signed128{negative: true, hi: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := mulInt64ByUint64ToSigned128(tt.value, tt.factor)

			assertSigned128Equal(t, tt.expected, actual)
		})
	}
}

// Test_addSigned128 verifies signed 128-bit addition across sign and magnitude branches.
func Test_addSigned128(t *testing.T) {
	tests := []struct {
		name     string
		left     signed128
		right    signed128
		expected signed128
		ok       bool
	}{
		{"same sign positive", signed128{lo: 10}, signed128{lo: 5}, signed128{lo: 15}, true},
		{"same sign negative", signed128{negative: true, lo: 10}, signed128{negative: true, lo: 5}, signed128{negative: true, lo: 15}, true},
		{"same sign overflow", signed128{hi: math.MaxUint64, lo: math.MaxUint64}, signed128{lo: 1}, signed128{}, false},
		{"opposite signs equal magnitude", signed128{lo: 10}, signed128{negative: true, lo: 10}, signed128{}, true},
		{"opposite signs left larger", signed128{lo: 10}, signed128{negative: true, lo: 3}, signed128{lo: 7}, true},
		{"opposite signs right larger", signed128{lo: 3}, signed128{negative: true, lo: 10}, signed128{negative: true, lo: 7}, true},
		{"opposite signs with borrow", signed128{hi: 1}, signed128{negative: true, lo: 1}, signed128{lo: math.MaxUint64}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := addSigned128(tt.left, tt.right)

			assert.Equal(t, tt.ok, ok, "success mismatch")
			assertSigned128Equal(t, tt.expected, actual)
		})
	}
}

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

// Test_divInt64ByUint64Exact verifies exact signed division without unsafe MinInt64 casts.
func Test_divInt64ByUint64Exact(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		divisor  uint64
		expected int64
		ok       bool
	}{
		{"positive exact", 42, 7, 6, true},
		{"negative exact", -42, 7, -6, true},
		{"MinInt64 exact", math.MinInt64, 2, math.MinInt64 / 2, true},
		{"not exact", 10, 3, 0, false},
		{"zero divisor", 10, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, ok := divInt64ByUint64Exact(tt.value, tt.divisor)

			assert.Equal(t, tt.ok, ok, "success mismatch")
			if tt.ok {
				assert.Equal(t, tt.expected, actual, "quotient mismatch")
			}
		})
	}
}

// Test_signed128DivisionHelpers verifies modulo, gcd, and int64 conversion for 128-bit values.
func Test_signed128DivisionHelpers(t *testing.T) {
	t.Run("modulo handles small divisor", func(t *testing.T) {
		var actual uint64
		assert.NotPanics(t, func() {
			actual = modUint128ByUint64(3, 5, 2)
		}, "modulo should handle high words larger than the divisor")

		assert.Equal(t, uint64(1), actual, "modulo mismatch")
	})

	t.Run("modulo handles high word", func(t *testing.T) {
		actual := modUint128ByUint64(1, 5, 3)

		assert.Equal(t, uint64(0), actual, "modulo mismatch")
	})

	t.Run("gcd with zero divisor returns zero", func(t *testing.T) {
		actual := gcdSigned128Uint64(signed128{lo: 5}, 0)

		assert.Equal(t, uint64(0), actual, "gcd mismatch")
	})

	t.Run("gcd uses 128-bit modulo", func(t *testing.T) {
		value := signed128{hi: 1, lo: 8}
		actual := gcdSigned128Uint64(value, 12)

		assert.Equal(t, uint64(12), actual, "gcd mismatch")
	})

	t.Run("gcd of zero returns divisor", func(t *testing.T) {
		actual := gcdSigned128Uint64(signed128{}, 12)

		assert.Equal(t, uint64(12), actual, "gcd mismatch")
	})

	t.Run("division rejects zero divisor", func(t *testing.T) {
		_, ok := divSigned128ByUint64ToInt64(signed128{lo: 5}, 0)

		assert.False(t, ok, "zero divisor should not fit")
	})

	t.Run("division returns positive int64", func(t *testing.T) {
		actual, ok := divSigned128ByUint64ToInt64(signed128{hi: 1, lo: 0}, 4)

		assert.True(t, ok, "division should fit")
		assert.Equal(t, int64(1<<63/2), actual, "quotient mismatch")
	})

	t.Run("division returns MinInt64", func(t *testing.T) {
		actual, ok := divSigned128ByUint64ToInt64(signed128{negative: true, hi: 1, lo: 0}, 2)

		assert.True(t, ok, "division should fit")
		assert.Equal(t, int64(math.MinInt64), actual, "quotient mismatch")
	})

	t.Run("division rejects positive quotient above MaxInt64", func(t *testing.T) {
		_, ok := divSigned128ByUint64ToInt64(signed128{hi: 1, lo: 0}, 2)

		assert.False(t, ok, "positive quotient should not exceed MaxInt64")
	})

	t.Run("division rejects oversized quotient", func(t *testing.T) {
		_, ok := divSigned128ByUint64ToInt64(signed128{hi: 2}, 2)

		assert.False(t, ok, "oversized quotient should not fit")
	})

	t.Run("division rejects remainder", func(t *testing.T) {
		_, ok := divSigned128ByUint64ToInt64(signed128{lo: 5}, 2)

		assert.False(t, ok, "non-exact quotient should not fit")
	})
}

// Test_reduceIfLarge verifies that the size threshold controls gcd work.
func Test_reduceIfLarge(t *testing.T) {
	t.Run("small unreduced value stays unchanged", func(t *testing.T) {
		r := Rat{numerator: 2, denominator: 4}
		r.reduceIfLarge()

		assert.Equal(t, int64(2), r.numerator, "numerator should stay unchanged")
		assert.Equal(t, uint64(4), r.denominator, "denominator should stay unchanged")
	})

	t.Run("large unreduced value is compacted", func(t *testing.T) {
		r := Rat{numerator: 2, denominator: autoReduceThreshold * 2}
		r.reduceIfLarge()

		assert.Equal(t, int64(1), r.numerator, "numerator mismatch")
		assert.Equal(t, autoReduceThreshold, r.denominator, "denominator mismatch")
	})

	t.Run("invalid value remains invalid", func(t *testing.T) {
		r := Rat{numerator: 1, denominator: 0}
		r.reduceIfLarge()

		assert.True(t, r.IsInvalid(), "invalid value should remain invalid")
		assert.Equal(t, uint64(0), r.denominator, "invalid denominator should stay zero")
	})
}

// Test_scaleMagnitude verifies signed scale normalization before decimal scaling.
func Test_scaleMagnitude(t *testing.T) {
	t.Run("positive scale", func(t *testing.T) {
		magnitude, negative, ok := scaleMagnitude(7)

		assert.True(t, ok, "positive scale should be accepted")
		assert.False(t, negative, "positive scale should not be marked negative")
		assert.Equal(t, 7, magnitude, "magnitude mismatch")
	})

	t.Run("negative scale", func(t *testing.T) {
		magnitude, negative, ok := scaleMagnitude(-7)

		assert.True(t, ok, "negative scale should be accepted")
		assert.True(t, negative, "negative scale should be marked negative")
		assert.Equal(t, 7, magnitude, "magnitude mismatch")
	})

	t.Run("minimum int scale", func(t *testing.T) {
		_, _, ok := scaleMagnitude(-int(^uint(0)>>1) - 1)

		assert.False(t, ok, "minimum int scale cannot be represented as a positive int")
	})
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
