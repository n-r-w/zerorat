package money

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

// hasSameCurrency checks if two Money values have the same currency.
// Returns true only if both Money values are valid and have matching currencies.
// This is a centralized helper function for currency mismatch checks.
func hasSameCurrency(a, b Money) bool {
	return a.IsValid() && b.IsValid() && a.currency == b.currency
}

// RegisterValidationFunc registers a custom validation function for Money types.
func RegisterValidationFunc(v *validator.Validate) {
	v.RegisterCustomTypeFunc(func(field reflect.Value) any {
		if m, ok := field.Interface().(Money); ok {
			return m.IsValid()
		}
		return nil
	}, Money{})
}
