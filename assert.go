package assert

import (
	"encoding/json"
	"errors"
)

// Testing is an interface wrapper around *testing.T
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// Fail reports a failure message through
func Fail(t Testing, message string) bool {
	t.Helper()

	t.Errorf(message)

	return false
}

// Failf reports a failure formatted message through
func Failf(t Testing, format string, args ...any) bool {
	t.Helper()

	t.Errorf(format, args...)

	return false
}

// True asserts that the specified value is true.
func True(t Testing, condition bool) bool {
	t.Helper()

	if !condition {
		return Fail(t, "Should be true")
	}

	return true
}

// False asserts that the specified value is false.
func False(t Testing, condition bool) bool {
	t.Helper()

	if condition {
		return Fail(t, "Should be false")
	}

	return true
}

// Equal asserts that two objects are equal.
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func Equal[T Comparable](t Testing, actual, expected T) bool {
	t.Helper()

	if !equal(actual, expected) {
		return Failf(t, "Should be equal:\n  actual: %s\nexpected: %s", print(actual), print(expected))
	}

	return true
}

// NotEqual asserts that the specified values are NOT equal.
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqual[T Comparable](t Testing, actual, expected T) bool {
	t.Helper()

	if equal(actual, expected) {
		return Failf(t, "Should not be equal\n  actual: %s", print(actual))
	}

	return true
}

// EqualDelta asserts that two numeric values difference is less then delta.
//
// Method panics if delta value is not positive.
func EqualDelta[T Numeric](t Testing, actual, expected, delta T) bool {
	t.Helper()

	if !equalDelta(actual, expected, delta) {
		return Failf(t, "Should be equal in delta:\n  actual: %s\nexpected: %s", print(actual), print(expected))
	}

	return true
}

// Same asserts that two pointers reference the same object.
//
// Both arguments must be pointer variables. Pointer variable equality is
// determined based on the equality of both type and value.
func Same[T Reference](t Testing, actual, expected T) bool {
	t.Helper()

	if valid, ok := same[T](expected, actual); !ok {
		return Failf(t, "Should be pointers\n  actual: %s\nexpected: %s", print(actual), print(expected))
	} else if !valid {
		return Failf(t, "Should be same\n  actual: %s\nexpected: %s", print(actual), print(expected))
	}

	return true
}

// NotSame asserts that two pointers do NOT reference the same object.
//
// Both arguments must be pointer variables. Pointer variable equality is
// determined based on the equality of both type and value.
func NotSame[T Reference](t Testing, actual, expected T) bool {
	t.Helper()

	if valid, ok := same(expected, actual); !ok {
		Failf(t, "Should be pointers\n  actual: %s\nexpected: %s", print(actual), print(expected))
	} else if valid {
		return Failf(t, "Should not be same\n  actual: %s", print(actual))
	}

	return true
}

// Length asserts that object have given length.
func Length[S Iterable[any]](t Testing, object S, expected int) bool {
	t.Helper()

	if actual := length(object); actual != expected {
		return Failf(t, "Should have element length\n  object: %#v\n  actual: %d\nexpected: %d", object, actual, expected)
	}

	return true
}

// Contains asserts that object contains given element
//
// Works with strings, arrays, slices, maps values and channels
func Contains[S Iterable[E], E Comparable](t Testing, object S, element E) bool {
	t.Helper()

	if found, ok := contains(object, element); !ok {
		return Failf(t, "Should be iterable\n  object: %#v", object)
	} else if !found {
		return Failf(t, "Should contain element\n  object: %#v\n element: %#v", object, element)
	}

	return true
}

// NotContains asserts that object do NOT contains given element
//
// Works with strings, arrays, slices, maps values and channels
func NotContains[S Iterable[E], E Comparable](t Testing, object S, element E) bool {
	t.Helper()

	if found, ok := contains(object, element); !ok {
		Failf(t, "Should be iterable\n  object: %#v", object)
	} else if found {
		return Failf(t, "Should not contain element\n  object: %#v\n element: %#v", object, element)
	}

	return true
}

// Error asserts that error is NOT nil
func Error(t Testing, err error) bool {
	t.Helper()

	if err == nil {
		return Failf(t, "Should be error")
	}

	return true
}

// NoError asserts that error is nil
func NoError(t Testing, err error) bool {
	t.Helper()

	if err != nil {
		return Failf(t, "Should not be error\n   error: %#v", err)
	}

	return true
}

// ErrorIs asserts that error is unwrappable to given target
func ErrorIs(t Testing, err error, target error) bool {
	t.Helper()

	if !errors.Is(err, target) {
		return Failf(t, "Should be same error\n   error: %#v\n  target: %#v", err, target)
	}

	return true
}

// NotErrorIs asserts that error is NOT unwrappable to given target
func NotErrorIs(t Testing, err error, target error) bool {
	t.Helper()

	if errors.Is(err, target) {
		return Failf(t, "Should not be same error\n   error: %#v", err, target)
	}

	return true
}

// EqualJSON asserts that JSON strings are equal
func EqualJSON(t Testing, actual, expected string) bool {
	t.Helper()

	var actualJSON, expectedJSON any

	if err := json.Unmarshal([]byte(actual), &actualJSON); err != nil {
		return Failf(t, "Should be valid JSON\n  actual: %s\n     err: %v", expected, err)
	}

	if err := json.Unmarshal([]byte(expected), &expectedJSON); err != nil {
		return Failf(t, "Should be valid JSON\nexpected: %s\n     err: %v", expected, err)
	}

	return Equal(t, expectedJSON, actualJSON)
}

// JSON asserts that object can be marshall to expected JSON string
func JSON(t Testing, actual any, expected string) bool {
	t.Helper()

	s, err := json.Marshal(actual)

	return NoError(t, err) && EqualJSON(t, string(s), expected)
}
