# Review Results: `NewFromFloat64` and float conversion path
Review scope: `NewFromFloat64`, `float64ToRatExact`, direct float arithmetic entry points, and related tests. Out of scope: unrelated rational operations.

## Summary
Outcome: 🟥 Reject

- FND-001: the positive-exponent overflow path returns a valid but numerically wrong rational instead of invalid overflow.
- FND-002: the negative-exponent path silently rounds exact finite `float64` values to zero or to a nearby rational while still marking them valid.
- FND-003: the current tests permit ambiguous outcomes, so the suite stays green while the documented contract is broken.

## Findings

### Blocker
- **⛔ FND-001**
  Severity: Blocker.
  Location: [helper contract](rat_float_convert.go#L8-L10), [positive-exponent branch](rat_float_convert.go#L52-L99), [public constructor call site](rat_construct.go#L93-L98), [float arithmetic entry points](rat_arithmetic.go#L301-L342), [money float entry points](money/arithmetic.go#L172-L173), [money float arithmetic](money/arithmetic.go#L262-L518).
  Issue: when `e >= 0` and the fast path no longer fits, the helper computes `newShift := e - neededDenPow` and returns `(mant << newShift) / 2^neededDenPow`. If $d = neededDenPow$, the returned value is $mant \times 2^{e-2d}$, while the original float value is $mant \times 2^e$. These are equal only when $d = 0`, which means the offload branch changes the number instead of preserving it.
  Impact: large finite floats near or above the `int64` boundary can be silently corrupted and still look valid to callers. The problem is not limited to `NewFromFloat64`; it also affects `AddFloat`, `SubFloat`, `MulFloat`, `DivFloat`, and the `money` package float APIs that depend on the same conversion path.
  Recommendation: do not "offload" exponent bits into the denominator for an integer-valued input. If the exact value does not fit the `(int64 numerator, uint64 denominator)` model, return invalid overflow, or redefine the public contract and helper naming so approximation is explicit.
  Verification: code inspection of the branch arithmetic is sufficient to prove the mismatch.

- **⛔ FND-002**
  Severity: Blocker.
  Location: [helper contract](rat_float_convert.go#L8-L10), [negative-exponent branch](rat_float_convert.go#L103-L133), [public constructor comments](rat_construct.go#L79-L98).
  Issue: the branch starts with `denPow := min(-e, 63)`, which caps the denominator at $2^{63}$ and rounds the numerator. That makes exact conversion impossible for many exactly representable floats whose exact rational needs a larger power-of-two denominator. The behavior can be derived directly from the code: `2^-64` becomes `0 / 2^63`, then reduces to `0/1`; `3 * 2^-64` becomes `2 / 2^63`, then reduces to `1 / 2^62`.
  Impact: tiny non-zero finite inputs can become zero or a different nearby value while `IsValid()` still returns true. Downstream arithmetic has no way to detect that corruption.
  Recommendation: choose one contract and implement it consistently. Either make the helper exact-or-invalid, or document it as an approximation routine and rename it accordingly. The current code and comments promise one behavior and implement another.
  Verification: code inspection of the denominator cap and rounding logic is sufficient to prove the loss of exactness.

### Major
- **⚠️ FND-003**
  Severity: Major.
  Location: [permissive overflow test](rat_float_convert_test.go#L64-L67), [discarded assertion](rat_float_convert_test.go#L81-L83), [tautological assertion](rat_float_convert_test.go#L260-L264), [smallest non-zero float edge case](rat_construct_test.go#L734-L749).
  Issue: several tests validate execution paths rather than behavior. Examples include accepting either valid or invalid results, computing a value and not asserting on it, and using a tautological check that is always true for any `int64`.
  Impact: the current suite passes even though the float conversion contract is broken. That makes regressions likely in the most branch-heavy part of the package.
  Recommendation: replace permissive assertions with behavior assertions. Add exact or invalid expectations around `2^63`, `-2^63`, `2^-64`, `3 * 2^-64`, and public wrappers that delegate to `NewFromFloat64`.
  Verification: the repository test suite passes with `go test ./...`, which shows the current tests do not catch the contract failures above.

## Assumptions
- **ASM-001:** The intended contract is exact conversion or invalid overflow. This is based on the helper name and comments in [rat_float_convert.go](rat_float_convert.go#L8-L10) and the constructor comments in [rat_construct.go](rat_construct.go#L79-L98).

## Open Questions
- **QST-001:** Should float conversion be exact-or-invalid, or approximate-within-a-documented-bound? This decision is required before implementation work starts, because the algorithm, naming, comments, and tests must all follow the same contract.

## Next Steps
- **NXT-001:** Decide the float-conversion contract that the package should expose. This resolves QST-001.
- **NXT-002:** Align the implementation with that contract in `float64ToRatExact` and all public callers.
- **NXT-003:** Rewrite the float-conversion tests so they assert the chosen contract on boundary values and wrapper APIs.

## References
- **REF-001:** [rat_construct.go](rat_construct.go#L79-L121) - public float constructors and documented behavior.
- **REF-002:** [rat_float_convert.go](rat_float_convert.go#L8-L133) - internal float-to-rational conversion algorithm.
- **REF-003:** [rat_arithmetic.go](rat_arithmetic.go#L301-L342) - direct float arithmetic entry points.
- **REF-004:** [money/arithmetic.go](money/arithmetic.go#L172-L173) - percentage path that uses `NewFromFloat64`.
- **REF-005:** [money/arithmetic.go](money/arithmetic.go#L262-L518) - `money` float arithmetic entry points.
- **REF-006:** [rat_float_convert_test.go](rat_float_convert_test.go#L55-L359) - helper-focused float conversion tests.
- **REF-007:** [rat_construct_test.go](rat_construct_test.go#L485-L752) - public constructor float tests.
- **REF-008:** `go test ./...` - current test suite passes despite the findings.
