package zerorat

import (
	"math"
	"math/bits"
	"reflect"

	"github.com/go-playground/validator/v10"
)

const (
	autoReduceThreshold uint64 = 1 << 62
	uint64BitSize              = 64
	uint64TopBitShift          = uint64BitSize - 1
	halfDivisor                = 2
)

// signed128 stores a signed 128-bit integer as sign plus unsigned magnitude.
// It is used only for overflow recovery paths where a temporary numerator may
// exceed int64 before reduction brings it back into the Rat range.
type signed128 struct {
	negative bool
	hi       uint64
	lo       uint64
}

// isZero reports whether the signed 128-bit magnitude is zero.
func (v signed128) isZero() bool {
	return v.hi == 0 && v.lo == 0
}

// normalize clears the sign on zero so comparisons and division stay stable.
func (v signed128) normalize() signed128 {
	if v.isZero() {
		v.negative = false
	}
	return v
}

// negateSigned128 flips the sign without creating a negative zero value.
func negateSigned128(v signed128) signed128 {
	if v.isZero() {
		return v.normalize()
	}
	v.negative = !v.negative
	return v
}

// mulInt64ByUint64ToSigned128 multiplies a signed int64 by uint64 and keeps the
// exact 128-bit magnitude for later reduction.
func mulInt64ByUint64ToSigned128(a int64, b uint64) signed128 {
	aAbs := absInt64ToUint64(a)
	hi, lo := bits.Mul64(aAbs, b)
	return signed128{negative: a < 0, hi: hi, lo: lo}.normalize()
}

// addSigned128 adds two signed 128-bit values and reports overflow outside the
// internal 128-bit recovery range.
func addSigned128(a, b signed128) (signed128, bool) {
	if a.negative == b.negative {
		lo, carry := bits.Add64(a.lo, b.lo, 0)
		hi, carry := bits.Add64(a.hi, b.hi, carry)
		if carry != 0 {
			return signed128{}, false
		}
		return signed128{negative: a.negative, hi: hi, lo: lo}.normalize(), true
	}

	cmp := compare128Bit(a.hi, a.lo, b.hi, b.lo)
	if cmp == 0 {
		return signed128{}, true
	}

	negative := a.negative
	leftHi, leftLo := a.hi, a.lo
	rightHi, rightLo := b.hi, b.lo
	if cmp < 0 {
		negative = b.negative
		leftHi, leftLo = b.hi, b.lo
		rightHi, rightLo = a.hi, a.lo
	}

	lo, borrow := bits.Sub64(leftLo, rightLo, 0)
	hi, _ := bits.Sub64(leftHi, rightHi, borrow)
	return signed128{negative: negative, hi: hi, lo: lo}.normalize(), true
}

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
		if lo > uint64(math.MaxInt64) {
			if lo == uint64(math.MaxInt64)+1 {
				return math.MinInt64, true
			}
			return 0, false
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
		if u > uint64(math.MaxInt64) {
			if u == uint64(math.MaxInt64)+1 {
				return math.MinInt64, true
			}
			return 0, false
		}

		return -int64(u), true
	}
	if u > uint64(math.MaxInt64) {
		return 0, false
	}
	return int64(u), true
}

// divInt64ByUint64Exact divides a signed numerator by a known exact unsigned
// divisor without unsafe casts around math.MinInt64.
func divInt64ByUint64Exact(value int64, divisor uint64) (int64, bool) {
	if divisor == 0 {
		return 0, false
	}
	if divisor == 1 {
		return value, true
	}

	absValue := absInt64ToUint64(value)
	if absValue%divisor != 0 {
		return 0, false
	}
	absValue /= divisor
	return uint64ToInt64WithSign(absValue, value < 0)
}

// addInt64AndUint64ToInt64 adds a signed int64 value and a positive uint64 value.
// It returns ok=false if the exact result cannot be represented as int64.
func addInt64AndUint64ToInt64(signed int64, positive uint64) (int64, bool) {
	if signed >= 0 {
		positiveInt, ok := uint64ToInt64WithSign(positive, false)
		if !ok || willOverflowInt64Add(signed, positiveInt) {
			return 0, false
		}

		return signed + positiveInt, true
	}

	negativeMagnitude := absInt64ToUint64(signed)
	if positive >= negativeMagnitude {
		return uint64ToInt64WithSign(positive-negativeMagnitude, false)
	}

	return uint64ToInt64WithSign(negativeMagnitude-positive, true)
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

// scaleMagnitude returns the absolute scale.
// It rejects the minimum int value because its absolute value cannot be represented as int.
func scaleMagnitude(scale int) (magnitude int, ok bool) {
	if scale >= 0 {
		return scale, true
	}

	// The minimum int value has no positive counterpart on two's-complement
	// platforms, so callers must treat it as an unsupported scale.
	minInt := -int(^uint(0)>>1) - 1
	if scale == minInt {
		return 0, false
	}

	return -scale, true
}

// appendBitModulo appends one high-to-low binary digit to a running remainder.
// It handles the temporary 65-bit value produced by remainder*2+bit without
// using division operations that panic for large high words.
func appendBitModulo(remainder, bit, divisor uint64) uint64 {
	// Keep the remainder below the divisor while appending a bit from the 128-bit
	// input. If the left shift overflows, the implicit high bit contributes one
	// divisor subtraction in binary long division.
	highBit := remainder >> uint64TopBitShift
	next := (remainder << 1) | bit
	if highBit != 0 || next >= divisor {
		return next - divisor
	}
	return next
}

// modUint128ByUint64 returns a 128-bit unsigned magnitude modulo a uint64 divisor.
func modUint128ByUint64(hi, lo, divisor uint64) uint64 {
	if divisor == 0 {
		// GCD helpers use zero as an invalid divisor marker. Returning zero keeps
		// that marker explicit and avoids a divide-by-zero panic.
		return 0
	}
	if divisor == 1 {
		return 0
	}

	// bits.Div64 requires divisor > high word. Binary long division has no such
	// precondition, so it is safe for any divisor and any 128-bit magnitude.
	remainder := uint64(0)
	for bitIndex := range uint64BitSize {
		bit := (hi >> (uint64TopBitShift - bitIndex)) & 1
		remainder = appendBitModulo(remainder, bit, divisor)
	}
	for bitIndex := range uint64BitSize {
		bit := (lo >> (uint64TopBitShift - bitIndex)) & 1
		remainder = appendBitModulo(remainder, bit, divisor)
	}
	return remainder
}

// gcdSigned128Uint64 calculates gcd(abs(value), divisor) without converting the
// temporary 128-bit numerator back to int64 first.
func gcdSigned128Uint64(value signed128, divisor uint64) uint64 {
	if divisor == 0 {
		// No valid denominator factor exists. The caller will fail the conversion
		// path instead of building a wrapped result.
		return 0
	}
	if value.isZero() {
		// gcd(0, d) is d, which lets zero numerators normalize to 0/1 upstream.
		return divisor
	}
	// Euclid only needs value modulo divisor, so the full 128-bit numerator never
	// has to be converted back to int64 for gcd calculation.
	return gcdUint64(modUint128ByUint64(value.hi, value.lo, divisor), divisor)
}

// divSigned128ByUint64ToInt64 divides an exact 128-bit numerator by divisor and
// returns the signed int64 result when the reduced value fits into Rat.
func divSigned128ByUint64ToInt64(value signed128, divisor uint64) (int64, bool) {
	if divisor == 0 {
		return 0, false
	}
	if value.hi >= divisor {
		// bits.Div64 would panic here, and the quotient would not fit into one word.
		return 0, false
	}

	quotient, remainder := bits.Div64(value.hi, value.lo, divisor)
	if remainder != 0 {
		// Recovery paths call this only when reduction should be exact.
		return 0, false
	}

	// The final signed magnitude must fit int64, including the MinInt64 edge when
	// the result is negative.
	return uint64ToInt64WithSign(quotient, value.negative)
}

// compareRemainderToHalf compares remainder/denominator with one half without
// doubling the remainder, which can overflow uint64.
func compareRemainderToHalf(remainder, denominator uint64) int {
	// Compare with denominator/2 directly. Multiplying the remainder by two would
	// wrap for remainders in the top half of uint64.
	half := denominator / halfDivisor
	if remainder > half {
		return 1
	}
	if denominator%halfDivisor == 0 && remainder == half {
		return 0
	}
	return -1
}

// shouldAutoReduce reports whether a valid Rat is large enough to justify gcd
// work after an arithmetic operation.
func shouldAutoReduce(numerator int64, denominator uint64) bool {
	// A high threshold keeps small hot-path operations free of gcd work while
	// leaving enough headroom before int64 and uint64 limits.
	return absInt64ToUint64(numerator) >= autoReduceThreshold || denominator >= autoReduceThreshold
}

// reduceIfLarge keeps intermediate arithmetic values compact without paying the
// gcd cost on the small-value path.
func (r *Rat) reduceIfLarge() {
	// Invalid values must stay invalid; Reduce would otherwise treat denominator
	// zero as part of a gcd calculation.
	if r.IsValid() && shouldAutoReduce(r.numerator, r.denominator) {
		r.Reduce()
	}
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

// Pre-computed powers of 10 up to 10^19 (max that fits in uint64)
// 10^20 = 100000000000000000000 > 2^64-1 = 18446744073709551615.
var powersOf10 = [...]uint64{ //nolint:gochecknoglobals // pre-computed constants
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

// powerOf10 calculates 10^exp as uint64, returning (result, overflow).
// Returns overflow=true if the result would exceed uint64 capacity.
func powerOf10(exp int) (uint64, bool) {
	if exp < 0 {
		return 0, true // Invalid input
	}
	if exp == 0 {
		return 1, false
	}

	if exp >= len(powersOf10) {
		return 0, true // Overflow
	}

	return powersOf10[exp], false
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

// divSigned128ByUint64 divides a signed 128-bit value by a uint64 divisor.
// The quotient keeps the input sign, and the remainder is always a magnitude.
func divSigned128ByUint64(value signed128, divisor uint64) (signed128, uint64, bool) {
	if divisor == 0 {
		return signed128{}, 0, false
	}
	if value.isZero() {
		return signed128{}, 0, true
	}

	quotientHi := value.hi / divisor
	highRemainder := value.hi % divisor
	quotientLo, remainder := bits.Div64(highRemainder, value.lo, divisor)

	quotient := signed128{negative: value.negative, hi: quotientHi, lo: quotientLo}.normalize()
	return quotient, remainder, true
}

// roundSigned128Division rounds a signed 128-bit numerator divided by a uint64 denominator.
func roundSigned128Division(numerator signed128, denominator uint64, roundType RoundType) (signed128, bool) {
	quotient, remainder, ok := divSigned128ByUint64(numerator, denominator)
	if !ok {
		return signed128{}, false
	}

	return applyRoundingToQuotient(quotient, remainder, denominator, numerator.negative, roundType)
}

// roundInt64ByUint128Denominator rounds an int64 numerator divided by a 128-bit
// denominator. It is used when negative-scale rounding makes the denominator exceed uint64.
func roundInt64ByUint128Denominator(
	numerator int64,
	denominatorHi, denominatorLo uint64,
	roundType RoundType,
) (int64, bool) {
	if denominatorHi == 0 {
		if denominatorLo == 0 {
			return 0, false
		}
		return roundDivision(numerator, denominatorLo, roundType), true
	}
	if numerator == 0 {
		return 0, true
	}

	negative := numerator < 0
	remainder := absInt64ToUint64(numerator)
	quotient := int64(0)

	if roundType != RoundDown && roundType != RoundUp && roundType != RoundHalfUp {
		return quotient, true
	}
	if roundType == RoundDown {
		return quotient, true
	}
	if roundType == RoundUp {
		if negative {
			return -1, true
		}
		return 1, true
	}

	halfComparison := compareUint64RemainderToHalf128(remainder, denominatorHi, denominatorLo)
	if halfComparison > 0 {
		if negative {
			return -1, true
		}
		return 1, true
	}
	if halfComparison == 0 && !negative {
		return 1, true
	}
	return quotient, true
}

// compareUint64RemainderToHalf128 compares a uint64 remainder with half of a
// 128-bit denominator.
func compareUint64RemainderToHalf128(remainder, denominatorHi, denominatorLo uint64) int {
	doubleHi, doubleLo := bits.Mul64(remainder, halfDivisor)
	return compare128Bit(doubleHi, doubleLo, denominatorHi, denominatorLo)
}

// applyRoundingToQuotient applies RoundType to a truncated quotient and a uint64
// remainder.
func applyRoundingToQuotient(
	quotient signed128,
	remainder, denominator uint64,
	negative bool,
	roundType RoundType,
) (signed128, bool) {
	if remainder == 0 {
		return quotient, true
	}

	if roundType != RoundDown && roundType != RoundUp && roundType != RoundHalfUp {
		return quotient, true
	}

	if roundType == RoundDown {
		return quotient, true
	}

	incrementAwayFromZero := func() (signed128, bool) {
		one := signed128{negative: negative, lo: 1}
		return addSigned128(quotient, one)
	}

	if roundType == RoundUp {
		return incrementAwayFromZero()
	}

	halfComparison := compareRemainderToHalf(remainder, denominator)
	if halfComparison > 0 {
		return incrementAwayFromZero()
	}
	if halfComparison == 0 && !negative {
		return incrementAwayFromZero()
	}
	return quotient, true
}

// roundDivision performs integer division with rounding according to RoundType.
// It computes round(numerator / denominator) using the specified rounding strategy.
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
		quotient = int64(uint64(numerator) / denominator) //nolint:gosec // quotient cannot exceed numerator
		remainder = uint64(numerator) % denominator
	} else {
		// Handle negative numerator
		absNum := absInt64ToUint64(numerator)
		var ok bool
		quotient, ok = uint64ToInt64WithSign(absNum/denominator, true)
		if !ok {
			return 0
		}
		remainder = absNum % denominator
	}

	// If no remainder, return exact quotient
	if remainder == 0 {
		return quotient
	}

	if roundType != RoundDown && roundType != RoundUp && roundType != RoundHalfUp {
		return quotient
	}

	// Apply rounding strategy
	switch roundType {
	case RoundDown:
	case RoundUp:
		// Round away from zero
		if numerator > 0 {
			return quotient + 1
		}
		return quotient - 1

	case RoundHalfUp:
		// Round half values toward positive infinity
		// This means: for positive numbers, round up; for negative numbers, round toward zero

		halfComparison := compareRemainderToHalf(remainder, denominator)
		if halfComparison > 0 {
			// More than half - round away from zero
			if numerator > 0 {
				return quotient + 1
			}
			return quotient - 1
		}

		if halfComparison == 0 {
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

	return quotient
}

// RegisterValidationFunc registers a custom validation function for Rat types.
func RegisterValidationFunc(v *validator.Validate) {
	v.RegisterCustomTypeFunc(func(field reflect.Value) any {
		if r, ok := field.Interface().(Rat); ok {
			return r.IsValid()
		}
		return nil
	}, Rat{})
}
