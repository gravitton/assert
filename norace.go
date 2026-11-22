//go:build !race

package assert

import "testing"

func allocationsOf(fn func()) (allocations int, counted bool) {
	return int(testing.AllocsPerRun(100, fn)), true
}
