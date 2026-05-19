package zerorat

import (
	"math"
)

// addSubCommon implements common logic for addition and subtraction.
// isAdd=true for addition, isAdd=false for subtraction.
func (r *Rat) addSubCommon(other Rat, isAdd bool) {
	// If any operand is invalid, the result is invalid
	if r.IsInvalid() || other.IsInvalid() {
		r.Invalidate()
		return
	}

	// Optimization for same denominators
	if r.denominator == other.denominator {
		newNum, ok := addSubInt64(r.numerator, other.numerator, isAdd)
		if !ok {
			r.addSubReduced(other, isAdd)
			return
		}

		// If the result is zero, normalize to 0/1
		if newNum == 0 {
			r.numerator = 0
			r.denominator = 1
			return
		}

		r.numerator = newNum
		// denominator remains the same
		r.reduceIfLarge()
		return
	}

	// General case: different denominators
	// Check for denominator multiplication overflow
	if willOverflowUint64Mul(r.denominator, other.denominator) {
		r.addSubReduced(other, isAdd)
		return
	}

	// Calculate new denominator
	newDenom := r.denominator * other.denominator

	// Check for numerator calculation overflow
	// Check a*d overflow and compute safely
	prod1, ok := mulInt64ByUint64ToInt64(r.numerator, other.denominator)
	if !ok {
		r.addSubReduced(other, isAdd)
		return
	}

	// Check c*b overflow and compute safely
	prod2, ok := mulInt64ByUint64ToInt64(other.numerator, r.denominator)
	if !ok {
		r.addSubReduced(other, isAdd)
		return
	}

	term1 := prod1
	term2 := prod2

	newNum, ok := addSubInt64(term1, term2, isAdd)
	if !ok {
		r.addSubReduced(other, isAdd)
		return
	}

	// If the result is zero, normalize to 0/1
	if newNum == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// Store the raw result first so small values avoid gcd work.
	r.numerator = newNum
	r.denominator = newDenom
	r.reduceIfLarge()
}

// addSubInt64 performs the selected signed operation and reports whether the
// exact result fits into int64.
func addSubInt64(left, right int64, isAdd bool) (int64, bool) {
	if isAdd {
		if willOverflowInt64Add(left, right) {
			return 0, false
		}
		return left + right, true
	}

	if willOverflowInt64Sub(left, right) {
		return 0, false
	}
	return left - right, true
}

// addSubReduced retries addition or subtraction with reduced operands and a
// 128-bit temporary numerator so representable results are not invalidated only
// because the unreduced intermediate terms overflowed int64.
func (r *Rat) addSubReduced(other Rat, isAdd bool) {
	// Work on reduced copies so the receiver is changed only after the retry is
	// known to fit into the Rat representation.
	left := *r
	right := other
	left.Reduce()
	right.Reduce()

	// Use the common denominator factor before cross multiplication. This keeps
	// the temporary numerator smaller and leaves a denominator factor available
	// for the final reduction step.
	denominatorGCD := gcdUint64(left.denominator, right.denominator)
	leftScale := right.denominator / denominatorGCD
	rightScale := left.denominator / denominatorGCD

	// The numerator terms can exceed int64 before cancellation. Keep them as a
	// signed 128-bit magnitude until the shared denominator factor is removed.
	leftTerm := mulInt64ByUint64ToSigned128(left.numerator, leftScale)
	rightTerm := mulInt64ByUint64ToSigned128(right.numerator, rightScale)
	if !isAdd {
		// Subtraction is addition with the right term sign flipped.
		rightTerm = negateSigned128(rightTerm)
	}

	// If the 128-bit numerator itself overflows, no valid Rat result can be built
	// from this retry path.
	numerator128, ok := addSigned128(leftTerm, rightTerm)
	if !ok {
		r.Invalidate()
		return
	}

	// A zero result must always use 0/1 so later operations do not carry a large
	// denominator that has no numeric meaning.
	if numerator128.isZero() {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// Remove the part of the common denominator factor that also divides the
	// temporary numerator before converting the numerator back to int64.
	commonFactor := gcdSigned128Uint64(numerator128, denominatorGCD)
	newNum, ok := divSigned128ByUint64ToInt64(numerator128, commonFactor)
	if !ok {
		r.Invalidate()
		return
	}

	// Build the reduced denominator as (left.denominator / gcd) *
	// (right.denominator / commonFactor). This is equivalent to the full common
	// denominator after the numerator cancellation above.
	leftDenominator := left.denominator / denominatorGCD
	rightDenominator := right.denominator / commonFactor
	if willOverflowUint64Mul(leftDenominator, rightDenominator) {
		r.Invalidate()
		return
	}

	r.numerator = newNum
	r.denominator = leftDenominator * rightDenominator
}

// Add adds another rational number to the current one (mutable operation).
// Formula: a/b + c/d = (a*d + c*b) / (b*d)
// Large intermediate results may be reduced when the exact result fits Rat.
// Sets invalid state on overflow or with invalid operands.
func (r *Rat) Add(other Rat) {
	r.addSubCommon(other, true)
}

// Sub subtracts another rational number from the current one (mutable operation).
// Formula: a/b - c/d = (a*d - c*b) / (b*d)
// Large intermediate results may be reduced when the exact result fits Rat.
// Sets invalid state on overflow or with invalid operands.
func (r *Rat) Sub(other Rat) {
	r.addSubCommon(other, false)
}

// Mul multiplies the current rational number by another (mutable operation).
// Formula: a/b * c/d = (a*c) / (b*d)
// Large intermediate results may be reduced when the exact result fits Rat.
// Sets invalid state on overflow or with invalid operands.
func (r *Rat) Mul(other Rat) {
	// If any operand is invalid, result is invalid
	if r.IsInvalid() || other.IsInvalid() {
		r.Invalidate()
		return
	}

	// Check numerator multiplication overflow
	if willOverflowInt64Mul(r.numerator, other.numerator) {
		r.mulReduced(other)
		return
	}

	// Check denominator multiplication overflow
	if willOverflowUint64Mul(r.denominator, other.denominator) {
		r.mulReduced(other)
		return
	}

	newNum := r.numerator * other.numerator
	newDenom := r.denominator * other.denominator

	// If result is zero, normalize to 0/1
	if newNum == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// Store the raw result first so small values avoid gcd work.
	r.numerator = newNum
	r.denominator = newDenom
	r.reduceIfLarge()
}

// mulReduced retries multiplication after reducing operands and cancelling
// cross factors, which prevents overflow when the exact product still fits Rat.
func (r *Rat) mulReduced(other Rat) {
	// Reduce operands first so same-side common factors do not hide cheaper cross
	// cancellation opportunities.
	left := *r
	right := other
	left.Reduce()
	right.Reduce()

	// Multiplication by zero has a fixed compact representation and does not need
	// any denominator work.
	if left.numerator == 0 || right.numerator == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// Cancel numerator factors against the opposite denominator before multiplying.
	// This prevents overflow for products such as (2/MaxUint64) * (MaxInt64/2).
	leftWithRightDenominator := gcdInt64Uint64(left.numerator, right.denominator)
	rightWithLeftDenominator := gcdInt64Uint64(right.numerator, left.denominator)

	// The gcd values are exact divisors by construction. A failed division means
	// an internal invariant was broken, so the public result becomes invalid.
	leftNum, ok := divInt64ByUint64Exact(left.numerator, leftWithRightDenominator)
	if !ok {
		r.Invalidate()
		return
	}
	rightNum, ok := divInt64ByUint64Exact(right.numerator, rightWithLeftDenominator)
	if !ok {
		r.Invalidate()
		return
	}

	leftDenominator := left.denominator / rightWithLeftDenominator
	rightDenominator := right.denominator / leftWithRightDenominator
	if willOverflowInt64Mul(leftNum, rightNum) || willOverflowUint64Mul(leftDenominator, rightDenominator) {
		// Cross cancellation was not enough to keep the exact result in Rat limits.
		r.Invalidate()
		return
	}

	r.numerator = leftNum * rightNum
	r.denominator = leftDenominator * rightDenominator
}

// Div divides the current rational number by another (mutable operation).
// Formula: a/b ÷ c/d = a/b * d/c = (a*d) / (b*c)
// Large intermediate results may be reduced when the exact result fits Rat.
// Sets invalid state on overflow, division by zero, or with invalid operands.
func (r *Rat) Div(other Rat) {
	// If any operand is invalid, result is invalid
	if r.IsInvalid() || other.IsInvalid() {
		r.Invalidate()
		return
	}

	// Check for division by zero
	if other.numerator == 0 {
		r.Invalidate()
		return
	}

	// Division is equivalent to multiplying by reciprocal
	// a/b ÷ c/d = a/b * d/c = (a*d) / (b*c)

	// Get absolute value of other.numerator for unsigned arithmetic
	otherNumAbs := absInt64ToUint64(other.numerator)

	// Apply the divisor sign before multiplying so the MinInt64 magnitude remains
	// representable when the exact result is negative.
	signedNumerator := r.numerator
	if other.numerator < 0 {
		if signedNumerator == math.MinInt64 {
			r.divReduced(other)
			return
		}
		signedNumerator = -signedNumerator
	}

	// Check for numerator * denominator overflow and compute safely
	prodNum, ok := mulInt64ByUint64ToInt64(signedNumerator, other.denominator)
	if !ok {
		r.divReduced(other)
		return
	}

	// Check for denominator * numerator overflow
	if willOverflowUint64Mul(r.denominator, otherNumAbs) {
		r.divReduced(other)
		return
	}

	newNum := prodNum
	newDenom := r.denominator * otherNumAbs

	// If result is zero, normalize to 0/1
	if newNum == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// Store the raw result first so small values avoid gcd work.
	r.numerator = newNum
	r.denominator = newDenom
	r.reduceIfLarge()
}

// divReduced retries division after reducing operands and cancelling cross
// factors from a/b * d/c, including the sign carried by c.
func (r *Rat) divReduced(other Rat) {
	// Division is multiplication by the reciprocal. Reduce both operands before
	// cancellation so the reciprocal product uses the smallest available factors.
	left := *r
	right := other
	left.Reduce()
	right.Reduce()

	// Zero divided by any valid non-zero value is always the compact zero value.
	if left.numerator == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// For a/b divided by c/d, cancel a with |c| and d with b before multiplying
	// the remaining terms.
	rightNumeratorAbs := absInt64ToUint64(right.numerator)
	leftWithRightNumerator := gcdInt64Uint64(left.numerator, rightNumeratorAbs)
	rightDenominatorWithLeftDenominator := gcdUint64(right.denominator, left.denominator)

	// The first cancellation removes the part of the divisor numerator that would
	// otherwise enlarge the result denominator.
	leftNum, ok := divInt64ByUint64Exact(left.numerator, leftWithRightNumerator)
	if !ok {
		r.Invalidate()
		return
	}
	rightNumeratorAbs /= leftWithRightNumerator
	leftDenominator := left.denominator / rightDenominatorWithLeftDenominator
	rightDenominator := right.denominator / rightDenominatorWithLeftDenominator

	// Apply the divisor sign before multiplying so a negative MinInt64 result is
	// kept when the exact value fits the Rat numerator range.
	if right.numerator < 0 {
		if leftNum == math.MinInt64 {
			r.Invalidate()
			return
		}
		leftNum = -leftNum
	}

	newNum, ok := mulInt64ByUint64ToInt64(leftNum, rightDenominator)
	if !ok {
		r.Invalidate()
		return
	}

	if willOverflowUint64Mul(leftDenominator, rightNumeratorAbs) {
		// The remaining denominator terms still exceed uint64 after cancellation.
		r.Invalidate()
		return
	}

	r.numerator = newNum
	r.denominator = leftDenominator * rightNumeratorAbs
}

// Divided returns the quotient of current divided by another rational number (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) Divided(other Rat) Rat {
	result := r // make a copy
	result.Div(other)
	return result
}

// Added returns the sum of current and another rational number (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) Added(other Rat) Rat {
	result := r // make a copy
	result.Add(other)
	return result
}

// Subtracted returns the difference of current and another rational number (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) Subtracted(other Rat) Rat {
	result := r // make a copy
	result.Sub(other)
	return result
}

// Multiplied returns the product of current and another rational number (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) Multiplied(other Rat) Rat {
	result := r // make a copy
	result.Mul(other)
	return result
}

// AddInt adds an int64 value to the current rational number (mutable operation).
func (r *Rat) AddInt(value int64) {
	r.Add(NewFromInt64(value))
}

// AddedInt returns the sum of current and an int64 value (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) AddedInt(value int64) Rat {
	result := r // make a copy
	result.AddInt(value)
	return result
}

// SubInt subtracts an int64 value from the current rational number (mutable operation).
func (r *Rat) SubInt(value int64) {
	r.Sub(NewFromInt64(value))
}

// SubtractedInt returns the difference of current and an int64 value (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) SubtractedInt(value int64) Rat {
	result := r // make a copy
	result.SubInt(value)
	return result
}

// MulInt multiplies the current rational number by an int64 value (mutable operation).
func (r *Rat) MulInt(value int64) {
	r.Mul(NewFromInt64(value))
}

// MultipliedInt returns the product of current and an int64 value (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) MultipliedInt(value int64) Rat {
	result := r // make a copy
	result.MulInt(value)
	return result
}

// DivInt divides the current rational number by an int64 value (mutable operation).
func (r *Rat) DivInt(value int64) {
	r.Div(NewFromInt64(value))
}

// DividedInt returns the quotient of current divided by an int64 value (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) DividedInt(value int64) Rat {
	result := r // make a copy
	result.DivInt(value)
	return result
}

// Invert inverts the current rational number (mutable operation).
// Formula: a/b -> b/a (with sign moved to numerator)
// Sets invalid state on zero inversion or overflow.
func (r *Rat) Invert() {
	// If already invalid, remain invalid
	if r.IsInvalid() {
		return
	}

	// Check for inversion of zero (division by zero)
	if r.numerator == 0 {
		r.Invalidate()
		return
	}

	// For inversion: a/b -> b/a
	// We need to handle the sign correctly since numerator is signed and denominator is unsigned

	// Get the sign from the numerator
	isNegative := r.numerator < 0

	// Convert denominator to signed int64 for new numerator
	newNum, ok := uint64ToInt64WithSign(r.denominator, isNegative)
	if !ok {
		// Overflow when converting denominator to signed numerator
		r.Invalidate()
		return
	}

	// Convert absolute value of numerator to uint64 for new denominator
	newDenom := absInt64ToUint64(r.numerator)

	// Store the result
	r.numerator = newNum
	r.denominator = newDenom
}

// Inverted returns the inverse of the current rational number (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) Inverted() Rat {
	result := r // make a copy
	result.Invert()
	return result
}

// ScaleDown scales the rational number down by n decimal places (mutable operation).
// Equivalent to dividing by 10^n, moving the decimal point left.
// For negative n, calls ScaleUp with |n|.
// Sets invalid state when overflow remains after cancellation, the exact result
// is outside Rat limits, or operands are invalid.
func (r *Rat) ScaleDown(n int) {
	// If already invalid, remain invalid
	if r.IsInvalid() {
		return
	}

	// Handle zero scale - no operation needed
	if n == 0 {
		return
	}

	// Handle negative scale by calling ScaleUp
	if n < 0 {
		magnitude, ok := scaleMagnitude(n)
		if !ok {
			r.Invalidate()
			return
		}
		r.ScaleUp(magnitude)
		return
	}

	// Get power of 10
	powerOf10, overflow := powerOf10(n)
	if overflow {
		r.scaleDownLarge(n)
		return
	}

	// ScaleDown: divide by 10^n = multiply denominator by 10^n
	// Check for denominator overflow
	if willOverflowUint64Mul(r.denominator, powerOf10) {
		r.scaleDownReduced(powerOf10)
		return
	}

	r.denominator *= powerOf10
	r.reduceIfLarge()
}

// scaleDownLarge scales down when 10^n does not fit uint64 by processing decimal
// chunks. It cancels numerator factors before growing the denominator.
func (r *Rat) scaleDownLarge(n int) {
	value := *r
	value.Reduce()
	if value.numerator == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	remaining := n
	for remaining > 0 {
		chunk := remaining
		if chunk >= len(powersOf10) {
			chunk = len(powersOf10) - 1
		}

		scaleFactor := powersOf10[chunk]
		commonFactor := gcdInt64Uint64(value.numerator, scaleFactor)

		var ok bool
		value.numerator, ok = divInt64ByUint64Exact(value.numerator, commonFactor)
		if !ok {
			r.Invalidate()
			return
		}
		scaleFactor /= commonFactor

		if willOverflowUint64Mul(value.denominator, scaleFactor) {
			r.Invalidate()
			return
		}
		value.denominator *= scaleFactor
		remaining -= chunk
	}

	r.numerator = value.numerator
	r.denominator = value.denominator
	r.reduceIfLarge()
}

func (r *Rat) scaleDownReduced(scaleFactor uint64) {
	// Reduce first so existing common factors do not multiply into the denominator.
	value := *r
	value.Reduce()

	// ScaleDown divides by 10^n. Cancelling numerator factors against 10^n keeps
	// the denominator from growing when the numeric value allows it.
	commonFactor := gcdInt64Uint64(value.numerator, scaleFactor)
	newNum, ok := divInt64ByUint64Exact(value.numerator, commonFactor)
	if !ok {
		r.Invalidate()
		return
	}
	scaleFactor /= commonFactor

	if willOverflowUint64Mul(value.denominator, scaleFactor) {
		// The scaled denominator is still outside uint64 after all available
		// cancellation, so the result is not representable.
		r.Invalidate()
		return
	}

	r.numerator = newNum
	r.denominator = value.denominator * scaleFactor
	r.reduceIfLarge()
}

// ScaledDown returns a new rational number scaled down by n decimal places (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) ScaledDown(n int) Rat {
	result := r // make a copy
	result.ScaleDown(n)
	return result
}

// ScaleUp scales the rational number up by n decimal places (mutable operation).
// Equivalent to multiplying by 10^n, moving the decimal point right.
// For negative n, calls ScaleDown with |n|.
// Sets invalid state when overflow remains after cancellation, the exact result
// is outside Rat limits, or operands are invalid.
func (r *Rat) ScaleUp(n int) {
	// If already invalid, remain invalid
	if r.IsInvalid() {
		return
	}

	// Handle zero scale - no operation needed
	if n == 0 {
		return
	}

	// Handle negative scale by calling ScaleDown
	if n < 0 {
		magnitude, ok := scaleMagnitude(n)
		if !ok {
			r.Invalidate()
			return
		}
		r.ScaleDown(magnitude)
		return
	}

	// Get power of 10
	powerOf10, overflow := powerOf10(n)
	if overflow {
		r.scaleUpLarge(n)
		return
	}

	// ScaleUp: multiply by 10^n = multiply numerator by 10^n
	// Check for numerator overflow
	if willOverflowInt64MulUint64(r.numerator, powerOf10) {
		r.scaleUpReduced(powerOf10)
		return
	}

	// Handle multiplication with proper sign handling
	if r.numerator >= 0 {
		r.numerator *= int64(powerOf10) //nolint:gosec // overflow checked above
	} else {
		// Handle negative case carefully to avoid overflow
		absNum := uint64(-r.numerator)
		r.numerator = -int64(absNum * powerOf10) //nolint:gosec // overflow checked above
	}
	r.reduceIfLarge()
}

// scaleUpLarge scales up when 10^n does not fit uint64 by processing decimal
// chunks. It cancels denominator factors before growing the numerator.
func (r *Rat) scaleUpLarge(n int) {
	value := *r
	value.Reduce()
	if value.numerator == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	remaining := n
	for remaining > 0 {
		chunk := remaining
		if chunk >= len(powersOf10) {
			chunk = len(powersOf10) - 1
		}

		scaleFactor := powersOf10[chunk]
		commonFactor := gcdUint64(value.denominator, scaleFactor)
		value.denominator /= commonFactor
		scaleFactor /= commonFactor

		newNum, ok := mulInt64ByUint64ToInt64(value.numerator, scaleFactor)
		if !ok {
			r.Invalidate()
			return
		}
		value.numerator = newNum
		remaining -= chunk
	}

	r.numerator = value.numerator
	r.denominator = value.denominator
	r.reduceIfLarge()
}

func (r *Rat) scaleUpReduced(scaleFactor uint64) {
	// Reduce first so denominator factors that already match the numerator do not
	// hide cancellation against 10^n.
	value := *r
	value.Reduce()

	// ScaleUp multiplies by 10^n. Cancelling denominator factors against 10^n
	// prevents numerator growth when the denominator already contains powers of ten.
	commonFactor := gcdUint64(value.denominator, scaleFactor)
	value.denominator /= commonFactor
	scaleFactor /= commonFactor

	// If the remaining scale factor still makes the numerator overflow int64, the
	// exact scaled value is outside Rat limits.
	newNum, ok := mulInt64ByUint64ToInt64(value.numerator, scaleFactor)
	if !ok {
		r.Invalidate()
		return
	}

	r.numerator = newNum
	r.denominator = value.denominator
	r.reduceIfLarge()
}

// ScaledUp returns a new rational number scaled up by n decimal places (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) ScaledUp(n int) Rat {
	result := r // make a copy
	result.ScaleUp(n)
	return result
}
