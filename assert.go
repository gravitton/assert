package assert

import (
	"errors"
)

// Testing is an interface wrapper around *testing.T
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// True asserts that the specified value is true.
//
//	assert.True(t, condition)
func True(t Testing, condition bool) bool {
	t.Helper()

	if !condition {
		return Fail(t, "Should be true")
	}

	return true
}

// False asserts that the specified value is false.
//
//	assert.False(t, condition)
func False(t Testing, condition bool) bool {
	t.Helper()

	if condition {
		return Fail(t, "Should be false")
	}

	return true
}

// Equal asserts that two objects are equal.
//
//	assert.Equal(t, actual, expected)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func Equal(t Testing, actual, expected any) bool {
	t.Helper()

	if !equal(actual, expected) {
		return Failf(t, "Should be equal:\n  actual: %#v\nexpected: %#v", actual, expected)
	}

	return true
}

// NotEqual asserts that the specified values are NOT equal.
//
//	assert.NotEqual(t, actual, expected)
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqual(t Testing, actual, expected any) bool {
	t.Helper()

	if equal(actual, expected) {
		return Failf(t, "Should not be equal\n  actual: %#v", actual)
	}

	return true
}

// Same asserts that two pointers reference the same object.
//
//	assert.Same(t, actual, expected)
//
// Both arguments must be pointer variables. Pointer variable equality is
// determined based on the equality of both type and value.
func Same(t Testing, actual, expected any) bool {
	t.Helper()

	if !same(expected, actual) {
		return Failf(t, "Should be same\n  actual: %[1]p %#[1]v\nexpected: %[2]p %#[2]v", actual, expected)
	}

	return true
}

// NotSame asserts that two pointers do NOT reference the same object.
//
//	assert.NotSame(t, actual, expected)
//
// Both arguments must be pointer variables. Pointer variable equality is
// determined based on the equality of both type and value.
func NotSame(t Testing, actual, expected any) bool {
	t.Helper()

	if same(expected, actual) {
		return Failf(t, "Should not be same\n  actual: %[1]p %#[1]v", actual)
	}

	return true
}

// Length asserts that object have given length.
//
//	assert.Length(t, object, expected)
func Length(t Testing, object any, expected int) bool {
	t.Helper()

	if actual := length(object); actual != expected {
		return Failf(t, "Should have element length\n  object: %#v\n  actual: %d\nexpected: %d", object, actual, expected)
	}

	return true
}

// Contains asserts that object contains given element
//
//	assert.Contains(t, object, element)
//
// Works with strings, arrays, slices, maps values and channels
func Contains(t Testing, object, element any) bool {
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
//	assert.NotContains(t, object, element)
//
// Works with strings, arrays, slices, maps values and channels
func NotContains(t Testing, object any, element any) bool {
	t.Helper()

	if found, ok := contains(object, element); !ok {
		return Failf(t, "Should be iterable\n  object: %#v", object)
	} else if found {
		return Failf(t, "Should not contain element\n  object: %#v\n element: %#v", object, element)
	}

	return true
}

// Error asserts that error is NOT nil
//
//	assert.Error(t, err)
func Error(t Testing, err error) bool {
	t.Helper()

	if err == nil {
		return Failf(t, "Should be error")
	}

	return true
}

// NoError asserts that error is nil
//
//	assert.NoError(t, err)
func NoError(t Testing, err error) bool {
	t.Helper()

	if err != nil {
		return Failf(t, "Should not be error\n   error: %#v", err)
	}

	return true
}

// ErrorIs asserts that error is unwrappable to given target
//
//	assert.ErrorIs(t, err)
func ErrorIs(t Testing, err error, target error) bool {
	t.Helper()

	if !errors.Is(err, target) {
		return Failf(t, "Should be same error\n   error: %#v\n  target: %#v", err, target)
	}

	return true
}

// NotErrorIs asserts that error is NOT unwrappable to given target
//
//	assert.NotErrorIs(t, err)
func NotErrorIs(t Testing, err error, target error) bool {
	t.Helper()

	if errors.Is(err, target) {
		return Failf(t, "Should not be same error\n   error: %#v", err, target)
	}

	return true
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
