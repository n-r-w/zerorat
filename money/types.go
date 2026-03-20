// Package money provides a Money type built on zerorat.Rat that supports exact arithmetic,
// rounding, currency safety, and both mutable/immutable APIs with multi-arg operations.
package money

import (
	"errors"

	"github.com/n-r-w/zerorat"
)

// Currency is an alias for string to provide clarity in documentation and function signatures.
type Currency = string

// Money represents a monetary value with a currency and exact rational amount.
// Uses zerorat.Rat for zero-allocation exact arithmetic.
// Money is invalid if currency is empty or amount is invalid.
type Money struct {
	currency Currency    // Currency code (case-sensitive, must be non-empty for valid Money)
	amount   zerorat.Rat // Exact rational amount
}

// Error definitions for Money operations.
var (
	// ErrMoneyInvalid indicates that a Money value is in an invalid state.
	ErrMoneyInvalid = errors.New("invalid money")

	// ErrMoneyCurrencyMismatch indicates that an operation was attempted between Money values with different currencies.
	ErrMoneyCurrencyMismatch = errors.New("money currency mismatch")
)

// NewInvalid creates a new invalid Money value.
func NewInvalid() Money {
	return Money{}
}

// NewMoney creates a new Money with the given currency and amount.
// Returns a value, not a pointer, following project preferences.
// The Money is invalid if currency is empty or amount is invalid.
// The amount is automatically reduced if needed by the zerorat constructor.
func NewMoney(currency Currency, amount zerorat.Rat) Money {
	m, _ := NewMoneyErr(currency, amount)
	return m
}

// NewMoneyErr creates a new Money with the given currency and amount.
// Returns a value, not a pointer, following project preferences.
// The Money is invalid if currency is empty or amount is invalid.
// The amount is automatically reduced if needed by the zerorat constructor.
func NewMoneyErr(currency Currency, amount zerorat.Rat) (Money, error) {
	// If currency is empty, return invalid Money
	if currency == "" {
		return Money{}, ErrMoneyInvalid // invalid: empty currency and zero-value Rat (which has denominator=0)
	}

	// If amount is invalid, return invalid Money
	if amount.IsInvalid() {
		return Money{}, ErrMoneyInvalid // invalid: empty currency and invalid Rat
	}

	// Create valid Money - amount is already reduced by zerorat constructor
	return Money{
		currency: currency,
		amount:   amount,
	}, nil
}

// NewMoneyInt creates a Money from an integer value.
// Equivalent to NewMoney(currency, zerorat.NewFromInt(value)).
func NewMoneyInt(currency Currency, value int64) Money {
	m, _ := NewMoneyIntErr(currency, value)
	return m
}

// NewMoneyInt64Ptr creates a Money from an integer pointer.
func NewMoneyInt64Ptr(currency Currency, value *int64) Money {
	if value == nil {
		return Money{}
	}
	return NewMoneyInt(currency, *value)
}

// NewMoneyIntPtr creates a Money from an integer pointer.
func NewMoneyIntPtr(currency Currency, value *int) Money {
	if value == nil {
		return Money{}
	}
	return NewMoneyInt(currency, int64(*value))
}

// NewMoneyIntErr creates a Money from an integer value.
// Equivalent to NewMoney(currency, zerorat.NewFromInt(value)).
func NewMoneyIntErr(currency Currency, value int64) (Money, error) {
	// If currency is empty, return invalid Money
	if currency == "" {
		return Money{}, nil // invalid: empty currency and zero-value Rat (which has denominator=0)
	}

	amount := zerorat.NewFromInt64(value)
	return NewMoneyErr(currency, amount)
}

// NewMoneyFloat64Ptr creates Money from a float64 pointer.
// Nil returns the zero-value Money and nil error.
func NewMoneyFloat64Ptr(currency Currency, value *float64) (Money, error) {
	if value == nil {
		return Money{}, nil
	}
	return NewMoneyFloat(currency, *value)
}

// NewMoneyFloat32Ptr creates Money from a float32 pointer.
// Nil returns the zero-value Money and nil error.
func NewMoneyFloat32Ptr(currency Currency, value *float32) (Money, error) {
	if value == nil {
		return Money{}, nil
	}
	return NewMoneyFloat(currency, float64(*value))
}

// NewMoneyFloat creates a Money from a float64 value.
// The amount preserves the exact IEEE-754 binary float value.
// It does not normalize decimal literals to currency precision.
// Use integer minor units or fractions when decimal money semantics are intended.
// Returns an error if currency is empty or float conversion fails.
func NewMoneyFloat(currency Currency, value float64) (Money, error) {
	amount, err := zerorat.NewFromFloat64(value)
	if err != nil {
		return Money{}, err
	}

	return NewMoneyErr(currency, amount)
}

// NewMoneyFromFraction creates a Money from a fraction (numerator/denominator).
// Returns invalid Money if currency is empty or denominator is zero.
// The fraction is automatically reduced by the zerorat constructor.
func NewMoneyFromFraction(numerator int64, denominator uint64, currency Currency) Money {
	m, _ := NewMoneyFromFractionErr(numerator, denominator, currency)
	return m
}

// NewMoneyFromFractionErr creates a Money from a fraction (numerator/denominator).
// Returns invalid Money if currency is empty or denominator is zero.
// The fraction is automatically reduced by the zerorat constructor.
func NewMoneyFromFractionErr(numerator int64, denominator uint64, currency Currency) (Money, error) {
	if currency == "" {
		return Money{}, ErrMoneyInvalid
	}

	amount := zerorat.New(numerator, denominator)
	return NewMoneyErr(currency, amount)
}

// ZeroMoney creates a Money representing zero in the given currency.
// Returns invalid Money if currency is empty.
func ZeroMoney(currency Currency) Money {
	m, _ := ZeroMoneyErr(currency)
	return m
}

// ZeroMoneyErr creates a Money representing zero in the given currency.
// Returns invalid Money if currency is empty.
func ZeroMoneyErr(currency Currency) (Money, error) {
	amount := zerorat.Zero()
	return NewMoneyErr(currency, amount)
}

// IsValid checks if the Money is in a valid state.
// Returns true if currency is non-empty and amount is valid.
func (m Money) IsValid() bool {
	return m.currency != "" && m.amount.IsValid()
}

// IsInvalid checks if the Money is in an invalid state.
// Returns true if currency is empty or amount is invalid.
func (m Money) IsInvalid() bool {
	return !m.IsValid()
}

// Invalidate marks the Money as invalid by clearing currency and invalidating amount.
// Uses pointer receiver as this is a mutable operation.
func (m *Money) Invalidate() {
	m.currency = ""
	m.amount.Invalidate()
}

// Currency returns the currency code of the Money.
// Returns empty string for invalid Money.
func (m Money) Currency() string {
	return m.currency
}

// Amount returns the underlying zerorat.Rat amount.
// Returns invalid Rat for invalid Money.
func (m Money) Amount() zerorat.Rat {
	return m.amount
}

// SameCurrency checks if this Money has the same currency as another Money.
// Returns true only if both Money values are valid and have matching currencies.
// Uses value receiver as this is an immutable operation.
func (m Money) SameCurrency(other Money) bool {
	return hasSameCurrency(m, other)
}

// SameCurrency is a convenience function that checks if two Money values have the same currency.
// Returns true only if both Money values are valid and have matching currencies.
func SameCurrency(a, b Money) bool {
	return hasSameCurrency(a, b)
}

// SameCurrencies is a convenience function that checks if all Money values have the same currency.
// Returns true if all Money values are valid and have matching currencies.
func SameCurrencies(moneys ...Money) bool {
	ok, err := SameCurrenciesErr(moneys...)
	return err == nil && ok
}

// SameCurrenciesErr is a convenience function that checks if all Money values have the same currency.
func SameCurrenciesErr(moneys ...Money) (bool, error) {
	//  check for invalid operands
	for _, m := range moneys {
		if m.IsInvalid() {
			return false, ErrMoneyInvalid
		}
	}

	if len(moneys) <= 1 {
		return true, nil
	}

	for i := 1; i < len(moneys); i++ {
		if !hasSameCurrency(moneys[0], moneys[i]) {
			return false, nil
		}
	}
	return true, nil
}

// IsNegative checks if the Money represents a negative value.
// Returns true if Money is valid and amount is less than zero.
// Uses value receiver as this is an immutable operation.
func (m Money) IsNegative() bool {
	return m.IsValid() && m.amount.Sign() < 0
}

// IsPositive checks if the Money represents a positive value.
// Returns true if Money is valid and amount is greater than zero.
// Uses value receiver as this is an immutable operation.
func (m Money) IsPositive() bool {
	return m.IsValid() && m.amount.Sign() > 0
}

// IsEmpty checks if the Money is in an empty (invalid) state.
// This is an alias for IsInvalid() for semantic clarity.
// Uses value receiver as this is an immutable operation.
func (m Money) IsEmpty() bool {
	return m.IsInvalid()
}

// ToInt64Err converts the monetary amount to an int64 with error handling.
// Returns ErrMoneyInvalid if the Money value is invalid.
func (m Money) ToInt64Err() (int64, error) {
	if m.IsInvalid() {
		return 0, ErrMoneyInvalid
	}
	return m.amount.ToInt64Err()
}

// ToInt64 converts the monetary amount to an int64.
// Returns 0 if the Money value is invalid.
func (m Money) ToInt64() int64 {
	result, _ := m.ToInt64Err()
	return result
}

// ToIntErr converts the monetary amount to an int with error handling.
// Returns ErrMoneyInvalid if the Money value is invalid.
func (m Money) ToIntErr() (int, error) {
	if m.IsInvalid() {
		return 0, ErrMoneyInvalid
	}
	return m.amount.ToIntErr()
}

// ToInt converts the monetary amount to an int.
// Returns 0 if the Money value is invalid.
func (m Money) ToInt() int {
	result, _ := m.ToIntErr()
	return result
}

// ToInt64Ptr converts the monetary amount to an int64 pointer.
func (m Money) ToInt64Ptr() *int64 {
	if m.IsInvalid() {
		return nil
	}
	result := m.ToInt64()
	return &result
}

// ToIntPtr converts the monetary amount to an int pointer.
func (m Money) ToIntPtr() *int {
	if m.IsInvalid() {
		return nil
	}
	result := m.ToInt()
	return &result
}

// ToFloat64Err converts the monetary amount to a float64 with error handling.
// Returns ErrMoneyInvalid if the Money value is invalid.
func (m Money) ToFloat64Err() (float64, error) {
	if m.IsInvalid() {
		return 0, ErrMoneyInvalid
	}
	return m.amount.ToFloat64Err()
}

// ToFloat64 converts the monetary amount to a float64.
// Returns 0 if the Money value is invalid.
func (m Money) ToFloat64() float64 {
	result, _ := m.ToFloat64Err()
	return result
}

// ToFloat32Err converts the monetary amount to a float32 with error handling.
// Returns ErrMoneyInvalid if the Money value is invalid.
func (m Money) ToFloat32Err() (float32, error) {
	if m.IsInvalid() {
		return 0, ErrMoneyInvalid
	}
	return m.amount.ToFloat32Err()
}

// ToFloat32 converts the monetary amount to a float32.
// Returns 0 if the Money value is invalid.
func (m Money) ToFloat32() float32 {
	result, _ := m.ToFloat32Err()
	return result
}

// ToFloat32Ptr converts the monetary amount to a float32 pointer.
func (m Money) ToFloat32Ptr() *float32 {
	if m.IsInvalid() {
		return nil
	}
	result := m.ToFloat32()
	return &result
}

// ToFloat64Ptr converts the monetary amount to a float64 pointer.
func (m Money) ToFloat64Ptr() *float64 {
	if m.IsInvalid() {
		return nil
	}
	result := m.ToFloat64()
	return &result
}
