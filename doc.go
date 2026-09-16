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
//	// point.X: Should be equal
//	//   actual: 1
//	// expected: 2
//
// # Available assertions
//
//   - [True] / [False] — boolean condition
//   - [Nil] / [NotNil] — value is nil / non-nil, including typed nil references
//   - [Equal] / [NotEqual] — deep value equality
//   - [EqualDelta] / [NotEqualDelta] — numeric equality within a tolerance
//   - [Greater] / [GreaterOrEqual] / [Less] / [LessOrEqual] — numeric ordering
//   - [Same] / [NotSame] — reference identity (same type and address)
//   - [Length] / [Empty] / [NotEmpty] — length of a string, array, slice, map, or channel
//   - [Contains] / [NotContains] — membership in a string, array, slice, or map values
//   - [Error] / [NoError] — error is non-nil / nil
//   - [ErrorIs] / [NotErrorIs] — error matches / does not match a target with errors.Is
//   - [ErrorAs] — error is assignable to a target with errors.As
//   - [Matches] / [NotMatches] — string matches / does not match a regular expression
//   - [EqualJSON] — semantic JSON string equality
//   - [JSON] — object marshals to an expected JSON string
//   - [Panics] / [NotPanics] — function panics / does not panic
//   - [PanicsWith] — function panics with an expected value or error
//   - [Fail] / [Failf] — unconditional failure with a custom message
package assert
