package zerorat

import (
	"math"
	"math/bits"
)

// float64Parts describes a finite float64 as mantissa * 2^exponent.
type float64Parts struct {
	negative bool
	mantissa uint64
	exponent int
}

const (
	float64SignBit          = 63
	float64FractionBits     = 52
	float64ExponentMask     = 0x7FF
	float64ExponentBias     = 1023
	float64SubnormalExp     = -1074
	float64DenominatorLimit = 63
	uint64BitWidth          = 64
)

// ratFromFloat64PartsExact converts already decomposed finite float parts into an exact Rat.
func ratFromFloat64PartsExact(parts float64Parts) (Rat, error) {
	if parts.exponent >= 0 {
		if parts.exponent >= uint64BitWidth {
			return Rat{}, ErrNotRepresentable
		}

		absLimit := uint64(math.MaxInt64)
		if parts.negative {
			absLimit++
		}
		if parts.mantissa > (absLimit >> uint(parts.exponent)) {
			return Rat{}, ErrNotRepresentable
		}

		shifted := parts.mantissa << uint(parts.exponent)
		// The absLimit guard above guarantees that this conversion fits.
		numerator, _ := uint64ToInt64WithSign(shifted, parts.negative)
		return Rat{numerator: numerator, denominator: 1}, nil
	}

	denPow := -parts.exponent
	if denPow > float64DenominatorLimit {
		return Rat{}, ErrNotRepresentable
	}

	// Finite float64 mantissa fits safely into int64 after decomposition.
	numerator, _ := uint64ToInt64WithSign(parts.mantissa, parts.negative)
	return Rat{numerator: numerator, denominator: uint64(1) << uint(denPow)}, nil
}

// decomposeFiniteFloat64 normalizes a finite float64 into mantissa * 2^exponent.
func decomposeFiniteFloat64(value float64) (float64Parts, bool, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return float64Parts{}, false, ErrNonFiniteFloat
	}
	if value == 0 {
		return float64Parts{}, true, nil
	}

	bits64 := math.Float64bits(value)
	parts := float64Parts{negative: (bits64 >> float64SignBit) != 0}
	expBits := int((bits64 >> float64FractionBits) & float64ExponentMask)
	frac := bits64 & ((uint64(1) << float64FractionBits) - 1)

	if expBits == 0 {
		parts.mantissa = frac
		parts.exponent = float64SubnormalExp
	} else {
		parts.mantissa = (uint64(1) << float64FractionBits) | frac
		parts.exponent = expBits - float64ExponentBias - float64FractionBits
	}

	if parts.exponent < 0 {
		if tz := bits.TrailingZeros64(parts.mantissa); tz > 0 {
			shift := min(tz, -parts.exponent)
			parts.mantissa >>= uint(shift)
			parts.exponent += shift
		}
	}

	return parts, false, nil
}

// float64ToRatExact converts a float64 to an exact rational using IEEE-754 decomposition.
// It returns ErrNonFiniteFloat for NaN and infinities.
// It returns ErrNotRepresentable when the exact value does not fit into Rat.
func float64ToRatExact(value float64) (Rat, error) {
	parts, isZero, err := decomposeFiniteFloat64(value)
	if err != nil {
		return Rat{}, err
	}
	if isZero {
		return Rat{numerator: 0, denominator: 1}, nil
	}

	return ratFromFloat64PartsExact(parts)
}

// float64ToRatApprox converts a float64 to the nearest representable Rat.
// It returns ErrNonFiniteFloat for NaN and infinities.
// It returns ErrNotRepresentable when the integer part does not fit into Rat.
func float64ToRatApprox(value float64) (Rat, error) {
	parts, isZero, err := decomposeFiniteFloat64(value)
	if err != nil {
		return Rat{}, err
	}
	if isZero {
		return Rat{numerator: 0, denominator: 1}, nil
	}

	if parts.exponent >= 0 {
		return ratFromFloat64PartsExact(parts)
	}

	denPow := -parts.exponent
	if denPow <= float64DenominatorLimit {
		return ratFromFloat64PartsExact(parts)
	}

	shift := uint(denPow - float64DenominatorLimit)
	var rounded uint64
	if shift < uint64BitWidth {
		base := parts.mantissa >> shift
		mask := (uint64(1) << shift) - 1
		rem := parts.mantissa & mask
		half := uint64(1) << (shift - 1)
		if rem > half || (rem == half && (base&1) == 1) {
			base++
		}
		rounded = base
	}

	// rounded is derived from the finite float64 mantissa and therefore fits into int64.
	numerator, _ := uint64ToInt64WithSign(rounded, parts.negative)
	if numerator == 0 {
		return Rat{numerator: 0, denominator: 1}, nil
	}
	r := Rat{numerator: numerator, denominator: uint64(1) << float64DenominatorLimit}
	r.Reduce()
	return r, nil
}
