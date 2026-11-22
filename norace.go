//go:build !race

package assert

import "testing"

func allocationsOf(runs int, fn func()) (allocations int, counted bool) {
	return int(testing.AllocsPerRun(runs, fn)), true
}
