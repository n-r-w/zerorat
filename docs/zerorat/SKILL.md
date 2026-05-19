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
- Use `zerorat.FactorFromPercent(percent)` when the caller already has a percentage and needs the exact multiplicative factor $1 + p / 100$. Example: `20 -> 6/5`.
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
  - `Rat`: `Add`, `Sub`, `Mul`, `Div`, `Round`, `ScaleUp`, `ScaleDown`, `Invalidate`
  - `Money`: `Add`, `Sub`, `Profit`, `Percent`, `PercentInt`, `PercentMoney`, `MulRat`, `DivRat`, `Convert`, `MulInt`, `DivInt`, `Round`, `ScaleUp`, `ScaleDown`, `Invalidate`
- Immutable methods use value receivers and return a new value:
  - `Rat`: `Added`, `Subtracted`, `Multiplied`, `Divided`, `Rounded`, `ScaledUp`, `ScaledDown`
  - `Money`: `Added`, `Subtracted`, `Profited`, `Percented`, `PercentedInt`, `PercentedMoney`, `MultipliedRat`, `DividedRat`, `Converted`, `MultipliedInt`, `DividedInt`, `Rounded`, `ScaledUp`, `ScaledDown`
- Error-returning immutable variants add `Err`.
- Package-level multi-value helpers complement methods when the operation works on many values at once, such as `SumErr` / `Sum` and `RedistributeRoundedErr` / `RedistributeRounded`.

Prefer mutable methods when the caller clearly wants in-place updates. Prefer immutable methods when chaining or preserving the original value matters.

## Arithmetic, reduction, and error semantics

### `Rat`
- `New` reduces immediately. `Add`, `Sub`, `Mul`, `Div`, `ScaleUp`, `ScaleDown`, and `Round` keep small results on the fast path and reduce or cancel large intermediates when that avoids preventable overflow.
- Call `Reduce()` / `Reduced()` when lowest terms are required by the caller, not just for overflow safety.
- Unrecoverable overflow, exact results outside `Rat` limits, and invalid operations produce an invalid value instead of widening precision.
- `Rat` can compare directly with primitive values:
  - `CompareInt64`, `EqualInt64`, `LessInt64`, `GreaterInt64`
  - `CompareFloat64`, `EqualFloat64`, `LessFloat64`, `GreaterFloat64`
- Use `ApproxEqualFloat64(value, eps)` when you need tolerance-based comparison in `float64` space. Keep using `EqualFloat64` when you need exact `float64` semantics.
- `float64` comparison methods use the exact IEEE-754 value, just like `NewFromFloat64`. For example, `EqualFloat64(0.5)` can be true, while `EqualFloat64(0.1)` is false for `New(1, 10)`.
- If a `float64` is finite but still outside the `Rat` model, such as `math.MaxFloat64`, the exact `float64` comparison helpers treat it as an invalid operand and return neutral results (`CompareFloat64` returns `0`; `EqualFloat64`, `LessFloat64`, and `GreaterFloat64` return `false`).
- Non-`Err` comparisons return neutral values for invalid inputs:
  - `Compare()` returns `0`
  - `Equal`, `Less`, `Greater` return `false`
- Use `*Err` variants when the caller must distinguish invalid input from a valid comparison result.

### `Money`
- Money methods that accept another `Money` require matching currencies.
- Mutable money methods invalidate the receiver on mismatch, overflow, or invalid input.
- Use `Rat` overloads such as `MulRat`, `DivRat`, `AddRat`, `SubRat` for rate-based math.
- Use `Convert` / `ConvertedErr` / `Converted` when the caller already has an exact ratio and needs a target-currency `Money` value without manual unwrap-and-rewrap code.
- `Convert` is not an FX-policy engine. It applies only the explicit ratio passed by the caller.
- If `Convert` receives the same source and target currency, it returns the original amount unchanged.
- Use `Percent`, `PercentInt`, `PercentMoney`, and `Profit` helpers when the task is expressed in those terms.
- Use `RedistributeRoundedErr` / `RedistributeRounded` when rounded line items must preserve input order and still sum to the rounded total.
- `RedistributeRoundedErr` requires same-currency valid inputs. Uses the same `scale` semantics as `Money.Round`: `0` rounds to integers, positive values round to decimal places, and negative values round to powers of ten. Order-sensitive because it rounds cumulative totals in input order.

## Formatting and parsing

- `Rat.String()` returns rational text such as `23/4`.
- `Rat.ToDecimalString()` returns an exact finite decimal string such as `0.35` or `12.5`.
- `Rat.ToDecimalString()` returns `ErrNonTerminatingDecimal` for values like `1/3` and `ErrInvalid` for invalid `Rat` values.
- `Rat` implements `encoding/json` support directly.
- `Rat.MarshalJSON()` writes an unquoted JSON number only for exact finite decimals and returns `ErrInvalid` or `ErrNonTerminatingDecimal` otherwise.
- `Rat.UnmarshalJSON()` accepts JSON numbers and quoted decimal/scientific strings such as `0.35`, `3.5e-1`, `"0.35"`, and `"3.5e-1"`.
- `Rat.UnmarshalJSON()` rejects quoted rational text such as `"1/3"`; quoted input must still match `NewFromDecimalString`.
- `Rat.UnmarshalJSON()` treats `null` as an invalid state for value fields. For optional fields, prefer `*Rat`, where JSON `null` maps to `nil` by standard `encoding/json` behavior.
- `zerorat.NewFromBigRat()` returns `ErrNilBigRat` for `nil` input and `ErrNotRepresentable` when the exact value does not fit package limits.
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
6. If the caller needs the factor $1 + p / 100$, use `zerorat.FactorFromPercent` instead of reimplementing percent-to-factor math at the call site.
7. If the caller already has the exact conversion ratio, use `Money.Convert` or `Money.ConvertedErr` instead of extracting `Amount()` and rebuilding `Money` manually.
8. If rounded line items must add up to the rounded total, use `money.RedistributeRoundedErr` or `money.RedistributeRounded`.

## Example patterns

```go
import "math/big"

// Exact ratio.
discountRate := zerorat.New(15, 100)

// Exact percent-to-factor helper.
vatFactor := zerorat.FactorFromPercent(zerorat.NewFromInt64(20))

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

// Exact conversion when the ratio is already known.
convertedPrice, err := price.ConvertedErr("EUR", zerorat.New(9, 10))
if err != nil {
  return err
}

// Cumulative redistribution after rounding.
roundedLineItems, err := money.RedistributeRoundedErr([]money.Money{
  money.NewMoneyFromFraction(104, 10, "EUR"),
  money.NewMoneyFromFraction(104, 10, "EUR"),
}, zerorat.RoundHalfUp, 0)
if err != nil {
  return err
}

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

_ = vatFactor
_ = convertedPrice
_ = roundedLineItems
_ = ratioBig
```

## Common pitfalls

- Treating `Rat{}` or `Money{}` as valid zero.
- Confusing a factor like `6/5` with a percentage result like `20/100`.
- Assuming every `Rat` can be formatted as a finite decimal string.
- Assuming every `*big.Rat` can be converted into `Rat` without checking package limits.
- Assuming arithmetic always returns lowest terms.
- Using `NewMoneyFloat` for decimal currency semantics.
- Assuming `Money.String()` is user-facing display formatting.
- Forgetting that money methods taking another `Money` require matching currencies.
- Treating `Convert` as a place to discover or compose FX policy.
- Expecting same-currency `Convert` to rescale the amount; it keeps the original value unchanged.
- Forgetting that `RedistributeRoundedErr` is order-sensitive.
- Ignoring `IsValid()` after operations that may invalidate a value.
