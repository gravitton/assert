# Assert

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Report Card][ico-go-report-card]][link-go-report-card]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Simple and lightweight testing assertion library for Go.


## Installation

```bash
go get github.com/gravitton/assert
```


## Usage

```go
import (
	"testing"

	"github.com/gravitton/assert"
)

func TestFoo(t *testing.T) {
	assert.Equal(t, actual, expected)
	assert.Contains(t, []int{1, 2, 3}, 2)
	assert.NoError(t, err)
	assert.JSON(t, obj, `{"data":1}`)
}
```


## Assertions

| Function | Description |
|---|---|
| `True(t, condition)` | condition is true |
| `False(t, condition)` | condition is false |
| `Same(t, actual, expected)` | pointers reference the same object |
| `NotSame(t, actual, expected)` | pointers reference different objects |
| `Equal(t, actual, expected)` | values are equal (deep) |
| `NotEqual(t, actual, expected)` | values are not equal |
| `EqualDelta(t, actual, expected, delta)` | numeric values differ by at most delta |
| `NotEqualDelta(t, actual, expected, delta)` | numeric values differ by more than delta |
| `Greater(t, actual, expected)` | actual > expected |
| `GreaterOrEqual(t, actual, expected)` | actual >= expected |
| `Less(t, actual, expected)` | actual < expected |
| `LessOrEqual(t, actual, expected)` | actual <= expected |
| `Length(t, object, n)` | string/slice/map/channel has length n |
| `Empty(t, object)` | string/slice/map/channel has zero length |
| `NotEmpty(t, object)` | string/slice/map/channel has non-zero length |
| `Contains(t, object, element)` | string/slice/map/channel contains element |
| `NotContains(t, object, element)` | string/slice/map/channel does not contain element |
| `Error(t, err)` | error is not nil |
| `NoError(t, err)` | error is nil |
| `ErrorIs(t, err, target)` | error unwraps to target |
| `NotErrorIs(t, err, target)` | error does not unwrap to target |
| `Matches(t, actual, pattern)` | string matches regular expression |
| `NotMatches(t, actual, pattern)` | string does not match regular expression |
| `EqualJSON(t, actual, expected)` | JSON strings are semantically equal |
| `JSON(t, object, expected)` | object marshals to expected JSON string |
| `Fail(t, message)` | always fails with message |
| `Failf(t, format, args...)` | always fails with formatted message |

All assertions return `bool`: `true` on success, `false` on failure.

An optional trailing `messages ...string` argument is available on every assertion. The strings are concatenated and prepended to the failure message, which is useful for building custom assertions.


## Custom assertions

The return value and optional message prefix make it straightforward to compose reusable assertion helpers:

```go
func TestRect(t *testing.T) {
	assertRect(t, image.Rect(1, 2, 3, 4), image.Rect(1, 3, 2, 4))
	// test.go:4: Min.Y: Should be equal:
	//       actual: 2
	//     expected: 3
	// test.go:4: Max.X: Should be equal:
	//       actual: 3
	//     expected: 2
}

func assertRect(t *testing.T, actual, expected image.Rectangle) bool {
	t.Helper()

	ok := assertPoint(t, actual.Min, expected.Min, "Min.")
	ok = assertPoint(t, actual.Max, expected.Max, "Max.") && ok

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
[ico-go-report-card]:       https://goreportcard.com/badge/github.com/gravitton/assert?style=flat-square
[ico-coverage]:             https://img.shields.io/coverallsCoverage/github/gravitton/assert?style=flat-square

[link-author]:              https://github.com/gravitton
[link-release]:             https://github.com/gravitton/assert/releases
[link-contributors]:        https://github.com/gravitton/assert/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/assert/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/assert
[link-go-report-card]:      https://goreportcard.com/report/github.com/gravitton/assert
[link-coverage]:            https://coveralls.io/github/gravitton/assert
