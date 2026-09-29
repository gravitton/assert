package assert

import (
	"encoding/json"
	"errors"
	"strings"
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
		return Failf(t, "%sShould be true", join(messages))
	}

	return true
}

// False asserts that the specified value is false.
func False(t Testing, condition bool, messages ...string) bool {
	t.Helper()

	if condition {
		return Failf(t, "%sShould be false", join(messages))
	}

	return true
}

// Nil asserts that actual is nil.
//
// A typed nil pointer, slice, map, channel, function or unsafe.Pointer is treated as nil.
func Nil(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if !isNil(actual) {
		return Failf(t, "%sShould be nil\n  actual: %s", join(messages), format(actual))
	}

	return true
}

// NotNil asserts that actual is NOT nil.
//
// A typed nil pointer, slice, map, channel, function or unsafe.Pointer is treated as nil.
func NotNil(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if isNil(actual) {
		return Failf(t, "%sShould not be nil\n  actual: %s", join(messages), format(actual))
	}

	return true
}

// Zero asserts that actual is the zero value of its type.
//
// A type's own IsZero() bool method is used when it has one, so a zero time.Time
// in any location is zero. A nil pointer is zero without calling the method.
// Unlike [Empty], a non-nil empty slice or map is not zero.
func Zero(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if !isZero(actual) {
		return Failf(t, "%sShould be zero\n  actual: %s", join(messages), format(actual))
	}

	return true
}

// NotZero asserts that actual is NOT the zero value of its type.
//
// A type's own IsZero() bool method is used when it has one, so a zero time.Time
// in any location is zero. A nil pointer is zero without calling the method.
// Unlike [Empty], a non-nil empty slice or map is not zero.
func NotZero(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if isZero(actual) {
		return Failf(t, "%sShould not be zero\n  actual: %s", join(messages), format(actual))
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

	if identical, reason := same(actual, expected); reason != valid {
		return Failf(t, "%s%s\n  actual: %s\nexpected: %s", join(messages), reason, format(actual), format(expected))
	} else if !identical {
		return Failf(t, "%sShould be same\n  actual: %s\nexpected: %s", join(messages), format(actual), format(expected))
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

	if identical, reason := same(actual, expected); reason != valid {
		return Failf(t, "%s%s\n  actual: %s\nexpected: %s", join(messages), reason, format(actual), format(expected))
	} else if identical {
		return Failf(t, "%sShould not be same\n  actual: %s", join(messages), format(actual))
	}

	return true
}

// Equal asserts that two objects are equal.
//
// Equality is determined with reflect.DeepEqual: pointers are compared by the
// values they reference, and two non-nil functions are never equal. A typed nil
// stored in an interface is not equal to an untyped nil; use [Nil] for that.
func Equal[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if !equal(actual, expected) {
		return Failf(t, "%sShould be equal\n  actual: %s\nexpected: %s", join(messages), format(actual), format(expected))
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
		return Failf(t, "%sShould not be equal\n  actual: %s", join(messages), format(actual))
	}

	return true
}

// EqualDelta asserts that two numeric values differ by at most delta.
//
// Fails when delta is negative or NaN. NaN is only equal to NaN.
func EqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if within, reason := equalDelta(actual, expected, delta); reason != valid {
		return Failf(t, "%s%s\n   delta: %s", join(messages), reason, format(delta))
	} else if !within {
		return Failf(t, "%sShould be equal in delta\n  actual: %s\nexpected: %s\n   delta: %s\n    diff: %s", join(messages), format(actual), format(expected), format(delta), formatDistance(actual, expected))
	}

	return true
}

// NotEqualDelta asserts that two numeric values differ by more than delta.
//
// Fails when delta is negative or NaN. NaN is only equal to NaN.
func NotEqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if within, reason := equalDelta(actual, expected, delta); reason != valid {
		return Failf(t, "%s%s\n   delta: %s", join(messages), reason, format(delta))
	} else if within {
		return Failf(t, "%sShould not be equal in delta\n  actual: %s\nexpected: %s\n   delta: %s\n    diff: %s", join(messages), format(actual), format(expected), format(delta), formatDistance(actual, expected))
	}

	return true
}

// Greater asserts that actual is greater than bound.
//
// Fails when either value is NaN.
func Greater[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, ok := compare(actual, bound); !ok || order <= 0 {
		return Failf(t, "%sShould be greater\n  actual: %s\n   bound: %s", join(messages), format(actual), format(bound))
	}

	return true
}

// GreaterOrEqual asserts that actual is greater than or equal to bound.
//
// Fails when either value is NaN.
func GreaterOrEqual[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, ok := compare(actual, bound); !ok || order < 0 {
		return Failf(t, "%sShould be greater or equal\n  actual: %s\n   bound: %s", join(messages), format(actual), format(bound))
	}

	return true
}

// Less asserts that actual is less than bound.
//
// Fails when either value is NaN.
func Less[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, ok := compare(actual, bound); !ok || order >= 0 {
		return Failf(t, "%sShould be less\n  actual: %s\n   bound: %s", join(messages), format(actual), format(bound))
	}

	return true
}

// LessOrEqual asserts that actual is less than or equal to bound.
//
// Fails when either value is NaN.
func LessOrEqual[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, ok := compare(actual, bound); !ok || order > 0 {
		return Failf(t, "%sShould be less or equal\n  actual: %s\n   bound: %s", join(messages), format(actual), format(bound))
	}

	return true
}

// Length asserts that actual has given length.
//
// Works with strings, arrays, pointers to arrays, slices, maps and channels. String length is measured in bytes, not runes.
func Length[S Iterable](t Testing, actual S, expected int, messages ...string) bool {
	t.Helper()

	if size, reason := length(actual); reason != valid {
		return Failf(t, "%s%s\n  actual: %s", join(messages), reason, format(actual))
	} else if size != expected {
		return Failf(t, "%sShould have length\n  actual: %s\n  length: %d\nexpected: %d", join(messages), format(actual), size, expected)
	}

	return true
}

// Empty asserts that actual has zero length.
//
// Works with strings, arrays, pointers to arrays, slices, maps and channels. String length is measured in bytes, not runes.
func Empty[S Iterable](t Testing, actual S, messages ...string) bool {
	t.Helper()

	if size, reason := length(actual); reason != valid {
		return Failf(t, "%s%s\n  actual: %s", join(messages), reason, format(actual))
	} else if size != 0 {
		return Failf(t, "%sShould be empty\n  actual: %s", join(messages), format(actual))
	}

	return true
}

// NotEmpty asserts that actual has non-zero length.
//
// Works with strings, arrays, pointers to arrays, slices, maps and channels. String length is measured in bytes, not runes.
func NotEmpty[S Iterable](t Testing, actual S, messages ...string) bool {
	t.Helper()

	if size, reason := length(actual); reason != valid {
		return Failf(t, "%s%s\n  actual: %s", join(messages), reason, format(actual))
	} else if size == 0 {
		return Failf(t, "%sShould not be empty\n  actual: %s", join(messages), format(actual))
	}

	return true
}

// Contains asserts that actual contains given element.
//
// Works with strings, arrays, slices and map values.
func Contains[S Iterable, E Comparable](t Testing, actual S, element E, messages ...string) bool {
	t.Helper()

	if found, reason := contains(actual, element); reason != valid {
		return Failf(t, "%s%s\n  actual: %s\n element: %s", join(messages), reason, format(actual), format(element))
	} else if !found {
		return Failf(t, "%sShould contain element\n  actual: %s\n element: %s", join(messages), format(actual), format(element))
	}

	return true
}

// NotContains asserts that actual does NOT contain given element.
//
// Works with strings, arrays, slices and map values.
func NotContains[S Iterable, E Comparable](t Testing, actual S, element E, messages ...string) bool {
	t.Helper()

	if found, reason := contains(actual, element); reason != valid {
		return Failf(t, "%s%s\n  actual: %s\n element: %s", join(messages), reason, format(actual), format(element))
	} else if found {
		return Failf(t, "%sShould not contain element\n  actual: %s\n element: %s", join(messages), format(actual), format(element))
	}

	return true
}

// EqualUnordered asserts that two arrays or slices have the same elements in any order.
//
// Elements are compared with reflect.DeepEqual and each must occur the same number
// of times in both. A nil slice and an empty slice are equal.
func EqualUnordered[S Iterable](t Testing, actual, expected S, messages ...string) bool {
	t.Helper()

	if extra, missing, reason := unorderedDifference(actual, expected); reason != valid {
		return Failf(t, "%s%s\n  actual: %s\nexpected: %s", join(messages), reason, format(actual), format(expected))
	} else if extra.Len() > 0 || missing.Len() > 0 {
		return Failf(t, "%sShould be equal in any order\n  actual: %s\nexpected: %s\n   extra: %s\n missing: %s", join(messages), format(actual), format(expected), formatGoSyntax(extra), formatGoSyntax(missing))
	}

	return true
}

// HasPrefix asserts that actual begins with prefix.
//
// Works with strings and slices. Slice elements are compared with reflect.DeepEqual.
func HasPrefix[S Iterable](t Testing, actual, prefix S, messages ...string) bool {
	t.Helper()

	if found, reason := hasPrefix(actual, prefix); reason != valid {
		return Failf(t, "%s%s\n  actual: %s\n  prefix: %s", join(messages), reason, format(actual), format(prefix))
	} else if !found {
		return Failf(t, "%sShould have prefix\n  actual: %s\n  prefix: %s", join(messages), format(actual), format(prefix))
	}

	return true
}

// HasSuffix asserts that actual ends with suffix.
//
// Works with strings and slices. Slice elements are compared with reflect.DeepEqual.
func HasSuffix[S Iterable](t Testing, actual, suffix S, messages ...string) bool {
	t.Helper()

	if found, reason := hasSuffix(actual, suffix); reason != valid {
		return Failf(t, "%s%s\n  actual: %s\n  suffix: %s", join(messages), reason, format(actual), format(suffix))
	} else if !found {
		return Failf(t, "%sShould have suffix\n  actual: %s\n  suffix: %s", join(messages), format(actual), format(suffix))
	}

	return true
}

// Error asserts that error is NOT nil.
//
// A typed nil stored in err (e.g. (*MyError)(nil)) is an error, as err != nil sees it.
func Error(t Testing, err error, messages ...string) bool {
	t.Helper()

	if err == nil {
		return Failf(t, "%sShould be error", join(messages))
	}

	return true
}

// NoError asserts that error is nil.
//
// A typed nil stored in err (e.g. (*MyError)(nil)) is an error, as err != nil sees it.
func NoError(t Testing, err error, messages ...string) bool {
	t.Helper()

	if err != nil {
		return Failf(t, "%sShould not be error\n%s", join(messages), formatError(err))
	}

	return true
}

// ErrorIs asserts that error matches given target according to errors.Is.
func ErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if !errors.Is(err, target) {
		return Failf(t, "%sShould match error\n%s\n  target: %s", join(messages), formatError(err), formatGoSyntax(target))
	}

	return true
}

// NotErrorIs asserts that error does NOT match given target according to errors.Is.
func NotErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if errors.Is(err, target) {
		return Failf(t, "%sShould not match error\n%s\n  target: %s", join(messages), formatError(err), formatGoSyntax(target))
	}

	return true
}

// ErrorAs asserts that error can be assigned to target according to errors.As.
//
// As with errors.As, target must be a non-nil pointer to a type implementing
// error or to any interface type; otherwise ErrorAs fails.
func ErrorAs(t Testing, err error, target any, messages ...string) bool {
	t.Helper()

	if assignable, reason := errorAs(err, target); reason != valid {
		return Failf(t, "%s%s\n  target: %T", join(messages), reason, target)
	} else if !assignable {
		return Failf(t, "%sShould be assignable to target\n%s\n  target: %T", join(messages), formatError(err), target)
	}

	return true
}

// ErrorContains asserts that error is NOT nil and its message contains substr.
//
// A typed nil stored in err (e.g. (*MyError)(nil)) is an error, as err != nil sees it.
func ErrorContains(t Testing, err error, substr string, messages ...string) bool {
	t.Helper()

	if err == nil {
		return Failf(t, "%sShould be error\n  substr: %s", join(messages), format(substr))
	} else if !strings.Contains(err.Error(), substr) {
		return Failf(t, "%sShould contain substring\n%s\n  substr: %s", join(messages), formatError(err), format(substr))
	}

	return true
}

// Matches asserts that a string matches the given regular expression.
func Matches(t Testing, actual, pattern string, messages ...string) bool {
	t.Helper()

	if matched, err := matches(actual, pattern); err != nil {
		return Failf(t, "%sShould be valid regexp\n pattern: %s\n     err: %v", join(messages), truncate(pattern), err)
	} else if !matched {
		return Failf(t, "%sShould match regexp\n  actual: %s\n pattern: %s", join(messages), format(actual), truncate(pattern))
	}

	return true
}

// NotMatches asserts that a string does NOT match the given regular expression.
func NotMatches(t Testing, actual, pattern string, messages ...string) bool {
	t.Helper()

	if matched, err := matches(actual, pattern); err != nil {
		return Failf(t, "%sShould be valid regexp\n pattern: %s\n     err: %v", join(messages), truncate(pattern), err)
	} else if matched {
		return Failf(t, "%sShould not match regexp\n  actual: %s\n pattern: %s", join(messages), format(actual), truncate(pattern))
	}

	return true
}

// EqualJSON asserts that JSON strings are semantically equal.
//
// Numbers are compared by exact decimal value, so 1.0 equals 1, 1e2 equals 100
// and large integers and exponents keep their precision.
func EqualJSON(t Testing, actual, expected string, messages ...string) bool {
	t.Helper()

	actualJSON, err := decodeJSON(actual)
	if err != nil {
		return Failf(t, "%sShould be valid JSON\n  actual: %s\n     err: %v", join(messages), truncate(actual), err)
	}

	expectedJSON, err := decodeJSON(expected)
	if err != nil {
		return Failf(t, "%sShould be valid JSON\nexpected: %s\n     err: %v", join(messages), truncate(expected), err)
	}

	if !equal(actualJSON, expectedJSON) {
		return Failf(t, "%sShould be equal JSON\n  actual: %s\nexpected: %s", join(messages), truncate(actual), truncate(expected))
	}

	return true
}

// JSON asserts that actual can be marshaled to expected JSON string.
func JSON(t Testing, actual any, expected string, messages ...string) bool {
	t.Helper()

	s, err := json.Marshal(actual)
	if err != nil {
		return Failf(t, "%sShould be marshalable\n  actual: %s\n     err: %v", join(messages), format(actual), err)
	}

	return EqualJSON(t, string(s), expected, messages...)
}

// Panics asserts that fn panics.
//
// A panic(nil) is recognised regardless of the GODEBUG panicnil setting.
func Panics(t Testing, fn func(), messages ...string) bool {
	t.Helper()

	if panicked, _ := panics(fn); !panicked {
		return Failf(t, "%sShould panic", join(messages))
	}

	return true
}

// PanicsWith asserts that fn panics with the expected value.
//
// The panic value must be deeply equal to expected. When expected is an error,
// a panic value matching it according to errors.Is is accepted as well.
// A panic(nil) is reported as a nil value regardless of the GODEBUG panicnil setting.
func PanicsWith(t Testing, fn func(), expected any, messages ...string) bool {
	t.Helper()

	panicked, value := panics(fn)
	if !panicked {
		return Failf(t, "%sShould panic\nexpected: %s", join(messages), format(expected))
	}

	if !panicsWith(value, expected) {
		return Failf(t, "%sShould panic with value\n  actual: %s\nexpected: %s", join(messages), format(value), format(expected))
	}

	return true
}

// NotPanics asserts that fn does NOT panic.
//
// A panic(nil) is recognised regardless of the GODEBUG panicnil setting.
func NotPanics(t Testing, fn func(), messages ...string) bool {
	t.Helper()

	if panicked, value := panics(fn); panicked {
		return Failf(t, "%sShould not panic\n  value: %s", join(messages), format(value))
	}

	return true
}
