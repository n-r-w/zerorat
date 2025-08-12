package money

// MulMoney multiplies this Money by another Money (mutable operation).
// Requires same currency. Result currency remains the same as operands.
func (m *Money) MulMoney(other Money) error {
	// Check if either operand is invalid
	if m.IsInvalid() || other.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Check currency match
	if !m.SameCurrency(other) {
		m.Invalidate()
		return ErrMoneyCurrencyMismatch
	}

	// Delegate to Rat arithmetic
	m.amount.Mul(other.amount)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// MultipliedMoneyErr returns the product of this Money and another Money (immutable operation with error).
func (m Money) MultipliedMoneyErr(other Money) (Money, error) {
	result := m // copy
	err := result.MulMoney(other)
	return result, err
}

// MultipliedMoney returns the product of this Money and another Money (immutable operation without error).
func (m Money) MultipliedMoney(other Money) Money {
	result, _ := m.MultipliedMoneyErr(other)
	return result
}

// DivMoney divides this Money by another Money (mutable operation).
// Requires same currency. Result currency remains the same as operands.
func (m *Money) DivMoney(other Money) error {
	// Check if either operand is invalid
	if m.IsInvalid() || other.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Check currency match
	if !m.SameCurrency(other) {
		m.Invalidate()
		return ErrMoneyCurrencyMismatch
	}

	// Check for division by zero (zero amount)
	if other.amount.IsZero() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Delegate to Rat arithmetic
	m.amount.Div(other.amount)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// DividedMoneyErr returns the quotient of this Money and another Money (immutable operation with error).
func (m Money) DividedMoneyErr(other Money) (Money, error) {
	result := m // copy
	err := result.DivMoney(other)
	return result, err
}

// DividedMoney returns the quotient of this Money and another Money (immutable operation without error).
func (m Money) DividedMoney(other Money) Money {
	result, _ := m.DividedMoneyErr(other)
	return result
}

// AddMany adds multiple Money values to this Money (mutable varargs operation).
// Requires all operands to have same currency. Early termination on first error.
func (m *Money) AddMany(others ...Money) error {
	for _, other := range others {
		if err := m.Add(other); err != nil {
			return err
		}
	}
	return nil
}

// AddedManyErr returns the sum of this Money and multiple Money values (immutable varargs operation with error).
func (m Money) AddedManyErr(others ...Money) (Money, error) {
	result := m // copy
	err := result.AddMany(others...)
	return result, err
}

// AddedMany returns the sum of this Money and multiple Money values (immutable varargs operation without error).
func (m Money) AddedMany(others ...Money) Money {
	result, _ := m.AddedManyErr(others...)
	return result
}

// SubMany subtracts multiple Money values from this Money (mutable varargs operation).
// Requires all operands to have same currency. Early termination on first error.
func (m *Money) SubMany(others ...Money) error {
	for _, other := range others {
		if err := m.Sub(other); err != nil {
			return err
		}
	}
	return nil
}

// SubtractedManyErr returns the difference of this Money and multiple Money values
// (immutable varargs operation with error).
func (m Money) SubtractedManyErr(others ...Money) (Money, error) {
	result := m // copy
	err := result.SubMany(others...)
	return result, err
}

// SubtractedMany returns the difference of this Money and multiple Money values
// (immutable varargs operation without error).
func (m Money) SubtractedMany(others ...Money) Money {
	result, _ := m.SubtractedManyErr(others...)
	return result
}

// MulManyInt multiplies this Money by multiple int64 values (mutable varargs operation).
// Multiplies sequentially: m = m * v1 * v2 * ... * vN.
func (m *Money) MulManyInt(values ...int64) error {
	for _, value := range values {
		if err := m.MulInt(value); err != nil {
			return err
		}
	}
	return nil
}

// MultipliedManyIntErr returns the product of this Money and multiple int64 values
// (immutable varargs operation with error).
func (m Money) MultipliedManyIntErr(values ...int64) (Money, error) {
	result := m // copy
	err := result.MulManyInt(values...)
	return result, err
}

// MultipliedManyInt returns the product of this Money and multiple int64 values
// (immutable varargs operation without error).
func (m Money) MultipliedManyInt(values ...int64) Money {
	result, _ := m.MultipliedManyIntErr(values...)
	return result
}

// MulManyFloat multiplies this Money by multiple float64 values (mutable varargs operation).
// Multiplies sequentially: m = m * v1 * v2 * ... * vN.
func (m *Money) MulManyFloat(values ...float64) error {
	for _, value := range values {
		if err := m.MulFloat(value); err != nil {
			return err
		}
	}
	return nil
}

// MultipliedManyFloatErr returns the product of this Money and multiple float64 values
// (immutable varargs operation with error).
func (m Money) MultipliedManyFloatErr(values ...float64) (Money, error) {
	result := m // copy
	err := result.MulManyFloat(values...)
	return result, err
}

// MultipliedManyFloat returns the product of this Money and multiple float64 values
// (immutable varargs operation without error).
func (m Money) MultipliedManyFloat(values ...float64) Money {
	result, _ := m.MultipliedManyFloatErr(values...)
	return result
}
