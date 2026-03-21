# ZeroRat

[![CI](https://github.com/n-r-w/zerorat/actions/workflows/ci.yml/badge.svg)](https://github.com/n-r-w/zerorat/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/n-r-w/zerorat/branch/main/graph/badge.svg)](https://codecov.io/gh/n-r-w/zerorat)
[![Go Report Card](https://goreportcard.com/badge/github.com/n-r-w/zerorat)](https://goreportcard.com/report/github.com/n-r-w/zerorat)
[![GoDoc](https://pkg.go.dev/badge/github.com/n-r-w/zerorat)](https://pkg.go.dev/github.com/n-r-w/zerorat)

Zero-allocation rational number library for Go, providing big.Rat-like functionality without heap allocations.

## Features

- **Zero heap allocations** for arithmetic and comparison operations
- **16-byte struct** with perfect memory alignment
- **Mutable and immutable APIs** for flexible usage
- **Overflow-safe comparisons** using 128-bit arithmetic via math/bits package
- **Overflow detection** with graceful invalid state handling
- **Optional fraction reduction** using greatest common divisor algorithm (call Reduce() when needed)
- **Value semantics** with pointer receivers as requested

## Installation

```bash
go get github.com/n-r-w/zerorat@latest
```

## Basic Usage

```go
import "math/big"

// Construction
a := New(3, 4)        // 3/4
b := NewFromInt(5)    // 5/1
c, err := NewFromFloat64(0.125) // exact float -> 1/8
if err != nil {
    // handle ErrNonFiniteFloat / ErrNotRepresentable
}

decimal, err := NewFromDecimalString("3.5e-1") // exact decimal/scientific notation -> 7/20
if err != nil {
    // handle ErrInvalidDecimalString / ErrNotRepresentable
}

fromBig, err := NewFromBigRat(big.NewRat(7, 20)) // exact big.Rat -> 7/20
if err != nil {
    // handle ErrNilBigRat / ErrNotRepresentable
}

d, err := NewApproxFromFloat64(3.0 / (1 << 64)) // nearest representable Rat on the 1/2^63 grid
if err != nil {
    // handle ErrNonFiniteFloat / ErrNotRepresentable
}

// Mutable operations (results not auto-reduced)
a.Add(b)                 // a = 3/4 + 5/1 = (3*1 + 5*4)/(4*1) = 23/4
a.Add(c)                 // explicit conversion first, then exact arithmetic
a.Add(d)                 // approximate path is opt-in and named explicitly
a.Reduce()               // a = 23/4 (manually reduce if needed)

// Immutable operations  
result := a.Added(b)     // returns new value, a unchanged

// Comparisons
equal := a.Equal(b)      // true/false
less := a.Less(b)        // true/false

// Utilities
str := a.String()        // "23/4"
decimalStr, err := decimal.ToDecimalString() // "0.35"
if err != nil {
    // handle ErrInvalid / ErrNonTerminatingDecimal
}
backToBig, err := fromBig.ToBigRatErr() // exact Rat -> *big.Rat
if err != nil {
    // handle ErrInvalid
}
valid := a.IsValid()     // true/false

_ = backToBig
```

## JSON support

`Rat` can be used directly with `encoding/json`.

- `MarshalJSON` writes a JSON number without quotes when the value has an exact finite decimal form
- `MarshalJSON` returns `ErrInvalid` for invalid `Rat` values
- `MarshalJSON` returns `ErrNonTerminatingDecimal` for values like `1/3` that cannot be written as an exact JSON number
- `UnmarshalJSON` accepts JSON numbers such as `0.35` and `3.5e-1`
- `UnmarshalJSON` also accepts quoted decimal text such as `"0.35"` and `"3.5e-1"`
- `UnmarshalJSON` rejects quoted rational text such as `"1/3"`; quoted input must still match `NewFromDecimalString`
- `UnmarshalJSON` treats `null` as an invalid `Rat` state for value fields
- Failed `UnmarshalJSON` calls also leave the receiver invalid, so bad payloads do not preserve stale values
- For optional fields, prefer `*Rat`; JSON `null` maps to `nil` by standard `encoding/json` behavior

```go
type payload struct {
    Rate zerorat.Rat `json:"rate"`
}

var p payload

err := json.Unmarshal([]byte(`{"rate":"3.5e-1"}`), &p)
if err != nil {
    // handle invalid decimal input
}

data, err := json.Marshal(payload{Rate: zerorat.New(7, 20)})
if err != nil {
    // handle invalid or non-terminating values like 1/3
}

_ = data // {"rate":0.35}
```

## Money package

Currency-aware monetary calculations built on top of ZeroRat. Provides type-safe money operations with automatic currency validation, rounding modes, and formatting support. Ensures operations only occur between compatible currencies while maintaining ZeroRat's zero-allocation performance characteristics.

```go
import (
    "github.com/n-r-w/zerorat"
    "github.com/n-r-w/zerorat/money"
)

// Construction
price := money.NewMoneyInt("USD", 1299)        // $12.99 (in cents)
tax, err := money.NewMoneyFloat("USD", 0.85) // exact float input
discount := money.NewMoney("USD", zerorat.New(15, 100)) // 15% as fraction
discount1 := money.NewMoneyFromFraction("USD", 20, 100) // 20% as fraction

if err != nil {
    // handle ErrNonFiniteFloat / ErrNotRepresentable
}

// NewMoneyFloat preserves the exact IEEE-754 binary value.
// For decimal money semantics, prefer integer minor units or explicit fractions.

// Mutable operations (error handling for currency mismatches)
err := price.Add(tax)                          // price = $13.84
err = price.PercentInt(10)                     // 10% of price

// Immutable operations
total := price.Added(tax)                      // returns new Money, price unchanged
discountRate, err := zerorat.NewFromFloat64(0.85)
afterDiscount := total.MultipliedRat(discountRate)  // 15% discount

// Currency validation
usd := money.NewMoneyInt("USD", 100)
eur := money.NewMoneyInt("EUR", 85)
err := usd.Add(eur)                           // returns ErrMoneyCurrencyMismatch

// Formatting
display := money.NewMoneyFromFraction(1299, 100, "USD")
fmt.Println(display.String())                 // "USD/1299/100"
```

## Validation support

zerorat and money packages support validation via github.com/go-playground/validator/v10.
There are two validation functions available:
- zerorat.RegisterValidationFunc
- money.RegisterValidationFunc

## Current benchmark snapshot

Benchmarks below were run in this session on **darwin/arm64 (Apple M4 Max)** with `go test -run '^$' -bench ... -benchmem ./...`.

| Benchmark | Time | Memory | Allocations |
|-----------|------|--------|-------------|
| `BenchmarkZeroRat_NewFromFloat64Exact` | 3.160 ns/op | 0 B/op | 0 allocs/op |
| `BenchmarkZeroRat_NewApproxFromFloat64` | 4.659 ns/op | 0 B/op | 0 allocs/op |
| `BenchmarkZeroRat_NewFromFloat32Exact` | 3.672 ns/op | 0 B/op | 0 allocs/op |
| `BenchmarkNewMoneyFloat` | 4.364 ns/op | 0 B/op | 0 allocs/op |
| `BenchmarkMoneyMulRat` | 4.034 ns/op | 0 B/op | 0 allocs/op |
| `BenchmarkMoneyDivRat` | 4.835 ns/op | 0 B/op | 0 allocs/op |

## Limitations

- Numbers must fit within int64 (numerator) and uint64 (denominator) ranges
- No arbitrary precision support
- Arithmetic overflow results in invalid state rather than expanding precision
- Exact float conversion returns `ErrNonFiniteFloat` or `ErrNotRepresentable` instead of changing the value silently
- `NewFromDecimalString` returns `ErrInvalidDecimalString` for malformed input and `ErrNotRepresentable` when the exact value does not fit in `Rat`
- `ToDecimalString` returns `ErrNonTerminatingDecimal` for values that do not have a finite decimal form, such as `1/3`
- Approximate float conversion is opt-in through `NewApproxFromFloat64` / `NewApproxFromFloat32` and rounds to the nearest representable Rat on the `1/2^63` grid

## Use Cases

- High-frequency trading systems
- Scientific computing with rational arithmetic  
- Game engines using rational coordinates
- Financial calculations requiring exact fractions
- Any performance-critical rational number operations

## Agents Skill

[ZeroRat Agents Skill](docs/zerorat/SKILL.md)