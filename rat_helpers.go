package zerorat

import (
	"math"
	"math/bits"
)

// willOverflowUint64Mul checks if multiplying two uint64 values would overflow.
// Uses math/bits for improved clarity and performance.
func willOverflowUint64Mul(a, b uint64) bool {
	if a == 0 || b == 0 {
		return false
	}
	// Use bits.Mul64 to detect overflow
	hi, _ := bits.Mul64(a, b)
	return hi != 0
}

// willOverflowInt64Mul checks if multiplying two int64 values would overflow.
// Uses math/bits for improved clarity and performance.
func willOverflowInt64Mul(a, b int64) bool {
	if a == 0 || b == 0 {
		return false
	}

	// Convert to unsigned for bits.Mul64, handling signs separately
	aAbs := absInt64ToUint64(a)
	bAbs := absInt64ToUint64(b)

	// Use bits.Mul64 to detect overflow
	hi, lo := bits.Mul64(aAbs, bAbs)

	// Check if result fits in int64 range
	sameSign := (a > 0) == (b > 0)
	if sameSign {
		// Positive result: must fit in [0, MaxInt64]
		return hi != 0 || lo > uint64(math.MaxInt64)
	}
	// Negative result: must fit in [MinInt64, -1]
	// MaxInt64 + 1 = 9223372036854775808 (absolute value of MinInt64)
	return hi != 0 || lo > 9223372036854775808
}

// mulInt64ByUint64ToInt64 multiplies an int64 by a uint64 and returns the int64 result if it fits.
// The function performs 128-bit multiplication on absolute values and then reapplies the sign safely.
// Returns (0, false) if the product would overflow int64.
func mulInt64ByUint64ToInt64(a int64, b uint64) (int64, bool) {
	if a == 0 || b == 0 {
		return 0, true
	}
	neg := a < 0
	aAbs := absInt64ToUint64(a)
	hi, lo := bits.Mul64(aAbs, b)
	if hi != 0 {
		// product >= 2^64, definitely exceeds int64 range
		return 0, false
	}
	if neg {
		limit := uint64(math.MaxInt64) + 1 // allow MinInt64 magnitude
		if lo > limit {
			return 0, false
		}
		if lo == limit {
			return math.MinInt64, true
		}
		return -int64(lo), true
	}
	// positive result
	if lo > uint64(math.MaxInt64) {
		return 0, false
	}
	return int64(lo), true
}

// uint64ToInt64WithSign converts an unsigned magnitude to a signed int64, given desired sign.
// Returns ok=false if magnitude cannot be represented in int64 with the given sign.
func uint64ToInt64WithSign(u uint64, neg bool) (int64, bool) {
	if neg {
		limit := uint64(math.MaxInt64) + 1
		if u > limit {
			return 0, false
		}
		if u == limit {
			return math.MinInt64, true
		}
		return -int64(u), true
	}
	if u > uint64(math.MaxInt64) {
		return 0, false
	}
	return int64(u), true
}

// willOverflowInt64Add checks if adding two int64 values would overflow.
// Uses simple range checking for clarity and correctness.
func willOverflowInt64Add(a, b int64) bool {
	if b > 0 {
		return a > math.MaxInt64-b
	}
	return a < math.MinInt64-b
}

// willOverflowInt64Sub checks if subtracting two int64 values would overflow.
// Uses simple range checking for clarity and correctness.
func willOverflowInt64Sub(a, b int64) bool {
	if b > 0 {
		return a < math.MinInt64+b
	}
	return a > math.MaxInt64+b
}

// gcdInt64Uint64 calculates the greatest common divisor of int64 and uint64.
func gcdInt64Uint64(a int64, b uint64) uint64 {
	// Use absolute value for int64
	absA := absInt64ToUint64(a)
	return gcdUint64(absA, b)
}

// gcdUint64 calculates the greatest common divisor of two uint64 values using Euclid's algorithm.
func gcdUint64(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// absInt64ToUint64 converts an int64 to its absolute value as uint64.
// Handles the special case where math.MinInt64 cannot be negated within int64 range.
func absInt64ToUint64(value int64) uint64 {
	if value < 0 {
		if value == math.MinInt64 {
			// Special case: absolute value of MinInt64 doesn't fit in int64
			return uint64(math.MaxInt64) + 1
		}
		return uint64(-value)
	}
	return uint64(value)
}

// compare128Bit compares two 128-bit numbers represented as (hi, lo) pairs.
// Returns -1 if first < second, 0 if equal, 1 if first > second.
func compare128Bit(hi1, lo1, hi2, lo2 uint64) int {
	if hi1 < hi2 {
		return -1
	}
	if hi1 > hi2 {
		return 1
	}
	// High parts are equal, compare low parts
	if lo1 < lo2 {
		return -1
	}
	if lo1 > lo2 {
		return 1
	}
	return 0
}

// compareRationalsCrossMul compares two rational numbers using 128-bit cross-multiplication.
// Returns -1 if a/b < c/d, 0 if a/b == c/d, 1 if a/b > c/d.
// Uses math/bits to handle potential overflow in intermediate calculations.
func compareRationalsCrossMul(aNum int64, aDenom uint64, cNum int64, cDenom uint64) int {
	// Handle signs separately to work with unsigned arithmetic
	aSign := 1
	if aNum < 0 {
		aSign = -1
	}
	cSign := 1
	if cNum < 0 {
		cSign = -1
	}

	// Get absolute values for unsigned arithmetic
	var aAbs, cAbs uint64
	aAbs = absInt64ToUint64(aNum)

	cAbs = absInt64ToUint64(cNum)

	// Calculate a*d and c*b using 128-bit arithmetic
	aTimesDHi, aTimesDLo := bits.Mul64(aAbs, cDenom)
	cTimesBHi, cTimesBLo := bits.Mul64(cAbs, aDenom)

	// Compare the 128-bit results
	cmpResult := compare128Bit(aTimesDHi, aTimesDLo, cTimesBHi, cTimesBLo)

	// Apply sign logic - simplified
	if aSign != cSign {
		// Different signs: negative < positive
		if aSign < 0 {
			return -1
		}
		return 1
	}
	// Same signs: if both negative, reverse magnitude comparison
	if aSign < 0 {
		return -cmpResult
	}
	// Both positive: use direct magnitude comparison
	return cmpResult
}

// powerOf10 calculates 10^exp as uint64, returning (result, overflow).
// Returns overflow=true if the result would exceed uint64 capacity.
func powerOf10(exp int) (uint64, bool) {
	if exp < 0 {
		return 0, true // Invalid input
	}
	if exp == 0 {
		return 1, false
	}

	// Pre-computed powers of 10 up to 10^19 (max that fits in uint64)
	// 10^20 = 100000000000000000000 > 2^64-1 = 18446744073709551615
	powers := []uint64{
		1,                    // 10^0
		10,                   // 10^1
		100,                  // 10^2
		1000,                 // 10^3
		10000,                // 10^4
		100000,               // 10^5
		1000000,              // 10^6
		10000000,             // 10^7
		100000000,            // 10^8
		1000000000,           // 10^9
		10000000000,          // 10^10
		100000000000,         // 10^11
		1000000000000,        // 10^12
		10000000000000,       // 10^13
		100000000000000,      // 10^14
		1000000000000000,     // 10^15
		10000000000000000,    // 10^16
		100000000000000000,   // 10^17
		1000000000000000000,  // 10^18
		10000000000000000000, // 10^19
	}

	if exp >= len(powers) {
		return 0, true // Overflow
	}

	return powers[exp], false
}

// willOverflowInt64MulUint64 checks if multiplying int64 by uint64 would overflow int64 range.
func willOverflowInt64MulUint64(a int64, b uint64) bool {
	if a == 0 || b == 0 {
		return false
	}

	if a > 0 {
		// Positive case: check if a * b > MaxInt64
		return uint64(a) > uint64(math.MaxInt64)/b
	}
	// Negative case: check if a * b < MinInt64
	// Since a < 0, we need |a| * b <= |MinInt64| = 2^63
	absA := uint64(-a)
	// Special case for MinInt64: -MinInt64 would overflow, but we can handle it
	if a == math.MinInt64 {
		// MinInt64 * b should not overflow if b == 1
		return b > 1
	}
	// For other negative values, check if |a| * b > 2^63
	return absA > (uint64(math.MaxInt64)+1)/b
}

// roundDivision performs integer division with rounding according to RoundType.
// Computes round(numerator / denominator) using the specified rounding strategy.
func roundDivision(numerator int64, denominator uint64, roundType RoundType) int64 {
	if denominator == 0 {
		return 0 // Should not happen, but be safe
	}

	if numerator == 0 {
		return 0
	}

	// Get the quotient and remainder
	var quotient int64
	var remainder uint64

	if numerator >= 0 {
		quotient = numerator / int64(denominator) //nolint:gosec // safe conversion
		remainder = uint64(numerator) % denominator
	} else {
		// Handle negative numerator
		absNum := uint64(-numerator)
		quotient = -int64(absNum / denominator) //nolint:gosec // safe conversion
		remainder = absNum % denominator
	}

	// If no remainder, return exact quotient
	if remainder == 0 {
		return quotient
	}

	// Apply rounding strategy
	switch roundType {
	case RoundDown:
		// Truncate toward zero (no adjustment needed)
		return quotient

	case RoundUp:
		// Round away from zero
		if numerator > 0 {
			return quotient + 1
		}
		return quotient - 1

	case RoundHalfUp:
		// Round half values toward positive infinity
		// This means: for positive numbers, round up; for negative numbers, round toward zero

		// Check if remainder * 2 compared to denominator
		// remainder * 2 > denominator: more than half
		// remainder * 2 = denominator: exactly half
		// remainder * 2 < denominator: less than half
		doubleRemainder := remainder * 2

		if doubleRemainder > denominator {
			// More than half - round away from zero
			if numerator > 0 {
				return quotient + 1
			}
			return quotient - 1
		}

		if doubleRemainder == denominator {
			// Exactly half - round toward positive infinity
			if numerator > 0 {
				// Positive: round up (away from zero)
				return quotient + 1
			}
			// Negative: round toward positive (toward zero)
			return quotient
		}

		// Less than half - no adjustment
		return quotient

	default:
		return quotient
	}
}