package zerorat

// Reduce reduces the rational number to its lowest terms (mutable operation).
// Uses the Euclidean algorithm to find the GCD and divides both numerator and denominator by it.
// If the rational number is invalid, it remains invalid.
func (r *Rat) Reduce() {
	if r.IsInvalid() {
		return
	}

	// Zero is already in reduced form (0/1)
	if r.numerator == 0 {
		r.denominator = 1
		return
	}

	// Find GCD and reduce
	gcd := gcdInt64Uint64(r.numerator, r.denominator)
	if gcd > 1 {
		// Reduce numerator using unsigned magnitude to avoid unsafe casts
		absNum := absInt64ToUint64(r.numerator)
		absNum /= gcd
		newNum, ok := uint64ToInt64WithSign(absNum, r.numerator < 0)
		if !ok {
			// Should not happen, but be safe: mark invalid
			r.Invalidate()
			return
		}
		r.numerator = newNum
		r.denominator /= gcd
	}
}

// Reduced returns a new rational number reduced to its lowest terms (immutable operation).
// Does not modify the original rational number.
func (r Rat) Reduced() Rat {
	result := r // create copy
	result.Reduce()
	return result
}

// Round rounds the rational number to the specified scale (mutable operation).
// Uses the given rounding strategy (RoundType) to determine the rounding behavior.
// If the rational number is invalid, it remains invalid.
//
// Scale interpretation:
// - scale = 0: round to integer (1.23 -> 1)
// - scale > 0: round to decimal places (1.234 with scale=2 -> 1.23)
// - scale < 0: round to powers of 10 (1234 with scale=-2 -> 1200).
func (r *Rat) Round(roundType RoundType, scale int) {
	// If invalid, remain invalid
	if r.IsInvalid() {
		return
	}

	// If zero, no rounding needed
	if r.numerator == 0 {
		return
	}

	// Calculate the scaling factor (10^|scale|)
	var scaleFactor uint64
	var scaleFactorOverflow bool

	if scale >= 0 {
		scaleFactor, scaleFactorOverflow = powerOf10(scale)
	} else {
		scaleFactor, scaleFactorOverflow = powerOf10(-scale)
	}

	// Handle overflow in scale factor calculation
	if scaleFactorOverflow {
		if scale >= 0 {
			// Very large positive scale - result would be extremely precise
			// For practical purposes, mark as invalid
			r.Invalidate()
			return
		}
		// Very large negative scale - round to zero
		r.numerator = 0
		r.denominator = 1
		return
	}

	if scale == 0 {
		// Round to integer: compute round(numerator/denominator)
		roundedInt := roundDivision(r.numerator, r.denominator, roundType)
		r.numerator = roundedInt
		r.denominator = 1
		return
	}

	if scale > 0 {
		r.roundToDecimalPlaces(scaleFactor, roundType)
	} else {
		r.roundToPowersOfTen(scaleFactor, roundType)
	}

	// Handle zero result
	if r.numerator == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// Reduce to lowest terms
	r.Reduce()
}

// roundToDecimalPlaces handles rounding to a positive number of decimal places.
func (r *Rat) roundToDecimalPlaces(scaleFactor uint64, roundType RoundType) {
	// Round to decimal places
	// To round a/b to scale decimal places:
	// 1. Multiply by 10^scale: (a * 10^scale) / b
	// 2. Round to integer: round((a * 10^scale) / b)
	// 3. Result is rounded_value / 10^scale

	// First check if the number is already exact at the requested scale
	// This happens when the denominator divides 10^scale evenly
	if scaleFactor%r.denominator == 0 {
		// Already exact - convert to standard scale format
		// Convert to the requested scale: a/b = (a * (10^scale / b)) / 10^scale
		multiplier := scaleFactor / r.denominator

		// Check for overflow in the multiplication
		if willOverflowInt64MulUint64(r.numerator, multiplier) {
			// If we can't represent at the requested scale due to overflow,
			// but the value is already exact, just leave it as-is
			// This handles cases like MaxInt64 with scale 1
			return
		}

		var newNumerator int64
		if r.numerator >= 0 {
			newNumerator = r.numerator * int64(multiplier) //nolint:gosec // overflow checked above
		} else {
			absNum := uint64(-r.numerator)
			newNumerator = -int64(absNum * multiplier) //nolint:gosec // overflow checked above
		}

		r.numerator = newNumerator
		r.denominator = scaleFactor
		return
	}

	// Check for overflow in numerator multiplication
	if willOverflowInt64MulUint64(r.numerator, scaleFactor) {
		r.Invalidate()
		return
	}

	// Multiply numerator by scale factor
	var scaledNumerator int64
	if r.numerator >= 0 {
		scaledNumerator = r.numerator * int64(scaleFactor) //nolint:gosec // overflow checked above
	} else {
		// Handle negative case carefully to avoid overflow
		absNum := uint64(-r.numerator)
		scaledNumerator = -int64(absNum * scaleFactor) //nolint:gosec // overflow checked above
	}

	// Round the scaled value to integer
	roundedInt := roundDivision(scaledNumerator, r.denominator, roundType)

	// Set result as roundedInt / scaleFactor
	r.numerator = roundedInt
	r.denominator = scaleFactor
}

// roundToPowersOfTen handles rounding to powers of 10 (negative scale).
func (r *Rat) roundToPowersOfTen(scaleFactor uint64, roundType RoundType) {
	// scale < 0: Round to powers of 10
	// To round a/b to nearest multiple of 10^(-scale):
	// 1. Compute a/b
	// 2. Divide by 10^(-scale): (a/b) / 10^(-scale) = a / (b * 10^(-scale))
	// 3. Round to integer: round(a / (b * 10^(-scale)))
	// 4. Multiply back: rounded_value * 10^(-scale)

	// Check for overflow in denominator multiplication
	if willOverflowUint64Mul(r.denominator, scaleFactor) {
		r.Invalidate()
		return
	}

	// Scale the denominator
	scaledDenominator := r.denominator * scaleFactor

	// Round to integer
	roundedInt := roundDivision(r.numerator, scaledDenominator, roundType)

	// Multiply back by scale factor
	if willOverflowInt64MulUint64(roundedInt, scaleFactor) {
		r.Invalidate()
		return
	}

	var finalNumerator int64
	if roundedInt >= 0 {
		finalNumerator = roundedInt * int64(scaleFactor) //nolint:gosec // overflow checked above
	} else {
		// Handle negative case carefully
		absRounded := uint64(-roundedInt)
		finalNumerator = -int64(absRounded * scaleFactor) //nolint:gosec // overflow checked above
	}

	r.numerator = finalNumerator
	r.denominator = 1
}

// Rounded returns a new rational number rounded to the nearest integer (immutable operation).
func (r Rat) Rounded(roundType RoundType, scale int) Rat {
	result := r // create copy
	result.Round(roundType, scale)
	return result
}
