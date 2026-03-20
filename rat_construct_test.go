package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew_ValidInputs tests creation of valid rational numbers
func TestNew_ValidInputs(t *testing.T) {
	tests := []struct {
		name        string
		numerator   int64
		denominator uint64
		wantNum     int64
		wantDenom   uint64
	}{
		{
			name:        "positive fraction",
			numerator:   3,
			denominator: 4,
			wantNum:     3,
			wantDenom:   4,
		},
		{
			name:        "negative fraction",
			numerator:   -5,
			denominator: 7,
			wantNum:     -5,
			wantDenom:   7,
		},
		{
			name:        "positive integer",
			numerator:   42,
			denominator: 1,
			wantNum:     42,
			wantDenom:   1,
		},
		{
			name:        "negative integer",
			numerator:   -42,
			denominator: 1,
			wantNum:     -42,
			wantDenom:   1,
		},
		{
			name:        "zero numerator",
			numerator:   0,
			denominator: 5,
			wantNum:     0,
			wantDenom:   1, // should normalize to 0/1
		},
		{
			name:        "large values",
			numerator:   9223372036854775807,  // MaxInt64
			denominator: 18446744073709551615, // MaxUint64
			wantNum:     9223372036854775807,
			wantDenom:   18446744073709551615,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.numerator, tt.denominator)
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestNew_InvalidInputs tests creation of invalid rational numbers
func TestNew_InvalidInputs(t *testing.T) {
	tests := []struct {
		name        string
		numerator   int64
		denominator uint64
	}{
		{
			name:        "zero denominator positive numerator",
			numerator:   5,
			denominator: 0,
		},
		{
			name:        "zero denominator negative numerator",
			numerator:   -3,
			denominator: 0,
		},
		{
			name:        "zero denominator zero numerator",
			numerator:   0,
			denominator: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.numerator, tt.denominator)
			assert.Equal(t, uint64(0), r.denominator, "should be invalid (denominator = 0)")
		})
	}
}

// TestNew_SignNormalization tests sign normalization
func TestNew_SignNormalization(t *testing.T) {
	// Note: in Go uint64 is always positive, so sign is always in numerator
	// This test verifies sign is handled correctly
	tests := []struct {
		name        string
		numerator   int64
		denominator uint64
		wantNum     int64
		wantDenom   uint64
	}{
		{
			name:        "positive/positive",
			numerator:   3,
			denominator: 4,
			wantNum:     3,
			wantDenom:   4,
		},
		{
			name:        "negative/positive",
			numerator:   -3,
			denominator: 4,
			wantNum:     -3,
			wantDenom:   4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.numerator, tt.denominator)
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestNew_AutomaticReduction tests that New automatically reduces fractions
func TestNew_AutomaticReduction(t *testing.T) {
	tests := []struct {
		name      string
		num       int64
		denom     uint64
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "reduce 6/9 to 2/3",
			num:       6,
			denom:     9,
			wantNum:   2,
			wantDenom: 3,
		},
		{
			name:      "reduce 10/15 to 2/3",
			num:       10,
			denom:     15,
			wantNum:   2,
			wantDenom: 3,
		},
		{
			name:      "reduce -12/18 to -2/3",
			num:       -12,
			denom:     18,
			wantNum:   -2,
			wantDenom: 3,
		},
		{
			name:      "already reduced 3/4 stays 3/4",
			num:       3,
			denom:     4,
			wantNum:   3,
			wantDenom: 4,
		},
		{
			name:      "reduce 100/200 to 1/2",
			num:       100,
			denom:     200,
			wantNum:   1,
			wantDenom: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rat := New(tt.num, tt.denom)
			assert.Equal(t, tt.wantNum, rat.Numerator(), "Numerator should be reduced")
			assert.Equal(t, tt.wantDenom, rat.Denominator(), "Denominator should be reduced")
		})
	}
}

// TestNewFromInt tests creation of rational number from integer
func TestNewFromInt(t *testing.T) {
	tests := []struct {
		name      string
		value     int64
		wantNum   int64
		wantDenom uint64
	}{
		{
			name:      "positive integer",
			value:     42,
			wantNum:   42,
			wantDenom: 1,
		},
		{
			name:      "negative integer",
			value:     -17,
			wantNum:   -17,
			wantDenom: 1,
		},
		{
			name:      "zero",
			value:     0,
			wantNum:   0,
			wantDenom: 1,
		},
		{
			name:      "max int64",
			value:     9223372036854775807,
			wantNum:   9223372036854775807,
			wantDenom: 1,
		},
		{
			name:      "min int64",
			value:     -9223372036854775808,
			wantNum:   -9223372036854775808,
			wantDenom: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewFromInt64(tt.value)
			assert.Equal(t, tt.wantNum, r.numerator, "numerator mismatch")
			assert.Equal(t, tt.wantDenom, r.denominator, "denominator mismatch")
		})
	}
}

// TestZero tests creation of zero rational number
func TestZero(t *testing.T) {
	r := Zero()
	assert.Equal(t, int64(0), r.numerator, "zero numerator expected")
	assert.Equal(t, uint64(1), r.denominator, "denominator should be 1")
}

// TestOne tests creation of one rational number
func TestOne(t *testing.T) {
	r := One()
	assert.Equal(t, int64(1), r.numerator, "numerator should be 1")
	assert.Equal(t, uint64(1), r.denominator, "denominator should be 1")
}

// TestRat_FieldAccess tests struct field access
func TestRat_FieldAccess(t *testing.T) {
	r := New(3, 4)

	// Verify fields are accessible (this compiles)
	assert.Equal(t, int64(3), r.numerator)
	assert.Equal(t, uint64(4), r.denominator)
}

// TestRat_IsValid tests validity check of rational number
func TestRat_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected bool
	}{
		{
			name:     "valid positive fraction",
			rat:      New(3, 4),
			expected: true,
		},
		{
			name:     "valid negative fraction",
			rat:      New(-5, 7),
			expected: true,
		},
		{
			name:     "valid zero",
			rat:      New(0, 1),
			expected: true,
		},
		{
			name:     "valid integer",
			rat:      New(42, 1),
			expected: true,
		},
		{
			name:     "invalid zero denominator",
			rat:      New(5, 0),
			expected: false,
		},
		{
			name:     "invalid zero denominator with zero numerator",
			rat:      New(0, 0),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.IsValid()
			assert.Equal(t, tt.expected, result, "IsValid() result mismatch")
		})
	}
}

// TestRat_IsInvalid tests invalidity check of rational number
func TestRat_IsInvalid(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected bool
	}{
		{
			name:     "valid positive fraction",
			rat:      New(3, 4),
			expected: false,
		},
		{
			name:     "valid negative fraction",
			rat:      New(-5, 7),
			expected: false,
		},
		{
			name:     "valid zero",
			rat:      New(0, 1),
			expected: false,
		},
		{
			name:     "invalid zero denominator",
			rat:      New(5, 0),
			expected: true,
		},
		{
			name:     "invalid zero denominator with zero numerator",
			rat:      New(0, 0),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.IsInvalid()
			assert.Equal(t, tt.expected, result, "IsInvalid() result mismatch")
		})
	}
}

// TestRat_Invalidate tests forced invalidation
func TestRat_Invalidate(t *testing.T) {
	r := New(3, 4)
	assert.True(t, r.IsValid(), "should be valid initially")

	r.Invalidate()
	assert.True(t, r.IsInvalid(), "should be invalid after Invalidate()")
	assert.False(t, r.IsValid(), "should not be valid after Invalidate()")
}

// TestRat_UtilityMethods tests utility methods
func TestRat_UtilityMethods(t *testing.T) {
	// Numerator, Denominator
	r := New(3, 4)
	assert.Equal(t, int64(3), r.Numerator())
	assert.Equal(t, uint64(4), r.Denominator())

	// Sign
	assert.Equal(t, 1, New(3, 4).Sign())
	assert.Equal(t, -1, New(-3, 4).Sign())
	assert.Equal(t, 0, New(0, 1).Sign())
	assert.Equal(t, 0, New(1, 0).Sign()) // invalid

	// IsZero, IsOne
	assert.True(t, New(0, 1).IsZero())
	assert.False(t, New(1, 2).IsZero())
	assert.True(t, New(1, 1).IsOne())
	assert.False(t, New(2, 1).IsOne())

	// HasFractional
	assert.False(t, New(2, 1).HasFractional()) // whole number
	assert.True(t, New(1, 2).HasFractional())  // fractional
	assert.False(t, New(0, 1).HasFractional()) // zero

	// IntegerAndFraction
	intPart, fracPart := New(7, 3).IntegerAndFraction()
	assert.Equal(t, int64(2), intPart) // 7/3 = 2 + 1/3
	assert.Equal(t, int64(1), fracPart.numerator)
	assert.Equal(t, uint64(3), fracPart.denominator)
}

// TestRat_HasFractional tests the HasFractional method
func TestRat_HasFractional(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected bool
		desc     string
	}{
		// Whole numbers (should return false)
		{"whole positive", New(5, 1), false, "5/1 is a whole number"},
		{"whole negative", New(-5, 1), false, "-5/1 is a whole number"},
		{"zero", New(0, 1), false, "0/1 is a whole number"},
		{"large whole", New(1000000, 1), false, "1000000/1 is a whole number"},
		{"reduced whole", New(10, 5), false, "10/5 = 2/1 is a whole number"},
		{"negative reduced whole", New(-15, 3), false, "-15/3 = -5/1 is a whole number"},

		// Fractional numbers (should return true)
		{"simple fraction", New(1, 2), true, "1/2 has fractional part"},
		{"negative fraction", New(-1, 2), true, "-1/2 has fractional part"},
		{"complex fraction", New(7, 3), true, "7/3 has fractional part"},
		{"large fraction", New(1000001, 1000000), true, "1000001/1000000 has fractional part"},
		{"reduced fraction", New(3, 6), true, "3/6 = 1/2 has fractional part"},
		{"negative reduced fraction", New(-7, 14), true, "-7/14 = -1/2 has fractional part"},

		// Edge cases
		{"invalid state", Rat{numerator: 1, denominator: 0}, false, "invalid rational should return false"},
		{"zero invalid", Rat{numerator: 0, denominator: 0}, false, "invalid zero should return false"},
		{"max int numerator whole", New(9223372036854775807, 1), false, "MaxInt64/1 is whole"},
		{"min int numerator whole", New(-9223372036854775808, 1), false, "MinInt64/1 is whole"},
		{"max denominator fraction", New(1, 18446744073709551615), true, "1/MaxUint64 has fractional part"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.HasFractional()
			assert.Equal(t, tt.expected, result, tt.desc)
		})
	}
}

// TestRat_IntegerAndFraction tests the IntegerAndFraction method
func TestRat_IntegerAndFraction(t *testing.T) {
	tests := []struct {
		name         string
		rat          Rat
		expectedInt  int64
		expectedFrac Rat
		desc         string
	}{
		// Whole numbers (fractional part should be 0/1)
		{"positive whole", New(5, 1), 5, New(0, 1), "5/1 = 5 + 0/1"},
		{"negative whole", New(-5, 1), -5, New(0, 1), "-5/1 = -5 + 0/1"},
		{"zero", New(0, 1), 0, New(0, 1), "0/1 = 0 + 0/1"},
		{"reduced whole", New(10, 5), 2, New(0, 1), "10/5 = 2 + 0/1"},
		{"negative reduced whole", New(-15, 3), -5, New(0, 1), "-15/3 = -5 + 0/1"},

		// Mixed numbers (proper integer and fractional parts)
		{"simple mixed positive", New(7, 3), 2, New(1, 3), "7/3 = 2 + 1/3"},
		{"simple mixed negative", New(-7, 3), -2, New(-1, 3), "-7/3 = -2 + (-1/3)"},
		{"complex mixed", New(22, 7), 3, New(1, 7), "22/7 = 3 + 1/7"},
		{"negative complex mixed", New(-22, 7), -3, New(-1, 7), "-22/7 = -3 + (-1/7)"},
		{"large mixed", New(1000001, 1000000), 1, New(1, 1000000), "1000001/1000000 = 1 + 1/1000000"},

		// Proper fractions (integer part should be 0)
		{"proper fraction positive", New(1, 2), 0, New(1, 2), "1/2 = 0 + 1/2"},
		{"proper fraction negative", New(-1, 2), 0, New(-1, 2), "-1/2 = 0 + (-1/2)"},
		{"small proper fraction", New(3, 7), 0, New(3, 7), "3/7 = 0 + 3/7"},
		{"negative small proper", New(-3, 7), 0, New(-3, 7), "-3/7 = 0 + (-3/7)"},

		// Edge cases
		{"invalid state", Rat{numerator: 1, denominator: 0}, 0, Rat{numerator: 0, denominator: 0}, "invalid rational should return 0 and invalid fraction"},
		{"max int numerator", New(9223372036854775807, 1), 9223372036854775807, New(0, 1), "MaxInt64/1 = MaxInt64 + 0/1"},
		{"min int numerator", New(-9223372036854775808, 1), -9223372036854775808, New(0, 1), "MinInt64/1 = MinInt64 + 0/1"},
		{"large denominator", New(1, 18446744073709551615), 0, New(1, 18446744073709551615), "1/MaxUint64 = 0 + 1/MaxUint64"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			intPart, fracPart := tt.rat.IntegerAndFraction()

			assert.Equal(t, tt.expectedInt, intPart, "integer part mismatch: %s", tt.desc)
			assert.Equal(t, tt.expectedFrac.numerator, fracPart.numerator, "fractional numerator mismatch: %s", tt.desc)
			assert.Equal(t, tt.expectedFrac.denominator, fracPart.denominator, "fractional denominator mismatch: %s", tt.desc)
		})
	}
}

// TestNewFromFloat64 tests exact float64 construction.
func TestNewFromFloat64(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  Rat
	}{
		{name: "zero", value: 0.0, want: Zero()},
		{name: "negative zero", value: math.Copysign(0.0, -1), want: Zero()},
		{name: "positive integer", value: 42.0, want: NewFromInt64(42)},
		{name: "negative integer", value: -17.0, want: NewFromInt64(-17)},
		{name: "one half", value: 0.5, want: New(1, 2)},
		{name: "three quarters", value: 0.75, want: New(3, 4)},
		{name: "one tenth exact binary", value: 0.1, want: New(3602879701896397, 36028797018963968)},
		{name: "negative mixed", value: -3.75, want: New(-15, 4)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := NewFromFloat64(tt.value)
			require.NoError(t, err)
			assert.True(t, r.Equal(tt.want))
		})
	}
}

// TestNewFromFloat64_SpecialValues tests special float64 values.
func TestNewFromFloat64_SpecialValues(t *testing.T) {
	_, err := NewFromFloat64(math.Inf(1))
	require.ErrorIs(t, err, ErrNonFiniteFloat)

	_, err = NewFromFloat64(math.Inf(-1))
	require.ErrorIs(t, err, ErrNonFiniteFloat)

	_, err = NewFromFloat64(math.NaN())
	require.ErrorIs(t, err, ErrNonFiniteFloat)
}

// TestNewFromFloat64_EdgeCases tests exact edge cases and helper wrappers.
func TestNewFromFloat64_EdgeCases(t *testing.T) {
	_, err := NewFromFloat64(math.MaxFloat64)
	require.ErrorIs(t, err, ErrNotRepresentable)

	_, err = NewFromFloat64(math.SmallestNonzeroFloat64)
	require.ErrorIs(t, err, ErrNotRepresentable)

	_, err = NewFromFloat64(math.Ldexp(1, -64))
	require.ErrorIs(t, err, ErrNotRepresentable)

	_, err = NewFromFloat64(3 * math.Ldexp(1, -64))
	require.ErrorIs(t, err, ErrNotRepresentable)

	v := 0.25
	r, err := NewFromFloat64Ptr(&v)
	require.NoError(t, err)
	assert.True(t, r.Equal(New(1, 4)))

	values, err := FromFloat64Slice([]float64{0.5, 0.25})
	require.NoError(t, err)
	require.Len(t, values, 2)
	assert.True(t, values[0].Equal(New(1, 2)))
	assert.True(t, values[1].Equal(New(1, 4)))
}

// TestNewFromFloat32 tests float32 exact and approximate wrappers.
func TestNewFromFloat32(t *testing.T) {
	r, err := NewFromFloat32(0.5)
	require.NoError(t, err)
	assert.True(t, r.Equal(New(1, 2)))

	_, err = NewFromFloat32(math.SmallestNonzeroFloat32)
	require.ErrorIs(t, err, ErrNotRepresentable)

	r, err = NewApproxFromFloat32(math.SmallestNonzeroFloat32)
	require.NoError(t, err)
	assert.True(t, r.Equal(Zero()))

	v := float32(0.25)
	r, err = NewFromFloat32Ptr(&v)
	require.NoError(t, err)
	assert.True(t, r.Equal(New(1, 4)))

	r, err = NewFromFloat32Ptr(nil)
	require.NoError(t, err)
	assert.True(t, r.IsInvalid())
}

// TestFromFloat64SlicePtr tests pointer slice conversion.
func TestFromFloat64SlicePtr(t *testing.T) {
	a := 0.5
	b := 0.25
	values, err := FromFloat64SlicePtr([]*float64{&a, nil, &b})
	require.NoError(t, err)
	require.Len(t, values, 3)
	assert.True(t, values[0].Equal(New(1, 2)))
	assert.True(t, values[1].IsInvalid())
	assert.True(t, values[2].Equal(New(1, 4)))

	nan := math.NaN()
	_, err = FromFloat64SlicePtr([]*float64{&a, &nan})
	require.ErrorIs(t, err, ErrNonFiniteFloat)
}

// TestRat_ToInt64Err tests the ToInt64Err method
func TestRat_ToInt64Err(t *testing.T) {
	tests := []struct {
		name        string
		rat         Rat
		expected    int64
		expectError bool
	}{
		{
			name:        "valid positive integer",
			rat:         New(5, 1),
			expected:    5,
			expectError: false,
		},
		{
			name:        "valid negative integer",
			rat:         New(-3, 1),
			expected:    -3,
			expectError: false,
		},
		{
			name:        "valid fraction",
			rat:         New(7, 3),
			expected:    2,
			expectError: false,
		},
		{
			name:        "valid negative fraction",
			rat:         New(-7, 3),
			expected:    -2,
			expectError: false,
		},
		{
			name:        "zero",
			rat:         New(0, 1),
			expected:    0,
			expectError: false,
		},
		{
			name:        "invalid rat",
			rat:         Rat{numerator: 1, denominator: 0},
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.rat.ToInt64Err()

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, ErrInvalid, err)
				assert.Equal(t, int64(0), result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestRat_ToInt64 tests the ToInt64 method
func TestRat_ToInt64(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected int64
	}{
		{
			name:     "valid positive integer",
			rat:      New(5, 1),
			expected: 5,
		},
		{
			name:     "valid negative integer",
			rat:      New(-3, 1),
			expected: -3,
		},
		{
			name:     "valid fraction",
			rat:      New(7, 3),
			expected: 2,
		},
		{
			name:     "valid negative fraction",
			rat:      New(-7, 3),
			expected: -2,
		},
		{
			name:     "zero",
			rat:      New(0, 1),
			expected: 0,
		},
		{
			name:     "invalid rat",
			rat:      Rat{numerator: 1, denominator: 0},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.ToInt64()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRat_ToIntErr tests the ToIntErr method
func TestRat_ToIntErr(t *testing.T) {
	tests := []struct {
		name        string
		rat         Rat
		expected    int
		expectError bool
	}{
		{
			name:        "valid positive integer",
			rat:         New(5, 1),
			expected:    5,
			expectError: false,
		},
		{
			name:        "valid negative integer",
			rat:         New(-3, 1),
			expected:    -3,
			expectError: false,
		},
		{
			name:        "valid fraction",
			rat:         New(7, 3),
			expected:    2,
			expectError: false,
		},
		{
			name:        "valid negative fraction",
			rat:         New(-7, 3),
			expected:    -2,
			expectError: false,
		},
		{
			name:        "zero",
			rat:         New(0, 1),
			expected:    0,
			expectError: false,
		},
		{
			name:        "invalid rat",
			rat:         Rat{numerator: 1, denominator: 0},
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.rat.ToIntErr()

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, ErrInvalid, err)
				assert.Equal(t, 0, result)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// TestRat_ToInt tests the ToInt method
func TestRat_ToInt(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected int
	}{
		{
			name:     "valid positive integer",
			rat:      New(5, 1),
			expected: 5,
		},
		{
			name:     "valid negative integer",
			rat:      New(-3, 1),
			expected: -3,
		},
		{
			name:     "valid fraction",
			rat:      New(7, 3),
			expected: 2,
		},
		{
			name:     "valid negative fraction",
			rat:      New(-7, 3),
			expected: -2,
		},
		{
			name:     "zero",
			rat:      New(0, 1),
			expected: 0,
		},
		{
			name:     "invalid rat",
			rat:      Rat{numerator: 1, denominator: 0},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.ToInt()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRat_ToFloatErr tests the ToFloatErr method
func TestRat_ToFloatErr(t *testing.T) {
	tests := []struct {
		name        string
		rat         Rat
		expected    float64
		expectError bool
	}{
		{
			name:        "valid positive integer",
			rat:         New(5, 1),
			expected:    5.0,
			expectError: false,
		},
		{
			name:        "valid negative integer",
			rat:         New(-3, 1),
			expected:    -3.0,
			expectError: false,
		},
		{
			name:        "valid fraction",
			rat:         New(1, 2),
			expected:    0.5,
			expectError: false,
		},
		{
			name:        "valid negative fraction",
			rat:         New(-3, 4),
			expected:    -0.75,
			expectError: false,
		},
		{
			name:        "zero",
			rat:         New(0, 1),
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "invalid rat",
			rat:         Rat{numerator: 1, denominator: 0},
			expected:    0.0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.rat.ToFloat64Err()

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, ErrInvalid, err)
				assert.InDelta(t, 0.0, result, 1e-15)
			} else {
				require.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-15)
			}
		})
	}
}

// TestRat_ToFloat tests the ToFloat method
func TestRat_ToFloat(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected float64
	}{
		{
			name:     "valid positive integer",
			rat:      New(5, 1),
			expected: 5.0,
		},
		{
			name:     "valid negative integer",
			rat:      New(-3, 1),
			expected: -3.0,
		},
		{
			name:     "valid fraction",
			rat:      New(1, 2),
			expected: 0.5,
		},
		{
			name:     "valid negative fraction",
			rat:      New(-3, 4),
			expected: -0.75,
		},
		{
			name:     "zero",
			rat:      New(0, 1),
			expected: 0.0,
		},
		{
			name:     "invalid rat",
			rat:      Rat{numerator: 1, denominator: 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.ToFloat64()
			assert.InDelta(t, tt.expected, result, 1e-15)
		})
	}
}

// TestRat_ToFloat32Err tests the ToFloat32Err method
func TestRat_ToFloat32Err(t *testing.T) {
	tests := []struct {
		name        string
		rat         Rat
		expected    float32
		expectError bool
	}{
		{
			name:        "valid positive integer",
			rat:         New(5, 1),
			expected:    5.0,
			expectError: false,
		},
		{
			name:        "valid negative integer",
			rat:         New(-3, 1),
			expected:    -3.0,
			expectError: false,
		},
		{
			name:        "valid fraction",
			rat:         New(1, 2),
			expected:    0.5,
			expectError: false,
		},
		{
			name:        "valid negative fraction",
			rat:         New(-3, 4),
			expected:    -0.75,
			expectError: false,
		},
		{
			name:        "zero",
			rat:         New(0, 1),
			expected:    0.0,
			expectError: false,
		},
		{
			name:        "invalid rat",
			rat:         Rat{numerator: 1, denominator: 0},
			expected:    0.0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.rat.ToFloat32Err()

			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, ErrInvalid, err)
				assert.InDelta(t, float32(0.0), result, 1e-7)
			} else {
				require.NoError(t, err)
				assert.InDelta(t, tt.expected, result, 1e-7)
			}
		})
	}
}

// TestRat_ToFloat32Err_Overflow tests overflow scenarios for ToFloat32Err
func TestRat_ToFloat32Err_Overflow(t *testing.T) {
	// Test with a value that should overflow float32
	// Create a rational number that represents a value larger than MaxFloat32
	largeRat := New(1, 1)
	largeRat.numerator = 1000000000000000000 // 1e18
	largeRat.denominator = 1

	// Convert to float64 first to check the actual value
	result64 := float64(largeRat.numerator) / float64(largeRat.denominator)

	// If the value is within float32 range, the conversion should succeed
	// If it's outside the range, it should return an error
	result32, err := largeRat.ToFloat32Err()

	// For this specific large value, it should be within float32 range
	// so we expect no error
	if result64 <= math.MaxFloat32 && result64 >= -math.MaxFloat32 {
		require.NoError(t, err)
		assert.InDelta(t, float32(result64), result32, 1e-7)
	}
}

// TestRat_ToFloat32 tests the ToFloat32 method
func TestRat_ToFloat32(t *testing.T) {
	tests := []struct {
		name     string
		rat      Rat
		expected float32
	}{
		{
			name:     "valid positive integer",
			rat:      New(5, 1),
			expected: 5.0,
		},
		{
			name:     "valid negative integer",
			rat:      New(-3, 1),
			expected: -3.0,
		},
		{
			name:     "valid fraction",
			rat:      New(1, 2),
			expected: 0.5,
		},
		{
			name:     "valid negative fraction",
			rat:      New(-3, 4),
			expected: -0.75,
		},
		{
			name:     "zero",
			rat:      New(0, 1),
			expected: 0.0,
		},
		{
			name:     "invalid rat",
			rat:      Rat{numerator: 1, denominator: 0},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rat.ToFloat32()
			assert.InDelta(t, tt.expected, result, 1e-7)
		})
	}
}
