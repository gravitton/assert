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
- **Generic** – both sides of a comparison share one type, so untyped constants take the type of the value under test;
  named types and `time.Duration` just work.
- **Readable failures** – aligned `actual` / `expected` output, pointing at the call site in your test; every
  assertion's example on pkg.go.dev shows its failure messages.
- **Composable** – every assertion returns `bool` and takes an optional message prefix, so building your own is trivial.
- **Exact** – integers compared without float rounding, NaN handled explicitly.
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
assert.Zero(t, user.DeletedAt)      // zero value of any type, structs and time.Time included
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
assert.Length(t, items, 3)          // string, array, array pointer, slice, map or channel
assert.Empty(t, errs)
assert.NotEmpty(t, name)
assert.Contains(t, items, item)     // array or slice element, map value
assert.Contains(t, body, "<html>")  // substring
assert.EqualUnordered(t, ids, []int{3, 1, 2})
assert.HasPrefix(t, path, "/api/")
assert.HasSuffix(t, events, []string{"commit", "close"}) // strings or slices
assert.Matches(t, id, `^[a-f0-9]{8}$`)
```

Errors:

```go
assert.NoError(t, err)
assert.ErrorIs(t, err, io.EOF)      // errors.Is
assert.ErrorAs(t, err, &pathErr)    // errors.As
assert.ErrorContains(t, err, "permission denied")
```

Panics:

```go
assert.Panics(t, func() { div(1, 0) })
assert.PanicsWith(t, func() { stack.Pop() }, "pop from empty stack")
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

| Function                                    | Description                                                                                            |
|---------------------------------------------|--------------------------------------------------------------------------------------------------------|
| `True(t, condition)`                        | condition is true                                                                                      |
| `False(t, condition)`                       | condition is false                                                                                     |
| `Nil(t, actual)`                            | actual is nil, including a typed nil pointer, slice, map, channel, function or unsafe.Pointer          |
| `NotNil(t, actual)`                         | actual is not nil                                                                                      |
| `Zero(t, actual)`                           | actual is the zero value of its type                                                                   |
| `NotZero(t, actual)`                        | actual is not the zero value of its type                                                               |
| `Same(t, actual, expected)`                 | references have the same type and address; slices also the same length and capacity                    |
| `NotSame(t, actual, expected)`              | references differ in type or address, or slices in length or capacity                                  |
| `Equal(t, actual, expected)`                | values are deeply equal                                                                                |
| `NotEqual(t, actual, expected)`             | values are not deeply equal                                                                            |
| `EqualDelta(t, actual, expected, delta)`    | numeric values differ by at most delta                                                                 |
| `NotEqualDelta(t, actual, expected, delta)` | numeric values differ by more than delta                                                               |
| `Greater(t, actual, bound)`                 | actual > bound, for any ordered type                                                                   |
| `GreaterOrEqual(t, actual, bound)`          | actual >= bound                                                                                        |
| `Less(t, actual, bound)`                    | actual < bound                                                                                         |
| `LessOrEqual(t, actual, bound)`             | actual <= bound                                                                                        |
| `Length(t, actual, n)`                      | string, array, array pointer, slice, map or channel has length n                                       |
| `Empty(t, actual)`                          | string, array, array pointer, slice, map or channel has zero length                                    |
| `NotEmpty(t, actual)`                       | string, array, array pointer, slice, map or channel has non-zero length                                |
| `Contains(t, actual, element)`              | string contains substring, or array, array pointer, slice or map values contain element                |
| `NotContains(t, actual, element)`           | string does not contain substring, or array, array pointer, slice or map values do not contain element |
| `EqualUnordered(t, actual, expected)`       | arrays, array pointers or slices have the same elements, each occurring as often, in any order         |
| `HasPrefix(t, actual, prefix)`              | string or slice begins with prefix                                                                     |
| `HasSuffix(t, actual, suffix)`              | string or slice ends with suffix                                                                       |
| `Error(t, err)`                             | error is not nil                                                                                       |
| `NoError(t, err)`                           | error is nil                                                                                           |
| `ErrorIs(t, err, target)`                   | error matches target with `errors.Is`                                                                  |
| `NotErrorIs(t, err, target)`                | error does not match target with `errors.Is`                                                           |
| `ErrorAs(t, err, target)`                   | error is assignable to target with `errors.As`; fails on an invalid target                             |
| `ErrorContains(t, err, substr)`             | error is not nil and its message contains substr                                                       |
| `Matches(t, actual, pattern)`               | string matches regular expression                                                                      |
| `NotMatches(t, actual, pattern)`            | string does not match regular expression                                                               |
| `Panics(t, fn)`                             | fn panics                                                                                              |
| `PanicsWith(t, fn, expected)`               | fn panics with a value deeply equal to expected, or matching it with `errors.Is` for errors            |
| `NotPanics(t, fn)`                          | fn does not panic                                                                                      |
| `EqualJSON(t, actual, expected)`            | JSON strings are semantically equal                                                                    |
| `JSON(t, actual, expected)`                 | actual marshals to JSON semantically equal to expected                                                 |
| `Fail(t, message)`                          | always fails with message                                                                              |
| `Failf(t, format, args...)`                 | always fails with formatted message                                                                    |

## Conventions

**Return value:** Every assertion returns `bool`, `true` on success and `false` on failure. Failures are reported
with `t.Errorf`, so the test continues; return early yourself when a later assertion depends on an earlier one.

**Messages:** Every assertion except `Fail` and `Failf` accepts a trailing `messages ...string`. The strings are
concatenated and prepended to the failure message. Use them as a prefix, `"user.Name: "`, rather than a sentence. No
separator is inserted, neither between the strings nor before the failure message, so each carries its own: a nested
helper appends `"Min."`, the leaf appends `"X: "`, and the failure reads `Min.X: Should be equal`.

**Nil:** `Nil` and `NotNil` treat a typed nil stored in their `any` argument as nil, unlike a plain `== nil`
comparison. Pointers, slices, maps, channels, functions and `unsafe.Pointer` all qualify. `Error`, `NoError` and
`ErrorContains` do not: a `(*MyErr)(nil)` returned as an `error` is an error, exactly as `if err != nil` sees it.

**Equality:** `Equal` uses `reflect.DeepEqual`: pointers are compared by the values they reference, a typed nil stored
in an interface is not equal to an untyped `nil`, and NaN is not equal to itself, as with `==`. Non-nil functions are
never deeply equal, so `Equal` and `NotEqual` reject them as invalid arguments; use `Nil` or `NotNil` instead. `Same`
compares identity: same type and address, and for slices also the same length and capacity. Two nil references of the
same type are the same. Pointers to zero-size values and slices with zero capacity may share an address even when
allocated separately.

**Numbers:** `EqualDelta` compares integer types exactly, without a detour through `float64`, and fails on a negative
or NaN delta. NaN is only equal to NaN. Its failure shows the delta and the actual difference. `Greater`, `Less` and
friends accept any `cmp.Ordered` type, strings included, and report a NaN on either side as `Should not be NaN`.

**Length and contents:** String length is measured in bytes, not runes. A pointer to an array works like the array,
except that `Contains` and `EqualUnordered` reject a nil one. `HasPrefix` and `HasSuffix` take strings and slices only,
since two arrays of one type always have the same length. **`Contains` on a map searches the values, not the keys**,
unlike testify. The element must have the container's element type, or be assignable to it when that is an interface
type, and a substring must have the same type as the string; a mismatch is reported as a failure, not silently `false`.
An untyped constant takes its default type, so write `Contains(t, ids, int64(1))` for an `[]int64`.
`EqualUnordered`, `HasPrefix` and `HasSuffix` require both arguments to have the same type, so a slice and an array, or
a string and a named string type, are reported as `Should have same type`. `EqualUnordered` treats a nil slice and an
empty one as equal, unlike `Equal`. `Zero` does the opposite of `Empty` there: a non-nil empty slice or map is not zero.

**Zero:** `Zero` and `NotZero` call the value's own `IsZero() bool` method when it has one, so a zero `time.Time` in
any location is zero. A method promoted from an embedded field counts too, so a struct embedding `time.Time` is zero
whenever that time is, whatever its other fields hold. A pointer is zero only when it is nil, even when it points to a
zero value, unlike `encoding/json`'s `omitzero`.

**JSON:** `EqualJSON` and `JSON` compare structure, not text. Both sides are decoded as `encoding/json` decodes into
`any`, so key order and whitespace are ignored and the last of duplicate keys wins. Numbers become `float64`: `1.0`
equals `1` and `1e2` equals `100`, but integers beyond 2^53 may compare equal when they differ, and numbers beyond the
`float64` range are reported as invalid JSON. To compare such numbers exactly, decode into a struct and use `Equal`.
`JSON` marshals its argument with `json.Marshal`, so a `[]byte` becomes a base64 string; compare JSON held in bytes with
`EqualJSON(t, string(b), expected)`.

**Panics:** `Panics`, `PanicsWith` and `NotPanics` recognise `panic(nil)` and report its value as `nil` whatever the
`GODEBUG=panicnil` setting. `PanicsWith` matches with `reflect.DeepEqual`, or with `errors.Is` when the expected value
is an `error`. All three fail on a nil function.

**Output:** Values are printed with `%#v`, with two exceptions: a number prints in decimal or through its `String`
method, so `uint(5)` prints as `5` and a `time.Duration` as `1.5s`, and a pointer prints as `&` followed by the value it
points to, one level deep. Everything else prints exactly as `%#v` prints it, so a pointer to a pointer, and pointers,
channels and functions nested in a value show their addresses. When two values of different types print the same, or a
check fails on a type mismatch, each value is printed with its type, e.g. `int(1)` and `int64(1)`. When `Equal` or
`PanicsWith` fails on two values of one type that print the same, such as NaN or a `time.Time` differing only in its
monotonic clock, a `hint` line says so. The error under test prints its `Error()` text as `msg` and its type as `error`,
and a panic value that is an error also prints its text as `msg`; regexp patterns, JSON and error details print as plain
text, and multi-line text is indented under its label. Values are printed in full.

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
