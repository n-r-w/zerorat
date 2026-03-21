package zerorat

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRat_UnmarshalJSON_AcceptsNumbersAndStrings verifies that Rat accepts both JSON numbers and quoted decimal strings.
func TestRat_UnmarshalJSON_AcceptsNumbersAndStrings(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  Rat
	}{
		{
			name:  "integer number",
			input: `42`,
			want:  NewFromInt64(42),
		},
		{
			name:  "decimal number",
			input: `12.5`,
			want:  New(25, 2),
		},
		{
			name:  "scientific notation number",
			input: `3.5e-1`,
			want:  New(7, 20),
		},
		{
			name:  "decimal string",
			input: `"0.35"`,
			want:  New(7, 20),
		},
		{
			name:  "scientific notation string",
			input: `"3.5e-1"`,
			want:  New(7, 20),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got Rat

			err := json.Unmarshal([]byte(tt.input), &got)
			require.NoError(t, err)
			assert.True(t, got.Equal(tt.want), "decoded value mismatch")
		})
	}
}

// TestRat_UnmarshalJSON_Null verifies that JSON null invalidates a value Rat field.
func TestRat_UnmarshalJSON_Null(t *testing.T) {
	t.Parallel()

	type payload struct {
		Value Rat `json:"value"`
	}

	got := payload{Value: One()}

	err := json.Unmarshal([]byte(`{"value":null}`), &got)
	require.NoError(t, err)
	assert.True(t, got.Value.IsInvalid(), "null should invalidate value Rat fields")
}

// TestRat_UnmarshalJSON_PointerNull verifies that JSON null keeps pointer Rat fields nil.
func TestRat_UnmarshalJSON_PointerNull(t *testing.T) {
	t.Parallel()

	type payload struct {
		Value *Rat `json:"value"`
	}

	var got payload

	err := json.Unmarshal([]byte(`{"value":null}`), &got)
	require.NoError(t, err)
	assert.Nil(t, got.Value)
}

// TestRat_UnmarshalJSON_InvalidInputs verifies that Rat rejects unsupported JSON token kinds and unsupported numeric text.
func TestRat_UnmarshalJSON_InvalidInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "rational string",
			input:   `"1/3"`,
			wantErr: ErrInvalidDecimalString,
		},
		{
			name:    "boolean",
			input:   `true`,
			wantErr: ErrInvalidDecimalString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := One()

			err := json.Unmarshal([]byte(tt.input), &got)
			require.ErrorIs(t, err, tt.wantErr)
			assert.True(t, got.IsInvalid(), "failed unmarshal should leave Rat invalid")
		})
	}
}

// TestRat_UnmarshalJSON_MalformedQuotedToken verifies that broken quoted JSON leaves the receiver invalid.
func TestRat_UnmarshalJSON_MalformedQuotedToken(t *testing.T) {
	t.Parallel()

	got := One()

	err := got.UnmarshalJSON([]byte(`"unterminated`))
	require.Error(t, err)
	assert.True(t, got.IsInvalid(), "failed quoted-token decode should leave Rat invalid")
}

// TestRat_MarshalJSON verifies that Rat writes normalized JSON numbers and rejects unsupported states.
func TestRat_MarshalJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		rat     Rat
		want    string
		wantErr error
	}{
		{
			name: "zero",
			rat:  Zero(),
			want: `0`,
		},
		{
			name: "integer",
			rat:  NewFromInt64(42),
			want: `42`,
		},
		{
			name: "decimal",
			rat:  New(25, 2),
			want: `12.5`,
		},
		{
			name: "normalized decimal",
			rat:  New(350, 1000),
			want: `0.35`,
		},
		{
			name:    "invalid rat",
			rat:     Rat{},
			wantErr: ErrInvalid,
		},
		{
			name:    "non terminating decimal",
			rat:     New(1, 3),
			wantErr: ErrNonTerminatingDecimal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := json.Marshal(tt.rat)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, string(got))
		})
	}
}

// TestRat_JSONRoundTrip verifies that Rat normalizes supported JSON inputs to normalized decimal output.
func TestRat_JSONRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		input      string
		wantOutput string
	}{
		{
			name:       "number input stays numeric",
			input:      `0.35`,
			wantOutput: `0.35`,
		},
		{
			name:       "string input becomes normalized number",
			input:      `"3.5e-1"`,
			wantOutput: `0.35`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var rat Rat

			err := json.Unmarshal([]byte(tt.input), &rat)
			require.NoError(t, err)

			got, err := json.Marshal(rat)
			require.NoError(t, err)
			assert.Equal(t, tt.wantOutput, string(got))
		})
	}
}
