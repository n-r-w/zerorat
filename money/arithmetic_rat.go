package money

import "github.com/n-r-w/zerorat"

// validateRatOperation performs common validation for all Rat operations.
// Returns true if the operation should proceed, false if it should fail.
// If validation fails, the Money is invalidated and the function returns false.
func (m *Money) validateRatOperation(value zerorat.Rat) bool {
	if m.IsInvalid() {
		return false
	}

	if value.IsInvalid() {
		m.Invalidate()
		return false
	}

	return true
}

// finalizeRatOperation performs post-operation validation and cleanup.
// Returns ErrMoneyInvalid if the operation resulted in an invalid state.
func (m *Money) finalizeRatOperation() error {
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}
	return nil
}

// AddRat adds a zerorat.Rat value to this Money (mutable operation).
// Sets invalid state on invalid operands or arithmetic overflow.
// Uses pointer receiver for mutable operation.
func (m *Money) AddRat(value zerorat.Rat) error {
	if !m.validateRatOperation(value) {
		return ErrMoneyInvalid
	}

	m.amount.Add(value)
	return m.finalizeRatOperation()
}

// AddedRatErr returns the sum of this Money and a zerorat.Rat value (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) AddedRatErr(value zerorat.Rat) (Money, error) {
	result := m // copy
	err := result.AddRat(value)
	return result, err
}

// AddedRat returns the sum of this Money and a zerorat.Rat value (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) AddedRat(value zerorat.Rat) Money {
	result, _ := m.AddedRatErr(value)
	return result
}

// SubRat subtracts a zerorat.Rat value from this Money (mutable operation).
// Sets invalid state on invalid operands or arithmetic overflow.
// Uses pointer receiver for mutable operation.
func (m *Money) SubRat(value zerorat.Rat) error {
	if !m.validateRatOperation(value) {
		return ErrMoneyInvalid
	}

	m.amount.Sub(value)
	return m.finalizeRatOperation()
}

// SubtractedRatErr returns the difference of this Money and a zerorat.Rat value (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) SubtractedRatErr(value zerorat.Rat) (Money, error) {
	result := m // copy
	err := result.SubRat(value)
	return result, err
}

// SubtractedRat returns the difference of this Money and a zerorat.Rat value (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) SubtractedRat(value zerorat.Rat) Money {
	result, _ := m.SubtractedRatErr(value)
	return result
}

// MulRat multiplies this Money by a zerorat.Rat value (mutable operation).
// Sets invalid state on invalid operands or arithmetic overflow.
// Uses pointer receiver for mutable operation.
func (m *Money) MulRat(value zerorat.Rat) error {
	if !m.validateRatOperation(value) {
		return ErrMoneyInvalid
	}

	m.amount.Mul(value)
	return m.finalizeRatOperation()
}

// MultipliedRatErr returns the product of this Money and a zerorat.Rat value (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) MultipliedRatErr(value zerorat.Rat) (Money, error) {
	result := m // copy
	err := result.MulRat(value)
	return result, err
}

// MultipliedRat returns the product of this Money and a zerorat.Rat value (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) MultipliedRat(value zerorat.Rat) Money {
	result, _ := m.MultipliedRatErr(value)
	return result
}

// DivRat divides this Money by a zerorat.Rat value (mutable operation).
// Sets invalid state on invalid operands, division by zero, or arithmetic overflow.
// Uses pointer receiver for mutable operation.
func (m *Money) DivRat(value zerorat.Rat) error {
	if !m.validateRatOperation(value) {
		return ErrMoneyInvalid
	}

	// Check for division by zero
	if value.IsZero() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	m.amount.Div(value)
	return m.finalizeRatOperation()
}

// DividedRatErr returns the quotient of this Money and a zerorat.Rat value (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) DividedRatErr(value zerorat.Rat) (Money, error) {
	result := m // copy
	err := result.DivRat(value)
	return result, err
}

// DividedRat returns the quotient of this Money and a zerorat.Rat value (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) DividedRat(value zerorat.Rat) Money {
	result, _ := m.DividedRatErr(value)
	return result
}
