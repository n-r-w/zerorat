package zerorat

import (
	"math"
)

// RoundType defines rounding strategies for rounding rationals to integers.
type RoundType int

const (
	// RoundDown rounds toward zero.
	RoundDown RoundType = iota
	// RoundUp rounds away from zero.
	RoundUp
	// RoundHalfUp (half up).
	RoundHalfUp
)

// Rat represents a rational number without heap allocation.
// Uses denominator = 0 to represent an invalid state.
type Rat struct {
	numerator   int64  // Signed numerator
	denominator uint64 // Denominator (always positive, 0 = invalid state)
}

// New creates a new rational number with given numerator and denominator.
// Returns a value, not a pointer.
func New(numerator int64, denominator uint64) (r Rat) {
	// If denominator is 0, return invalid state
	if denominator == 0 {
		return Rat{numerator: numerator, denominator: 0}
	}

	// If numerator is 0, normalize to 0/1
	if numerator == 0 {
		return Rat{numerator: 0, denominator: 1}
	}

	// Construct and explicitly reduce (hot path without defer)
	r = Rat{numerator: numerator, denominator: denominator}
	r.Reduce()
	return r
}

// NewFromInt creates a rational number from an integer.
// Equivalent to New(value, 1).
func NewFromInt(value int64) Rat {
	return Rat{numerator: value, denominator: 1}
}

// NewFromFloat64 creates a rational number from a float64 with minimum precision loss.
// Returns invalid state (denominator = 0) for special values: NaN, +Inf, -Inf.
// Returns invalid state if the conversion would overflow int64/uint64 limits.
func NewFromFloat64(value float64) (r Rat) {
	// Handle special cases
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return Rat{numerator: 0, denominator: 0} // invalid state
	}

	// Handle zero (including negative zero)
	if value == 0.0 {
		return Rat{numerator: 0, denominator: 1}
	}

	// Use IEEE 754 decomposition for exact conversion.
	// Note: NewFromFloat64 must invalidate on overflow; float64ToRatExact returns Rat{}
	// when representation would exceed int64/uint64 bounds.
	r = float64ToRatExact(value)
	if r.IsValid() {
		r.Reduce()
	}
	return r
}

// Zero returns a rational number representing zero (0/1).
func Zero() Rat {
	return Rat{numerator: 0, denominator: 1}
}

// One returns a rational number representing one (1/1).
func One() Rat {
	return Rat{numerator: 1, denominator: 1}
}

// IsValid checks if the rational number is valid.
// Returns true if denominator > 0.
func (r Rat) IsValid() bool {
	return r.denominator > 0
}

// IsInvalid checks if the rational number is invalid.
// Returns true if denominator == 0.
func (r Rat) IsInvalid() bool {
	return r.denominator == 0
}

// Invalidate marks the rational number as invalid,
// by setting denominator to 0.
func (r *Rat) Invalidate() {
	r.denominator = 0
}

// Numerator returns the numerator of rational number.
func (r Rat) Numerator() int64 {
	return r.numerator
}

// Denominator returns the denominator of rational number.
func (r Rat) Denominator() uint64 {
	return r.denominator
}

// Sign returns the sign of rational number.
// Returns -1 for negative, 0 for zero or invalid, 1 for positive.
func (r Rat) Sign() int {
	if r.IsInvalid() {
		return 0
	}

	if r.numerator < 0 {
		return -1
	} else if r.numerator > 0 {
		return 1
	}
	return 0
}

// IsZero checks if rational number equals zero.
func (r Rat) IsZero() bool {
	return r.IsValid() && r.numerator == 0
}

// IsOne checks if rational number equals one.
func (r Rat) IsOne() bool {
	return r.IsValid() && r.numerator == 1 && r.denominator == 1
}
