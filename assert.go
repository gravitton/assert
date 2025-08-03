package assert

import (
	"reflect"
)

// Testing is an interface wrapper around *testing.T
type Testing interface {
	Errorf(format string, args ...any)
}

// helper is an interface wrapper around *testing.T
type helper interface {
	Helper()
}

// True asserts that the specified value is true.
//
//	assert.True(t, condition)
func True(t Testing, actual bool) bool {
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	if !actual {
		return Fail(t, "Should be true")
	}

	return true
}

// False asserts that the specified value is false.
//
//	assert.False(t, condition)
func False(t Testing, actual bool) bool {
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	if actual {
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
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	if !equal(actual, expected) {
		return Failf(t, "Should be equal: \nactual: %#v\n  expected: %#v", actual, expected)
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
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	if equal(actual, expected) {
		return Failf(t, "Should not be equal: %#v", actual)
	}

	return true
}

func equal(actual, expected any) bool {
	if actual == nil || expected == nil {
		return actual == expected
	}

	return reflect.DeepEqual(actual, expected)
}

// Same asserts that two pointers reference the same object.
//
//	assert.Same(t, actual, expected)
//
// Both arguments must be pointer variables. Pointer variable sameness is
// determined based on the equality of both type and value.
func Same(t Testing, actual, expected any) bool {
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	if !same(expected, actual) {
		return Failf(t, "Should be same: \nactual: %[1]p %#[1]v\n  expected: %[2]p %#[2]v", actual, expected)
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
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	if same(expected, actual) {
		return Failf(t, "Should not be same: %[1]p %#[1]v", actual)
	}

	return true
}

func same(actual, expected any) bool {
	actualPtr, expectedPtr := reflect.ValueOf(actual), reflect.ValueOf(expected)
	if actualPtr.Kind() != reflect.Ptr || expectedPtr.Kind() != reflect.Ptr {
		return false
	}

	actualType, expectedType := reflect.TypeOf(actual), reflect.TypeOf(expected)
	if actualType != expectedType {
		return false
	}

	// compare pointer addresses
	return actual == expected
}

// Fail reports a failure through
func Fail(t Testing, message string) bool {
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	t.Errorf(message)

	return false
}

// Failf reports a failure through
func Failf(t Testing, format string, args ...any) bool {
	if h, ok := t.(helper); ok {
		h.Helper()
	}

	t.Errorf(format, args...)

	return false
}
