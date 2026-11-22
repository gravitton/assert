//go:build race

package assert

func allocationsOf(fn func()) (allocations int, counted bool) {
	return 0, false
}
