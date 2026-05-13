package assert

import (
	"encoding/json"
	"errors"
	"regexp"
)

// Testing is an interface wrapper around *testing.T
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// Fail reports a failure message via t.Errorf.
func Fail(t Testing, message string) bool {
	t.Helper()

	t.Errorf("%s", message)

	return false
}

// Failf reports a formatted failure message via t.Errorf.
func Failf(t Testing, format string, args ...any) bool {
	t.Helper()

	t.Errorf(format, args...)

	return false
}

// True asserts that the specified value is true.
func True(t Testing, condition bool, messages ...string) bool {
	t.Helper()

	if !condition {
		return Failf(t, "%sShould be true", message(messages))
	}

	return true
}

// False asserts that the specified value is false.
func False(t Testing, condition bool, messages ...string) bool {
	t.Helper()

	if condition {
		return Failf(t, "%sShould be false", message(messages))
	}

	return true
}

// Same asserts that two pointers reference the same object.
//
// Both arguments must be pointer variables. Pointer variable equality is
// determined based on the equality of both type and value.
func Same[T Reference](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if valid, ok := same[T](expected, actual); !ok {
		return Failf(t, "%sShould be pointers\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	} else if !valid {
		return Failf(t, "%sShould be same\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// NotSame asserts that two pointers do NOT reference the same object.
//
// Both arguments must be pointer variables. Pointer variable equality is
// determined based on the equality of both type and value.
func NotSame[T Reference](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if valid, ok := same(expected, actual); !ok {
		Failf(t, "%sShould be pointers\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	} else if valid {
		return Failf(t, "%sShould not be same\n  actual: %s", message(messages), print(actual))
	}

	return true
}

// Equal asserts that two objects are equal.
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func Equal[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if !equal(actual, expected) {
		return Failf(t, "%sShould be equal:\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// NotEqual asserts that the specified values are NOT equal.
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqual[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if equal(actual, expected) {
		return Failf(t, "%sShould not be equal\n  actual: %s", message(messages), print(actual))
	}

	return true
}

// EqualDelta asserts that two numeric values difference is less than delta.
//
// Panics if delta is negative.
func EqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if !equalDelta(actual, expected, delta) {
		return Failf(t, "%sShould be equal in delta:\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// NotEqualDelta asserts that two numeric values difference is greater than delta.
//
// Panics if delta is negative.
func NotEqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if equalDelta(actual, expected, delta) {
		return Failf(t, "%sShould not be equal in delta:\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// Greater asserts that actual is greater than expected.
func Greater[T Numeric](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if actual <= expected {
		return Failf(t, "%sShould be greater\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// GreaterOrEqual asserts that actual is greater than or equal to expected.
func GreaterOrEqual[T Numeric](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if actual < expected {
		return Failf(t, "%sShould be greater or equal\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// Less asserts that actual is less than expected.
func Less[T Numeric](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if actual >= expected {
		return Failf(t, "%sShould be less\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// LessOrEqual asserts that actual is less than or equal to expected.
func LessOrEqual[T Numeric](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if actual > expected {
		return Failf(t, "%sShould be less or equal\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// Length asserts that object has given length.
func Length[S Iterable[any]](t Testing, object S, expected int, messages ...string) bool {
	t.Helper()

	if actual := length(object); actual != expected {
		return Failf(t, "%sShould have element length\n  object: %#v\n  actual: %d\nexpected: %d", message(messages), object, actual, expected)
	}

	return true
}

// Empty asserts that object has zero length.
func Empty[S Iterable[any]](t Testing, object S, messages ...string) bool {
	t.Helper()

	if length(object) != 0 {
		return Failf(t, "%sShould be empty\n  object: %#v", message(messages), object)
	}

	return true
}

// NotEmpty asserts that object has non-zero length.
func NotEmpty[S Iterable[any]](t Testing, object S, messages ...string) bool {
	t.Helper()

	if length(object) == 0 {
		return Failf(t, "%sShould not be empty\n  object: %#v", message(messages), object)
	}

	return true
}

// Contains asserts that object contains given element.
//
// Works with strings, arrays, slices, maps values and channels.
func Contains[S Iterable[E], E Comparable](t Testing, object S, element E, messages ...string) bool {
	t.Helper()

	if found, ok := contains(object, element); !ok {
		return Failf(t, "%sShould be iterable\n  object: %#v", message(messages), object)
	} else if !found {
		return Failf(t, "%sShould contain element\n  object: %#v\n element: %#v", message(messages), object, element)
	}

	return true
}

// NotContains asserts that object does NOT contain given element.
//
// Works with strings, arrays, slices, maps values and channels.
func NotContains[S Iterable[E], E Comparable](t Testing, object S, element E, messages ...string) bool {
	t.Helper()

	if found, ok := contains(object, element); !ok {
		Failf(t, "%sShould be iterable\n  object: %#v", message(messages), object)
	} else if found {
		return Failf(t, "%sShould not contain element\n  object: %#v\n element: %#v", message(messages), object, element)
	}

	return true
}

// Error asserts that error is NOT nil.
//
// A typed nil error (e.g. (*MyError)(nil)) is treated as nil.
func Error(t Testing, err error, messages ...string) bool {
	t.Helper()

	if isNilError(err) {
		return Failf(t, "%sShould be error", message(messages))
	}

	return true
}

// NoError asserts that error is nil.
//
// A typed nil error (e.g. (*MyError)(nil)) is treated as nil.
func NoError(t Testing, err error, messages ...string) bool {
	t.Helper()

	if !isNilError(err) {
		return Failf(t, "%sShould not be error\n     msg: %[2]v\n   error: %#[2]v", message(messages), err)
	}

	return true
}

// ErrorIs asserts that error is unwrappable to given target.
func ErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if !errors.Is(err, target) {
		return Failf(t, "%sShould be same error\n     msg: %[2]v\n   error: %#[2]v\n  target: %#v", message(messages), err, target)
	}

	return true
}

// NotErrorIs asserts that error is NOT unwrappable to given target.
func NotErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if errors.Is(err, target) {
		return Failf(t, "%sShould not be same error\n     msg: %[2]v\n   error: %#[2]v\n  target: %#v", message(messages), err, target)
	}

	return true
}

// Matches asserts that a string matches the given regular expression.
func Matches(t Testing, actual, pattern string, messages ...string) bool {
	t.Helper()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return Failf(t, "%sShould be valid regexp\n pattern: %s\n     err: %v", message(messages), pattern, err)
	}

	if !re.MatchString(actual) {
		return Failf(t, "%sShould match regexp\n  actual: %s\n pattern: %s", message(messages), actual, pattern)
	}

	return true
}

// NotMatches asserts that a string does NOT match the given regular expression.
func NotMatches(t Testing, actual, pattern string, messages ...string) bool {
	t.Helper()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return Failf(t, "%sShould be valid regexp\n pattern: %s\n     err: %v", message(messages), pattern, err)
	}

	if re.MatchString(actual) {
		return Failf(t, "%sShould not match regexp\n  actual: %s\n pattern: %s", message(messages), actual, pattern)
	}

	return true
}

// EqualJSON asserts that JSON strings are equal.
func EqualJSON(t Testing, actual, expected string, messages ...string) bool {
	t.Helper()

	var actualJSON, expectedJSON any

	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		return Failf(t, "%sShould be valid JSON\n  actual: %s\n     err: %v", message(messages), actual, err)
	}

	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		return Failf(t, "%sShould be valid JSON\nexpected: %s\n     err: %v", message(messages), expected, err)
	}

	return Equal(t, actualJSON, expectedJSON)
}

// JSON asserts that object can be marshaled to expected JSON string.
func JSON(t Testing, actual any, expected string, messages ...string) bool {
	t.Helper()

	s, err := json.Marshal(actual)

	return NoError(t, err, messages...) && EqualJSON(t, string(s), expected, messages...)
}

// Panics asserts that fn panics.
//
// When expected is nil, only the presence of a panic is checked.
// When expected is an error, the panic value is validated with errors.Is.
// Otherwise, the panic value is validated with deep equality.
func Panics(t Testing, fn func(), expected any, messages ...string) bool {
	t.Helper()

	panicked, value := panics(fn)
	if !panicked {
		return Failf(t, "%sShould panic", message(messages))
	}

	if expected == nil {
		return true
	}

	if target, ok := expected.(error); ok {
		if err, ok := value.(error); ok {
			return ErrorIs(t, err, target, messages...)
		}

		return Failf(t, "%sShould panic with error\n  actual: %s\nexpected: %s", message(messages), print(value), print(expected))
	}

	if !equal(value, expected) {
		return Failf(t, "%sShould panic with value\n  actual: %s\nexpected: %s", message(messages), print(value), print(expected))
	}

	return true
}

// NotPanics asserts that fn does NOT panic.
func NotPanics(t Testing, fn func(), messages ...string) bool {
	t.Helper()

	if panicked, value := panics(fn); panicked {
		return Failf(t, "%sShould not panic\n  value: %s", message(messages), print(value))
	}

	return true
}
