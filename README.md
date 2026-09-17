<div align="center" width="100%">

<a href="https://github.com/gravitton">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/assert/refs/heads/main/docs/images/logo-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/assert/refs/heads/main/docs/images/logo-light.svg">
  <img alt="Gravitton assert" src="https://raw.githubusercontent.com/gravitton/assert/refs/heads/main/docs/images/logo-light.svg" width="300">
</picture>
</a>

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Simple and lightweight testing assertion library for Go

<hr>

</div>


## Features

- **Zero dependencies** – standard library only.
- **Generic** – type-safe arguments, no `interface{}` juggling; named types and `time.Duration` just work.
- **Readable failures** – aligned `actual` / `expected` output, pointing at the call site in your test.
- **Composable** – every assertion returns `bool` and takes an optional message prefix, so building your own is trivial.
- **Exact** – integers compared without float rounding, JSON numbers compared as rationals, NaN handled explicitly.
- **Works with anything** that has `Helper()` and `Errorf()`: `*testing.T`, `*testing.B`, `testing.TB`, or your own.

## Installation

```shell
go get github.com/gravitton/assert
```

## Usage

```go
import "github.com/gravitton/assert"
```

Values:

```go
assert.True(t, ok)
assert.Equal(t, got, want)          // reflect.DeepEqual
assert.NotEqual(t, got, 0)
assert.Nil(t, ptr)                  // typed nil pointers, slices, maps, channels and funcs count
assert.Same(t, a, b)                // same type and address
```

Numbers and ordering:

```go
assert.EqualDelta(t, got, 3.14, 0.01)
assert.EqualDelta(t, elapsed, 100*time.Millisecond, 20*time.Millisecond)
assert.Greater(t, len(items), 0)
assert.LessOrEqual(t, ratio, 1.0)
assert.Less(t, "apple", "banana")    // any ordered type, strings included
```

Collections and strings:

```go
assert.Length(t, items, 3)          // string, array, slice, map or channel
assert.Empty(t, errs)
assert.NotEmpty(t, name)
assert.Contains(t, items, item)     // array or slice element, map value
assert.Contains(t, body, "<html>")  // substring
assert.Matches(t, id, `^[a-f0-9]{8}$`)
```

Errors:

```go
assert.NoError(t, err)
assert.ErrorIs(t, err, io.EOF)      // errors.Is
assert.ErrorAs(t, err, &pathErr)    // errors.As
```

Panics:

```go
assert.Panics(t, func() { div(1, 0) })
assert.PanicsWith(t, func() { div(1, 0) }, "division by zero")
assert.PanicsWith(t, func() { must(err) }, err) // errors.Is on the panic value
assert.NotPanics(t, func() { div(1, 1) })
```

JSON, compared semantically:

```go
assert.EqualJSON(t, body, `{"id": 1, "tags": ["a", "b"]}`)
assert.JSON(t, user, `{"id":1,"name":"Ann"}`) // marshals user first
```

Failure output is aligned and shows both sides:

```
user_test.go:12: Should be equal
      actual: "Ann"
    expected: "Bob"
```

Full reference: [pkg.go.dev][link-go-dev-reference].

## Assertions

| Function                                    | Description                                                                              |
|---------------------------------------------|------------------------------------------------------------------------------------------|
| `True(t, condition)`                        | condition is true                                                                        |
| `False(t, condition)`                       | condition is false                                                                       |
| `Nil(t, object)`                            | object is nil, including a typed nil pointer, slice, map, channel, function or unsafe.Pointer |
| `NotNil(t, object)`                         | object is not nil                                                                        |
| `Same(t, actual, expected)`                 | references have the same type and address; slices also the same length and capacity     |
| `NotSame(t, actual, expected)`              | references have a different type or address                                              |
| `Equal(t, actual, expected)`                | values are deeply equal                                                                  |
| `NotEqual(t, actual, expected)`             | values are not deeply equal                                                              |
| `EqualDelta(t, actual, expected, delta)`    | numeric values differ by at most delta                                                   |
| `NotEqualDelta(t, actual, expected, delta)` | numeric values differ by more than delta                                                 |
| `Greater(t, actual, expected)`              | actual > expected, for any ordered type                                                  |
| `GreaterOrEqual(t, actual, expected)`       | actual >= expected                                                                       |
| `Less(t, actual, expected)`                 | actual < expected                                                                        |
| `LessOrEqual(t, actual, expected)`          | actual <= expected                                                                       |
| `Length(t, object, n)`                      | string, array, slice, map or channel has length n                                        |
| `Empty(t, object)`                          | string, array, slice, map or channel has zero length                                     |
| `NotEmpty(t, object)`                       | string, array, slice, map or channel has non-zero length                                 |
| `Contains(t, object, element)`              | string contains substring, or array, slice or map values contain element                 |
| `NotContains(t, object, element)`           | string does not contain substring, or array, slice or map values do not contain element  |
| `Error(t, err)`                             | error is not nil                                                                         |
| `NoError(t, err)`                           | error is nil                                                                             |
| `ErrorIs(t, err, target)`                   | error matches target with `errors.Is`                                                    |
| `NotErrorIs(t, err, target)`                | error does not match target with `errors.Is`                                             |
| `ErrorAs(t, err, target)`                   | error is assignable to target with `errors.As`                                           |
| `Matches(t, actual, pattern)`               | string matches regular expression                                                        |
| `NotMatches(t, actual, pattern)`            | string does not match regular expression                                                 |
| `Panics(t, fn)`                             | fn panics                                                                                |
| `PanicsWith(t, fn, expected)`               | fn panics with a value deeply equal to expected; an error expected is matched with `errors.Is` |
| `NotPanics(t, fn)`                          | fn does not panic                                                                        |
| `EqualJSON(t, actual, expected)`            | JSON strings are semantically equal                                                      |
| `JSON(t, object, expected)`                 | object marshals to JSON semantically equal to expected                                   |
| `Fail(t, message)`                          | always fails with message                                                                |
| `Failf(t, format, args...)`                 | always fails with formatted message                                                      |

## Conventions

**Return value:** Every assertion returns `bool`, `true` on success and `false` on failure. Failures are reported
with `t.Errorf`, so the test continues; return early yourself when a later assertion depends on an earlier one.

**Messages:** Every assertion except `Fail` and `Failf` accepts a trailing `messages ...string`. The strings are
concatenated and prepended to the failure message. Use them as a prefix, `"user.Name: "`, rather than a sentence.

**Nil:** `Nil`, `NotNil`, `Error` and `NoError` treat a typed nil (`(*MyErr)(nil)` stored in an `error`) as nil,
unlike a plain `== nil` comparison. Pointers, slices, maps, channels, functions and `unsafe.Pointer` all qualify.

**Equality:** `Equal` uses `reflect.DeepEqual`: pointers are compared by the values they reference and two non-nil
functions are never equal. `Same` compares identity instead: same type and address, and for slices also the same
length and capacity. Pointers to zero-size values and slices with zero capacity may share an address even when
allocated separately.

**Numbers:** `EqualDelta` compares integer types exactly, without a detour through `float64`, and panics on a negative
or NaN delta. NaN is only equal to NaN. `Greater`, `Less` and friends accept any `cmp.Ordered` type, strings
included, and fail when either side is NaN.

**Length and contents:** String length is measured in bytes, not runes. **`Contains` on a map searches the values,
not the keys**, unlike testify. The element must be assignable to the container's element type; a mismatch is reported as a failure, not
silently `false`.

**JSON:** `EqualJSON` and `JSON` compare structure, not text. Key order and whitespace are ignored, numbers are compared
exactly, so `1.0` equals `1` and integers beyond 2^53 keep their precision.

**Panics:** `Panics` recognises `panic(nil)`. `PanicsWith` matches an `error` expected value with `errors.Is` and
anything else with `reflect.DeepEqual`.

## Custom assertions

The return value and message prefix compose into reusable helpers:

```go
func TestRect(t *testing.T) {
	assertRect(t, image.Rect(1, 2, 3, 4), image.Rect(1, 3, 2, 4))
	// rect_test.go:4: Min.Y: Should be equal
	//       actual: 2
	//     expected: 3
	// rect_test.go:4: Max.X: Should be equal
	//       actual: 3
	//     expected: 2
}

func assertRect(t *testing.T, actual, expected image.Rectangle, messages ...string) bool {
	t.Helper()

	ok := assertPoint(t, actual.Min, expected.Min, append(messages, "Min.")...)
	ok = assertPoint(t, actual.Max, expected.Max, append(messages, "Max.")...) && ok

	return ok
}

func assertPoint(t *testing.T, actual, expected image.Point, messages ...string) bool {
	t.Helper()

	ok := assert.Equal(t, actual.X, expected.X, append(messages, "X: ")...)
	ok = assert.Equal(t, actual.Y, expected.Y, append(messages, "Y: ")...) && ok

	return ok
}
```

## Credits

- [Tomáš Novotný](https://github.com/tomas-novotny)
- [All Contributors][link-contributors]

## License

The MIT License (MIT). Please see [License File][link-licence] for more information.


[ico-license]:              https://img.shields.io/github/license/gravitton/assert.svg?style=flat-square&colorB=blue
[ico-workflow]:             https://img.shields.io/github/actions/workflow/status/gravitton/assert/main.yml?branch=main&style=flat-square
[ico-release]:              https://img.shields.io/github/v/release/gravitton/assert?style=flat-square&colorB=blue
[ico-go-dev-reference]:     https://img.shields.io/badge/go.dev-reference-blue?style=flat-square
[ico-coverage]:             https://img.shields.io/coverallsCoverage/github/gravitton/assert?style=flat-square

[link-author]:              https://github.com/gravitton
[link-release]:             https://github.com/gravitton/assert/releases
[link-contributors]:        https://github.com/gravitton/assert/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/assert/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/assert
[link-coverage]:            https://coveralls.io/github/gravitton/assert
