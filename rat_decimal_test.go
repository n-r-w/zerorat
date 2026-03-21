package zerorat

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFromDecimalString tests exact decimal parsing, including scientific notation.
func TestNewFromDecimalString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  Rat
	}{
		{
			name:  "zero",
			input: "0",
			want:  Zero(),
		},
		{
			name:  "negative zero with fraction",
			input: "-0.0",
			want:  Zero(),
		},
		{
			name:  "integer",
			input: "12",
			want:  NewFromInt64(12),
		},
		{
			name:  "signed integer",
			input: "+12",
			want:  NewFromInt64(12),
		},
		{
			name:  "fractional",
			input: "12.5",
			want:  New(25, 2),
		},
		{
			name:  "leading zero fraction",
			input: "0.35",
			want:  New(7, 20),
		},
		{
			name:  "trailing zeros are reduced",
			input: "1.2300",
			want:  New(123, 100),
		},
		{
			name:  "fraction without integer part",
			input: ".5",
			want:  New(1, 2),
		},
		{
			name:  "fraction without fractional part",
			input: "5.",
			want:  NewFromInt64(5),
		},
		{
			name:  "scientific notation positive exponent",
			input: "1.25e2",
			want:  NewFromInt64(125),
		},
		{
			name:  "scientific notation negative exponent",
			input: "3.5e-1",
			want:  New(7, 20),
		},
		{
			name:  "scientific notation uppercase exponent",
			input: "-7.5E+1",
			want:  NewFromInt64(-75),
		},
		{
			name:  "min int64",
			input: "-9223372036854775808",
			want:  NewFromInt64(math.MinInt64),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewFromDecimalString(tt.input)
			require.NoError(t, err)
			assert.True(t, got.Equal(tt.want), "parsed value mismatch")
		})
	}
}

// TestNewFromDecimalString_InvalidInputs tests malformed and non-representable decimal strings.
func TestNewFromDecimalString_InvalidInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: ErrInvalidDecimalString,
		},
		{
			name:    "sign only",
			input:   "+",
			wantErr: ErrInvalidDecimalString,
		},
		{
			name:    "bare dot",
			input:   ".",
			wantErr: ErrInvalidDecimalString,
		},
		{
			name:    "multiple dots",
			input:   "1.2.3",
			wantErr: ErrInvalidDecimalString,
		},
		{
			name:    "embedded spaces",
			input:   "1 2",
			wantErr: ErrInvalidDecimalString,
		},
		{
			name:    "missing exponent digits",
			input:   "1e",
			wantErr: ErrInvalidDecimalString,
		},
		{
			name:    "positive overflow",
			input:   "9223372036854775808",
			wantErr: ErrNotRepresentable,
		},
		{
			name:    "too much decimal scale",
			input:   "1e-20",
			wantErr: ErrNotRepresentable,
		},
		{
			name:    "min int64 exponent overflow boundary",
			input:   "1e-9223372036854775808",
			wantErr: ErrNotRepresentable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewFromDecimalString(tt.input)
			require.ErrorIs(t, err, tt.wantErr)
			assert.True(t, got.IsInvalid(), "invalid input should return invalid Rat")
		})
	}
}

// TestRat_ToDecimalString tests exact decimal formatting for terminating decimals.
func TestRat_ToDecimalString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		rat  Rat
		want string
	}{
		{
			name: "zero",
			rat:  Zero(),
			want: "0",
		},
		{
			name: "integer",
			rat:  NewFromInt64(42),
			want: "42",
		},
		{
			name: "simple fraction",
			rat:  New(1, 2),
			want: "0.5",
		},
		{
			name: "two decimal digits",
			rat:  New(3, 4),
			want: "0.75",
		},
		{
			name: "mixed value",
			rat:  New(25, 2),
			want: "12.5",
		},
		{
			name: "negative fraction",
			rat:  New(-1, 8),
			want: "-0.125",
		},
		{
			name: "leading zero after decimal point",
			rat:  New(1, 20),
			want: "0.05",
		},
		{
			name: "multiple leading zeroes after decimal point",
			rat:  New(1, 40),
			want: "0.025",
		},
		{
			name: "maximum supported decimal scale",
			rat:  New(1, 10000000000000000000),
			want: "0.0000000000000000001",
		},
		{
			name: "canonical without trailing zeros",
			rat:  New(12, 10),
			want: "1.2",
		},
		{
			name: "min int64",
			rat:  NewFromInt64(math.MinInt64),
			want: "-9223372036854775808",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.rat.ToDecimalString()
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestRat_ToDecimalString_InvalidInputs tests decimal formatting failures.
func TestRat_ToDecimalString_InvalidInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rat     Rat
		wantErr error
	}{
		{
			name:    "invalid rat",
			rat:     Rat{},
			wantErr: ErrInvalid,
		},
		{
			name:    "one third",
			rat:     New(1, 3),
			wantErr: ErrNonTerminatingDecimal,
		},
		{
			name:    "two sevenths",
			rat:     New(2, 7),
			wantErr: ErrNonTerminatingDecimal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.rat.ToDecimalString()
			require.ErrorIs(t, err, tt.wantErr)
			assert.Empty(t, got)
		})
	}
}

// TestRat_DecimalRoundTrip tests canonical parse and exact format round trips.
func TestRat_DecimalRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         string
		wantCanonical string
	}{
		{
			name:          "plain decimal",
			input:         "1.2300",
			wantCanonical: "1.23",
		},
		{
			name:          "scientific notation",
			input:         "3.5e-1",
			wantCanonical: "0.35",
		},
		{
			name:          "integer scientific notation",
			input:         "1.25e2",
			wantCanonical: "125",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rat, err := NewFromDecimalString(tt.input)
			require.NoError(t, err)

			got, err := rat.ToDecimalString()
			require.NoError(t, err)
			assert.Equal(t, tt.wantCanonical, got)
		})
	}
}
