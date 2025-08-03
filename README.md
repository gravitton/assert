# Assert

[![Build Status][ico-workflow]][link-workflow]
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

## Contributing

Please see [CONTRIBUTING][link-contributing] and [CODE_OF_CONDUCT][link-code-of-conduct] for details.


## Credits

- [Tomáš Novotný](https://github.com/tomas-novotny)
- [All Contributors][link-contributors]


## License

The MIT License (MIT). Please see [License File][link-licence] for more information.


[ico-license]:              https://img.shields.io/github/license/gravitton/assert.svg?style=flat-square&colorB=blue
[ico-workflow]:             https://img.shields.io/github/actions/workflow/status/gravitton/assert/master.yml?branch=master&style=flat-square

[link-author]:              https://github.com/gravitton
[link-contributors]:        https://github.com/gravitton/assert/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/assert/actions
