package assert

import (
	"github.com/gravitton/assert/internal"
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

	if !internal.Assert(condition) {
		return Fail(t, "Should be true")
	}

	return true
}

// False asserts that the specified value is false.
//
//	assert.False(t, condition)
func False(t Testing, condition bool) bool {
	t.Helper()

	if internal.Assert(condition) {
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

	if !internal.Equal(actual, expected) {
		return Failf(t, "Should be equal:\n  object: %#v\nelement: %#v", actual, expected)
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

	if internal.Equal(actual, expected) {
		return Failf(t, "Should not be equal\n  object: %#v", actual)
	}

	return true
}

// Same asserts that two pointers reference the same object.
//
//	assert.Same(t, actual, expected)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func Same(t Testing, actual, expected any) bool {
	t.Helper()

	if !internal.Same(expected, actual) {
		return Failf(t, "Should be same\n  object: %[1]p %#[1]v\nelement: %[2]p %#[2]v", actual, expected)
	}

	return true
}

// NotSame asserts that two pointers do NOT reference the same object.
//
//	assert.NotSame(t, actual, expected)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func NotSame(t Testing, actual, expected any) bool {
	t.Helper()

	if internal.Same(expected, actual) {
		return Failf(t, "Should not be same\n  object: %[1]p %#[1]v", actual)
	}

	return true
}

func Length(t Testing, object any, expected int) bool {
	t.Helper()

	if actual := internal.Length(object); actual != expected {
		return Failf(t, "Should have element length\n  object: %#v\n  object: %d\nelement: %d", object, actual, expected)
	}

	return true
}

func Contains(t Testing, object, element any) bool {
	t.Helper()

	if found, ok := internal.Contains(object, element); !ok {
		return Failf(t, "Should be iterable\n  object: %#v", object)
	} else if !found {
		return Failf(t, "Should contain element\n  object: %#v\n element: %#v", object, element)
	}

	return true
}

func NotContains(t Testing, object any, element any) bool {
	t.Helper()

	if found, ok := internal.Contains(object, element); !ok {
		return Failf(t, "Should be iterable\n  object: %#v", object)
	} else if found {
		return Failf(t, "Should not contain element\n  object: %#v\n element: %#v", object, element)
	}

	return true
}

func Error(t Testing, err error) bool {
	t.Helper()

	if !internal.Error(err) {
		return Failf(t, "Should be error")
	}

	return true
}

func NoError(t Testing, err error) bool {
	t.Helper()

	if internal.Error(err) {
		return Failf(t, "Should not be error\n   error: %#v", err)
	}

	return true
}

func ErrorIs(t Testing, err error, target error) bool {
	t.Helper()

	if !internal.ErrorIs(err, target) {
		return Failf(t, "Should be same error\n   error: %#v\n  target: %#v", err, target)
	}

	return true
}

func NotErrorIs(t Testing, err error, target error) bool {
	t.Helper()

	if internal.ErrorIs(err, target) {
		return Failf(t, "Should not be same error\n   error: %#v", err, target)
	}

	return true
}

// Fail reports a failure through
func Fail(t Testing, message string) bool {
	t.Helper()

	t.Errorf(message)

	return false
}

// Failf reports a failure through
func Failf(t Testing, format string, args ...any) bool {
	t.Helper()

	t.Errorf(format, args...)

	return false
}
