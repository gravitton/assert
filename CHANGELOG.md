# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/assert/compare/v1.5.0...main)
### Added
- `Zero` and `NotZero`, using the type's `IsZero` method when present; a non-nil empty slice or map is not zero
- `EqualUnordered`, comparing arrays or slices as multisets and reporting extra and missing elements
- `HasPrefix` and `HasSuffix` for strings and slices
- `ErrorContains`, matching a substring of the error message
- `Length`, `Empty` and `NotEmpty` accept a pointer to an array

### Changed
- `EqualDelta` and `NotEqualDelta` failures show the delta and the actual difference
- `Greater`, `GreaterOrEqual`, `Less` and `LessOrEqual` name their second argument `bound` and print it as `bound`
- `ErrorAs` fails on an invalid target instead of panicking
- The value under test is named `actual` in every assertion's signature and failure message; `Length` prints the measured length as `length`
- Failure messages: `Matches` quotes the actual string, channels and functions print their address once, and a trailing value after JSON is reported as `unexpected "…" after top-level value`

### Fixed
- `Error`, `NoError` and `ErrorContains` treated a typed nil stored in an `error` as nil, hiding the bug `err != nil` would hit
- `EqualJSON` and `JSON` compared numbers with exponents beyond 1e6 as text and took long to expand large exponents
- `EqualJSON`, `JSON`, `Matches`, `NotMatches` and the error assertions printed values without the 1024-byte cut

## [v1.5.0 (2026-09-17)](https://github.com/gravitton/assert/compare/v1.4.0...v1.5.0)
### Added
- `Nil` and `NotNil`; typed nil pointers, slices, maps, channels, functions and `unsafe.Pointer` count as nil
- `ErrorAs`, mirroring `errors.As`
- `PanicsWith`, matching the panic value by `errors.Is` or deep equality
- `Ordered` constraint; `Greater`, `GreaterOrEqual`, `Less` and `LessOrEqual` accept any `cmp.Ordered` type, strings included

### Changed
- **Breaking:** `Panics` no longer takes an expected value; use `PanicsWith`
- **Breaking:** `Iterable` has no type parameter; `Comparable`, `Reference` and `Iterable` are descriptive only
- `Same` and `NotSame` require the same type and address, and slices the same length and capacity; `unsafe.Pointer` and functions are not references
- `Greater`, `GreaterOrEqual`, `Less` and `LessOrEqual` fail on NaN
- `EqualDelta` and `NotEqualDelta` compare integers exactly, accept `uintptr`, and fail on a negative or NaN delta instead of panicking or never matching
- `EqualJSON` and `JSON` compare numbers exactly instead of through `float64`
- `Contains` and `NotContains` report a mismatched element type as `Should have element of same type`
- Failure messages: no trailing colon, `Should match error`, `Should be equal JSON`, `Should be marshalable`, objects printed like `Equal` with reference addresses, values cut at 1024 bytes
- Documented byte-length semantics of `Length`, `Empty` and `NotEmpty`, and typed nil and untyped constant behaviour of `Equal` and `Contains`

### Fixed
- `NotSame` and `NotContains` returned `true` after reporting an invalid argument
- `Length`, `Empty`, `NotEmpty`, `Contains` and `NotContains` panicked on unsupported types
- `Equal`, `NotEqual`, `Same`, `NotSame` and `Panics` panicked when printing a nil pointer
- `EqualJSON` and `JSON` dropped the custom message prefix
- `Panics` and `NotPanics` missed `panic(nil)` under `GODEBUG=panicnil=1`
- `Error` and `NoError` treat a typed nil error as nil

## [v1.4.0 (2026-08-25)](https://github.com/gravitton/assert/compare/v1.3.0...v1.4.0)
### Changed
- Require Go 1.27

## [v1.3.0 (2026-05-13)](https://github.com/gravitton/assert/compare/v1.2.1...v1.3.0)
### Added
- Added `Panics` and `NotPanics` assert methods for panic detection; `Panics` accepts an optional expected value — `nil` skips value validation, an `error` validates with `errors.Is`, anything else uses deep equality


## [v1.2.1 (2026-04-27)](https://github.com/gravitton/assert/compare/v1.2.0...v1.2.1)
### Changed
- Improve error formatting in `Error` assert methods for better readability


## [v1.2.0 (2026-04-22)](https://github.com/gravitton/assert/compare/v1.1.0...v1.2.0)
### Added
- Added `Greater`, `GreaterOrEqual`, `Less`, and `LessOrEqual` assert methods for numeric comparison
- Added `Empty` and `NotEmpty` assert methods for zero-length checks on strings, slices, maps, and channels
- Added `NotEqualDelta` assert method, complementing `EqualDelta`

### Fixed
- `Error` and `NoError` now correctly treat typed nil errors (e.g. `(*MyError)(nil)`) as nil


## [v1.1.0 (2026-04-22)](https://github.com/gravitton/assert/compare/v1.0.0...v1.1.0)
### Added
- Added `Matches` and `NotMatches` assert methods for regexp matching on strings


## v1.0.0 (2025-10-17)
### Added
- Added new assert methods
  - `True`
  - `False`
  - `Equal`
  - `NotEqual`
  - `EqualDelta`
  - `Same`
  - `NotSame`
  - `Length`
  - `Contains`
  - `NotContains`
  - `Error`
  - `NoError`
  - `ErrorIs`
  - `NotErrorIs`
  - `EqualJSON`
  - `JSON`
  - `Fail`
  - `Failf`
