package money

// ScaleDown scales the Money down by n decimal places (mutable operation).
// Equivalent to dividing by 10^n, moving the decimal point left.
// For negative n, calls ScaleUp with |n|.
// Sets invalid state on overflow or with invalid operands.
// Uses pointer receiver for mutable operation.
func (m *Money) ScaleDown(n int) error {
	// Check if Money is invalid
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Delegate to Rat scaling
	m.amount.ScaleDown(n)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// ScaledDownErr returns a new Money scaled down by n decimal places (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) ScaledDownErr(n int) (Money, error) {
	result := m // copy
	err := result.ScaleDown(n)
	return result, err
}

// ScaledDown returns a new Money scaled down by n decimal places (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) ScaledDown(n int) Money {
	result, _ := m.ScaledDownErr(n)
	return result
}

// ScaleUp scales the Money up by n decimal places (mutable operation).
// Equivalent to multiplying by 10^n, moving the decimal point right.
// For negative n, calls ScaleDown with |n|.
// Sets invalid state on overflow or with invalid operands.
// Uses pointer receiver for mutable operation.
func (m *Money) ScaleUp(n int) error {
	// Check if Money is invalid
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Delegate to Rat scaling
	m.amount.ScaleUp(n)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// ScaledUpErr returns a new Money scaled up by n decimal places (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) ScaledUpErr(n int) (Money, error) {
	result := m // copy
	err := result.ScaleUp(n)
	return result, err
}

// ScaledUp returns a new Money scaled up by n decimal places (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) ScaledUp(n int) Money {
	result, _ := m.ScaledUpErr(n)
	return result
}
