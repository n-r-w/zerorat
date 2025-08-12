package money

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// invalidMoneyString is the string representation for invalid Money.
	invalidMoneyString = "invalid"
	// minParseParts is the minimum number of parts required for parsing.
	minParseParts = 2
	// maxParseParts is the maximum number of parts allowed for parsing.
	maxParseParts = 3
	// fractionParts is the number of parts for full fraction format.
	fractionParts = 3
)

// String returns string representation of Money.
// Format: "currency/amount" where amount uses zerorat.Rat.String() format.
// Returns "invalid" for invalid Money.
// Uses value receiver as this is an immutable operation.
func (m Money) String() string {
	if m.IsInvalid() {
		return invalidMoneyString
	}

	// Use the underlying Rat's string representation for the amount part
	amountStr := m.amount.String()

	// Format: currency/amount (using whatever format Rat.String() produces)
	return fmt.Sprintf("%s/%s", m.currency, amountStr)
}

// ParseMoney parses a string representation of Money.
// Expected format: "currency/numerator/denominator" or "currency/numerator" (implies denominator=1).
// Returns error for invalid format, empty currency, or invalid rational number.
func ParseMoney(s string) (Money, error) {
	if s == "" {
		return Money{}, errors.New("empty string")
	}

	if s == invalidMoneyString {
		return Money{}, errors.New("invalid money string")
	}

	// Split by '/' to get currency and amount parts
	parts := strings.Split(s, "/")

	// We need at least currency/numerator, optionally currency/numerator/denominator
	if len(parts) < minParseParts {
		return Money{}, errors.New("invalid format: missing currency or amount")
	}

	if len(parts) > maxParseParts {
		return Money{}, errors.New("invalid format: too many parts")
	}

	// Extract currency
	currency := parts[0]
	if currency == "" {
		return Money{}, errors.New("invalid format: empty currency")
	}

	// Parse numerator
	numeratorStr := parts[1]
	numerator, err := strconv.ParseInt(numeratorStr, 10, 64)
	if err != nil {
		return Money{}, fmt.Errorf("invalid numerator: %w", err)
	}

	// Parse denominator (default to 1 if not provided)
	var denominator uint64 = 1
	if len(parts) == fractionParts {
		denominatorStr := parts[2]
		denominatorParsed, parseErr := strconv.ParseUint(denominatorStr, 10, 64)
		if parseErr != nil {
			return Money{}, fmt.Errorf("invalid denominator: %w", parseErr)
		}
		denominator = denominatorParsed
	}

	// Create the Money using the fraction constructor
	money := NewMoneyFromFraction(numerator, denominator, currency)

	// Check if the resulting Money is valid
	if money.IsInvalid() {
		return Money{}, errors.New("invalid money: fraction or currency invalid")
	}

	return money, nil
}
