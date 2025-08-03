# Assert

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
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
	"github.com/gravitton/assert"
	"testing"
)

func TestSomething(t *testing.T) {
	// assert equality
	assert.Equal(t, 123, 123)
	// assert inequality
	assert.NotEqual(t, 123, 456)
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

[link-author]:              https://github.com/gravitton
[link-release]:             https://github.com/gravitton/assert/releases
[link-contributors]:        https://github.com/gravitton/assert/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/assert/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/assert

