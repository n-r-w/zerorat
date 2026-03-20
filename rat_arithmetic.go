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
		var newNum int64
		var overflowCheck bool

		if isAdd {
			overflowCheck = willOverflowInt64Add(r.numerator, other.numerator)
			newNum = r.numerator + other.numerator
		} else {
			overflowCheck = willOverflowInt64Sub(r.numerator, other.numerator)
			newNum = r.numerator - other.numerator
		}

		if overflowCheck {
			r.Invalidate()
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
		return
	}

	// General case: different denominators
	// Check for denominator multiplication overflow
	if willOverflowUint64Mul(r.denominator, other.denominator) {
		r.Invalidate()
		return
	}

	// Calculate new denominator
	newDenom := r.denominator * other.denominator

	// Check for numerator calculation overflow
	// Check a*d overflow and compute safely
	prod1, ok := mulInt64ByUint64ToInt64(r.numerator, other.denominator)
	if !ok {
		r.Invalidate()
		return
	}

	// Check c*b overflow and compute safely
	prod2, ok := mulInt64ByUint64ToInt64(other.numerator, r.denominator)
	if !ok {
		r.Invalidate()
		return
	}

	term1 := prod1
	term2 := prod2

	var newNum int64
	var overflowCheck bool

	if isAdd {
		overflowCheck = willOverflowInt64Add(term1, term2)
		newNum = term1 + term2
	} else {
		overflowCheck = willOverflowInt64Sub(term1, term2)
		newNum = term1 - term2
	}

	if overflowCheck {
		r.Invalidate()
		return
	}

	// If the result is zero, normalize to 0/1
	if newNum == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// Store result without automatic reduction
	r.numerator = newNum
	r.denominator = newDenom
}

// Add adds another rational number to the current one (mutable operation).
// Formula: a/b + c/d = (a*d + c*b) / (b*d)
// Result is not automatically reduced to lowest terms - use Reduce() if needed.
// Sets invalid state on overflow or with invalid operands.
func (r *Rat) Add(other Rat) {
	r.addSubCommon(other, true)
}

// Sub subtracts another rational number from the current one (mutable operation).
// Formula: a/b - c/d = (a*d - c*b) / (b*d)
// Result is not automatically reduced to lowest terms - use Reduce() if needed.
// Sets invalid state on overflow or with invalid operands.
func (r *Rat) Sub(other Rat) {
	r.addSubCommon(other, false)
}

// Mul multiplies the current rational number by another (mutable operation).
// Formula: a/b * c/d = (a*c) / (b*d)
// Result is not automatically reduced to lowest terms - use Reduce() if needed.
// Sets invalid state on overflow or with invalid operands.
func (r *Rat) Mul(other Rat) {
	// If any operand is invalid, result is invalid
	if r.IsInvalid() || other.IsInvalid() {
		r.Invalidate()
		return
	}

	// Check numerator multiplication overflow
	if willOverflowInt64Mul(r.numerator, other.numerator) {
		r.Invalidate()
		return
	}

	// Check denominator multiplication overflow
	if willOverflowUint64Mul(r.denominator, other.denominator) {
		r.Invalidate()
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

	// Store result without automatic reduction
	r.numerator = newNum
	r.denominator = newDenom
}

// Div divides the current rational number by another (mutable operation).
// Formula: a/b ÷ c/d = a/b * d/c = (a*d) / (b*c)
// Result is not automatically reduced to lowest terms - use Reduce() if needed.
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

	// Check for numerator * denominator overflow and compute safely
	prodNum, ok := mulInt64ByUint64ToInt64(r.numerator, other.denominator)
	if !ok {
		r.Invalidate()
		return
	}

	// Check for denominator * numerator overflow
	if willOverflowUint64Mul(r.denominator, otherNumAbs) {
		r.Invalidate()
		return
	}

	newNum := prodNum
	newDenom := r.denominator * otherNumAbs

	// Apply sign: if other.numerator was negative, negate result
	if other.numerator < 0 {
		if newNum == math.MinInt64 {
			// cannot negate MinInt64 safely; treat as overflow
			r.Invalidate()
			return
		}
		newNum = -newNum
	}

	// If result is zero, normalize to 0/1
	if newNum == 0 {
		r.numerator = 0
		r.denominator = 1
		return
	}

	// Store result without automatic reduction
	r.numerator = newNum
	r.denominator = newDenom
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
// Sets invalid state on overflow or with invalid operands.
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
		r.ScaleUp(-n)
		return
	}

	// Get power of 10
	powerOf10, overflow := powerOf10(n)
	if overflow {
		r.Invalidate()
		return
	}

	// ScaleDown: divide by 10^n = multiply denominator by 10^n
	// Check for denominator overflow
	if willOverflowUint64Mul(r.denominator, powerOf10) {
		r.Invalidate()
		return
	}

	r.denominator *= powerOf10
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
// Sets invalid state on overflow or with invalid operands.
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
		r.ScaleDown(-n)
		return
	}

	// Get power of 10
	powerOf10, overflow := powerOf10(n)
	if overflow {
		r.Invalidate()
		return
	}

	// ScaleUp: multiply by 10^n = multiply numerator by 10^n
	// Check for numerator overflow
	if willOverflowInt64MulUint64(r.numerator, powerOf10) {
		r.Invalidate()
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
}

// ScaledUp returns a new rational number scaled up by n decimal places (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) ScaledUp(n int) Rat {
	result := r // make a copy
	result.ScaleUp(n)
	return result
}
