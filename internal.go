package assert

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
)

type reason string

const (
	valid        reason = ""
	notReference reason = "Should be reference"
	notIterable  reason = "Should be iterable"
	elementType  reason = "Should have element of same type"
)

type jsonNumber string

func equal[T Comparable](actual, expected T) bool {
	return reflect.DeepEqual(actual, expected)
}

func equalDelta[T Numeric](actual, expected, delta T) bool {
	if delta < 0 || delta != delta {
		panic("delta must be non-negative")
	}

	if isFloat[T]() {
		return equalDeltaFloat(float64(actual), float64(expected), float64(delta))
	}

	return integerDistance(actual, expected) <= uint64(delta)
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
	return T(1)/T(2) != 0
}

func compare[T Numeric](actual, expected T) (int, bool) {
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
		value = recover()
	}()

	panicked = true
	fn()

	return false, nil
}

func isNil(object any) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)

	switch value.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return value.IsNil()
	default:
		return false
	}
}

func message(message []string) string {
	return strings.Join(message, "")
}

func print(object any) string {
	valueOf := reflect.ValueOf(object)

	switch valueOf.Kind() {
	case reflect.Pointer:
		if valueOf.IsNil() {
			return fmt.Sprintf("%#v", object)
		}

		return fmt.Sprintf("[%p] %#v", object, valueOf.Elem().Interface())
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return fmt.Sprintf("[%[1]p] %#[1]v", object)
	default:
		return fmt.Sprintf("%#v", object)
	}
}
