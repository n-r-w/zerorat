package zerorat

import "math/bits"

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

	r.reduceIfLarge()

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
			reduced := *r
			reduced.Reduce()
			if denominatorDividesPowerOf10(reduced.denominator, scale) {
				r.numerator = reduced.numerator
				r.denominator = reduced.denominator
				return
			}

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

	// Keep rounded small values on the fast path and compact only large results.
	r.reduceIfLarge()
}

// reduceForRound reduces the receiver only after a rounding overflow is detected.
// It reports whether the representation changed enough to justify a retry.
func (r *Rat) reduceForRound() bool {
	oldNumerator := r.numerator
	oldDenominator := r.denominator
	r.Reduce()
	return r.numerator != oldNumerator || r.denominator != oldDenominator
}

// denominatorDividesPowerOf10 reports whether denominator needs no more than scale decimal places.
func denominatorDividesPowerOf10(denominator uint64, scale int) bool {
	if denominator == 0 || scale < 0 {
		return false
	}

	twoCount := 0
	for denominator%2 == 0 {
		denominator /= 2
		twoCount++
	}

	fiveCount := 0
	for denominator%5 == 0 {
		denominator /= 5
		fiveCount++
	}

	requiredScale := twoCount
	if fiveCount > requiredScale {
		requiredScale = fiveCount
	}

	return denominator == 1 && requiredScale <= scale
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
		scaledNumerator := mulInt64ByUint64ToSigned128(r.numerator, scaleFactor)
		roundedScaled, ok := roundSigned128Division(scaledNumerator, r.denominator, roundType)
		if !ok {
			r.Invalidate()
			return
		}

		roundedInt, ok := divSigned128ByUint64ToInt64(roundedScaled, 1)
		if ok {
			r.numerator = roundedInt
			r.denominator = scaleFactor
			return
		}

		commonDivisor := gcdSigned128Uint64(roundedScaled, scaleFactor)
		reducedNumerator, ok := divSigned128ByUint64ToInt64(roundedScaled, commonDivisor)
		if !ok {
			r.Invalidate()
			return
		}

		r.numerator = reducedNumerator
		r.denominator = scaleFactor / commonDivisor
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
		scaledDenominatorHi, scaledDenominatorLo := bits.Mul64(r.denominator, scaleFactor)
		roundedInt, ok := roundInt64ByUint128Denominator(
			r.numerator,
			scaledDenominatorHi,
			scaledDenominatorLo,
			roundType,
		)
		if !ok {
			r.Invalidate()
			return
		}
		newNum, ok := mulInt64ByUint64ToInt64(roundedInt, scaleFactor)
		if !ok {
			r.Invalidate()
			return
		}

		r.numerator = newNum
		r.denominator = 1
		return
	}

	// Scale the denominator
	scaledDenominator := r.denominator * scaleFactor

	// Round to integer
	roundedInt := roundDivision(r.numerator, scaledDenominator, roundType)

	// Multiply back by scale factor
	if willOverflowInt64MulUint64(roundedInt, scaleFactor) {
		if r.reduceForRound() {
			r.roundToPowersOfTen(scaleFactor, roundType)
			return
		}
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

// Round64Down rounds the rational number to an int64 using RoundDown strategy with scale 0.
func (r Rat) Round64Down() int64 {
	return r.Rounded(RoundDown, 0).ToInt64()
}

// Round64Up rounds the rational number to an int64 using RoundUp strategy with scale 0.
func (r Rat) Round64Up() int64 {
	return r.Rounded(RoundUp, 0).ToInt64()
}

// Round64HalfUp rounds the rational number to an int64 using RoundHalfUp strategy with scale 0.
func (r Rat) Round64HalfUp() int64 {
	return r.Rounded(RoundHalfUp, 0).ToInt64()
}

// RoundDown rounds the rational number to an int using RoundDown strategy with scale 0.
func (r Rat) RoundDown() int {
	return r.Rounded(RoundDown, 0).ToInt()
}

// RoundUp rounds the rational number to an int using RoundUp strategy with scale 0.
func (r Rat) RoundUp() int {
	return r.Rounded(RoundUp, 0).ToInt()
}

// RoundHalfUp rounds the rational number to an int using RoundHalfUp strategy with scale 0.
func (r Rat) RoundHalfUp() int {
	return r.Rounded(RoundHalfUp, 0).ToInt()
}

// RoundFloatDown rounds the rational number to a float64 using RoundDown strategy with scale 0.
func (r Rat) RoundFloatDown() float64 {
	return r.Rounded(RoundDown, 0).ToFloat64()
}

// RoundFloatUp rounds the rational number to a float64 using RoundUp strategy with scale 0.
func (r Rat) RoundFloatUp() float64 {
	return r.Rounded(RoundUp, 0).ToFloat64()
}

// RoundFloatHalfUp rounds the rational number to a float64 using RoundHalfUp strategy with scale 0.
func (r Rat) RoundFloatHalfUp() float64 {
	return r.Rounded(RoundHalfUp, 0).ToFloat64()
}

// RoundFloat32Down rounds the rational number to a float32 using RoundDown strategy with scale 0.
func (r Rat) RoundFloat32Down() float32 {
	return r.Rounded(RoundDown, 0).ToFloat32()
}

// RoundFloat32Up rounds the rational number to a float32 using RoundUp strategy with scale 0.
func (r Rat) RoundFloat32Up() float32 {
	return r.Rounded(RoundUp, 0).ToFloat32()
}

// RoundFloat32HalfUp rounds the rational number to a float32 using RoundHalfUp strategy with scale 0.
func (r Rat) RoundFloat32HalfUp() float32 {
	return r.Rounded(RoundHalfUp, 0).ToFloat32()
}
