package zerorat

import (
	"errors"
	"math"
)

// Error definitions for Rat operations.
var (
	// ErrInvalid indicates that a Rat value is in an invalid state.
	ErrInvalid = errors.New("invalid rat")
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

// NewFromInt64 creates a rational number from an integer.
// Equivalent to New(value, 1).
func NewFromInt64(value int64) Rat {
	return Rat{numerator: value, denominator: 1}
}

// NewFromInt creates a rational number from an integer.
func NewFromInt(value int) Rat {
	return NewFromInt64(int64(value))
}

// NewFromInt64Ptr creates a rational number from an integer pointer.
func NewFromInt64Ptr(value *int64) Rat {
	if value == nil {
		return Rat{}
	}
	return NewFromInt64(*value)
}

// NewFromIntPtr creates a rational number from an integer pointer.
func NewFromIntPtr(value *int) Rat {
	if value == nil {
		return Rat{}
	}
	return NewFromInt64(int64(*value))
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

// NewFromFloat64Ptr creates a rational number from a float64 pointer.
func NewFromFloat64Ptr(value *float64) Rat {
	if value == nil {
		return Rat{}
	}
	return NewFromFloat64(*value)
}

// NewFromFloat32 creates a rational number from a float32 with minimum precision loss.
func NewFromFloat32(value float32) Rat {
	return NewFromFloat64(float64(value))
}

// NewFromFloat32Ptr creates a rational number from a float32 pointer.
func NewFromFloat32Ptr(value *float32) Rat {
	if value == nil {
		return Rat{}
	}
	return NewFromFloat64(float64(*value))
}

// FromInt64Slice creates a rational slice from an integer slice.
func FromInt64Slice(values []int64) []Rat {
	if values == nil {
		return nil
	}

	result := make([]Rat, len(values))
	for i, v := range values {
		result[i] = NewFromInt64(v)
	}
	return result
}

// FromInt64SlicePtr creates a rational slice from an integer slice pointer.
func FromInt64SlicePtr(values []*int64) []Rat {
	if values == nil {
		return nil
	}

	result := make([]Rat, len(values))
	for i, v := range values {
		if v == nil {
			result[i] = Rat{}
		} else {
			result[i] = NewFromInt64(*v)
		}
	}
	return result
}

// FromIntSlice creates a rational slice from an integer slice.
func FromIntSlice(values []int) []Rat {
	if values == nil {
		return nil
	}

	result := make([]Rat, len(values))
	for i, v := range values {
		result[i] = NewFromInt64(int64(v))
	}
	return result
}

// FromIntSlicePtr creates a rational slice from an integer slice pointer.
func FromIntSlicePtr(values []*int) []Rat {
	if values == nil {
		return nil
	}

	result := make([]Rat, len(values))
	for i, v := range values {
		if v == nil {
			result[i] = Rat{}
		} else {
			result[i] = NewFromInt64(int64(*v))
		}
	}
	return result
}

// FromFloat64Slice creates a rational slice from a float64 slice.
func FromFloat64Slice(values []float64) []Rat {
	if values == nil {
		return nil
	}

	result := make([]Rat, len(values))
	for i, v := range values {
		result[i] = NewFromFloat64(v)
	}
	return result
}

// FromFloat64SlicePtr creates a rational slice from a float64 slice pointer.
func FromFloat64SlicePtr(values []*float64) []Rat {
	if values == nil {
		return nil
	}

	result := make([]Rat, len(values))
	for i, v := range values {
		if v == nil {
			result[i] = Rat{}
		} else {
			result[i] = NewFromFloat64(*v)
		}
	}
	return result
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

// HasFractional checks if rational number has a fractional part.
// Returns true if there is a remainder from division (numerator % denominator != 0).
// Returns false for whole numbers, zero, and invalid states.
func (r Rat) HasFractional() bool {
	// Invalid rationals return false
	if r.IsInvalid() {
		return false
	}

	// Zero has no fractional part
	if r.numerator == 0 {
		return false
	}

	// Check if there's a remainder when dividing numerator by denominator
	// For negative numerators, we need to use absolute value for modulo operation
	if r.numerator < 0 {
		return uint64(-r.numerator)%r.denominator != 0
	}
	return uint64(r.numerator)%r.denominator != 0
}

// IntegerAndFraction returns the integer and fractional parts of the rational number.
// For 7/3, returns (2, 1/3). For -7/3, returns (-2, -1/3).
// For whole numbers like 5/1, returns (5, 0/1).
// For proper fractions like 1/2, returns (0, 1/2).
// Returns (0, invalid) for invalid rational numbers.
func (r Rat) IntegerAndFraction() (int64, Rat) {
	// Invalid rationals return (0, invalid)
	if r.IsInvalid() {
		return 0, Rat{numerator: 0, denominator: 0}
	}

	// Zero returns (0, 0/1)
	if r.numerator == 0 {
		return 0, Rat{numerator: 0, denominator: 1}
	}

	// Calculate integer part and remainder
	// We need to handle the case where denominator > MaxInt64
	var integerPart int64
	var remainder int64

	if r.denominator <= uint64(math.MaxInt64) {
		// Safe to convert denominator to int64
		denom := int64(r.denominator)
		integerPart = r.numerator / denom
		remainder = r.numerator % denom
	} else {
		// Denominator is larger than MaxInt64
		// For any int64 numerator, integer part will be 0 (since |numerator| <= MaxInt64 < denominator)
		integerPart = 0
		remainder = r.numerator
	}

	// Create fractional part as remainder/denominator
	// The fractional part keeps the same sign as the original number
	fractionalPart := Rat{numerator: remainder, denominator: r.denominator}

	return integerPart, fractionalPart
}

// IntegerPart returns the integer part of the rational number.
func (r Rat) IntegerPart() int64 {
	intPart, _ := r.IntegerAndFraction()
	return intPart
}

// FractionalPart returns the fractional part of the rational number.
func (r Rat) FractionalPart() Rat {
	_, frac := r.IntegerAndFraction()
	return frac
}

// ToInt64Err converts the rational number to an int64 with error handling.
// Returns ErrInvalid if the rational number is invalid.
func (r Rat) ToInt64Err() (int64, error) {
	if r.IsInvalid() {
		return 0, ErrInvalid
	}
	return r.IntegerPart(), nil
}

// ToInt64 converts the rational number to an int64.
// Returns 0 if the rational number is invalid.
func (r Rat) ToInt64() int64 {
	result, _ := r.ToInt64Err()
	return result
}

// ToIntErr converts the rational number to an int with error handling.
// Returns ErrInvalid if the rational number is invalid.
// This is an alias for ToInt64Err for compatibility with int type.
func (r Rat) ToIntErr() (int, error) {
	result, err := r.ToInt64Err()
	if err != nil {
		return 0, err
	}
	return int(result), nil
}

// ToInt converts the rational number to an int.
// Returns 0 if the rational number is invalid.
// This is an alias for ToInt64 for compatibility with int type.
func (r Rat) ToInt() int {
	result, _ := r.ToIntErr()
	return result
}

// ToInt64Ptr converts the rational number to an int64 pointer.
func (r Rat) ToInt64Ptr() *int64 {
	if r.IsInvalid() {
		return nil
	}
	result := r.ToInt64()
	return &result
}

// ToIntPtr converts the rational number to an int pointer.
func (r Rat) ToIntPtr() *int {
	if r.IsInvalid() {
		return nil
	}
	result := r.ToInt()
	return &result
}

// ToInt64Slice converts a slice of rational numbers to a slice of int64.
func ToInt64Slice(values []Rat) []int64 {
	if values == nil {
		return nil
	}

	result := make([]int64, len(values))
	for i, v := range values {
		result[i] = v.ToInt64()
	}
	return result
}

// ToIntSlice converts a slice of rational numbers to a slice of int.
func ToIntSlice(values []Rat) []int {
	if values == nil {
		return nil
	}

	result := make([]int, len(values))
	for i, v := range values {
		result[i] = v.ToInt()
	}
	return result
}

// ToInt64SlicePtr converts a slice of rational numbers to a slice of int64 pointers.
func ToInt64SlicePtr(values []Rat) []*int64 {
	if values == nil {
		return nil
	}

	result := make([]*int64, len(values))
	for i, v := range values {
		result[i] = v.ToInt64Ptr()
	}
	return result
}

// ToIntSlicePtr converts a slice of rational numbers to a slice of int pointers.
func ToIntSlicePtr(values []Rat) []*int {
	if values == nil {
		return nil
	}

	result := make([]*int, len(values))
	for i, v := range values {
		result[i] = v.ToIntPtr()
	}
	return result
}

// ToFloat64Err converts the rational number to a float64 with error handling.
// Returns ErrInvalid if the rational number is invalid.
func (r Rat) ToFloat64Err() (float64, error) {
	if r.IsInvalid() {
		return 0, ErrInvalid
	}
	return float64(r.numerator) / float64(r.denominator), nil
}

// ToFloat64 converts the rational number to a float64.
// Returns 0 if the rational number is invalid.
func (r Rat) ToFloat64() float64 {
	result, _ := r.ToFloat64Err()
	return result
}

// ToFloat32Err converts the rational number to a float32 with error handling.
// Returns ErrInvalid if the rational number is invalid or if the conversion would overflow float32.
func (r Rat) ToFloat32Err() (float32, error) {
	if r.IsInvalid() {
		return 0, ErrInvalid
	}

	// Convert to float64 first to check for overflow
	result64 := float64(r.numerator) / float64(r.denominator)

	// Check for overflow to float32
	if math.IsInf(result64, 0) || math.IsNaN(result64) {
		return 0, ErrInvalid
	}

	// Check if the result is within float32 range
	if result64 > math.MaxFloat32 || result64 < -math.MaxFloat32 {
		return 0, ErrInvalid
	}

	return float32(result64), nil
}

// ToFloat32 converts the rational number to a float32.
// Returns 0 if the rational number is invalid or if the conversion would overflow float32.
func (r Rat) ToFloat32() float32 {
	result, _ := r.ToFloat32Err()
	return result
}

// ToFloat32Ptr converts the rational number to a float32 pointer.
func (r Rat) ToFloat32Ptr() *float32 {
	if r.IsInvalid() {
		return nil
	}
	result := r.ToFloat32()
	return &result
}

// ToFloat64Ptr converts the rational number to a float64 pointer.
func (r Rat) ToFloat64Ptr() *float64 {
	if r.IsInvalid() {
		return nil
	}
	result := r.ToFloat64()
	return &result
}

// ToFloat64Slice converts a slice of rational numbers to a slice of float64.
func ToFloat64Slice(values []Rat) []float64 {
	if values == nil {
		return nil
	}

	result := make([]float64, len(values))
	for i, v := range values {
		result[i] = v.ToFloat64()
	}
	return result
}

// ToFloat32Slice converts a slice of rational numbers to a slice of float32.
func ToFloat32Slice(values []Rat) []float32 {
	if values == nil {
		return nil
	}

	result := make([]float32, len(values))
	for i, v := range values {
		result[i] = v.ToFloat32()
	}
	return result
}

// ToFloat64SlicePtr converts a slice of rational numbers to a slice of float64 pointers.
func ToFloat64SlicePtr(values []Rat) []*float64 {
	if values == nil {
		return nil
	}

	result := make([]*float64, len(values))
	for i, v := range values {
		result[i] = v.ToFloat64Ptr()
	}
	return result
}

// ToFloat32SlicePtr converts a slice of rational numbers to a slice of float32 pointers.
func ToFloat32SlicePtr(values []Rat) []*float32 {
	if values == nil {
		return nil
	}

	result := make([]*float32, len(values))
	for i, v := range values {
		result[i] = v.ToFloat32Ptr()
	}
	return result
}
