// Package assert provides assertions for Go tests.
//
// An assertion takes a [Testing], which *testing.T satisfies, and calls t.Helper, so a failure
// points to the calling test. It reports a failure with t.Errorf, letting the test continue, and
// returns true on success and false on failure.
//
//	func TestFoo(t *testing.T) {
//		assert.True(t, ok)
//		assert.Equal(t, actual, expected)
//		assert.NoError(t, err)
//		assert.Contains(t, []int{1, 2, 3}, 2)
//		assert.JSON(t, obj, `{"key":"value"}`)
//	}
//
// # Rules
//
// Values are compared with reflect.DeepEqual unless an assertion says otherwise. An argument an
// assertion cannot check, such as a value of the wrong kind or two values of different types,
// fails it rather than passing or panicking. Each negation, NotX or NoError, follows the rules of
// its positive and fails on the same invalid arguments.
//
// # Messages
//
// Every assertion except [Fail] and [Failf] accepts trailing messages, joined without a separator
// and prepended to the failure, so each string carries its own: a nested helper appends "Min."
// and the leaf appends "X: ".
//
//	assert.Equal(t, actual.X, expected.X, "Min.", "X: ")
//	// Min.X: Should be equal
//	//   actual: 1
//	// expected: 2
//
// # Assertions
//
//   - [True] / [False] — boolean condition
//   - [Nil] / [NotNil] — nil, including a typed nil
//   - [Zero] / [NotZero] — zero value of its type
//   - [Equal] / [NotEqual] — deep equality
//   - [EqualDelta] / [NotEqualDelta] — numeric equality within a delta
//   - [Greater] / [GreaterOrEqual] / [Less] / [LessOrEqual] — ordering of numbers and strings
//   - [Same] / [NotSame] — reference identity
//   - [Length] / [Empty] / [NotEmpty] — length
//   - [Contains] / [NotContains] — substring, or element of an array, slice or map values
//   - [EqualUnordered] — same elements in any order
//   - [HasPrefix] / [HasSuffix] — string or slice begins / ends with another
//   - [Error] / [NoError] — non-nil / nil error
//   - [ErrorIs] / [NotErrorIs] — errors.Is
//   - [ErrorAs] — errors.As
//   - [ErrorContains] — error message contains a substring
//   - [Matches] / [NotMatches] — regular expression
//   - [EqualJSON] — equal JSON documents
//   - [JSON] — value marshals to a JSON document
//   - [Panics] / [NotPanics] — function panics
//   - [PanicsWith] — function panics with a value
//   - [Fail] / [Failf] — unconditional failure
package assert
