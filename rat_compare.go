package zerorat

import (
	"fmt"
	"math"
	"strconv"
)

// Equal checks equality of two rational numbers.
// Returns false for any invalid operands, consistent with comparison semantics.
func (r Rat) Equal(other Rat) bool {
	// Invalid operands are never equal to anything (including other invalid operands)
	if r.IsInvalid() || other.IsInvalid() {
		return false
	}
	return compareRationalsCrossMul(r.numerator, r.denominator, other.numerator, other.denominator) == 0
}

// Less checks if current rational number is less than another.
// Returns false for any invalid operands, consistent with comparison semantics.
func (r Rat) Less(other Rat) bool {
	// Invalid operands cannot be ordered
	if r.IsInvalid() || other.IsInvalid() {
		return false
	}
	return compareRationalsCrossMul(r.numerator, r.denominator, other.numerator, other.denominator) < 0
}

// Greater checks if current rational number is greater than another.
// Returns false for any invalid operands, consistent with comparison semantics.
func (r Rat) Greater(other Rat) bool {
	// Invalid operands cannot be ordered
	if r.IsInvalid() || other.IsInvalid() {
		return false
	}
	return compareRationalsCrossMul(r.numerator, r.denominator, other.numerator, other.denominator) > 0
}

// Compare performs three-way comparison of rational numbers.
// Returns -1 if r < other, 0 if r == other, 1 if r > other.
// Returns 0 for any invalid operands (cannot be meaningfully compared).
// Uses single 128-bit cross-multiplication for optimal performance.
func (r Rat) Compare(other Rat) int {
	// Invalid operands cannot be meaningfully compared - return equal
	if r.IsInvalid() || other.IsInvalid() {
		return 0
	}

	// Normalize zeros: 0/x == 0/y for any non-zero x, y
	if r.numerator == 0 && other.numerator == 0 {
		return 0
	}

	// Use single 128-bit cross-multiplication for optimal performance
	return compareRationalsCrossMul(r.numerator, r.denominator, other.numerator, other.denominator)
}

// EqualInt64 checks equality against an int64 value.
// Reuses Rat comparison semantics, including false for invalid operands.
func (r Rat) EqualInt64(other int64) bool {
	return r.Equal(NewFromInt64(other))
}

// LessInt64 checks whether this rational value is less than an int64 value.
// Reuses Rat comparison semantics, including false for invalid operands.
func (r Rat) LessInt64(other int64) bool {
	return r.Less(NewFromInt64(other))
}

// GreaterInt64 checks whether this rational value is greater than an int64 value.
// Reuses Rat comparison semantics, including false for invalid operands.
func (r Rat) GreaterInt64(other int64) bool {
	return r.Greater(NewFromInt64(other))
}

// CompareInt64 performs three-way comparison against an int64 value.
// Reuses Rat comparison semantics, including 0 for invalid operands.
func (r Rat) CompareInt64(other int64) int {
	return r.Compare(NewFromInt64(other))
}

// EqualFloat64 checks equality against a float64 value.
// Non-finite floats cannot be represented as Rat and therefore compare as false.
func (r Rat) EqualFloat64(other float64) bool {
	otherRat, err := NewFromFloat64(other)
	if err != nil {
		return false
	}

	return r.Equal(otherRat)
}

// LessFloat64 checks whether this rational value is less than a float64 value.
// Non-finite floats cannot be represented as Rat and therefore compare as false.
func (r Rat) LessFloat64(other float64) bool {
	otherRat, err := NewFromFloat64(other)
	if err != nil {
		return false
	}

	return r.Less(otherRat)
}

// GreaterFloat64 checks whether this rational value is greater than a float64 value.
// Non-finite floats cannot be represented as Rat and therefore compare as false.
func (r Rat) GreaterFloat64(other float64) bool {
	otherRat, err := NewFromFloat64(other)
	if err != nil {
		return false
	}

	return r.Greater(otherRat)
}

// CompareFloat64 performs three-way comparison against a float64 value.
// Non-finite floats cannot be represented as Rat and therefore compare as 0.
func (r Rat) CompareFloat64(other float64) int {
	otherRat, err := NewFromFloat64(other)
	if err != nil {
		return 0
	}

	return r.Compare(otherRat)
}

// ApproxEqualFloat64 checks whether this rational value is within eps of a float64 value.
// It compares in float64 space because epsilon semantics belong to float values, not exact rationals.
// Non-finite floats, invalid Rat values, and negative epsilon return false.
func (r Rat) ApproxEqualFloat64(other, eps float64) bool {
	if r.IsInvalid() || math.IsNaN(other) || math.IsInf(other, 0) || math.IsNaN(eps) || math.IsInf(eps, 0) || eps < 0 {
		return false
	}

	value, err := r.ToFloat64Err()
	if err != nil {
		return false
	}

	return math.Abs(value-other) <= eps
}

// String returns string representation of rational number.
// Format: "numerator/denominator" or "numerator" if denominator == 1.
// Returns "invalid" for invalid state.
func (r Rat) String() string {
	if r.IsInvalid() {
		return "invalid"
	}

	if r.numerator == 0 {
		return "0"
	}

	if r.denominator == 1 {
		return strconv.FormatInt(r.numerator, 10)
	}

	return fmt.Sprintf("%d/%d", r.numerator, r.denominator)
}
