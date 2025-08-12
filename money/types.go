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

// NewMoneyIntErr creates a Money from an integer value.
// Equivalent to NewMoney(currency, zerorat.NewFromInt(value)).
func NewMoneyIntErr(currency Currency, value int64) (Money, error) {
	// If currency is empty, return invalid Money
	if currency == "" {
		return Money{}, nil // invalid: empty currency and zero-value Rat (which has denominator=0)
	}

	amount := zerorat.NewFromInt(value)
	return NewMoneyErr(currency, amount)
}

// NewMoneyFloat creates a Money from a float64 value.
// Returns invalid Money if currency is empty or float conversion fails.
// Equivalent to NewMoney(currency, zerorat.NewFromFloat64(value)).
func NewMoneyFloat(currency Currency, value float64) Money {
	m, _ := NewMoneyFloatErr(currency, value)
	return m
}

// NewMoneyFloatErr creates a Money from a float64 value.
// Returns invalid Money if currency is empty or float conversion fails.
// Equivalent to NewMoney(currency, zerorat.NewFromFloat64(value)).
func NewMoneyFloatErr(currency Currency, value float64) (Money, error) {
	amount := zerorat.NewFromFloat64(value)
	if amount.IsInvalid() {
		return Money{}, ErrMoneyInvalid
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
// Returns true if there are less than 2 Money values, or if all Money values have the same currency.
func SameCurrencies(moneys ...Money) bool {
	if len(moneys) == 0 {
		return true
	}

	if len(moneys) == 1 {
		return moneys[0].IsValid()
	}

	for i := 1; i < len(moneys); i++ {
		if !hasSameCurrency(moneys[0], moneys[i]) {
			return false
		}
	}
	return true
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
