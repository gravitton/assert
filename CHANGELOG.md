# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/assert/compare/v1.5.0...main)
### Added
- `Zero` and `NotZero`, using the value's `IsZero` method when present; a non-nil empty slice or map is not zero, nor is a non-nil pointer
- `EqualUnordered`, comparing arrays or slices as multisets and reporting extra and missing elements
- `HasPrefix` and `HasSuffix` for strings and slices
- `ErrorContains`, matching a substring of the error message
- `Length`, `Empty`, `NotEmpty`, `Contains` and `NotContains` accept a pointer to an array
- Runnable examples showing the failure messages of each assertion

### Changed
- **Breaking:** `Error`, `NoError` and `ErrorContains` treat a typed nil stored in an `error` as an error, as `err != nil` does
- **Breaking:** `Equal` and `NotEqual` reject non-nil functions as `Should not be function` instead of treating them as never equal
- **Breaking:** `EqualJSON` and `JSON` compare numbers as `float64`, as `encoding/json` decodes them, instead of by exact decimal value; integers beyond 2^53 may compare equal and numbers beyond the `float64` range are invalid JSON
- **Breaking:** `Contains` and `NotContains` on a string require the substring to have the string's type, so a named string type and `string` are reported as `Should have element of same type`
- `Same` reports references of different types as `Should have same type` instead of `Should be same`
- `Greater`, `GreaterOrEqual`, `Less` and `LessOrEqual` report a NaN as `Should not be NaN`
- `ErrorAs` fails on an invalid target as `Should have pointer to error or interface target` instead of panicking
- The value under test is named `actual` in every signature and failure message, and the second argument of `Greater`, `GreaterOrEqual`, `Less` and `LessOrEqual` is named `bound`; `Length` prints the measured length as `length`
- Failure messages: `EqualDelta` and `NotEqualDelta` show the delta and the actual difference, and an invalid delta with `actual` and `expected`
- Failure messages: a number prints in decimal or through its `String` method, so `uint(5)` prints as `5` and a `time.Duration` as `1.5s`
- Failure messages: a pointer prints as `&` followed by the value it points to, one level deep
- Failure messages: the error under test prints its type as `error` instead of its `%#v` value, which showed the addresses inside a wrapped error
- Failure messages: values of different types that print the same, such as `1` and `int64(1)`, and the values in a type mismatch are printed with their types
- Failure messages: `Matches` quotes the actual string, multi-line text such as an error message or JSON is indented under its label, and all labels are right-aligned
- Failure messages: values are printed in full instead of being cut at 1024 bytes
- Failure messages: `NotPanics` and `PanicsWith` print the `Error()` text of a panic value that is an error as `msg`
- Failure messages: `Equal` and `PanicsWith` add a `hint` when two unequal values of one type print the same, such as NaN
- Documented that the `messages` strings are concatenated with no separator, neither between them nor before the failure message

### Fixed
- `Panics` passed and `NotPanics` reported a panic for a nil function; all three panic assertions now fail with `Should be non-nil function`
- `Contains` and `NotContains` accepted an element assignable to the element type but not identical to it, such as `[]int` in a slice of a named `[]int` type, which could never match; it is now reported as `Should have element of same type`
- `EqualJSON` and `JSON` accepted non-JSON whitespace such as `\f` after the top-level value

## [v1.5.0](https://github.com/gravitton/assert/compare/v1.4.0...v1.5.0) (2026-09-17)
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

## [v1.4.0](https://github.com/gravitton/assert/compare/v1.3.0...v1.4.0) (2026-08-25)
### Changed
- Require Go 1.27

## [v1.3.0](https://github.com/gravitton/assert/compare/v1.2.1...v1.3.0) (2026-05-13)
### Added
- Added `Panics` and `NotPanics` assert methods for panic detection; `Panics` accepts an optional expected value — `nil` skips value validation, an `error` validates with `errors.Is`, anything else uses deep equality


## [v1.2.1](https://github.com/gravitton/assert/compare/v1.2.0...v1.2.1) (2026-04-27)
### Changed
- Improve error formatting in `Error` assert methods for better readability


## [v1.2.0](https://github.com/gravitton/assert/compare/v1.1.0...v1.2.0) (2026-04-22)
### Added
- Added `Greater`, `GreaterOrEqual`, `Less`, and `LessOrEqual` assert methods for numeric comparison
- Added `Empty` and `NotEmpty` assert methods for zero-length checks on strings, slices, maps, and channels
- Added `NotEqualDelta` assert method, complementing `EqualDelta`

### Fixed
- `Error` and `NoError` now correctly treat typed nil errors (e.g. `(*MyError)(nil)`) as nil


## [v1.1.0](https://github.com/gravitton/assert/compare/v1.0.0...v1.1.0) (2026-04-22)
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
