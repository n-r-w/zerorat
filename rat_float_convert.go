package zerorat

import (
	"math"
	"math/bits"
)

// float64ToRatExact converts a float64 to a rational using IEEE-754 decomposition.
// It returns the exact rational m * 2^e when it fits into (int64 numerator, uint64 denominator).
// If exact representation would overflow the type bounds, it returns an invalid Rat (denominator=0).
//
//nolint:gocognit,nestif,mnd // complexity and magic constants are acceptable here due to bit-level IEEE-754 handling
func float64ToRatExact(value float64) Rat {
	// Reject non-finite
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return Rat{}
	}
	// Zero (incl. -0) handled here for completeness
	if value == 0 {
		return Rat{numerator: 0, denominator: 1}
	}

	bits64 := math.Float64bits(value)
	neg := (bits64 >> 63) != 0
	// exponent bits are 11-bit value in [0, 2047]
	expBits := int((bits64 >> 52) & 0x7FF) //nolint:gosec // bounded to 11 bits, safe int conversion
	frac := bits64 & ((uint64(1) << 52) - 1)

	var mant uint64
	var e int // power of two exponent in value = mant * 2^e
	if expBits == 0 {
		// Subnormal: no implicit 1, unbiased exp is 1-1023, but value = frac * 2^-1074
		// frac is non-zero here because value != 0
		mant = frac
		e = -1074
	} else {
		// Normalized: implicit leading 1
		mant = (uint64(1) << 52) | frac
		e = expBits - 1023 - 52
	}

	// Reduce common factors of 2 when e < 0 to keep denominator smaller
	if e < 0 {
		// Remove up to -e trailing zeros from mant
		if tz := bits.TrailingZeros64(mant); tz > 0 {
			if tz > -e {
				mant >>= uint(-e)
				e = 0
			} else {
				mant >>= uint(tz)
				e += tz
			}
		}
	}

	// Try to construct exact value; if too large, move some exponent into denominator (power-of-two)
	if e >= 0 {
		absLimit := uint64(math.MaxInt64)
		limitBits := 63
		if neg {
			absLimit = uint64(math.MaxInt64) + 1 // allow -2^63
			limitBits = 64
		}

		mantBits := bits.Len64(mant)
		// Fast exact path: mant << e must fit absLimit
		if e >= 0 && e < 64 {
			shifted := mant << uint(e)
			if mant <= (absLimit >> uint(e)) {
				// shifted is <= absLimit by the guard above, safe to cast
				n := int64(shifted) //nolint:gosec // guarded by absLimit check above
				if neg {
					n = -n
				}
				return Rat{numerator: n, denominator: 1}
			}
		}

		// Need to offload some exponent to denominator 2^d so numerator fits
		maxShiftAllowed := max(limitBits-mantBits, 0)
		neededDenPow := max(e-maxShiftAllowed, 0)
		if neededDenPow > 63 {
			// Even with denom capped at 2^63, numerator still too large → invalid (overflow)
			return Rat{}
		}
		newShift := e - neededDenPow
		// Now mant << newShift must fit
		if newShift < 0 || newShift >= 64 || mant > (absLimit>>uint(newShift)) {
			return Rat{}
		}
		shifted := mant << uint(newShift)
		// shifted is <= absLimit by the guard above, safe to cast
		n := int64(shifted) //nolint:gosec // guarded by absLimit check above
		if neg {
			n = -n
		}
		// neededDenPow in [0,63] by guard above; safe shift
		den := uint64(1) << uint(neededDenPow) //nolint:gosec // neededDenPow in [0,63] by guard above, safe shift
		return Rat{numerator: n, denominator: den}
	}

	// e < 0: choose denominator d = min(-e, 63) and compute nearest numerator by rounding
	denPow := min(-e, 63)
	shiftUp := e + denPow // <= 0
	var n64 uint64
	if shiftUp < 0 {
		shift := uint(-shiftUp)
		if shift >= 64 {
			n64 = 0
		} else {
			base := mant >> shift
			// Round-to-nearest-even on the truncated bits
			mask := (uint64(1) << shift) - 1
			rem := mant & mask
			half := uint64(1) << (shift - 1)
			if rem > half || (rem == half && (base&1) == 1) {
				base++
			}
			n64 = base
		}
	} else {
		// shiftUp == 0 → exact
		n64 = mant
	}
	// n64 now fits well within int64 range (<= mant)
	// n64 is derived from mant which fits in int64 after rounding; safe cast
	n := int64(n64) //nolint:gosec // derived from mant and bounded, safe cast
	if neg {
		n = -n
	}
	// denPow is in [0,63]; safe shift
	den := uint64(1) << uint(denPow) //nolint:gosec // denPow is bounded in [0,63], safe shift
	return Rat{numerator: n, denominator: den}
}
