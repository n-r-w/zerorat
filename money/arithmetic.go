package money

import (
	"github.com/n-r-w/zerorat"
)

const (
	// percentDivisor is used to convert percentage values to fractions (value/100).
	percentDivisor = 100
	// percentDivisorFloat is used to convert percentage values to decimals (value/100.0).
	percentDivisorFloat = 100.0
)

// Add adds another Money to this Money (mutable operation).
// Requires same currency. Sets invalid state on currency mismatch or invalid operands.
// Uses pointer receiver for mutable operation.
func (m *Money) Add(other Money) error {
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
	m.amount.Add(other.amount)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// AddedErr returns the sum of this Money and another Money (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) AddedErr(other Money) (Money, error) {
	result := m // copy
	err := result.Add(other)
	return result, err
}

// Added returns the sum of this Money and another Money (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) Added(other Money) Money {
	result, _ := m.AddedErr(other)
	return result
}

// Sub subtracts another Money from this Money (mutable operation).
// Requires same currency. Sets invalid state on currency mismatch or invalid operands.
// Uses pointer receiver for mutable operation.
func (m *Money) Sub(other Money) error {
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
	m.amount.Sub(other.amount)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// SubtractedErr returns the difference of this Money and another Money (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) SubtractedErr(other Money) (Money, error) {
	result := m // copy
	err := result.Sub(other)
	return result, err
}

// Subtracted returns the difference of this Money and another Money (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) Subtracted(other Money) Money {
	result, _ := m.SubtractedErr(other)
	return result
}

// Profit calculates profit by subtracting another Money from this Money (mutable operation).
// This is an alias for Sub operation. Requires same currency.
// Uses pointer receiver for mutable operation.
func (m *Money) Profit(other Money) error {
	return m.Sub(other)
}

// ProfitedErr returns the profit of this Money minus another Money (immutable operation with error).
// This is an alias for SubtractedErr operation. Uses value receiver for immutable operation.
func (m Money) ProfitedErr(other Money) (Money, error) {
	return m.SubtractedErr(other)
}

// Profited returns the profit of this Money minus another Money (immutable operation without error).
// This is an alias for Subtracted operation. Returns invalid Money on error.
// Uses value receiver for immutable operation.
func (m Money) Profited(other Money) Money {
	return m.Subtracted(other)
}

// PercentInt calculates percentage of this Money using an int64 value (mutable operation).
// Formula: m = m * (value / 100). Uses pointer receiver for mutable operation.
func (m *Money) PercentInt(value int64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Convert percentage to fraction: value/100
	percentRat := zerorat.New(value, percentDivisor)

	// Delegate to Rat arithmetic
	m.amount.Mul(percentRat)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// PercentIntErr returns percentage of this Money using an int64 value (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) PercentIntErr(value int64) (Money, error) {
	result := m // copy
	err := result.PercentInt(value)
	return result, err
}

// PercentedInt returns percentage of this Money using an int64 value (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) PercentedInt(value int64) Money {
	result, _ := m.PercentIntErr(value)
	return result
}

// PercentFloat calculates percentage of this Money using a float64 value (mutable operation).
// Formula: m = m * (value / 100). Uses pointer receiver for mutable operation.
func (m *Money) PercentFloat(value float64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Convert percentage to decimal: value/100
	percentFloat := value / percentDivisorFloat
	percentRat := zerorat.NewFromFloat64(percentFloat)

	// Check if float conversion was invalid
	if percentRat.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Delegate to Rat arithmetic
	m.amount.Mul(percentRat)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// PercentFloatErr returns percentage of this Money using a float64 value (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) PercentFloatErr(value float64) (Money, error) {
	result := m // copy
	err := result.PercentFloat(value)
	return result, err
}

// PercentedFloat returns percentage of this Money using a float64 value (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) PercentedFloat(value float64) Money {
	result, _ := m.PercentFloatErr(value)
	return result
}

// PercentOf calculates this Money as percentage of another Money (mutable operation).
// Formula: m = m * other (interpreted as proportion). Requires same currency.
// Uses pointer receiver for mutable operation.
func (m *Money) PercentOf(other Money) error {
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

// PercentOfErr returns this Money as percentage of another Money (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) PercentOfErr(other Money) (Money, error) {
	result := m // copy
	err := result.PercentOf(other)
	return result, err
}

// PercentedOf returns this Money as percentage of another Money (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) PercentedOf(other Money) Money {
	result, _ := m.PercentOfErr(other)
	return result
}

// AddInt adds an int64 value to this Money (mutable operation).
// Converts int64 to Rat and delegates to Money addition.
func (m *Money) AddInt(value int64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Convert int64 to Rat
	ratValue := zerorat.NewFromInt(value)

	// Delegate to Rat arithmetic
	m.amount.Add(ratValue)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// AddedIntErr returns the sum of this Money and an int64 value (immutable operation with error).
func (m Money) AddedIntErr(value int64) (Money, error) {
	result := m // copy
	err := result.AddInt(value)
	return result, err
}

// AddedInt returns the sum of this Money and an int64 value (immutable operation without error).
func (m Money) AddedInt(value int64) Money {
	result, _ := m.AddedIntErr(value)
	return result
}

// AddFloat adds a float64 value to this Money (mutable operation).
// Converts float64 to Rat and delegates to Money addition.
func (m *Money) AddFloat(value float64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Convert float64 to Rat
	ratValue := zerorat.NewFromFloat64(value)

	// Check if float conversion was invalid
	if ratValue.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Delegate to Rat arithmetic
	m.amount.Add(ratValue)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// AddedFloatErr returns the sum of this Money and a float64 value (immutable operation with error).
func (m Money) AddedFloatErr(value float64) (Money, error) {
	result := m // copy
	err := result.AddFloat(value)
	return result, err
}

// AddedFloat returns the sum of this Money and a float64 value (immutable operation without error).
func (m Money) AddedFloat(value float64) Money {
	result, _ := m.AddedFloatErr(value)
	return result
}

// SubInt subtracts an int64 value from this Money (mutable operation).
func (m *Money) SubInt(value int64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Convert int64 to Rat
	ratValue := zerorat.NewFromInt(value)

	// Delegate to Rat arithmetic
	m.amount.Sub(ratValue)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// SubtractedIntErr returns the difference of this Money and an int64 value (immutable operation with error).
func (m Money) SubtractedIntErr(value int64) (Money, error) {
	result := m // copy
	err := result.SubInt(value)
	return result, err
}

// SubtractedInt returns the difference of this Money and an int64 value (immutable operation without error).
func (m Money) SubtractedInt(value int64) Money {
	result, _ := m.SubtractedIntErr(value)
	return result
}

// SubFloat subtracts a float64 value from this Money (mutable operation).
func (m *Money) SubFloat(value float64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Convert float64 to Rat
	ratValue := zerorat.NewFromFloat64(value)

	// Check if float conversion was invalid
	if ratValue.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Delegate to Rat arithmetic
	m.amount.Sub(ratValue)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// SubtractedFloatErr returns the difference of this Money and a float64 value (immutable operation with error).
func (m Money) SubtractedFloatErr(value float64) (Money, error) {
	result := m // copy
	err := result.SubFloat(value)
	return result, err
}

// SubtractedFloat returns the difference of this Money and a float64 value (immutable operation without error).
func (m Money) SubtractedFloat(value float64) Money {
	result, _ := m.SubtractedFloatErr(value)
	return result
}

// MulInt multiplies this Money by an int64 value (mutable operation).
func (m *Money) MulInt(value int64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Convert int64 to Rat
	ratValue := zerorat.NewFromInt(value)

	// Delegate to Rat arithmetic
	m.amount.Mul(ratValue)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// MultipliedIntErr returns the product of this Money and an int64 value (immutable operation with error).
func (m Money) MultipliedIntErr(value int64) (Money, error) {
	result := m // copy
	err := result.MulInt(value)
	return result, err
}

// MultipliedInt returns the product of this Money and an int64 value (immutable operation without error).
func (m Money) MultipliedInt(value int64) Money {
	result, _ := m.MultipliedIntErr(value)
	return result
}

// MulFloat multiplies this Money by a float64 value (mutable operation).
func (m *Money) MulFloat(value float64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Convert float64 to Rat
	ratValue := zerorat.NewFromFloat64(value)

	// Check if float conversion was invalid
	if ratValue.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Delegate to Rat arithmetic
	m.amount.Mul(ratValue)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// MultipliedFloatErr returns the product of this Money and a float64 value (immutable operation with error).
func (m Money) MultipliedFloatErr(value float64) (Money, error) {
	result := m // copy
	err := result.MulFloat(value)
	return result, err
}

// MultipliedFloat returns the product of this Money and a float64 value (immutable operation without error).
func (m Money) MultipliedFloat(value float64) Money {
	result, _ := m.MultipliedFloatErr(value)
	return result
}

// DivInt divides this Money by an int64 value (mutable operation).
func (m *Money) DivInt(value int64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Check for division by zero
	if value == 0 {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Convert int64 to Rat
	ratValue := zerorat.NewFromInt(value)

	// Delegate to Rat arithmetic
	m.amount.Div(ratValue)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// DividedIntErr returns the quotient of this Money and an int64 value (immutable operation with error).
func (m Money) DividedIntErr(value int64) (Money, error) {
	result := m // copy
	err := result.DivInt(value)
	return result, err
}

// DividedInt returns the quotient of this Money and an int64 value (immutable operation without error).
func (m Money) DividedInt(value int64) Money {
	result, _ := m.DividedIntErr(value)
	return result
}

// DivFloat divides this Money by a float64 value (mutable operation).
func (m *Money) DivFloat(value float64) error {
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Check for division by zero
	if value == 0.0 {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Convert float64 to Rat
	ratValue := zerorat.NewFromFloat64(value)

	// Check if float conversion was invalid
	if ratValue.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	// Delegate to Rat arithmetic
	m.amount.Div(ratValue)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// DividedFloatErr returns the quotient of this Money and a float64 value (immutable operation with error).
func (m Money) DividedFloatErr(value float64) (Money, error) {
	result := m // copy
	err := result.DivFloat(value)
	return result, err
}

// DividedFloat returns the quotient of this Money and a float64 value (immutable operation without error).
func (m Money) DividedFloat(value float64) Money {
	result, _ := m.DividedFloatErr(value)
	return result
}

// Sum returns the sum of multiple Money values (immutable varargs operation).
func Sum(moneys ...Money) Money {
	result, _ := SumErr(moneys...)
	return result
}

// SumErr returns the sum of multiple Money values (immutable varargs operation).
func SumErr(moneys ...Money) (Money, error) {
	if len(moneys) == 0 {
		return Money{}, nil
	}
	if len(moneys) == 1 {
		return moneys[0], nil
	}

	result := moneys[0]
	for i := 1; i < len(moneys); i++ {
		if err := result.Add(moneys[i]); err != nil {
			return Money{}, err
		}
	}
	return result, nil
}
