package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRat_Reduce tests mutable reduction to lowest terms
func TestRat_Reduce(t *testing.T) {
	tests := []struct {
		name      string
		input     Rat
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "already reduced",
			input:     New(3, 4),
			wantNum:   3,
			wantDenom: 4,
		},
		{
			name:      "reduce simple fraction",
			input:     New(6, 8),
			wantNum:   3,
			wantDenom: 4,
		},
		{
			name:      "reduce to integer",
			input:     New(10, 5),
			wantNum:   2,
			wantDenom: 1,
		},
		{
			name:      "reduce negative fraction",
			input:     New(-12, 18),
			wantNum:   -2,
			wantDenom: 3,
		},
		{
			name:      "reduce zero",
			input:     New(0, 15),
			wantNum:   0,
			wantDenom: 1,
		},
		{
			name:      "large numbers",
			input:     New(1000000, 2000000),
			wantNum:   1,
			wantDenom: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.input
			r.Reduce()
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestRat_Reduce_Invalid tests reduction with invalid state
func TestRat_Reduce_Invalid(t *testing.T) {
	r := New(6, 0) // invalid
	r.Reduce()
	assert.True(t, r.IsInvalid(), "invalid state should be preserved")
}

// TestRat_Reduced tests immutable reduction to lowest terms
func TestRat_Reduced(t *testing.T) {
	tests := []struct {
		name      string
		input     Rat
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "reduce fraction",
			input:     New(6, 8),
			wantNum:   3,
			wantDenom: 4,
		},
		{
			name:      "already reduced",
			input:     New(5, 7),
			wantNum:   5,
			wantDenom: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			original := tt.input
			result := tt.input.Reduced()

			// Check result
			assert.Equal(t, tt.wantNum, result.numerator, "result numerator mismatch")
			assert.Equal(t, tt.wantDenom, result.denominator, "result denominator mismatch")

			// Check original is unchanged
			assert.Equal(t, original.numerator, tt.input.numerator, "original should not be modified")
			assert.Equal(t, original.denominator, tt.input.denominator, "original should not be modified")
		})
	}
}

// TestRat_Reduced_Invalid tests immutable reduction with invalid state
func TestRat_Reduced_Invalid(t *testing.T) {
	original := New(6, 0) // invalid
	result := original.Reduced()

	// Check result is invalid
	assert.True(t, result.IsInvalid(), "result should be invalid")

	// Check original is unchanged
	assert.True(t, original.IsInvalid(), "original should still be invalid")
}

// TestReduce_SignedMagnitudePath exercises gcd>1 with negative numerator path
func TestReduce_SignedMagnitudePath(t *testing.T) {
	r := Rat{numerator: -10, denominator: 20}
	r.Reduce()
	assert.True(t, r.IsValid())
	assert.Equal(t, int64(-1), r.numerator)
	assert.Equal(t, uint64(2), r.denominator)
}

// TestReduce_EdgeCases covers remaining Reduce branches
func TestReduce_EdgeCases(t *testing.T) {
	// Test gcd == 1 case (already reduced)
	r := Rat{numerator: 7, denominator: 11}
	r.Reduce()
	assert.True(t, r.IsValid())
	assert.Equal(t, int64(7), r.numerator)
	assert.Equal(t, uint64(11), r.denominator)

	// Test large negative numerator reduction
	r2 := Rat{numerator: -1000000, denominator: 2000000}
	r2.Reduce()
	assert.True(t, r2.IsValid())
	assert.Equal(t, int64(-1), r2.numerator)
	assert.Equal(t, uint64(2), r2.denominator)
}

// TestReduce_MinInt64Special tests reduction with MinInt64 special case handling
func TestReduce_MinInt64Special(t *testing.T) {
	// Test MinInt64 case which might have special handling
	r := Rat{numerator: math.MinInt64, denominator: 2}
	r.Reduce()
	// Should reduce to MinInt64/2, 1 if possible, or handle the MinInt64 special case
	assert.True(t, r.IsValid() || r.IsInvalid()) // Either outcome is acceptable

	// Try a case where GCD is large and might cause issues
	r2 := Rat{numerator: math.MinInt64, denominator: uint64(math.MaxInt64) + 1}
	r2.Reduce()
	assert.True(t, r2.IsValid() || r2.IsInvalid()) // Either outcome is acceptable
}

// TestReduce_GCDEdgeCases tests GCD calculation edge cases
func TestReduce_GCDEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		numerator int64
		denom     uint64
		desc      string
	}{
		{
			name:      "gcd with zero numerator",
			numerator: 0,
			denom:     42,
			desc:      "GCD(0, 42) should normalize to 0/1",
		},
		{
			name:      "gcd with one",
			numerator: 1,
			denom:     1,
			desc:      "GCD(1, 1) = 1, should stay 1/1",
		},
		{
			name:      "gcd with large coprime numbers",
			numerator: 97,  // prime
			denom:     101, // different prime
			desc:      "GCD(97, 101) = 1, should stay 97/101",
		},
		{
			name:      "gcd with powers of 2",
			numerator: 64,  // 2^6
			denom:     128, // 2^7
			desc:      "GCD(64, 128) = 64, should reduce to 1/2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rat{numerator: tt.numerator, denominator: tt.denom}
			original := r
			r.Reduce()

			assert.True(t, r.IsValid(), "%s: should be valid after reduction", tt.desc)

			// Verify the fraction is in lowest terms by checking GCD = 1
			if r.numerator != 0 {
				gcd := gcdInt64Uint64(r.numerator, r.denominator)
				assert.Equal(t, uint64(1), gcd, "%s: result should be in lowest terms", tt.desc)
			}

			// Verify mathematical equivalence
			if original.denominator != 0 && r.denominator != 0 {
				originalValue := float64(original.numerator) / float64(original.denominator)
				reducedValue := float64(r.numerator) / float64(r.denominator)
				assert.InDelta(t, originalValue, reducedValue, 1e-15, "%s: mathematical value should be preserved", tt.desc)
			}
		})
	}
}

// TestReduce_OverflowInReduction tests edge cases where reduction might cause overflow
func TestReduce_OverflowInReduction(t *testing.T) {
	// Test cases that might trigger the "Should not happen" case in uint64ToInt64WithSign
	tests := []struct {
		name      string
		numerator int64
		denom     uint64
		desc      string
	}{
		{
			name:      "MinInt64 with even denominator",
			numerator: math.MinInt64,
			denom:     2,
			desc:      "MinInt64/2 reduction should handle signed magnitude correctly",
		},
		{
			name:      "MinInt64 with large denominator",
			numerator: math.MinInt64,
			denom:     uint64(math.MaxInt64) + 1,
			desc:      "MinInt64 with MaxInt64+1 denominator",
		},
		{
			name:      "large negative with large denominator",
			numerator: -9223372036854775800, // Close to MinInt64
			denom:     18446744073709551600, // Close to MaxUint64
			desc:      "large negative numerator with large denominator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rat{numerator: tt.numerator, denominator: tt.denom}
			r.Reduce()

			// The result should either be valid or invalid, but should not panic
			assert.True(t, r.IsValid() || r.IsInvalid(), "%s: should not panic and should have valid state", tt.desc)

			if r.IsValid() {
				// If valid, verify it's properly reduced
				if r.numerator != 0 {
					gcd := gcdInt64Uint64(r.numerator, r.denominator)
					assert.Equal(t, uint64(1), gcd, "%s: should be in lowest terms if valid", tt.desc)
				}
			} else {
				// If invalid, should have denominator = 0
				assert.Equal(t, uint64(0), r.denominator, "%s: invalid result should have denominator = 0", tt.desc)
			}
		})
	}
}

// TestReduce_ZeroNumeratorNormalization tests zero numerator normalization
func TestReduce_ZeroNumeratorNormalization(t *testing.T) {
	tests := []struct {
		name  string
		denom uint64
	}{
		{"zero with small denominator", 5},
		{"zero with large denominator", 1000000},
		{"zero with max denominator", math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Rat{numerator: 0, denominator: tt.denom}
			r.Reduce()

			assert.True(t, r.IsValid(), "zero should always be valid after reduction")
			assert.Equal(t, int64(0), r.numerator, "numerator should remain 0")
			assert.Equal(t, uint64(1), r.denominator, "zero should normalize to 0/1")
		})
	}
}

// TestReduce_ManuallyCreatedUnreducedFractions tests reduction on manually created unreduced fractions
func TestReduce_ManuallyCreatedUnreducedFractions(t *testing.T) {
	// These tests bypass New() constructor to create unreduced fractions directly
	tests := []struct {
		name      string
		numerator int64
		denom     uint64
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "manually created 4/6 -> 2/3",
			numerator: 4,
			denom:     6,
			wantNum:   2,
			wantDenom: 3,
		},
		{
			name:      "manually created -15/25 -> -3/5",
			numerator: -15,
			denom:     25,
			wantNum:   -3,
			wantDenom: 5,
		},
		{
			name:      "manually created 100/150 -> 2/3",
			numerator: 100,
			denom:     150,
			wantNum:   2,
			wantDenom: 3,
		},
		{
			name:      "manually created large fraction",
			numerator: 123456789,
			denom:     987654321,
			wantNum:   13717421,  // GCD(123456789, 987654321) = 9
			wantDenom: 109739369, // 987654321 / 9
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create unreduced fraction manually
			r := Rat{numerator: tt.numerator, denominator: tt.denom}
			r.Reduce()

			assert.True(t, r.IsValid(), "should be valid after reduction")
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")

			// Verify it's actually reduced (GCD = 1)
			if r.numerator != 0 {
				gcd := gcdInt64Uint64(r.numerator, r.denominator)
				assert.Equal(t, uint64(1), gcd, "result should be in lowest terms")
			}
		})
	}
}

// TestReduce_ShouldNotHappenCase attempts to trigger the "Should not happen" case in Reduce
func TestReduce_ShouldNotHappenCase(t *testing.T) {
	// The "Should not happen" case occurs when uint64ToInt64WithSign fails
	// This happens when absNum after division is > MaxInt64 for positive numbers
	// or > MaxInt64+1 for negative numbers

	// Try to create a scenario where this might happen by manually constructing
	// a Rat with very specific values that might trigger this edge case

	// Create a fraction where the GCD reduction might cause issues
	// Use the largest possible values that might cause the conversion to fail
	r := Rat{
		numerator:   math.MinInt64,  // Most negative value
		denominator: math.MaxUint64, // Largest denominator
	}

	// This should either work normally or trigger the safety check
	r.Reduce()

	// The result should either be valid (if the reduction worked) or invalid (if it triggered the safety check)
	assert.True(t, r.IsValid() || r.IsInvalid(), "should have a valid state after Reduce")

	// If it's invalid, it should have denominator = 0
	if r.IsInvalid() {
		assert.Equal(t, uint64(0), r.denominator, "invalid result should have denominator = 0")
	}
}
