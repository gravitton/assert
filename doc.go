// Package assert provides a simple and lightweight testing assertion library for Go.
//
// It works with any type that implements the [Testing] interface, which *testing.T satisfies.
// All assertion functions call t.Helper() so that failures point to the call site in your
// test, not inside this package.
//
// # Basic usage
//
//	func TestFoo(t *testing.T) {
//		assert.True(t, ok)
//		assert.Equal(t, actual, expected)
//		assert.NoError(t, err)
//		assert.Contains(t, []int{1, 2, 3}, 2)
//		assert.JSON(t, obj, `{"key":"value"}`)
//	}
//
// # Return values
//
// Every assertion returns bool: true on success, false on failure. This lets you
// short-circuit or accumulate failures in custom assertion helpers.
//
// # Custom messages
//
// All assertions accept an optional trailing messages ...string argument. The strings
// are concatenated and prepended to the failure output, making it easy to add context
// or build composable helpers:
//
//	assert.Equal(t, actual.X, expected.X, "point.X: ")
//	// point.X: Should be equal:
//	//   actual: 1
//	// expected: 2
//
// # Available assertions
//
//   - [True] / [False] — boolean condition
//   - [Equal] / [NotEqual] — deep value equality
//   - [EqualDelta] — numeric equality within a tolerance
//   - [Same] / [NotSame] — pointer identity (same memory address)
//   - [Length] — length of a string, slice, map, or channel
//   - [Contains] / [NotContains] — membership in a string, slice, map, or channel
//   - [Error] / [NoError] — error is non-nil / nil
//   - [ErrorIs] / [NotErrorIs] — error wraps / does not wrap a target
//   - [EqualJSON] — semantic JSON string equality
//   - [JSON] — object marshals to an expected JSON string
//   - [Fail] / [Failf] — unconditional failure with a custom message
package assert
