package money

import (
	"github.com/n-r-w/zerorat"
)

// Round rounds the Money to the specified scale using the given rounding strategy (mutable operation).
// Scale interpretation:
// - scale = 0: round to integer (1.23 -> 1)
// - scale > 0: round to decimal places (1.234 with scale=2 -> 1.23)
// - scale < 0: round to powers of 10 (1234 with scale=-2 -> 1200)
// Uses pointer receiver for mutable operation.
func (m *Money) Round(roundType zerorat.RoundType, scale int) error {
	// Check if Money is invalid
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// Delegate to Rat rounding
	m.amount.Round(roundType, scale)

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// RoundedErr returns a new Money rounded to the specified scale (immutable operation with error).
// Uses value receiver for immutable operation.
func (m Money) RoundedErr(roundType zerorat.RoundType, scale int) (Money, error) {
	result := m // copy
	err := result.Round(roundType, scale)
	return result, err
}

// Rounded returns a new Money rounded to the specified scale (immutable operation without error).
// Returns invalid Money on error. Uses value receiver for immutable operation.
func (m Money) Rounded(roundType zerorat.RoundType, scale int) Money {
	result, _ := m.RoundedErr(roundType, scale)
	return result
}

// RedistributeRoundedErr rounds cumulative totals and emits per-item deltas in input order.
// Scale uses the same semantics as Money.Round: 0 rounds to integers, positive values round
// to decimal places, and negative values round to powers of ten.
// The returned values preserve order and sum to the rounded total.
func RedistributeRoundedErr(values []Money, roundType zerorat.RoundType, scale int) ([]Money, error) {
	if values == nil {
		return nil, nil
	}

	if len(values) == 0 {
		return values[:0], nil
	}

	sameCurrency, err := SameCurrenciesErr(values...)
	if err != nil {
		return nil, err
	}
	if !sameCurrency {
		return nil, ErrMoneyCurrencyMismatch
	}

	currency := values[0].Currency()
	runningTotal := ZeroMoney(currency)
	previousRoundedTotal := ZeroMoney(currency)
	result := make([]Money, 0, len(values))

	for _, value := range values {
		if addErr := runningTotal.Add(value); addErr != nil {
			return nil, addErr
		}

		roundedTotal, roundErr := runningTotal.RoundedErr(roundType, scale)
		if roundErr != nil {
			return nil, roundErr
		}

		delta, deltaErr := roundedTotal.SubtractedErr(previousRoundedTotal)
		if deltaErr != nil {
			return nil, deltaErr
		}

		result = append(result, delta)
		previousRoundedTotal = roundedTotal
	}

	return result, nil
}

// RedistributeRounded rounds cumulative totals and emits per-item deltas in input order.
// Scale uses the same semantics as Money.Round.
// It returns nil on error.
func RedistributeRounded(values []Money, roundType zerorat.RoundType, scale int) []Money {
	result, _ := RedistributeRoundedErr(values, roundType, scale)
	return result
}

// Ceil rounds the Money toward positive infinity to the specified scale (mutable operation).
// Mathematical ceiling function: always rounds up for positive numbers, truncates for negative numbers.
// Uses pointer receiver for mutable operation.
func (m *Money) Ceil(scale int) error {
	// Check if Money is invalid
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// For ceiling, we need to determine the correct rounding strategy based on sign
	if m.amount.Sign() >= 0 {
		// Positive or zero: use RoundUp (away from zero, which is toward positive infinity)
		m.amount.Round(zerorat.RoundUp, scale)
	} else {
		// Negative: use RoundDown (toward zero, which is toward positive infinity for negatives)
		m.amount.Round(zerorat.RoundDown, scale)
	}

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// CeiledErr returns a new Money rounded toward positive infinity to the specified scale
// (immutable operation with error). Uses value receiver for immutable operation.
func (m Money) CeiledErr(scale int) (Money, error) {
	result := m // copy
	err := result.Ceil(scale)
	return result, err
}

// Ceiled returns a new Money rounded toward positive infinity to the specified scale
// (immutable operation without error). Returns invalid Money on error.
// Uses value receiver for immutable operation.
func (m Money) Ceiled(scale int) Money {
	result, _ := m.CeiledErr(scale)
	return result
}

// Floor rounds the Money toward negative infinity to the specified scale (mutable operation).
// Mathematical floor function: truncates for positive numbers, always rounds down for negative numbers.
// Uses pointer receiver for mutable operation.
func (m *Money) Floor(scale int) error {
	// Check if Money is invalid
	if m.IsInvalid() {
		return ErrMoneyInvalid
	}

	// For floor, we need to determine the correct rounding strategy based on sign
	if m.amount.Sign() >= 0 {
		// Positive or zero: use RoundDown (toward zero, which is toward negative infinity)
		m.amount.Round(zerorat.RoundDown, scale)
	} else {
		// Negative: use RoundUp (away from zero, which is toward negative infinity for negatives)
		m.amount.Round(zerorat.RoundUp, scale)
	}

	// Check if Rat operation resulted in invalid state
	if m.amount.IsInvalid() {
		m.Invalidate()
		return ErrMoneyInvalid
	}

	return nil
}

// FlooredErr returns a new Money rounded toward negative infinity to the specified scale
// (immutable operation with error). Uses value receiver for immutable operation.
func (m Money) FlooredErr(scale int) (Money, error) {
	result := m // copy
	err := result.Floor(scale)
	return result, err
}

// Floored returns a new Money rounded toward negative infinity to the specified scale
// (immutable operation without error). Returns invalid Money on error.
// Uses value receiver for immutable operation.
func (m Money) Floored(scale int) Money {
	result, _ := m.FlooredErr(scale)
	return result
}
