package zerorat

import (
	"errors"
	"math/big"
)

// NewFromBigRat creates a rational number from a big.Rat value.
func NewFromBigRat(value *big.Rat) (Rat, error) {
	if value == nil {
		return Rat{}, errors.New("nil big.Rat")
	}

	numerator := value.Num()
	if !numerator.IsInt64() {
		return Rat{}, ErrNotRepresentable
	}

	denominator := value.Denom()
	if !denominator.IsUint64() {
		return Rat{}, ErrNotRepresentable
	}

	return New(numerator.Int64(), denominator.Uint64()), nil
}

// ToBigRatErr converts a Rat into an exact big.Rat value.
func (r Rat) ToBigRatErr() (*big.Rat, error) {
	if r.IsInvalid() {
		return nil, ErrInvalid
	}

	return new(big.Rat).SetFrac(big.NewInt(r.numerator), new(big.Int).SetUint64(r.denominator)), nil
}
