//go:build race

package assert

func allocationsOf(runs int, fn func()) (allocations int, counted bool) {
	return 0, false
}
