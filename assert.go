package assert

import (
	"reflect"
)

type Testing interface {
	Errorf(format string, args ...any)
}

type Helper interface {
	Helper()
}

func True(t Testing, actual bool) bool {
	if h, ok := t.(Helper); ok {
		h.Helper()
	}

	if !actual {
		return Fail(t, "Should be true")
	}

	return true
}

func False(t Testing, actual bool) bool {
	if h, ok := t.(Helper); ok {
		h.Helper()
	}

	if actual {
		return Fail(t, "Should be false")
	}

	return true
}

func Equal(t Testing, actual, expected any) bool {
	if h, ok := t.(Helper); ok {
		h.Helper()
	}

	if !equal(actual, expected) {
		return Failf(t, "Should be equal: \nactual: %#v\n  expected: %#v", actual, expected)
	}

	return true
}

func NotEqual(t Testing, actual, expected any) bool {
	if h, ok := t.(Helper); ok {
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

func Same(t Testing, actual, expected any) bool {
	if h, ok := t.(Helper); ok {
		h.Helper()
	}

	if !same(expected, actual) {
		return Failf(t, "Should be same: \nactual: %[1]p %#[1]v\n  expected: %[2]p %#[2]v", actual, expected)
	}

	return true
}

func NotSame(t Testing, actual, expected any) bool {
	if h, ok := t.(Helper); ok {
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

func Fail(t Testing, message string) bool {
	if h, ok := t.(Helper); ok {
		h.Helper()
	}

	t.Errorf(message)

	return false
}

func Failf(t Testing, format string, args ...any) bool {
	if h, ok := t.(Helper); ok {
		h.Helper()
	}

	t.Errorf(format, args...)

	return false
}
