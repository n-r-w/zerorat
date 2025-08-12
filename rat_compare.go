package zerorat

import (
	"fmt"
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