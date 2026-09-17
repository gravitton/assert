package assert

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"unicode/utf8"
)

type reason string

const (
	valid        reason = ""
	notReference reason = "Should be reference"
	notIterable  reason = "Should be iterable"
	elementType  reason = "Should have element of same type"
	invalidDelta reason = "Should have non-negative delta"
)

const formatLimit = 1024

type jsonNumber string

func equal[T Comparable](actual, expected T) bool {
	return reflect.DeepEqual(actual, expected)
}

func equalDelta[T Numeric](actual, expected, delta T) (bool, reason) {
	if delta < 0 || delta != delta {
		return false, invalidDelta
	}

	if isFloat[T]() {
		return equalDeltaFloat(float64(actual), float64(expected), float64(delta)), valid
	}

	return integerDistance(actual, expected) <= uint64(delta), valid
}

func equalDeltaFloat(actual, expected, delta float64) bool {
	if actual == expected {
		return true
	}

	if actual != actual || expected != expected {
		return actual != actual && expected != expected
	}

	diff := expected - actual

	return diff >= -delta && diff <= delta
}

func integerDistance[T Numeric](a, b T) uint64 {
	if a < b {
		a, b = b, a
	}

	return uint64(a) - uint64(b)
}

func isFloat[T Numeric]() bool {
	switch reflect.TypeFor[T]().Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func compare[T Ordered](actual, expected T) (int, bool) {
	if actual != actual || expected != expected {
		return 0, false
	}

	switch {
	case actual < expected:
		return -1, true
	case actual > expected:
		return 1, true
	default:
		return 0, true
	}
}

func same[T Reference](actual, expected T) (bool, reason) {
	valueOfActual := reflect.ValueOf(actual)
	valueOfExpected := reflect.ValueOf(expected)

	if !isReference(valueOfActual) || !isReference(valueOfExpected) {
		return false, notReference
	}

	if valueOfActual.Type() != valueOfExpected.Type() {
		return false, valid
	}

	if valueOfActual.Pointer() != valueOfExpected.Pointer() {
		return false, valid
	}

	if valueOfActual.Kind() == reflect.Slice {
		return valueOfActual.Len() == valueOfExpected.Len() && valueOfActual.Cap() == valueOfExpected.Cap(), valid
	}

	return true, valid
}

func isReference(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan:
		return true
	default:
		return false
	}
}

func length[S Iterable](object S) (int, reason) {
	valueOf := reflect.ValueOf(object)

	switch valueOf.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return valueOf.Len(), valid
	default:
		return 0, notIterable
	}
}

func contains[S Iterable, E Comparable](object S, element E) (bool, reason) {
	valueOf := reflect.ValueOf(object)

	switch valueOf.Kind() {
	case reflect.String:
		return containsSubstring(valueOf, reflect.ValueOf(element))
	case reflect.Array, reflect.Slice, reflect.Map:
		return containsValue(valueOf, element)
	default:
		return false, notIterable
	}
}

func containsSubstring(object, element reflect.Value) (bool, reason) {
	if element.Kind() != reflect.String {
		return false, elementType
	}

	return strings.Contains(object.String(), element.String()), valid
}

func containsValue(object reflect.Value, element any) (bool, reason) {
	if !assignable(reflect.TypeOf(element), object.Type().Elem()) {
		return false, elementType
	}

	for _, item := range object.Seq2() {
		if equal(item.Interface(), element) {
			return true, valid
		}
	}

	return false, valid
}

func assignable(from, to reflect.Type) bool {
	if from == nil {
		return to.Kind() == reflect.Interface
	}

	return from.AssignableTo(to)
}

func matches(actual, pattern string) (bool, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	return re.MatchString(actual), nil
}

func decodeJSON(s string) (any, error) {
	decoder := json.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()

	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}

	if rest := strings.TrimSpace(s[decoder.InputOffset():]); rest != "" {
		return nil, fmt.Errorf("invalid character %q after top-level value", rest[0])
	}

	return normalizeJSON(value), nil
}

func normalizeJSON(value any) any {
	switch v := value.(type) {
	case json.Number:
		return normalizeJSONNumber(v)
	case []any:
		for i, item := range v {
			v[i] = normalizeJSON(item)
		}

		return v
	case map[string]any:
		for key, item := range v {
			v[key] = normalizeJSON(item)
		}

		return v
	default:
		return value
	}
}

func normalizeJSONNumber(number json.Number) jsonNumber {
	if rat, ok := new(big.Rat).SetString(number.String()); ok {
		return jsonNumber(rat.RatString())
	}

	return jsonNumber(number)
}

func panics(fn func()) (panicked bool, value any) {
	defer func() {
		if panicked {
			value = normalizePanic(recover())
		}
	}()

	panicked = true
	fn()
	panicked = false

	return panicked, nil
}

func normalizePanic(value any) any {
	if _, ok := value.(*runtime.PanicNilError); ok {
		return nil
	}

	return value
}

func panicsWith(value, expected any) bool {
	if target, ok := expected.(error); ok {
		if err, ok := value.(error); ok && errors.Is(err, target) {
			return true
		}
	}

	return equal(value, expected)
}

func isNil(object any) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)

	switch value.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return value.IsNil()
	default:
		return false
	}
}

func join(messages []string) string {
	return strings.Join(messages, "")
}

func format(object any) string {
	valueOf := reflect.ValueOf(object)

	switch valueOf.Kind() {
	case reflect.Pointer:
		if valueOf.IsNil() {
			return truncate(fmt.Sprintf("%#v", object))
		}

		return truncate(fmt.Sprintf("[%p] %#v", object, valueOf.Elem().Interface()))
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return truncate(fmt.Sprintf("[%[1]p] %#[1]v", object))
	default:
		return truncate(fmt.Sprintf("%#v", object))
	}
}

func truncate(s string) string {
	if len(s) <= formatLimit {
		return s
	}

	end := formatLimit
	for end > 0 && !utf8.RuneStart(s[end]) {
		end--
	}

	return s[:end] + "…"
}
