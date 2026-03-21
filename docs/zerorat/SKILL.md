---
name: zerorat
description:  Guidelines for using `github.com/n-r-w/zerorat` package.
---

# ZeroRat usage guide

## When to use this skill

Use this skill when the user wants to:
- add exact rational arithmetic to a Go project;
- store ratios, percentages, rates, or fractions without float drift;
- represent money with currency safety;
- parse or format `zerorat` / `money` values;
- choose between mutable and immutable APIs;
- integrate the package with `validator/v10`.

## Typical imports

```go
import (
    "github.com/n-r-w/zerorat"
    "github.com/n-r-w/zerorat/money"
)
```

## Choose the right type

Use `zerorat.Rat` for:
- ratios, percentages, discounts, rates, proportions, and general exact fractions;
- code that needs exact rational math without `big.Rat` allocations;
- cases where currency is not part of the type.

Use `money.Money` for:
- amounts that must carry currency together with the value;
- operations that should reject currency mismatches;
- money-plus-rate workflows where you combine `Money` with `Rat`.

## Construct values correctly

### `zerorat.Rat`
- `Rat{}` is invalid. Do not use the zero-value struct as numeric zero.
- Use `zerorat.Zero()` for valid zero and `zerorat.One()` for valid one.
- Use `zerorat.New(numerator, denominator)` for exact fractions.
- Use `zerorat.NewFromInt64` / `zerorat.NewFromInt` for whole numbers.
- Use `zerorat.NewFromDecimalString` when the input is an exact decimal literal or scientific notation string such as `"0.35"` or `"3.5e-1"`.
- Use `zerorat.NewFromBigRat` when the caller already has a `*big.Rat` and wants an exact `Rat` if it fits the package limits.
- Use `zerorat.NewFromFloat64` / `zerorat.NewFromFloat32` when you need the exact IEEE-754 value.
- Use `zerorat.NewApproxFromFloat64` / `zerorat.NewApproxFromFloat32` only when approximation is acceptable and explicit.
- Pointer helper constructors such as `zerorat.NewFromInt64Ptr`, `zerorat.NewFromIntPtr`, `zerorat.NewFromFloat64Ptr`, and `zerorat.NewFromFloat32Ptr` return an invalid `Rat{}` on `nil` input.

### `money.Money`
- `Money{}` is invalid. Use `money.ZeroMoney(currency)` for a valid zero amount.
- Use `money.NewMoney(currency, amount)` when you already have a `zerorat.Rat`.
- Use `money.NewMoneyInt(currency, value)` when you want to store an integer amount exactly as given.
- Use `money.NewMoneyFromFraction(numerator, denominator, currency)` for exact fractional money values.
- Use `money.NewMoneyFloat(currency, value)` only when exact IEEE-754 semantics are desired.
- If your business logic uses decimal minor units such as cents, model that explicitly in your app. `NewMoneyFloat` does not normalize decimal inputs, and `Money.String()` is not a decimal currency formatter.

## Mutable vs immutable APIs

The package exposes both styles on purpose.

- Mutable methods use pointer receivers and change the receiver in place:
  - `Add`, `Sub`, `Mul`, `Div`, `Round`, `ScaleUp`, `ScaleDown`, `Invalidate`
- Immutable methods use value receivers and return a new value:
  - `Added`, `Subtracted`, `Multiplied`, `Divided`, `Rounded`, `ScaledUp`, `ScaledDown`
- Error-returning immutable variants add `Err`.

Prefer mutable methods when the caller clearly wants in-place updates. Prefer immutable methods when chaining or preserving the original value matters.

## Arithmetic, reduction, and error semantics

### `Rat`
- `New` reduces immediately, but arithmetic methods do not auto-reduce.
- Call `Reduce()` / `Reduced()` only when lowest terms are actually needed.
- Overflow or invalid operations usually produce an invalid value instead of widening precision.
- Non-`Err` comparisons return neutral values for invalid inputs:
  - `Compare()` returns `0`
  - `Equal`, `Less`, `Greater` return `false`
- Use `*Err` variants when the caller must distinguish invalid input from a valid comparison result.

### `Money`
- Money-vs-money operations require matching currencies.
- Mutable money methods invalidate the receiver on mismatch, overflow, or invalid input.
- Use `Rat` overloads such as `MulRat`, `DivRat`, `AddRat`, `SubRat` for rate-based math.
- Use `Percent`, `PercentInt`, `PercentMoney`, and `Profit` helpers when the task is expressed in those terms.

## Formatting and parsing

- `Rat.String()` returns rational text such as `23/4`.
- `Rat.ToDecimalString()` returns an exact finite decimal string such as `0.35` or `12.5`.
- `Rat.ToDecimalString()` returns `ErrNonTerminatingDecimal` for values like `1/3` and `ErrInvalid` for invalid `Rat` values.
- `Rat.ToBigRatErr()` returns an exact `*big.Rat` for valid values and `ErrInvalid` for invalid `Rat` values.
- `Money.String()` returns slash-separated rational text such as:
  - `USD/123/100`
  - `GBP/42`
  - `invalid`
- `money.ParseMoney()` expects `currency/numerator/denominator` or `currency/numerator`.
- If the user needs human-facing money display like `USD 12.99`, format it separately in application code. Do not rely on `Money.String()` for presentation formatting.

## Validation support

Both packages integrate with `github.com/go-playground/validator/v10`.

- Use `zerorat.RegisterValidationFunc(v)` for `zerorat.Rat` fields.
- Use `money.RegisterValidationFunc(v)` for `money.Money` fields.

This makes validator treat package values as valid or invalid based on their built-in `IsValid()` rules.

## Recommended decision rules

When helping a user:
1. If the value is a general fraction or rate, start with `zerorat.Rat`.
2. If the value must carry currency, use `money.Money`.
3. If the input is a decimal literal and exact decimal semantics matter, prefer integer units or `NewMoneyFromFraction` over `NewMoneyFloat`.
4. If the caller may receive invalid input or currency mismatches, prefer `*Err` variants.
5. If the code starts from a zero value, replace it with `zerorat.Zero()` or `money.ZeroMoney(currency)` as appropriate.

## Example patterns

```go
import "math/big"

// Exact ratio.
discountRate := zerorat.New(15, 100)

// Exact decimal input.
taxRate, err := zerorat.NewFromDecimalString("7.5e-2")
if err != nil {
  return err
}

// Exact interop with math/big.
ratioFromBig, err := zerorat.NewFromBigRat(big.NewRat(7, 20))
if err != nil {
  return err
}

// Exact money amount from a fraction.
price := money.NewMoneyFromFraction(1299, 100, "USD")

// Immutable money-plus-rate operation.
discountValue := price.MultipliedRat(discountRate)

// Validation-aware construction from floats.
rate, err := zerorat.NewFromFloat64(0.125)
if err != nil {
    return err
}

taxValue := price.MultipliedRat(rate)
if taxValue.IsInvalid() {
  // handle invalid result
}

taxRateText, err := taxRate.ToDecimalString()
if err != nil {
    return err
}

_ = taxRateText // "0.075"

ratioBig, err := ratioFromBig.ToBigRatErr()
if err != nil {
  return err
}

_ = ratioBig
```

## Common pitfalls

- Treating `Rat{}` or `Money{}` as valid zero.
- Assuming every `Rat` can be formatted as a finite decimal string.
- Assuming every `*big.Rat` can be converted into `Rat` without checking package limits.
- Assuming arithmetic auto-reduces.
- Using `NewMoneyFloat` for decimal currency semantics.
- Assuming `Money.String()` is user-facing display formatting.
- Forgetting that money-to-money operations require matching currencies.
- Ignoring `IsValid()` after operations that may invalidate a value.
