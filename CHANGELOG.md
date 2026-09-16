# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/assert/compare/v1.4.0...main)
### Added
- Added `Nil` and `NotNil` assert methods; typed nil pointers, slices, maps, channels and functions count as nil
- Added `ErrorAs` assert method, mirroring `errors.As`
- Added `PanicsWith` assert method, taking over the expected-value check from `Panics`

### Changed
- **Breaking:** `Panics` no longer takes an expected value; use `PanicsWith` for that
- **Breaking:** `Iterable` constraint no longer has a type parameter; `Comparable`, `Reference` and `Iterable` are documented as descriptive only
- Lowered the minimum Go version to 1.23; CI tests 1.23 through 1.27
- Failure messages: dropped the trailing colon from `Should be equal` and `Should be equal in delta`, `ErrorIs` reports `Should match error`, `EqualJSON` reports `Should be equal JSON`, `JSON` reports `Should be marshalable` on a marshal error
- `Same` and `NotSame` no longer accept `unsafe.Pointer` as a reference
- `Same` and `NotSame` now require both type and address to match, and report non-reference arguments as `Should be reference`; function values are no longer accepted as references
- `Contains` and `NotContains` report an element whose type does not fit the container as `Should have element of same type` instead of `Should be iterable`
- `Greater`, `GreaterOrEqual`, `Less` and `LessOrEqual` fail when either value is NaN
- `EqualDelta` and `NotEqualDelta` compare integer types exactly instead of through `float64`

### Fixed
- `Same` reported a slice and its resliced prefix (e.g. `s` and `s[:0]`) as the same; slices now also compare length and capacity
- `EqualJSON` and `JSON` conflated integers beyond 2^53 by decoding through `float64`; numbers are now compared exactly and the failure shows the original JSON strings
- `EqualDelta` and `NotEqualDelta` silently never matched on a NaN delta; they now panic like they do on a negative delta
- `NotSame` and `NotContains` returned `true` after reporting a non-reference or non-iterable argument
- `Equal`, `NotEqual`, `Same`, `NotSame` and `Panics` panicked when printing a nil pointer
- `Length`, `Empty`, `NotEmpty`, `Contains` and `NotContains` panicked on unsupported types instead of failing; `Contains` on a channel no longer panics but fails as not iterable
- `Contains` on a string with a non-string element compared against reflect's debug representation
- `EqualJSON` and `JSON` dropped the custom message prefix on the comparison failure
- `Error` and `NoError` now treat nil map, slice, func and channel errors as nil


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
