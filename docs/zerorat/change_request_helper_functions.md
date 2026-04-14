# Change Request: Add small helper functions for common exact-math combinations

This change request adds a small set of public helpers to reduce repeated call-site boilerplate in projects that already use `zerorat` and `money` as low-level arithmetic primitives. The goal is to remove repeated multi-step exact-math combinations without moving pricing policy, FX policy, VAT policy, or application DTO logic into the library.

## Scope
In scope:
- ISP-001: Add one public helper in the root `zerorat` package that builds an exact multiplicative factor from a percentage or rate using the formula $1 + p / 100$.
- ISP-002: Add one public helper in `money` that converts `Money` into a target currency when the caller already has an explicit exact conversion ratio.
- ISP-003: Add one public helper in `money` that redistributes rounding across an ordered list of same-currency values by cumulative rounding deltas.
- ISP-004: Add tests and user-facing documentation for the new helpers.

Out of scope:
- OSP-001: FX pair discovery from quote objects, hub-currency logic, or pricing-specific FX modes.
- OSP-002: VAT mode handling, price restriction rules, pricing pipeline state mutation, or response DTO materialization.
- OSP-003: Any helper that depends on external application structs, enums, feature flags, or business policy objects.

## Requested Changes
- FRQ-001: WHEN the caller has an exact percentage or rate as `Rat`, THE `zerorat` package SHALL provide a public helper that returns the exact multiplicative factor $1 + p / 100$.
- FRQ-002: WHEN the caller has a valid `Money` value, a target currency code, and an explicit exact conversion ratio, THE `money` package SHALL provide a public helper that returns the converted `Money` value without requiring the caller to unwrap `Amount()`, multiply manually, and re-wrap with `NewMoney`.
- FRQ-003: WHEN the caller provides an ordered collection of same-currency `Money` values, a `RoundType`, and a rounding scale, THE `money` package SHALL provide a public helper that rounds cumulative totals and returns per-item deltas so that the returned values preserve order and sum to the rounded total.
- FRQ-004: WHEN any of the new helpers receive invalid input, THE package SHALL follow the existing `zerorat` and `money` invalid-state and error conventions instead of introducing a new error model.
- FRQ-005: WHEN these helpers are added, THE public API SHALL remain math-centric and SHALL NOT encode pricing, VAT, FX-policy, or response-shape semantics.

Proposed public signatures:
- APC-001: `func FactorFromPercent(percent Rat) Rat`
- APC-002: `func (m *Money) Convert(targetCurrency Currency, rate zerorat.Rat) error`
- APC-003: `func (m Money) ConvertedErr(targetCurrency Currency, rate zerorat.Rat) (Money, error)`
- APC-004: `func (m Money) Converted(targetCurrency Currency, rate zerorat.Rat) Money`
- APC-005: `func RedistributeRoundedErr(values []Money, roundType zerorat.RoundType, scale int) ([]Money, error)`
- APC-006: `func RedistributeRounded(values []Money, roundType zerorat.RoundType, scale int) []Money`

Architecture and naming rationale:
- DEC-001: `FactorFromPercent` is a root-package constructor-style helper because it creates a `Rat` value from another `Rat` value, similar in role to `Zero()` and `One()`.
- DEC-002: `Convert` / `ConvertedErr` / `Converted` are instance methods because the operation transforms one `Money` value and matches the existing mutable and immutable method pairs such as `Round` / `RoundedErr` / `Rounded` and `MulRat` / `MultipliedRatErr` / `MultipliedRat`.
- DEC-003: `RedistributeRoundedErr` / `RedistributeRounded` are package-level helpers because they operate on many `Money` values at once, which matches the existing package-level multi-value API style such as `SumErr` / `Sum`.
- DEC-004: The proposed names are neutral and math-centric. They do not mention VAT, FX mode, restriction rules, supplier logic, or any `calc-meal` terminology.
- DEC-005: `FactorFromPercent` returns a multiplicative coefficient, not a percentage value. Example: `FactorFromPercent(20)` returns `6/5`, so `amount.MultipliedRat(FactorFromPercent(20))` means “increase amount by 20%”.
- DEC-006: `scale` in `RedistributeRoundedErr` / `RedistributeRounded` uses the same meaning as `Money.Round` and `Money.Rounded`: `0` rounds to integers, positive values round to decimal places, and negative values round to powers of ten.
- DEC-007: Example: with `RoundHalfUp` and `scale = 0`, values `10.4` and `10.4` are redistributed as `10` and `11`, because the helper rounds cumulative totals `10.4 -> 10` and `20.8 -> 21`, then returns the deltas.
- DEC-008: `RedistributeRounded` works in four steps: accumulate values in input order, round each cumulative total, subtract the previous rounded cumulative total, and emit that delta as the rounded value for the current item.
- DEC-009: The key invariant is: the sum of returned values is equal to rounding the total once, while each item keeps its original position in the sequence.

## Affected Areas
- CMP-001: `zerorat` root package - new `Rat` helper for percent-factor construction.
- CMP-002: `money` package - new conversion and rounding-redistribution helpers.
- CMP-003: [README.md](README.md) - document the new helpers and their intended boundaries.
- CMP-004: [docs/zerorat/SKILL.md](docs/zerorat/SKILL.md) - describe when to use the new helpers.
- CMP-005: Unit tests - verify valid paths, invalid inputs, same-currency rules, and rounding invariants.

## Constraints and Risks
- CNS-001: Names and signatures must stay neutral and reusable outside pricing domains.
- CNS-002: The conversion helper must apply an already-known ratio only. It must not discover or compose FX rates.
- CNS-003: The rounding redistribution helper must make order sensitivity, same-currency requirements, and scale explicit.
- CNS-004: Single-value conversion must follow the existing mutable and immutable method pattern of the `Money` API.
- CNS-005: Multi-value redistribution must stay a package-level helper and must not depend on a pricing-state container.
- RSK-001: A vague name could make the percent-factor helper look domain-specific or redundant.
- RSK-002: A poorly scoped conversion helper could leak FX policy into the library API.
- RSK-003: A redistribution helper without a precise invariant could create surprising rounding behavior for callers.

## Acceptance Criteria
- ACC-001: The library exposes one public `Rat` helper named `FactorFromPercent` for exact percent-factor construction.
- ACC-002: The library exposes `Convert`, `ConvertedErr`, and `Converted` on `Money` for conversion by an explicit exact ratio into a target currency.
- ACC-003: The library exposes `RedistributeRoundedErr` and `RedistributeRounded` in `money` with documented same-currency, order, and scale guarantees.
- ACC-004: Tests cover valid cases, invalid inputs, and the key redistribution invariant: returned items sum to the rounded total.
- ACC-005: Documentation explains the helpers and explicitly states that FX-policy logic, VAT-policy logic, and pricing-specific orchestration remain outside the library.

## Assumptions
- ASM-001: The implementation should follow the current mutable and immutable API style already used in `zerorat` and `money`.
- ASM-002: `FactorFromPercent`, `Convert`, `ConvertedErr`, `Converted`, `RedistributeRoundedErr`, and `RedistributeRounded` are the preferred public names unless repository maintainers identify a stronger conflict with existing naming conventions during implementation.

## Open Questions
- QST-001: Should `Money.Convert` short-circuit same-currency conversion before validating `rate`, or should it validate `rate` even when no conversion is needed? This affects edge-case behavior but does not change the public signature.
- QST-002: Should `RedistributeRoundedErr` return `nil` or an empty slice for empty input? This affects edge-case consistency but does not change the public signature.

## References
- REF-001: [/Users/rvnikulenk/dev/bronevik/calc-meal/internal/usecase/calculatemeals/pricing/prepare.go](/Users/rvnikulenk/dev/bronevik/calc-meal/internal/usecase/calculatemeals/pricing/prepare.go) - repeated construction of the VAT multiplier.
- REF-002: [/Users/rvnikulenk/dev/bronevik/calc-meal/internal/usecase/calculatemeals/pricing/runtime.go](/Users/rvnikulenk/dev/bronevik/calc-meal/internal/usecase/calculatemeals/pricing/runtime.go) - money conversion by explicit ratio and cumulative rounding redistribution.
- REF-003: [/Users/rvnikulenk/dev/bronevik/calc-meal/internal/usecase/calculatemeals/pricing/schema_our_net.go](/Users/rvnikulenk/dev/bronevik/calc-meal/internal/usecase/calculatemeals/pricing/schema_our_net.go) - repeated percent-based pricing combinations.
- REF-004: [/Users/rvnikulenk/dev/bronevik/calc-meal/internal/usecase/calculatemeals/pricing/schema_turnover_based_on_hotel.go](/Users/rvnikulenk/dev/bronevik/calc-meal/internal/usecase/calculatemeals/pricing/schema_turnover_based_on_hotel.go) - repeated percent-based pricing combinations.
- REF-005: [money/arithmetic.go](money/arithmetic.go) - existing `Money.Percented` API that already covers percent-of-money.
- REF-006: [money/arithmetic_rat.go](money/arithmetic_rat.go) - existing multiply and divide operations by `Rat`.
- REF-007: [money/rounding.go](money/rounding.go) - existing rounding primitives that the redistribution helper must complement, not duplicate.
- REF-008: [docs/zerorat/SKILL.md](docs/zerorat/SKILL.md) - current usage guidance for `zerorat` and `money`.
- REF-009: [README.md](README.md) - public library scope and user-facing API overview.
