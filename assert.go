package assert

import (
	"encoding/json"
)

// Testing is the part of *testing.T the assertions use.
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// Fail reports message with t.Errorf and returns false.
func Fail(t Testing, message string) bool {
	t.Helper()

	t.Errorf("%s", message)

	return false
}

// Failf reports the message formatted from format and args with t.Errorf and returns false.
func Failf(t Testing, format string, args ...any) bool {
	t.Helper()

	t.Errorf(format, args...)

	return false
}

// True asserts that condition is true.
func True(t Testing, condition bool, messages ...string) bool {
	t.Helper()

	if !condition {
		return fail(t, messages, "Should be true")
	}

	return true
}

// False asserts that condition is false.
func False(t Testing, condition bool, messages ...string) bool {
	t.Helper()

	if condition {
		return fail(t, messages, "Should be false")
	}

	return true
}

// Nil asserts that actual is nil. A typed nil stored in actual is nil.
func Nil(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if !isNil(actual) {
		return fail(t, messages, "Should be nil", field{"actual", actual})
	}

	return true
}

// NotNil asserts that actual is not nil.
func NotNil(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if isNil(actual) {
		return fail(t, messages, "Should not be nil", field{"actual", actual})
	}

	return true
}

// Zero asserts that actual is the zero value of its type.
//
// When actual is not a pointer and has an IsZero() bool method, the method decides, even when it is
// promoted from an embedded field and ignores the other fields. A pointer is zero only when nil,
// and a non-nil empty slice or map is not zero.
func Zero(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if !isZero(actual) {
		return fail(t, messages, "Should be zero", field{"actual", actual})
	}

	return true
}

// NotZero asserts that actual is not the zero value of its type.
func NotZero(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if isZero(actual) {
		return fail(t, messages, "Should not be zero", field{"actual", actual})
	}

	return true
}

// Same asserts that actual and expected are the same reference: equal type and address, and for
// slices equal length and capacity.
//
// Two nil references of one type are the same, and so may be separately allocated zero-size values
// or zero-capacity slices.
func Same[T Reference](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if identical, reason := same(actual, expected); reason != valid {
		return fail(t, messages, string(reason), field{"actual", typed(actual)}, field{"expected", typed(expected)})
	} else if !identical {
		return fail(t, messages, "Should be same", field{"actual", identity(actual)}, field{"expected", identity(expected)})
	}

	return true
}

// NotSame asserts that actual and expected are not the same reference. Unlike Same, it holds for
// references of different types.
func NotSame[T Reference](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if identical, reason := same(actual, expected); reason == notReference {
		return fail(t, messages, string(reason), field{"actual", typed(actual)}, field{"expected", typed(expected)})
	} else if identical {
		return fail(t, messages, "Should not be same", field{"actual", identity(actual)})
	}

	return true
}

// Equal asserts that actual and expected are deeply equal.
//
// Pointers are compared by what they point to, a typed nil stored in an interface differs from an
// untyped nil, NaN differs from itself, and functions are equal only when both are nil.
func Equal[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if !equal(actual, expected) {
		return fail(t, messages, "Should be equal", comparisonFields(actual, expected)...)
	}

	return true
}

// NotEqual asserts that actual and expected are not deeply equal.
func NotEqual[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if equal(actual, expected) {
		return fail(t, messages, "Should not be equal", field{"actual", actual})
	}

	return true
}

// EqualDelta asserts that actual and expected differ by at most delta.
//
// Integers are compared exactly. Two NaNs are equal; a NaN and a number never are. It fails for a
// negative or NaN delta.
func EqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if within, reason := equalDelta(actual, expected, delta); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"expected", expected}, field{"delta", delta})
	} else if !within {
		return fail(t, messages, "Should be equal in delta", field{"actual", actual}, field{"expected", expected}, field{"delta", delta}, field{"diff", distance(actual, expected)})
	}

	return true
}

// NotEqualDelta asserts that actual and expected differ by more than delta.
func NotEqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if within, reason := equalDelta(actual, expected, delta); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"expected", expected}, field{"delta", delta})
	} else if within {
		return fail(t, messages, "Should not be equal in delta", field{"actual", actual}, field{"expected", expected}, field{"delta", delta}, field{"diff", distance(actual, expected)})
	}

	return true
}

// Greater asserts that actual is greater than bound. It fails for a NaN on either side.
func Greater[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, reason := compare(actual, bound); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"bound", bound})
	} else if order <= 0 {
		return fail(t, messages, "Should be greater", field{"actual", actual}, field{"bound", bound})
	}

	return true
}

// GreaterOrEqual asserts that actual is greater than or equal to bound. It fails for a NaN on
// either side.
func GreaterOrEqual[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, reason := compare(actual, bound); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"bound", bound})
	} else if order < 0 {
		return fail(t, messages, "Should be greater or equal", field{"actual", actual}, field{"bound", bound})
	}

	return true
}

// Less asserts that actual is less than bound. It fails for a NaN on either side.
func Less[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, reason := compare(actual, bound); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"bound", bound})
	} else if order >= 0 {
		return fail(t, messages, "Should be less", field{"actual", actual}, field{"bound", bound})
	}

	return true
}

// LessOrEqual asserts that actual is less than or equal to bound. It fails for a NaN on either
// side.
func LessOrEqual[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, reason := compare(actual, bound); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"bound", bound})
	} else if order > 0 {
		return fail(t, messages, "Should be less or equal", field{"actual", actual}, field{"bound", bound})
	}

	return true
}

// Length asserts that actual has length expected, in bytes for a string.
func Length[S Iterable](t Testing, actual S, expected int, messages ...string) bool {
	t.Helper()

	if size, reason := length(actual); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual})
	} else if size != expected {
		return fail(t, messages, "Should have length", field{"actual", actual}, field{"length", size}, field{"expected", expected})
	}

	return true
}

// Empty asserts that actual has zero length.
func Empty[S Iterable](t Testing, actual S, messages ...string) bool {
	t.Helper()

	if size, reason := length(actual); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual})
	} else if size != 0 {
		return fail(t, messages, "Should be empty", field{"actual", actual})
	}

	return true
}

// NotEmpty asserts that actual has non-zero length.
func NotEmpty[S Iterable](t Testing, actual S, messages ...string) bool {
	t.Helper()

	if size, reason := length(actual); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual})
	} else if size == 0 {
		return fail(t, messages, "Should not be empty", field{"actual", actual})
	}

	return true
}

// Contains asserts that actual contains element: a substring of a string, or an element of an
// array, array pointer, slice or the values of a map, not its keys.
//
// element must have the type of the elements of actual, or be assignable to it when that is an
// interface; a substring must have the type of the string. It fails for a nil array pointer.
func Contains[S Iterable, E Comparable](t Testing, actual S, element E, messages ...string) bool {
	t.Helper()

	if found, reason := contains(actual, element); reason == elementMismatch {
		return fail(t, messages, string(reason), field{"actual", typed(actual)}, field{"element", typed(element)})
	} else if reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"element", element})
	} else if !found {
		return fail(t, messages, "Should contain element", field{"actual", actual}, field{"element", element})
	}

	return true
}

// NotContains asserts that actual does not contain element.
func NotContains[S Iterable, E Comparable](t Testing, actual S, element E, messages ...string) bool {
	t.Helper()

	if found, reason := contains(actual, element); reason == elementMismatch {
		return fail(t, messages, string(reason), field{"actual", typed(actual)}, field{"element", typed(element)})
	} else if reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"element", element})
	} else if found {
		return fail(t, messages, "Should not contain element", field{"actual", actual}, field{"element", element})
	}

	return true
}

// EqualUnordered asserts that the arrays, array pointers or slices actual and expected have the
// same elements in any order, each occurring as often in both.
//
// actual and expected must have the same type. A nil slice equals an empty one. It fails for a nil
// array pointer.
func EqualUnordered[S Iterable](t Testing, actual, expected S, messages ...string) bool {
	t.Helper()

	if extra, missing, reason := unorderedDifference(actual, expected); reason == typeMismatch {
		return fail(t, messages, string(reason), field{"actual", typed(actual)}, field{"expected", typed(expected)})
	} else if reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"expected", expected})
	} else if extra.Len() > 0 || missing.Len() > 0 {
		return fail(t, messages, "Should be equal in any order", field{"actual", actual}, field{"expected", expected}, field{"extra", extra.Interface()}, field{"missing", missing.Interface()})
	}

	return true
}

// HasPrefix asserts that the string or slice actual begins with prefix. actual and prefix must
// have the same type.
func HasPrefix[S Iterable](t Testing, actual, prefix S, messages ...string) bool {
	t.Helper()

	if found, reason := hasPrefix(actual, prefix); reason == typeMismatch {
		return fail(t, messages, string(reason), field{"actual", typed(actual)}, field{"prefix", typed(prefix)})
	} else if reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"prefix", prefix})
	} else if !found {
		return fail(t, messages, "Should have prefix", field{"actual", actual}, field{"prefix", prefix})
	}

	return true
}

// HasSuffix asserts that the string or slice actual ends with suffix. actual and suffix must have
// the same type.
func HasSuffix[S Iterable](t Testing, actual, suffix S, messages ...string) bool {
	t.Helper()

	if found, reason := hasSuffix(actual, suffix); reason == typeMismatch {
		return fail(t, messages, string(reason), field{"actual", typed(actual)}, field{"suffix", typed(suffix)})
	} else if reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"suffix", suffix})
	} else if !found {
		return fail(t, messages, "Should have suffix", field{"actual", actual}, field{"suffix", suffix})
	}

	return true
}

// Error asserts that err is not nil. A typed nil stored in err is not nil, as with err != nil.
func Error(t Testing, err error, messages ...string) bool {
	t.Helper()

	if err == nil {
		return fail(t, messages, string(notError))
	}

	return true
}

// NoError asserts that err is nil.
func NoError(t Testing, err error, messages ...string) bool {
	t.Helper()

	if err != nil {
		return fail(t, messages, "Should not be error", errorFields(err)...)
	}

	return true
}

// ErrorIs asserts that errors.Is(err, target) holds.
func ErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if !errorIs(err, target) {
		return fail(t, messages, "Should match error", errorFields(err, field{"target", target})...)
	}

	return true
}

// NotErrorIs asserts that errors.Is(err, target) does not hold.
func NotErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if errorIs(err, target) {
		return fail(t, messages, "Should not match error", errorFields(err, field{"target", target})...)
	}

	return true
}

// ErrorAs asserts that errors.As(err, target) holds, which sets target to the matched error.
//
// It fails, rather than panics, for a target that is not a non-nil pointer to an error type or an
// interface.
func ErrorAs(t Testing, err error, target any, messages ...string) bool {
	t.Helper()

	if assignable, reason := errorAs(err, target); reason != valid {
		return fail(t, messages, string(reason), field{"target", typeName(target)})
	} else if !assignable {
		return fail(t, messages, "Should be assignable to target", errorFields(err, field{"target", typeName(target)})...)
	}

	return true
}

// ErrorContains asserts that err is not nil and its message contains substr.
//
// A typed nil stored in err is not nil; its message is <nil> when its Error method dereferences
// the receiver.
func ErrorContains(t Testing, err error, substr string, messages ...string) bool {
	t.Helper()

	if found, reason := errorContains(err, substr); reason != valid {
		return fail(t, messages, string(reason), field{"substr", substr})
	} else if !found {
		return fail(t, messages, "Should contain substring", errorFields(err, field{"substr", substr})...)
	}

	return true
}

// Matches asserts that actual matches the regular expression pattern, anywhere unless the pattern
// is anchored. It fails for an invalid pattern.
func Matches(t Testing, actual, pattern string, messages ...string) bool {
	t.Helper()

	if matched, err := matches(actual, pattern); err != nil {
		return fail(t, messages, "Should be valid regexp", field{"pattern", text(pattern)}, field{"err", errorText(err)})
	} else if !matched {
		return fail(t, messages, "Should match regexp", field{"actual", actual}, field{"pattern", text(pattern)})
	}

	return true
}

// NotMatches asserts that actual does not match the regular expression pattern.
func NotMatches(t Testing, actual, pattern string, messages ...string) bool {
	t.Helper()

	if matched, err := matches(actual, pattern); err != nil {
		return fail(t, messages, "Should be valid regexp", field{"pattern", text(pattern)}, field{"err", errorText(err)})
	} else if matched {
		return fail(t, messages, "Should not match regexp", field{"actual", actual}, field{"pattern", text(pattern)})
	}

	return true
}

// EqualJSON asserts that the JSON documents actual and expected are deeply equal once decoded.
//
// Both are decoded into any, so key order and whitespace do not matter, the last duplicate key
// wins, and numbers become float64: integers beyond 2^53 may compare equal when they differ, and
// numbers beyond the float64 range are invalid JSON.
func EqualJSON(t Testing, actual, expected string, messages ...string) bool {
	t.Helper()

	actualJSON, err := decodeJSON(actual)
	if err != nil {
		return fail(t, messages, "Should be valid JSON", field{"actual", text(actual)}, field{"err", errorText(err)})
	}

	expectedJSON, err := decodeJSON(expected)
	if err != nil {
		return fail(t, messages, "Should be valid JSON", field{"expected", text(expected)}, field{"err", errorText(err)})
	}

	if !equal(actualJSON, expectedJSON) {
		return fail(t, messages, "Should be equal JSON", field{"actual", text(actual)}, field{"expected", text(expected)})
	}

	return true
}

// JSON asserts that actual marshals with json.Marshal to a document equal to expected under
// EqualJSON. A []byte marshals to a base64 string.
func JSON(t Testing, actual any, expected string, messages ...string) bool {
	t.Helper()

	s, err := json.Marshal(actual)
	if err != nil {
		return fail(t, messages, "Should be marshalable", field{"actual", actual}, field{"err", errorText(err)})
	}

	return EqualJSON(t, string(s), expected, messages...)
}

// Panics asserts that fn panics. A panic(nil) counts whatever the GODEBUG panicnil setting. It
// fails for a nil fn.
func Panics(t Testing, fn func(), messages ...string) bool {
	t.Helper()

	if panicked, _, reason := panics(fn); reason != valid {
		return fail(t, messages, string(reason))
	} else if !panicked {
		return fail(t, messages, "Should panic")
	}

	return true
}

// PanicsWith asserts that fn panics with a value deeply equal to expected or, when expected is an
// error, matching it under errors.Is.
//
// The value of a panic(nil) is nil whatever the GODEBUG panicnil setting. It fails for a nil fn.
func PanicsWith(t Testing, fn func(), expected any, messages ...string) bool {
	t.Helper()

	if panicked, value, reason := panics(fn); reason != valid {
		return fail(t, messages, string(reason))
	} else if !panicked {
		return fail(t, messages, "Should panic", field{"expected", expected})
	} else if !panicsWith(value, expected) {
		return fail(t, messages, "Should panic with value", panicFields(value, comparisonFields(value, expected)...)...)
	}

	return true
}

// NotPanics asserts that fn does not panic.
func NotPanics(t Testing, fn func(), messages ...string) bool {
	t.Helper()

	if panicked, value, reason := panics(fn); reason != valid {
		return fail(t, messages, string(reason))
	} else if panicked {
		return fail(t, messages, "Should not panic", panicFields(value, field{"value", value})...)
	}

	return true
}
