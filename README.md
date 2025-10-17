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
package main

import (
	"testing"

	"github.com/gravitton/assert"
)

func Test(t *testing.T) {
	// assert equality
	assert.Equal(t, 123, 123, "Custom message: ")
	// assert inequality
	assert.NotEqual(t, 123, 456)
	// assert object contains element
	assert.Contains(t, []int{1, 2, 3}, 2)
	// assert object can be marshall to given JSON string
	assert.JSON(t, &obj, `{"data":1}`)
}

```

All assertions return `false` if assertion is not successful 
and accepts a custom error message(s) as the last argument which are prefixed to the assert error message.

You can use this to build your own complex assertions.

```go
package main

import (
	"image"
	"testing"

	"github.com/gravitton/assert"
)

func Test(t *testing.T) {
	assertRect(t, image.Rect(1, 2, 3, 4), image.Rect(1, 2, 3, 4))
	// test.go:11: Min.Y: Should be equal:
	//       actual: 2
	//     expected: 3
	// test.go:11: Max.X: Should be equal:
	//       actual: 3
	//     expected: 2
}

func assertRect(t *testing.T, actual image.Rectangle, expected image.Rectangle) bool {
	t.Helper()

	ok := true
	
	if !assertPoint(t, actual.Min, expected.Min, "Min.") {
		ok = false
	}
	if !assertPoint(t, actual.Max, expected.Max, "Max.") {
		ok = false
	}
	
	return ok
}

func assertPoint(t *testing.T, actual image.Point, expected image.Point, messages ...string) bool {
	t.Helper()

	ok := true
	
	if !assert.Equal(t, actual.X, expected.X, append(messages, "X: ")...) {
		ok = false
	}
	if !assert.Equal(t, actual.Y, expected.Y, append(messages, "Y: ")...) {
		ok = false
	}

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
