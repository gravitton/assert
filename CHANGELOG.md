# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/assert/compare/v1.5.0...main)

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
