# ZeroRat

[![CI](https://github.com/n-r-w/zerorat/actions/workflows/ci.yml/badge.svg)](https://github.com/n-r-w/zerorat/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/n-r-w/zerorat/branch/main/graph/badge.svg)](https://codecov.io/gh/n-r-w/zerorat)
[![Go Report Card](https://goreportcard.com/badge/github.com/n-r-w/zerorat)](https://goreportcard.com/report/github.com/n-r-w/zerorat)
[![GoDoc](https://pkg.go.dev/badge/github.com/n-r-w/zerorat)](https://pkg.go.dev/github.com/n-r-w/zerorat)

Zero-allocation rational number library for Go, providing big.Rat-like functionality without heap allocations.

WARNING: This library is still in alpha and should not be used in production.

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
// Construction
a := New(3, 4)        // 3/4
b := NewFromInt(5)    // 5/1
c, err := NewFromFloat64(0.125) // exact float -> 1/8
if err != nil {
    // handle ErrNonFiniteFloat / ErrNotRepresentable
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
valid := a.IsValid()     // true/false
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
fmt.Println(price.String())                   // "USD 12.99"
```

## Validation support

zerorat and money packages support validation via github.com/go-playground/validator/v10.
There are two validation functions available:
- zerorat.RegisterValidationFunc
- money.RegisterValidationFunc

## Performance vs big.Rat

| Operation | ZeroRat | big.Rat | Speedup |
|-----------|---------|---------|---------|
| Construction | 0.24 ns | 60 ns | **254x** |
| Addition | 4.3 ns | 97 ns | **23x** |
| Multiplication | 2.5 ns | 68 ns | **27x** |
| Division | 3.0 ns | 68 ns | **23x** |
| Comparison | 2.7 ns | 32 ns | **12x** |
| Complex Expression | 13.9 ns | 303 ns | **22x** |
| Array Operations (100 items) | 0.19 μs | 34 μs | **172x** |

**Memory Allocations:**
- ZeroRat: **0 allocs/op** for all arithmetic operations
- big.Rat: 2-20 allocs/op depending on operation complexity

## Limitations

- Numbers must fit within int64 (numerator) and uint64 (denominator) ranges
- No arbitrary precision support
- Arithmetic overflow results in invalid state rather than expanding precision
- Exact float conversion returns `ErrNonFiniteFloat` or `ErrNotRepresentable` instead of changing the value silently
- Approximate float conversion is opt-in through `NewApproxFromFloat64` / `NewApproxFromFloat32` and rounds to the nearest representable Rat on the `1/2^63` grid

## Use Cases

Perfect for:
- High-frequency trading systems
- Scientific computing with rational arithmetic  
- Game engines using rational coordinates
- Financial calculations requiring exact fractions
- Any performance-critical rational number operations

## License

MIT License
