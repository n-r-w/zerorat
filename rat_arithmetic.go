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

// Divided returns the quotient of current divided by another rational number (immutable operation).
// Doesn't modify the original rational number.
func (r Rat) Divided(other Rat) Rat {
	result := r // make a copy
	result.Div(other)
	return result
}