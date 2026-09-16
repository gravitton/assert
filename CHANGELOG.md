# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/assert/compare/v1.4.0...main)
### Changed
- `Same` and `NotSame` now require both type and address to match, and report non-reference arguments as `Should be reference`; function values are no longer accepted as references
- `Contains` and `NotContains` report an element whose type does not fit the container as `Should have element of same type` instead of `Should be iterable`
- `Greater`, `GreaterOrEqual`, `Less` and `LessOrEqual` fail when either value is NaN
- `EqualDelta` and `NotEqualDelta` compare integer types exactly instead of through `float64`

### Fixed
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
