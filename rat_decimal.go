package zerorat

import (
	"errors"
	"math"
	"math/big"
	"math/bits"
	"regexp"
	"strconv"
	"strings"
)

const (
	decimalBase       = 10
	decimalFactorTwo  = 2
	decimalFactorFive = 5
	decimalMaxDigit   = 9
)

var (
	// ErrInvalidDecimalString indicates that a decimal string cannot be parsed.
	ErrInvalidDecimalString = errors.New("invalid decimal string")
	// ErrNonTerminatingDecimal indicates that a rational number has no finite decimal form.
	ErrNonTerminatingDecimal = errors.New("non-terminating decimal")
	// decimalNumberPattern matches plain decimal and scientific-notation inputs.
	decimalNumberPattern = regexp.MustCompile(`^(?P<int>[-+]?\d*)?(?:\.(?P<dec>\d*))?(?:[Ee](?P<exp>[-+]?\d+))?$`)
)

// NewFromDecimalString creates a rational number from an exact decimal string.
func NewFromDecimalString(s string) (Rat, error) {
	if s == "" {
		return Rat{}, ErrInvalidDecimalString
	}

	matches := decimalNumberPattern.FindStringSubmatch(s)
	if len(matches) == 0 {
		return Rat{}, ErrInvalidDecimalString
	}

	integerPart := matches[1]
	fractionalPart := matches[2]
	exponentPart := matches[3]

	isNegative := false
	if integerPart != "" && (integerPart[0] == '-' || integerPart[0] == '+') {
		isNegative = integerPart[0] == '-'
		integerPart = integerPart[1:]
	}

	if integerPart == "" && fractionalPart == "" {
		return Rat{}, ErrInvalidDecimalString
	}

	scale, err := parseDecimalScale(exponentPart, len(fractionalPart))
	if err != nil {
		return Rat{}, err
	}

	digits := strings.TrimLeft(integerPart+fractionalPart, "0")
	if digits == "" {
		return Zero(), nil
	}

	for scale > 0 && strings.HasSuffix(digits, "0") {
		digits = digits[:len(digits)-1]
		scale--
	}

	magnitude, ok := new(big.Int).SetString(digits, decimalBase)
	if !ok {
		return Rat{}, ErrInvalidDecimalString
	}

	numerator, denominator, err := buildDecimalRatParts(magnitude, scale, isNegative)
	if err != nil {
		return Rat{}, err
	}

	return New(numerator, denominator), nil
}

// ToDecimalString converts a rational number to an exact decimal string.
func (r Rat) ToDecimalString() (string, error) {
	if r.IsInvalid() {
		return "", ErrInvalid
	}

	if r.numerator == 0 {
		return "0", nil
	}

	reduced := r.Reduced()
	twoCount, fiveCount, terminates := countDecimalFactors(reduced.denominator)
	if !terminates {
		return "", ErrNonTerminatingDecimal
	}

	integerPart, remainder := decimalIntegerAndRemainder(reduced)

	var builder strings.Builder
	if reduced.numerator < 0 {
		_ = builder.WriteByte('-')
	}
	_, _ = builder.WriteString(strconv.FormatUint(integerPart, decimalBase))

	if remainder == 0 {
		return builder.String(), nil
	}

	_ = builder.WriteByte('.')
	maxDigits := twoCount
	if fiveCount > maxDigits {
		maxDigits = fiveCount
	}

	for range maxDigits {
		digit, nextRemainder, ok := nextDecimalDigit(remainder, reduced.denominator)
		if !ok {
			return "", ErrNonTerminatingDecimal
		}
		_ = builder.WriteByte(digit)
		remainder = nextRemainder
		if remainder == 0 {
			return builder.String(), nil
		}
	}

	return "", ErrNonTerminatingDecimal
}

// parseDecimalScale converts the exponent segment into a final decimal scale.
func parseDecimalScale(exponentPart string, fractionalDigits int) (int64, error) {
	if exponentPart == "" {
		return int64(fractionalDigits), nil
	}

	exponent, err := strconv.ParseInt(exponentPart, 10, 64)
	if err != nil {
		return 0, ErrNotRepresentable
	}

	scale, ok := checkedInt64Sub(int64(fractionalDigits), exponent)
	if !ok {
		return 0, ErrNotRepresentable
	}

	return scale, nil
}

// buildDecimalRatParts converts a parsed decimal magnitude and scale into Rat fields.
func buildDecimalRatParts(
	magnitude *big.Int,
	scale int64,
	isNegative bool,
) (numerator int64, denominator uint64, err error) {
	absMagnitude := new(big.Int).Set(magnitude)

	if scale <= 0 {
		numerator, err = scaleDecimalMagnitude(absMagnitude, -scale, isNegative)
		if err != nil {
			return 0, 0, err
		}

		return numerator, 1, nil
	}

	twoExp, fiveExp := reduceDecimalScale(absMagnitude, scale)
	denominator, ok := decimalDenominator(twoExp, fiveExp)
	if !ok {
		return 0, 0, ErrNotRepresentable
	}

	numerator, ok = signedBigIntToInt64(absMagnitude, isNegative)
	if !ok {
		return 0, 0, ErrNotRepresentable
	}

	return numerator, denominator, nil
}

// scaleDecimalMagnitude applies a positive decimal exponent while guarding int64 bounds.
func scaleDecimalMagnitude(magnitude *big.Int, factorExp int64, isNegative bool) (int64, error) {
	for range factorExp {
		if wouldDecimalScaleOverflow(magnitude) {
			return 0, ErrNotRepresentable
		}
		magnitude.Mul(magnitude, big.NewInt(decimalBase))
	}

	numerator, ok := signedBigIntToInt64(magnitude, isNegative)
	if !ok {
		return 0, ErrNotRepresentable
	}

	return numerator, nil
}

// reduceDecimalScale removes denominator factors that can be absorbed into the numerator.
func reduceDecimalScale(magnitude *big.Int, scale int64) (twoExp, fiveExp int64) {
	twoExp = scale
	fiveExp = scale

	for twoExp > 0 && divideBigIntIfDivisibleBySmallFactor(magnitude, decimalFactorTwo) {
		twoExp--
	}

	for fiveExp > 0 && divideBigIntIfDivisibleBySmallFactor(magnitude, decimalFactorFive) {
		fiveExp--
	}

	return twoExp, fiveExp
}

// wouldDecimalScaleOverflow reports whether multiplying by ten would exceed Rat numerator bounds.
func wouldDecimalScaleOverflow(magnitude *big.Int) bool {
	limit := new(big.Int).Quo(maxAbsInt64Big(), big.NewInt(decimalBase))
	return magnitude.Cmp(limit) > 0
}

// divideBigIntIfDivisibleWithBigFactor divides a big integer by a big factor only when division is exact.
func divideBigIntIfDivisibleWithBigFactor(value, factor *big.Int) bool {
	var remainder big.Int
	remainder.Mod(value, factor)
	if remainder.Sign() != 0 {
		return false
	}

	value.Quo(value, factor)
	return true
}

// divideBigIntIfDivisibleBySmallFactor divides a big integer by a small factor only when division is exact.
func divideBigIntIfDivisibleBySmallFactor(value *big.Int, factor int64) bool {
	return divideBigIntIfDivisibleWithBigFactor(value, big.NewInt(factor))
}

// decimalDenominator builds a decimal denominator from powers of two and five.
func decimalDenominator(twoExp, fiveExp int64) (uint64, bool) {
	denominator := uint64(1)

	for range twoExp {
		if willOverflowUint64Mul(denominator, decimalFactorTwo) {
			return 0, false
		}
		denominator *= decimalFactorTwo
	}

	for range fiveExp {
		if willOverflowUint64Mul(denominator, decimalFactorFive) {
			return 0, false
		}
		denominator *= decimalFactorFive
	}

	return denominator, true
}

// signedBigIntToInt64 converts a positive magnitude and sign into an int64 numerator.
func signedBigIntToInt64(magnitude *big.Int, isNegative bool) (int64, bool) {
	if isNegative {
		maxAbs := maxAbsInt64Big()
		if magnitude.Cmp(maxAbs) > 0 {
			return 0, false
		}
		if magnitude.Cmp(maxAbs) == 0 {
			return math.MinInt64, true
		}

		return -magnitude.Int64(), true
	}

	if !magnitude.IsInt64() {
		return 0, false
	}

	return magnitude.Int64(), true
}

// countDecimalFactors determines whether a denominator has a finite decimal expansion.
func countDecimalFactors(denominator uint64) (twoCount, fiveCount int, terminates bool) {
	twoCount = 0
	for denominator%decimalFactorTwo == 0 {
		twoCount++
		denominator /= decimalFactorTwo
	}

	fiveCount = 0
	for denominator%decimalFactorFive == 0 {
		fiveCount++
		denominator /= decimalFactorFive
	}

	return twoCount, fiveCount, denominator == 1
}

// decimalIntegerAndRemainder splits a rational into absolute integer and remainder parts.
func decimalIntegerAndRemainder(r Rat) (integerPart, remainder uint64) {
	absNumerator := absInt64ToUint64(r.numerator)
	return absNumerator / r.denominator, absNumerator % r.denominator
}

// nextDecimalDigit advances one exact decimal digit using 128-bit division.
func nextDecimalDigit(remainder, denominator uint64) (digit byte, nextRemainder uint64, ok bool) {
	hi, lo := bits.Mul64(remainder, decimalBase)
	quotient, nextRemainder := bits.Div64(hi, lo, denominator)
	if quotient > decimalMaxDigit {
		return 0, 0, false
	}

	return "0123456789"[quotient], nextRemainder, true
}

// maxAbsInt64Big returns the largest absolute numerator magnitude representable by Rat.
func maxAbsInt64Big() *big.Int {
	return new(big.Int).SetUint64(uint64(math.MaxInt64) + 1)
}

// checkedInt64Sub subtracts two int64 values and reports whether the result fits in int64.
func checkedInt64Sub(left, right int64) (result int64, ok bool) {
	if right > 0 {
		if left < math.MinInt64+right {
			return 0, false
		}
	} else if left > math.MaxInt64+right {
		return 0, false
	}

	return left - right, true
}
