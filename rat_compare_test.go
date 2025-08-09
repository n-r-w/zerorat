package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestRat_Compare tests three-way comparison
func TestRat_Compare(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected int
	}{
		{"1/2 vs 3/4", New(1, 2), New(3, 4), -1},
		{"3/4 vs 1/2", New(3, 4), New(1, 2), 1},
		{"equal", New(1, 2), New(2, 4), 0},
		{"invalid", New(1, 0), New(1, 2), 0}, // invalid numbers return 0
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Compare(tt.b))
		})
	}
}

// TestRat_Compare_ZeroComparison tests comparison when both numerators are zero
func TestRat_Compare_ZeroComparison(t *testing.T) {
	// This tests the specific line: if r.numerator == 0 && other.numerator == 0
	a := New(0, 3) // 0/3
	b := New(0, 7) // 0/7

	result := a.Compare(b)
	assert.Equal(t, 0, result, "0/3 should equal 0/7")
}

// TestRat_Equal tests equality comparison
func TestRat_Equal(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected bool
	}{
		{"equal fractions", New(1, 2), New(2, 4), true},
		{"different fractions", New(1, 2), New(3, 4), false},
		{"equal integers", New(5, 1), New(5, 1), true},
		{"different integers", New(5, 1), New(7, 1), false},
		{"zero comparisons", New(0, 1), New(0, 5), true},
		{"invalid vs valid", New(1, 0), New(1, 2), false},
		{"both invalid", New(1, 0), New(2, 0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Equal(tt.b))
		})
	}
}

// TestRat_Less tests less-than comparison
func TestRat_Less(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected bool
	}{
		{"1/2 < 3/4", New(1, 2), New(3, 4), true},
		{"3/4 < 1/2", New(3, 4), New(1, 2), false},
		{"equal", New(1, 2), New(2, 4), false},
		{"negative vs positive", New(-1, 2), New(1, 2), true},
		{"both negative", New(-3, 4), New(-1, 2), true},
		{"invalid", New(1, 0), New(1, 2), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Less(tt.b))
		})
	}
}

// TestRat_Greater tests greater-than comparison
func TestRat_Greater(t *testing.T) {
	cases := []struct {
		name     string
		a, b     Rat
		expected bool
	}{
		{"3/4 > 1/2", New(3, 4), New(1, 2), true},
		{"1/2 > 3/4", New(1, 2), New(3, 4), false},
		{"equal", New(1, 2), New(2, 4), false},
		{"invalid", New(1, 0), New(1, 2), false},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.a.Greater(tt.b))
		})
	}
}

// TestCompareRationalsCrossMul_SignHandling tests cross multiplication with different sign combinations
func TestCompareRationalsCrossMul_SignHandling(t *testing.T) {
	tests := []struct {
		name     string
		aNum     int64
		aDenom   uint64
		cNum     int64
		cDenom   uint64
		expected int
	}{
		{
			name:     "negative vs positive",
			aNum:     -3,
			aDenom:   4,
			cNum:     2,
			cDenom:   5,
			expected: -1, // -3/4 < 2/5
		},
		{
			name:     "positive vs negative",
			aNum:     3,
			aDenom:   4,
			cNum:     -2,
			cDenom:   5,
			expected: 1, // 3/4 > -2/5
		},
		{
			name:     "both negative with MinInt64",
			aNum:     math.MinInt64,
			aDenom:   1,
			cNum:     -1,
			cDenom:   1,
			expected: -1, // MinInt64/1 < -1/1
		},
		{
			name:     "MinInt64 special case in cNum",
			aNum:     -1,
			aDenom:   1,
			cNum:     math.MinInt64,
			cDenom:   1,
			expected: 1, // -1/1 > MinInt64/1
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareRationalsCrossMul(tt.aNum, tt.aDenom, tt.cNum, tt.cDenom)
			assert.Equal(t, tt.expected, result, "comparison result mismatch")
		})
	}
}

// TestRat_Compare_InvalidHandling tests comparison with invalid operands
func TestRat_Compare_InvalidHandling(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected int
	}{
		{"invalid vs valid", New(1, 0), New(1, 2), 0},
		{"valid vs invalid", New(1, 2), New(1, 0), 0},
		{"both invalid", New(1, 0), New(2, 0), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			assert.Equal(t, tt.expected, result, "invalid comparisons should return 0")
		})
	}
}

// TestRat_Compare_OverflowSafeComparison tests overflow-safe cross multiplication
func TestRat_Compare_OverflowSafeComparison(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected int
		desc     string
	}{
		{
			name:     "large values that would overflow naive multiplication",
			a:        New(math.MaxInt64/2, math.MaxUint64/2),
			b:        New(math.MaxInt64/3, math.MaxUint64/3),
			expected: 1, // (MaxInt64/2)/(MaxUint64/2) > (MaxInt64/3)/(MaxUint64/3) since 1/2 > 1/3
			desc:     "overflow-safe cross multiplication",
		},
		{
			name:     "mixed signs with large values",
			a:        New(-math.MaxInt64/2, math.MaxUint64/2),
			b:        New(math.MaxInt64/3, math.MaxUint64/3),
			expected: -1, // negative < positive
			desc:     "mixed signs with overflow-safe comparison",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			assert.Equal(t, tt.expected, result, "%s: %s", tt.name, tt.desc)
		})
	}
}

// TestRat_Compare_ZeroHandling tests comparison involving zero values
func TestRat_Compare_ZeroHandling(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected int
	}{
		{"zero vs positive", New(0, 1), New(1, 2), -1},
		{"positive vs zero", New(1, 2), New(0, 1), 1},
		{"zero vs negative", New(0, 1), New(-1, 2), 1},
		{"negative vs zero", New(-1, 2), New(0, 1), -1},
		{"both zero different denominators", New(0, 3), New(0, 7), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			assert.Equal(t, tt.expected, result, "comparison result mismatch")
		})
	}
}

// TestRat_Compare_SameSignComparison tests comparison when both operands have same sign
func TestRat_Compare_SameSignComparison(t *testing.T) {
	tests := []struct {
		name     string
		a, b     Rat
		expected int
		desc     string
	}{
		{
			name:     "both positive, a < b",
			a:        New(1, 3),
			b:        New(1, 2),
			expected: -1,
			desc:     "1/3 < 1/2",
		},
		{
			name:     "both positive, a > b",
			a:        New(2, 3),
			b:        New(1, 2),
			expected: 1,
			desc:     "2/3 > 1/2",
		},
		{
			name:     "both negative, a < b (more negative)",
			a:        New(-2, 3),
			b:        New(-1, 2),
			expected: -1,
			desc:     "-2/3 < -1/2",
		},
		{
			name:     "both negative, a > b (less negative)",
			a:        New(-1, 3),
			b:        New(-1, 2),
			expected: 1,
			desc:     "-1/3 > -1/2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.a.Compare(tt.b)
			assert.Equal(t, tt.expected, result, "%s: %s", tt.name, tt.desc)
		})
	}
}
