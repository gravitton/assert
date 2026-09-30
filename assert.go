package assert

import "encoding/json"

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

// Nil asserts that actual is nil.
//
// A typed nil pointer, slice, map, channel, function or unsafe.Pointer is treated as nil.
func Nil(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if !isNil(actual) {
		return fail(t, messages, "Should be nil", field{"actual", actual})
	}

	return true
}

// NotNil asserts that actual is NOT nil.
//
// A typed nil pointer, slice, map, channel, function or unsafe.Pointer is treated as nil.
func NotNil(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if isNil(actual) {
		return fail(t, messages, "Should not be nil", field{"actual", actual})
	}

	return true
}

// Zero asserts that actual is the zero value of its type.
//
// A type's own IsZero() bool method is used when it has one, so a zero time.Time
// in any location is zero. That includes a method promoted from an embedded field,
// so a struct embedding time.Time is zero whenever the embedded time is, whatever
// its other fields hold. A pointer is zero only when it is nil, even when it
// points to a zero value.
// Unlike [Empty], a non-nil empty slice or map is not zero.
func Zero(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if !isZero(actual) {
		return fail(t, messages, "Should be zero", field{"actual", actual})
	}

	return true
}

// NotZero asserts that actual is NOT the zero value of its type.
//
// A type's own IsZero() bool method is used when it has one, so a zero time.Time
// in any location is zero. That includes a method promoted from an embedded field,
// so a struct embedding time.Time is zero whenever the embedded time is, whatever
// its other fields hold. A pointer is zero only when it is nil, even when it
// points to a zero value.
// Unlike [Empty], a non-nil empty slice or map is not zero.
func NotZero(t Testing, actual any, messages ...string) bool {
	t.Helper()

	if isZero(actual) {
		return fail(t, messages, "Should not be zero", field{"actual", actual})
	}

	return true
}

// Same asserts that two references point to the same object.
//
// Both arguments must be references: pointers, slices, maps or channels.
// Two references are the same when they have the same type and address;
// slices must also have the same length and capacity.
//
// Two nil references of the same type are the same. Pointers to zero-size values
// and slices with zero capacity may share an address even when allocated
// separately, and are then reported as same.
func Same[T Reference](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if identical, reason := same(actual, expected); reason != valid {
		return fail(t, messages, string(reason), field{"actual", typeOf(actual)}, field{"expected", typeOf(expected)})
	} else if !identical {
		return fail(t, messages, "Should be same", field{"actual", identity(actual)}, field{"expected", identity(expected)})
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

	if identical, reason := same(actual, expected); reason == notReference {
		return fail(t, messages, string(reason), field{"actual", typeOf(actual)}, field{"expected", typeOf(expected)})
	} else if identical {
		return fail(t, messages, "Should not be same", field{"actual", identity(actual)})
	}

	return true
}

// Equal asserts that two objects are equal.
//
// Equality is determined with reflect.DeepEqual: pointers are compared by the
// values they reference. A typed nil stored in an interface is not equal to an
// untyped nil; use [Nil] for that. Non-nil functions are never equal and fail
// as an invalid argument; use [Nil] or [NotNil] for them.
func Equal[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if equals, reason := equal(actual, expected); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"expected", expected})
	} else if !equals {
		distinctActual, distinctExpected := distinct(actual, expected)

		return fail(t, messages, "Should be equal", field{"actual", distinctActual}, field{"expected", distinctExpected})
	}

	return true
}

// NotEqual asserts that the specified values are NOT equal.
//
// Equality is determined with reflect.DeepEqual: pointers are compared by the
// values they reference. Non-nil functions are never equal and fail as an
// invalid argument; use [Nil] or [NotNil] for them.
func NotEqual[T Comparable](t Testing, actual, expected T, messages ...string) bool {
	t.Helper()

	if equals, reason := equal(actual, expected); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"expected", expected})
	} else if equals {
		return fail(t, messages, "Should not be equal", field{"actual", actual})
	}

	return true
}

// EqualDelta asserts that two numeric values differ by at most delta.
//
// Fails when delta is negative or NaN. NaN is only equal to NaN.
func EqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if within, reason := equalDelta(actual, expected, delta); reason != valid {
		return fail(t, messages, string(reason), field{"delta", delta})
	} else if !within {
		return fail(t, messages, "Should be equal in delta", field{"actual", actual}, field{"expected", expected}, field{"delta", delta}, field{"diff", distance(actual, expected)})
	}

	return true
}

// NotEqualDelta asserts that two numeric values differ by more than delta.
//
// Fails when delta is negative or NaN. NaN is only equal to NaN.
func NotEqualDelta[T Numeric](t Testing, actual, expected, delta T, messages ...string) bool {
	t.Helper()

	if within, reason := equalDelta(actual, expected, delta); reason != valid {
		return fail(t, messages, string(reason), field{"delta", delta})
	} else if within {
		return fail(t, messages, "Should not be equal in delta", field{"actual", actual}, field{"expected", expected}, field{"delta", delta}, field{"diff", distance(actual, expected)})
	}

	return true
}

// Greater asserts that actual is greater than bound.
//
// Fails when either value is NaN.
func Greater[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, reason := compare(actual, bound); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"bound", bound})
	} else if order <= 0 {
		return fail(t, messages, "Should be greater", field{"actual", actual}, field{"bound", bound})
	}

	return true
}

// GreaterOrEqual asserts that actual is greater than or equal to bound.
//
// Fails when either value is NaN.
func GreaterOrEqual[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, reason := compare(actual, bound); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"bound", bound})
	} else if order < 0 {
		return fail(t, messages, "Should be greater or equal", field{"actual", actual}, field{"bound", bound})
	}

	return true
}

// Less asserts that actual is less than bound.
//
// Fails when either value is NaN.
func Less[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, reason := compare(actual, bound); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"bound", bound})
	} else if order >= 0 {
		return fail(t, messages, "Should be less", field{"actual", actual}, field{"bound", bound})
	}

	return true
}

// LessOrEqual asserts that actual is less than or equal to bound.
//
// Fails when either value is NaN.
func LessOrEqual[T Ordered](t Testing, actual, bound T, messages ...string) bool {
	t.Helper()

	if order, reason := compare(actual, bound); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"bound", bound})
	} else if order > 0 {
		return fail(t, messages, "Should be less or equal", field{"actual", actual}, field{"bound", bound})
	}

	return true
}

// Length asserts that actual has given length.
//
// Works with strings, arrays, pointers to arrays, slices, maps and channels. String length is measured in bytes, not runes.
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
//
// Works with strings, arrays, pointers to arrays, slices, maps and channels. String length is measured in bytes, not runes.
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
//
// Works with strings, arrays, pointers to arrays, slices, maps and channels. String length is measured in bytes, not runes.
func NotEmpty[S Iterable](t Testing, actual S, messages ...string) bool {
	t.Helper()

	if size, reason := length(actual); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual})
	} else if size == 0 {
		return fail(t, messages, "Should not be empty", field{"actual", actual})
	}

	return true
}

// Contains asserts that actual contains given element.
//
// Works with strings, arrays, array pointers, slices and map values.
func Contains[S Iterable, E Comparable](t Testing, actual S, element E, messages ...string) bool {
	t.Helper()

	if found, reason := contains(actual, element); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"element", element})
	} else if !found {
		return fail(t, messages, "Should contain element", field{"actual", actual}, field{"element", element})
	}

	return true
}

// NotContains asserts that actual does NOT contain given element.
//
// Works with strings, arrays, array pointers, slices and map values.
func NotContains[S Iterable, E Comparable](t Testing, actual S, element E, messages ...string) bool {
	t.Helper()

	if found, reason := contains(actual, element); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"element", element})
	} else if found {
		return fail(t, messages, "Should not contain element", field{"actual", actual}, field{"element", element})
	}

	return true
}

// EqualUnordered asserts that two arrays, array pointers or slices have the same elements in any order.
//
// Elements are compared with reflect.DeepEqual and each must occur the same number
// of times in both. A nil slice and an empty slice are equal.
func EqualUnordered[S Iterable](t Testing, actual, expected S, messages ...string) bool {
	t.Helper()

	if extra, missing, reason := unorderedDifference(actual, expected); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"expected", expected})
	} else if extra.Len() > 0 || missing.Len() > 0 {
		return fail(t, messages, "Should be equal in any order", field{"actual", actual}, field{"expected", expected}, field{"extra", extra.Interface()}, field{"missing", missing.Interface()})
	}

	return true
}

// HasPrefix asserts that actual begins with prefix.
//
// Works with strings and slices. Slice elements are compared with reflect.DeepEqual.
// Arrays are not accepted: both arguments share one type, so they would have the same length.
func HasPrefix[S Iterable](t Testing, actual, prefix S, messages ...string) bool {
	t.Helper()

	if found, reason := hasPrefix(actual, prefix); reason == typeMismatch {
		return fail(t, messages, string(reason), field{"actual", typeOf(actual)}, field{"prefix", typeOf(prefix)})
	} else if reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"prefix", prefix})
	} else if !found {
		return fail(t, messages, "Should have prefix", field{"actual", actual}, field{"prefix", prefix})
	}

	return true
}

// HasSuffix asserts that actual ends with suffix.
//
// Works with strings and slices. Slice elements are compared with reflect.DeepEqual.
// Arrays are not accepted: both arguments share one type, so they would have the same length.
func HasSuffix[S Iterable](t Testing, actual, suffix S, messages ...string) bool {
	t.Helper()

	if found, reason := hasSuffix(actual, suffix); reason == typeMismatch {
		return fail(t, messages, string(reason), field{"actual", typeOf(actual)}, field{"suffix", typeOf(suffix)})
	} else if reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"suffix", suffix})
	} else if !found {
		return fail(t, messages, "Should have suffix", field{"actual", actual}, field{"suffix", suffix})
	}

	return true
}

// Error asserts that error is NOT nil.
//
// A typed nil stored in err (e.g. (*MyError)(nil)) is an error, as err != nil sees it.
func Error(t Testing, err error, messages ...string) bool {
	t.Helper()

	if err == nil {
		return fail(t, messages, "Should be error")
	}

	return true
}

// NoError asserts that error is nil.
//
// A typed nil stored in err (e.g. (*MyError)(nil)) is an error, as err != nil sees it.
func NoError(t Testing, err error, messages ...string) bool {
	t.Helper()

	if err != nil {
		return fail(t, messages, "Should not be error", field{"error", err})
	}

	return true
}

// ErrorIs asserts that error matches given target according to errors.Is.
func ErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if !errorIs(err, target) {
		return fail(t, messages, "Should match error", field{"error", err}, field{"target", target})
	}

	return true
}

// NotErrorIs asserts that error does NOT match given target according to errors.Is.
func NotErrorIs(t Testing, err error, target error, messages ...string) bool {
	t.Helper()

	if errorIs(err, target) {
		return fail(t, messages, "Should not match error", field{"error", err}, field{"target", target})
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
		return fail(t, messages, string(reason), field{"target", target})
	} else if !assignable {
		return fail(t, messages, "Should be assignable to target", field{"error", err}, field{"target", target})
	}

	return true
}

// ErrorContains asserts that error is NOT nil and its message contains substr.
//
// A typed nil stored in err (e.g. (*MyError)(nil)) is an error, as err != nil sees it.
func ErrorContains(t Testing, err error, substr string, messages ...string) bool {
	t.Helper()

	if found, reason := errorContains(err, substr); reason != valid {
		return fail(t, messages, string(reason), field{"substr", substr})
	} else if !found {
		return fail(t, messages, "Should contain substring", field{"error", err}, field{"substr", substr})
	}

	return true
}

// Matches asserts that a string matches the given regular expression.
func Matches(t Testing, actual, pattern string, messages ...string) bool {
	t.Helper()

	if matched, err := matches(actual, pattern); err != nil {
		return fail(t, messages, "Should be valid regexp", field{"pattern", pattern}, field{"err", err})
	} else if !matched {
		return fail(t, messages, "Should match regexp", field{"actual", actual}, field{"pattern", pattern})
	}

	return true
}

// NotMatches asserts that a string does NOT match the given regular expression.
func NotMatches(t Testing, actual, pattern string, messages ...string) bool {
	t.Helper()

	if matched, err := matches(actual, pattern); err != nil {
		return fail(t, messages, "Should be valid regexp", field{"pattern", pattern}, field{"err", err})
	} else if matched {
		return fail(t, messages, "Should not match regexp", field{"actual", actual}, field{"pattern", pattern})
	}

	return true
}

// EqualJSON asserts that JSON strings are semantically equal.
//
// Both strings are decoded as encoding/json decodes into any: numbers become
// float64, so 1.0 equals 1 and 1e2 equals 100, but integers beyond 2^53 may
// compare equal when they differ, and numbers beyond the float64 range are
// reported as invalid JSON. To compare such numbers exactly, decode into a
// struct and use [Equal]. The last of duplicate keys wins.
func EqualJSON(t Testing, actual, expected string, messages ...string) bool {
	t.Helper()

	actualJSON, err := decodeJSON(actual)
	if err != nil {
		return fail(t, messages, "Should be valid JSON", field{"actual", actual}, field{"err", err})
	}

	expectedJSON, err := decodeJSON(expected)
	if err != nil {
		return fail(t, messages, "Should be valid JSON", field{"expected", expected}, field{"err", err})
	}

	if equals, reason := equal(actualJSON, expectedJSON); reason != valid {
		return fail(t, messages, string(reason), field{"actual", actual}, field{"expected", expected})
	} else if !equals {
		return fail(t, messages, "Should be equal JSON", field{"actual", actual}, field{"expected", expected})
	}

	return true
}

// JSON asserts that actual can be marshaled to expected JSON string.
//
// actual is marshaled with json.Marshal, so a []byte becomes a base64 string.
// To compare JSON held in a []byte, use EqualJSON(t, string(actual), expected).
func JSON(t Testing, actual any, expected string, messages ...string) bool {
	t.Helper()

	s, err := json.Marshal(actual)
	if err != nil {
		return fail(t, messages, "Should be marshalable", field{"actual", actual}, field{"err", err})
	}

	return EqualJSON(t, string(s), expected, messages...)
}

// Panics asserts that fn panics.
//
// A panic(nil) is recognised regardless of the GODEBUG panicnil setting.
// Fails when fn is nil.
func Panics(t Testing, fn func(), messages ...string) bool {
	t.Helper()

	if panicked, _, reason := panics(fn); reason != valid {
		return fail(t, messages, string(reason))
	} else if !panicked {
		return fail(t, messages, "Should panic")
	}

	return true
}

// PanicsWith asserts that fn panics with the expected value.
//
// The panic value must be deeply equal to expected. When expected is an error,
// a panic value matching it according to errors.Is is accepted as well.
// A panic(nil) is reported as a nil value regardless of the GODEBUG panicnil setting.
// Fails when fn is nil.
func PanicsWith(t Testing, fn func(), expected any, messages ...string) bool {
	t.Helper()

	if panicked, value, reason := panics(fn); reason != valid {
		return fail(t, messages, string(reason))
	} else if !panicked {
		return fail(t, messages, "Should panic", field{"expected", expected})
	} else if !panicsWith(value, expected) {
		distinctActual, distinctExpected := distinct(value, expected)

		return fail(t, messages, "Should panic with value", field{"actual", distinctActual}, field{"expected", distinctExpected})
	}

	return true
}

// NotPanics asserts that fn does NOT panic.
//
// A panic(nil) is recognised regardless of the GODEBUG panicnil setting.
// Fails when fn is nil.
func NotPanics(t Testing, fn func(), messages ...string) bool {
	t.Helper()

	if panicked, value, reason := panics(fn); reason != valid {
		return fail(t, messages, string(reason))
	} else if panicked {
		return fail(t, messages, "Should not panic", field{"value", value})
	}

	return true
}
