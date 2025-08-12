package money

// Compare performs three-way comparison of Money values.
// Returns -1 if m < other, 0 if m == other, 1 if m > other.
// Returns 0 for invalid operands or currency mismatch (cannot be meaningfully compared).
// Uses value receiver as this is an immutable operation.
func (m Money) Compare(other Money) int {
	// Invalid operands or currency mismatch cannot be meaningfully compared
	if m.IsInvalid() || other.IsInvalid() || !hasSameCurrency(m, other) {
		return 0
	}

	// Delegate to underlying Rat comparison
	return m.amount.Compare(other.amount)
}

// CompareErr performs three-way comparison of Money values with error handling.
// Returns -1 if m < other, 0 if m == other, 1 if m > other.
// Returns error for invalid operands or currency mismatch.
// Uses value receiver as this is an immutable operation.
func (m Money) CompareErr(other Money) (int, error) {
	// Check if either operand is invalid
	if m.IsInvalid() || other.IsInvalid() {
		return 0, ErrMoneyInvalid
	}

	// Check currency match
	if !hasSameCurrency(m, other) {
		return 0, ErrMoneyCurrencyMismatch
	}

	// Delegate to underlying Rat comparison
	return m.amount.Compare(other.amount), nil
}

// Equal checks if two Money values are equal.
// Returns true only if both are valid, have same currency, and equal amounts.
// Uses value receiver as this is an immutable operation.
func (m Money) Equal(other Money) bool {
	// Invalid operands or currency mismatch are not equal
	if m.IsInvalid() || other.IsInvalid() || !hasSameCurrency(m, other) {
		return false
	}

	// Delegate to underlying Rat equality
	return m.amount.Equal(other.amount)
}

// EqualErr checks if two Money values are equal with error handling.
// Returns true only if both are valid, have same currency, and equal amounts.
// Returns error for invalid operands or currency mismatch.
// Uses value receiver as this is an immutable operation.
func (m Money) EqualErr(other Money) (bool, error) {
	// Check if either operand is invalid
	if m.IsInvalid() || other.IsInvalid() {
		return false, ErrMoneyInvalid
	}

	// Check currency match
	if !hasSameCurrency(m, other) {
		return false, ErrMoneyCurrencyMismatch
	}

	// Delegate to underlying Rat equality
	return m.amount.Equal(other.amount), nil
}

// Less checks if this Money is less than another Money.
// Returns true only if both are valid, have same currency, and m < other.
// Uses value receiver as this is an immutable operation.
func (m Money) Less(other Money) bool {
	// Invalid operands or currency mismatch cannot be compared
	if m.IsInvalid() || other.IsInvalid() || !hasSameCurrency(m, other) {
		return false
	}

	// Delegate to underlying Rat comparison
	return m.amount.Compare(other.amount) < 0
}

// LessErr checks if this Money is less than another Money with error handling.
// Returns true only if both are valid, have same currency, and m < other.
// Returns error for invalid operands or currency mismatch.
// Uses value receiver as this is an immutable operation.
func (m Money) LessErr(other Money) (bool, error) {
	// Check if either operand is invalid
	if m.IsInvalid() || other.IsInvalid() {
		return false, ErrMoneyInvalid
	}

	// Check currency match
	if !hasSameCurrency(m, other) {
		return false, ErrMoneyCurrencyMismatch
	}

	// Delegate to underlying Rat comparison
	return m.amount.Compare(other.amount) < 0, nil
}

// Greater checks if this Money is greater than another Money.
// Returns true only if both are valid, have same currency, and m > other.
// Uses value receiver as this is an immutable operation.
func (m Money) Greater(other Money) bool {
	// Invalid operands or currency mismatch cannot be compared
	if m.IsInvalid() || other.IsInvalid() || !hasSameCurrency(m, other) {
		return false
	}

	// Delegate to underlying Rat comparison
	return m.amount.Compare(other.amount) > 0
}

// GreaterErr checks if this Money is greater than another Money with error handling.
// Returns true only if both are valid, have same currency, and m > other.
// Returns error for invalid operands or currency mismatch.
// Uses value receiver as this is an immutable operation.
func (m Money) GreaterErr(other Money) (bool, error) {
	// Check if either operand is invalid
	if m.IsInvalid() || other.IsInvalid() {
		return false, ErrMoneyInvalid
	}

	// Check currency match
	if !hasSameCurrency(m, other) {
		return false, ErrMoneyCurrencyMismatch
	}

	// Delegate to underlying Rat comparison
	return m.amount.Compare(other.amount) > 0, nil
}

// IsZero checks if the Money represents a zero value.
func (m Money) IsZero() bool {
	return m.IsValid() && m.amount.IsZero()
}
