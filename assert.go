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

// Nil asserts that the specified object is nil.
//
// A typed nil pointer, slice, map, channel or function is treated as nil.
func Nil(t Testing, object any, messages ...string) bool {
	t.Helper()

	if !isNil(object) {
		return Failf(t, "%sShould be nil\n  actual: %s", message(messages), print(object))
	}

	return true
}

// NotNil asserts that the specified object is NOT nil.
//
// A typed nil pointer, slice, map, channel or function is treated as nil.
func NotNil(t Testing, object any, messages ...string) bool {
	t.Helper()

	if isNil(object) {
		return Failf(t, "%sShould not be nil\n  actual: %s", message(messages), print(object))
	}

	return true
}

// Same asserts that two references point to the same object.
//
// Both arguments must be references: pointers, slices, maps or channels.
// Two references are the same when they have the same type and address;
// slices must also have the same length and capacity.
//
// Pointers to zero-size values and slices with zero capacity may share an
// address even when allocated separately, and are then reported as same.
func Same[T Reference](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if identical, why := same(actual, expected); why != valid {
		return Failf(t, "%s%s\n  actual: %s\nexpected: %s", message(messages), why, print(actual), print(expected))
	} else if !identical {
		return Failf(t, "%sShould be same\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// NotSame asserts that two references do NOT point to the same object.
//
// Both arguments must be references: pointers, slices, maps or channels.
// Two references are the same when they have the same type and address;
// slices must also have the same length and capacity.
func NotSame[T Reference](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if identical, why := same(actual, expected); why != valid {
		return Failf(t, "%s%s\n  actual: %s\nexpected: %s", message(messages), why, print(actual), print(expected))
	} else if identical {
		return Failf(t, "%sShould not be same\n  actual: %s", message(messages), print(actual))
	}

	return true
}

// Equal asserts that two objects are equal.
//
// Equality is determined with reflect.DeepEqual: pointers are compared by the
// values they reference, and two non-nil functions are never equal.
func Equal[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if !equal(actual, expected) {
		return Failf(t, "%sShould be equal\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// NotEqual asserts that the specified values are NOT equal.
//
// Equality is determined with reflect.DeepEqual: pointers are compared by the
// values they reference, and two non-nil functions are never equal.
func NotEqual[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if equal(actual, expected) {
		return Failf(t, "%sShould not be equal\n  actual: %s", message(messages), print(actual))
	}

	return true
}

// EqualDelta asserts that two numeric values differ by at most delta.
//
// Panics if delta is negative or NaN. NaN is only equal to NaN.
func EqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if !equalDelta(actual, expected, delta) {
		return Failf(t, "%sShould be equal in delta\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// NotEqualDelta asserts that two numeric values differ by more than delta.
//
// Panics if delta is negative or NaN. NaN is only equal to NaN.
func NotEqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if equalDelta(actual, expected, delta) {
		return Failf(t, "%sShould not be equal in delta\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// Greater asserts that actual is greater than expected.
//
// Fails when either value is NaN.
func Greater[T Numeric](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if order, ok := compare(actual, expected); !ok || order <= 0 {
		return Failf(t, "%sShould be greater\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// GreaterOrEqual asserts that actual is greater than or equal to expected.
//
// Fails when either value is NaN.
func GreaterOrEqual[T Numeric](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if order, ok := compare(actual, expected); !ok || order < 0 {
		return Failf(t, "%sShould be greater or equal\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// Less asserts that actual is less than expected.
//
// Fails when either value is NaN.
func Less[T Numeric](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if order, ok := compare(actual, expected); !ok || order >= 0 {
		return Failf(t, "%sShould be less\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// LessOrEqual asserts that actual is less than or equal to expected.
//
// Fails when either value is NaN.
func LessOrEqual[T Numeric](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if order, ok := compare(actual, expected); !ok || order > 0 {
		return Failf(t, "%sShould be less or equal\n  actual: %s\nexpected: %s", message(messages), print(actual), print(expected))
	}

	return true
}

// Length asserts that object has given length.
//
// Works with strings, arrays, slices, maps and channels. String length is measured in bytes, not runes.
func Length[S Iterable](t Testing, object S, expected int, messages ...string) bool {
	t.Helper()

	if actual, why := length(object); why != valid {
		return Failf(t, "%s%s\n  object: %#v", message(messages), why, object)
	} else if actual != expected {
		return Failf(t, "%sShould have length\n  object: %#v\n  actual: %d\nexpected: %d", message(messages), object, actual, expected)
	}

	return true
}

// Empty asserts that object has zero length.
//
// Works with strings, arrays, slices, maps and channels. String length is measured in bytes, not runes.
func Empty[S Iterable](t Testing, object S, messages ...string) bool {
	t.Helper()

	if actual, why := length(object); why != valid {
		return Failf(t, "%s%s\n  object: %#v", message(messages), why, object)
	} else if actual != 0 {
		return Failf(t, "%sShould be empty\n  object: %#v", message(messages), object)
	}

	return true
}

// NotEmpty asserts that object has non-zero length.
//
// Works with strings, arrays, slices, maps and channels. String length is measured in bytes, not runes.
func NotEmpty[S Iterable](t Testing, object S, messages ...string) bool {
	t.Helper()

	if actual, why := length(object); why != valid {
		return Failf(t, "%s%s\n  object: %#v", message(messages), why, object)
	} else if actual == 0 {
		return Failf(t, "%sShould not be empty\n  object: %#v", message(messages), object)
	}

	return true
}

// Contains asserts that object contains given element.
//
// Works with strings, arrays, slices and map values.
func Contains[S Iterable, E Comparable](t Testing, object S, element E, messages ...string) bool {
	t.Helper()

	if found, why := contains(object, element); why != valid {
		return Failf(t, "%s%s\n  object: %#v\n element: %#v", message(messages), why, object, element)
	} else if !found {
		return Failf(t, "%sShould contain element\n  object: %#v\n element: %#v", message(messages), object, element)
	}

	return true
}

// NotContains asserts that object does NOT contain given element.
//
// Works with strings, arrays, slices and map values.
func NotContains[S Iterable, E Comparable](t Testing, object S, element E, messages ...string) bool {
	t.Helper()

	if found, why := contains(object, element); why != valid {
		return Failf(t, "%s%s\n  object: %#v\n element: %#v", message(messages), why, object, element)
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

	if isNil(err) {
		return Failf(t, "%sShould be error", message(messages))
	}

	return true
}

// NoError asserts that error is nil.
//
// A typed nil error (e.g. (*MyError)(nil)) is treated as nil.
func NoError(t Testing, err error, messages ...string) bool {
	t.Helper()

	if !isNil(err) {
		return Failf(t, "%sShould not be error\n     msg: %[2]v\n   error: %#[2]v", message(messages), err)
	}

	return true
}

// ErrorIs asserts that error matches given target according to errors.Is.
func ErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if !errors.Is(err, target) {
		return Failf(t, "%sShould match error\n     msg: %[2]v\n   error: %#[2]v\n  target: %#v", message(messages), err, target)
	}

	return true
}

// NotErrorIs asserts that error does NOT match given target according to errors.Is.
func NotErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if errors.Is(err, target) {
		return Failf(t, "%sShould not match error\n     msg: %[2]v\n   error: %#[2]v\n  target: %#v", message(messages), err, target)
	}

	return true
}

// ErrorAs asserts that error can be assigned to target according to errors.As.
//
// As with errors.As, target must be a non-nil pointer to a type implementing
// error or to any interface type; otherwise ErrorAs panics.
func ErrorAs(t Testing, err error, target any, messages ...string) bool {
	t.Helper()

	if !errors.As(err, target) {
		return Failf(t, "%sShould be assignable to target\n     msg: %[2]v\n   error: %#[2]v\n  target: %T", message(messages), err, target)
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

// EqualJSON asserts that JSON strings are semantically equal.
//
// Numbers are compared exactly, so 1.0 equals 1 and large integers keep their precision.
func EqualJSON(t Testing, actual, expected string, messages ...string) bool {
	t.Helper()

	actualJSON, err := decodeJSON(actual)
	if err != nil {
		return Failf(t, "%sShould be valid JSON\n  actual: %s\n     err: %v", message(messages), actual, err)
	}

	expectedJSON, err := decodeJSON(expected)
	if err != nil {
		return Failf(t, "%sShould be valid JSON\nexpected: %s\n     err: %v", message(messages), expected, err)
	}

	if !equal(actualJSON, expectedJSON) {
		return Failf(t, "%sShould be equal JSON\n  actual: %s\nexpected: %s", message(messages), actual, expected)
	}

	return true
}

// JSON asserts that object can be marshaled to expected JSON string.
func JSON(t Testing, actual any, expected string, messages ...string) bool {
	t.Helper()

	s, err := json.Marshal(actual)
	if err != nil {
		return Failf(t, "%sShould be marshalable\n  actual: %s\n     err: %v", message(messages), print(actual), err)
	}

	return EqualJSON(t, string(s), expected, messages...)
}

// Panics asserts that fn panics.
func Panics(t Testing, fn func(), messages ...string) bool {
	t.Helper()

	if panicked, _ := panics(fn); !panicked {
		return Failf(t, "%sShould panic", message(messages))
	}

	return true
}

// PanicsWith asserts that fn panics with the expected value.
//
// When expected is an error, the panic value must be an error matching it
// according to errors.Is. Otherwise, the panic value must be deeply equal to expected.
func PanicsWith(t Testing, fn func(), expected any, messages ...string) bool {
	t.Helper()

	panicked, value := panics(fn)
	if !panicked {
		return Failf(t, "%sShould panic\nexpected: %s", message(messages), print(expected))
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
