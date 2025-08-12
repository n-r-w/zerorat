package money

// hasSameCurrency checks if two Money values have the same currency.
// Returns true only if both Money values are valid and have matching currencies.
// This is a centralized helper function for currency mismatch checks.
func hasSameCurrency(a, b Money) bool {
	return a.IsValid() && b.IsValid() && a.currency == b.currency
}
