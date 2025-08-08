//nolint:gosec // false positive (no change needed)
package zerorat

import (
	"fmt"
	"math"
	"math/bits"
	"strconv"
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
	defer r.Reduce()

	// If denominator is 0, return invalid state
	if denominator == 0 {
		return Rat{numerator: numerator, denominator: 0}
	}

	// If numerator is 0, normalize to 0/1
	if numerator == 0 {
		return Rat{numerator: 0, denominator: 1}
	}

	// Return the rational number as-is
	return Rat{numerator: numerator, denominator: denominator}
}

// NewFromInt creates a rational number from an integer.
// Equivalent to NewRat(value, 1).
func NewFromInt(value int64) Rat {
	return Rat{numerator: value, denominator: 1}
}

// NewFromFloat64 creates a rational number from a float64 with minimum precision loss.
// Returns invalid state (denominator = 0) for special values: NaN, +Inf, -Inf.
// Returns invalid state if the conversion would overflow int64/uint64 limits.
func NewFromFloat64(value float64) (r Rat) {
	defer r.Reduce()

	// Handle special cases
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return Rat{numerator: 0, denominator: 0} // invalid state
	}

	// Handle zero (including negative zero)
	if value == 0.0 {
		return Rat{numerator: 0, denominator: 1}
	}

	// Use IEEE 754 decomposition for exact conversion
	return float64ToRatExact(value)
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
func (r *Rat) IsValid() bool {
	return r.denominator > 0
}

// IsInvalid checks if the rational number is invalid.
// Returns true if denominator == 0.
func (r *Rat) IsInvalid() bool {
	return r.denominator == 0
}

// Invalidate marks the rational number as invalid,
// by setting denominator to 0.
func (r *Rat) Invalidate() {
	r.denominator = 0
}

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
	// Check a*d overflow
	if willOverflowInt64Mul(r.numerator, int64(other.denominator)) {
		r.Invalidate()
		return
	}

	// Check c*b overflow
	if willOverflowInt64Mul(other.numerator, int64(r.denominator)) {
		r.Invalidate()
		return
	}

	term1 := r.numerator * int64(other.denominator)
	term2 := other.numerator * int64(r.denominator)

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

	// Check for numerator * denominator overflow
	if willOverflowInt64Mul(r.numerator, int64(other.denominator)) {
		r.Invalidate()
		return
	}

	// Check for denominator * numerator overflow
	if willOverflowUint64Mul(r.denominator, otherNumAbs) {
		r.Invalidate()
		return
	}

	newNum := r.numerator * int64(other.denominator)
	newDenom := r.denominator * otherNumAbs

	// Apply sign: if other.numerator was negative, negate result
	if other.numerator < 0 {
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
func (r *Rat) Added(other Rat) Rat {
	result := *r // make a copy
	result.Add(other)
	return result
}

// Subtracted returns the difference of current and another rational number (immutable operation).
// Doesn't modify the original rational number.
func (r *Rat) Subtracted(other Rat) Rat {
	result := *r // make a copy
	result.Sub(other)
	return result
}

// Multiplied returns the product of current and another rational number (immutable operation).
// Doesn't modify the original rational number.
func (r *Rat) Multiplied(other Rat) Rat {
	result := *r // make a copy
	result.Mul(other)
	return result
}

// Divided returns the quotient of current divided by another rational number (immutable operation).
// Doesn't modify the original rational number.
func (r *Rat) Divided(other Rat) Rat {
	result := *r // make a copy
	result.Div(other)
	return result
}

// Equal checks equality of two rational numbers.
// Returns false for any invalid operands, consistent with comparison semantics.
func (r *Rat) Equal(other Rat) bool {
	// Invalid operands are never equal to anything (including other invalid operands)
	if r.IsInvalid() || other.IsInvalid() {
		return false
	}
	return compareRationalsCrossMul(r.numerator, r.denominator, other.numerator, other.denominator) == 0
}

// Less checks if current rational number is less than another.
// Returns false for any invalid operands, consistent with comparison semantics.
func (r *Rat) Less(other Rat) bool {
	// Invalid operands cannot be ordered
	if r.IsInvalid() || other.IsInvalid() {
		return false
	}
	return compareRationalsCrossMul(r.numerator, r.denominator, other.numerator, other.denominator) < 0
}

// Greater checks if current rational number is greater than another.
// Returns false for any invalid operands, consistent with comparison semantics.
func (r *Rat) Greater(other Rat) bool {
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
func (r *Rat) Compare(other Rat) int {
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
func (r *Rat) String() string {
	if r.IsInvalid() {
		return "invalid"
	}

	if r.denominator == 1 {
		return strconv.FormatInt(r.numerator, 10)
	}

	return fmt.Sprintf("%d/%d", r.numerator, r.denominator)
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
		r.numerator /= int64(gcd)
		r.denominator /= gcd
	}
}

// Reduced returns a new rational number reduced to its lowest terms (immutable operation).
// Does not modify the original rational number.
func (r *Rat) Reduced() Rat {
	result := *r // create copy
	result.Reduce()
	return result
}

// willOverflowUint64Mul checks if multiplying two uint64 values would overflow.
// Uses math/bits for improved clarity and performance.
func willOverflowUint64Mul(a, b uint64) bool {
	if a == 0 || b == 0 {
		return false
	}
	// Use bits.Mul64 to detect overflow
	hi, _ := bits.Mul64(a, b)
	return hi != 0
}

// willOverflowInt64Mul checks if multiplying two int64 values would overflow.
// Uses math/bits for improved clarity and performance.
func willOverflowInt64Mul(a, b int64) bool {
	if a == 0 || b == 0 {
		return false
	}

	// Convert to unsigned for bits.Mul64, handling signs separately
	aAbs := absInt64ToUint64(a)
	bAbs := absInt64ToUint64(b)

	// Use bits.Mul64 to detect overflow
	hi, lo := bits.Mul64(aAbs, bAbs)

	// Check if result fits in int64 range
	sameSign := (a > 0) == (b > 0)
	if sameSign {
		// Positive result: must fit in [0, MaxInt64]
		return hi != 0 || lo > uint64(math.MaxInt64)
	}
	// Negative result: must fit in [MinInt64, -1]
	// MaxInt64 + 1 = 9223372036854775808 (absolute value of MinInt64)
	return hi != 0 || lo > 9223372036854775808
}

// willOverflowInt64Add checks if adding two int64 values would overflow.
// Uses simple range checking for clarity and correctness.
func willOverflowInt64Add(a, b int64) bool {
	if b > 0 {
		return a > math.MaxInt64-b
	}
	return a < math.MinInt64-b
}

// willOverflowInt64Sub checks if subtracting two int64 values would overflow.
// Uses simple range checking for clarity and correctness.
func willOverflowInt64Sub(a, b int64) bool {
	if b > 0 {
		return a < math.MinInt64+b
	}
	return a > math.MaxInt64+b
}

// gcdInt64Uint64 calculates the greatest common divisor of int64 and uint64.
func gcdInt64Uint64(a int64, b uint64) uint64 {
	// Use absolute value for int64
	absA := absInt64ToUint64(a)
	return gcdUint64(absA, b)
}

// gcdUint64 calculates the greatest common divisor of two uint64 values using Euclid's algorithm.
func gcdUint64(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// absInt64ToUint64 converts an int64 to its absolute value as uint64.
// Handles the special case where math.MinInt64 cannot be negated within int64 range.
func absInt64ToUint64(value int64) uint64 {
	if value < 0 {
		if value == math.MinInt64 {
			// Special case: absolute value of MinInt64 doesn't fit in int64
			return uint64(math.MaxInt64) + 1
		}
		return uint64(-value)
	}
	return uint64(value)
}

// compare128Bit compares two 128-bit numbers represented as (hi, lo) pairs.
// Returns -1 if first < second, 0 if equal, 1 if first > second.
func compare128Bit(hi1, lo1, hi2, lo2 uint64) int {
	if hi1 < hi2 {
		return -1
	}
	if hi1 > hi2 {
		return 1
	}
	// High parts are equal, compare low parts
	if lo1 < lo2 {
		return -1
	}
	if lo1 > lo2 {
		return 1
	}
	return 0
}

// compareRationalsCrossMul compares two rational numbers using 128-bit cross-multiplication.
// Returns -1 if a/b < c/d, 0 if a/b == c/d, 1 if a/b > c/d.
// Uses math/bits to handle potential overflow in intermediate calculations.
func compareRationalsCrossMul(aNum int64, aDenom uint64, cNum int64, cDenom uint64) int {
	// Handle signs separately to work with unsigned arithmetic
	aSign := 1
	if aNum < 0 {
		aSign = -1
	}
	cSign := 1
	if cNum < 0 {
		cSign = -1
	}

	// Get absolute values for unsigned arithmetic
	var aAbs, cAbs uint64
	aAbs = absInt64ToUint64(aNum)

	cAbs = absInt64ToUint64(cNum)

	// Calculate a*d and c*b using 128-bit arithmetic
	aTimesDHi, aTimesDLo := bits.Mul64(aAbs, cDenom)
	cTimesBHi, cTimesBLo := bits.Mul64(cAbs, aDenom)

	// Compare the 128-bit results
	cmpResult := compare128Bit(aTimesDHi, aTimesDLo, cTimesBHi, cTimesBLo)

	// Apply sign logic - simplified
	if aSign != cSign {
		// Different signs: negative < positive
		if aSign < 0 {
			return -1
		}
		return 1
	}
	// Same signs: if both negative, reverse magnitude comparison
	if aSign < 0 {
		return -cmpResult
	}
	// Both positive: use direct magnitude comparison
	return cmpResult
}

// float64ToRatExact converts a float64 to its exact rational representation.
// Returns invalid state if conversion would overflow int64/uint64 limits.
//
//nolint:mnd // magic numbers are fine here
func float64ToRatExact(value float64) Rat {
	if math.IsNaN(value) {
		return Rat{}
	}
	if value < -9.223372036854775e+18 || value > 9.223372036854775e+18 {
		return Rat{}
	}
	if value > -2.168404344971009e-19 && value < 2.168404344971009e-19 {
		return Rat{}
	}

	// Decompose float64
	fBits := math.Float64bits(value)
	isNegative := fBits&(1<<63) != 0
	exp := int64((fBits>>52)&(1<<11-1)) - 1023
	mantissa := (fBits & (1<<52 - 1)) | 1<<52 // Since we discarded tiny values, it'll never be denormalized.

	// Amount of times to shift the mantissa to the right to compensate for the exponent
	shift := 52 - exp

	// Reduce shift and mantissa as far as we can
	for mantissa&1 == 0 && shift > 0 {
		mantissa >>= 1
		shift--
	}

	// Choose whether to shift the numerator or denominator
	var shiftN, shiftD int64 = 0, 0
	if shift > 0 {
		shiftD = shift
	} else {
		shiftN = shift
	}

	// Shift that require larger shifts that what an int64 can hold, or larger than the mantissa itself, will be
	// approximated splitting it between the numerator and denominator.
	if shiftD > 62 {
		shiftD = 62
		shiftN = shift - 62
	} else if shiftN > 52 {
		shiftN = 52
		shiftD = shift - 52
	}

	numerator, denominator := int64(mantissa), int64(1)
	denominator <<= shiftD
	if shiftN < 0 {
		numerator <<= -shiftN
	} else {
		numerator >>= shiftN
	}

	if isNegative {
		numerator *= -1
	}

	return Rat{numerator: numerator, denominator: uint64(denominator)}
}
