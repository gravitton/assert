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
- **Generic** – both sides share one type, so untyped constants take the type of the value under test.
- **Readable failures** – aligned `actual` / `expected` output at the call site.
- **Composable** – every assertion returns `bool` and takes a message prefix.
- **Strict** – an argument that cannot be checked fails the assertion instead of passing or panicking.
- **Exact** – integers compared without float rounding, NaN handled explicitly.
- **Works with anything** – `*testing.T`, `testing.TB` or any type with `Helper()` and `Errorf()`.

The rules behind them are in the [package documentation][link-go-dev-reference].

## Installation

```shell
go get github.com/gravitton/assert
```

```go
import "github.com/gravitton/assert"
```

## Usage

### Values

`Equal` compares with `reflect.DeepEqual`, `Same` compares identity. `Nil` counts a typed nil as nil, `Error` does not,
exactly as `if err != nil` sees it.

```go
assert.True(t, ok)
assert.Equal(t, got, want)
assert.NotEqual(t, got, 0)
assert.Nil(t, ptr)             // typed nil pointers, slices, maps, channels and funcs count
assert.Zero(t, user.DeletedAt) // calls IsZero when the type has one, time.Time included
assert.Same(t, a, b)           // same type and address
```

### Numbers and ordering

Integers are compared exactly, and NaN is only equal to NaN.

```go
assert.EqualDelta(t, got, 3.14, 0.01)
assert.EqualDelta(t, elapsed, 100*time.Millisecond, 20*time.Millisecond)
assert.Greater(t, len(items), 0)
assert.LessOrEqual(t, ratio, 1.0)
assert.Less(t, "apple", "banana") // any ordered type, strings included
```

### Collections and strings

`Contains` on a map searches the values, not the keys. The element must have the container's element type, so an
`[]int64` takes `int64(1)`.

```go
assert.Length(t, items, 3)                    // string, array, array pointer, slice, map or channel
assert.Empty(t, errs)
assert.NotEmpty(t, name)
assert.Contains(t, items, item)
assert.Contains(t, body, "<html>")            // substring
assert.EqualUnordered(t, ids, []int{3, 1, 2}) // each element as often, in any order
assert.HasPrefix(t, path, "/api/")
assert.HasSuffix(t, events, []string{"commit", "close"})
assert.Matches(t, id, `^[a-f0-9]{8}$`)
```

### Errors

```go
assert.NoError(t, err)
assert.ErrorIs(t, err, io.EOF)   // errors.Is
assert.ErrorAs(t, err, &pathErr) // errors.As
assert.ErrorContains(t, err, "permission denied")
```

### Panics

```go
assert.Panics(t, func() { div(1, 0) })
assert.PanicsWith(t, func() { stack.Pop() }, "pop from empty stack")
assert.PanicsWith(t, func() { must(err) }, err) // errors.Is on the panic value
assert.NotPanics(t, func() { div(1, 1) })
```

### JSON

Compared by structure, so key order and whitespace are ignored. Numbers decode as `float64`, so `1.0` equals `1`.

```go
assert.EqualJSON(t, body, `{"id": 1, "tags": ["a", "b"]}`)
assert.JSON(t, user, `{"id":1,"name":"Ann"}`) // marshals user first
```

### Failures

A failure is reported with `t.Errorf`, so the test continues; the assertion returns `false`. Both sides print aligned,
each with its type when they would otherwise print the same:

```
user_test.go:12: Should be equal
      actual: "Ann"
    expected: "Bob"
```

### Custom assertions

Trailing messages are joined without a separator and prepended to the failure, so helpers compose into a path:

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

Full reference: [pkg.go.dev][link-go-dev-reference].

## Credits

- [Tomáš Novotný](https://github.com/tomas-novotny)
- [All Contributors][link-contributors]

## License

The MIT License (MIT). Please see [License File][link-licence] for more information.


[ico-license]:              https://img.shields.io/github/license/gravitton/assert.svg?style=flat-square&color=blue
[ico-workflow]:             https://img.shields.io/github/actions/workflow/status/gravitton/assert/main.yml?branch=main&style=flat-square
[ico-release]:              https://img.shields.io/github/v/release/gravitton/assert?style=flat-square&color=blue
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
