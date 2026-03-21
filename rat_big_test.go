package zerorat

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFromBigRat tests exact conversion from big.Rat into Rat.
func TestNewFromBigRat(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input *big.Rat
		want  Rat
	}{
		{
			name:  "simple fraction",
			input: big.NewRat(7, 20),
			want:  New(7, 20),
		},
		{
			name:  "negative fraction",
			input: big.NewRat(-7, 20),
			want:  New(-7, 20),
		},
		{
			name:  "integer",
			input: big.NewRat(42, 1),
			want:  NewFromInt64(42),
		},
		{
			name:  "min int64 numerator",
			input: big.NewRat(math.MinInt64, 1),
			want:  NewFromInt64(math.MinInt64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewFromBigRat(tt.input)
			require.NoError(t, err)
			assert.True(t, got.Equal(tt.want), "converted value mismatch")
		})
	}
}

// TestNewFromBigRat_InvalidInputs tests rejected big.Rat inputs.
func TestNewFromBigRat_InvalidInputs(t *testing.T) {
	t.Parallel()

	tooLargeNumerator := new(big.Int).Add(big.NewInt(math.MaxInt64), big.NewInt(1))
	tooLargeDenominator := new(big.Int).Add(
		new(big.Int).SetUint64(math.MaxUint64),
		big.NewInt(1),
	)

	tests := []struct {
		name      string
		input     *big.Rat
		wantErr   error
		checkText string
	}{
		{
			name:      "nil input",
			input:     nil,
			checkText: "nil big.Rat",
		},
		{
			name:    "numerator overflow",
			input:   new(big.Rat).SetFrac(tooLargeNumerator, big.NewInt(1)),
			wantErr: ErrNotRepresentable,
		},
		{
			name:    "denominator overflow",
			input:   new(big.Rat).SetFrac(big.NewInt(1), tooLargeDenominator),
			wantErr: ErrNotRepresentable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewFromBigRat(tt.input)
			require.Error(t, err)
			assert.True(t, got.IsInvalid(), "invalid input should return invalid Rat")

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}

			if tt.checkText != "" {
				assert.ErrorContains(t, err, tt.checkText)
			}
		})
	}
}

// TestRat_ToBigRatErr tests exact conversion from Rat into big.Rat.
func TestRat_ToBigRatErr(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rat  Rat
		want *big.Rat
	}{
		{
			name: "simple fraction",
			rat:  New(7, 20),
			want: big.NewRat(7, 20),
		},
		{
			name: "negative fraction",
			rat:  New(-7, 20),
			want: big.NewRat(-7, 20),
		},
		{
			name: "integer",
			rat:  NewFromInt64(42),
			want: big.NewRat(42, 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.rat.ToBigRatErr()
			require.NoError(t, err)
			assert.Equal(t, 0, got.Cmp(tt.want))
		})
	}
}

// TestRat_ToBigRatErr_InvalidRat tests rejected Rat inputs for big.Rat conversion.
func TestRat_ToBigRatErr_InvalidRat(t *testing.T) {
	t.Parallel()

	got, err := Rat{}.ToBigRatErr()
	require.ErrorIs(t, err, ErrInvalid)
	assert.Nil(t, got)
}

// TestRat_BigRatRoundTrip tests exact round-trip conversion through big.Rat.
func TestRat_BigRatRoundTrip(t *testing.T) {
	t.Parallel()

	original := New(7, 20)
	asBig, err := original.ToBigRatErr()
	require.NoError(t, err)

	roundTrip, err := NewFromBigRat(asBig)
	require.NoError(t, err)
	assert.True(t, roundTrip.Equal(original))
}
