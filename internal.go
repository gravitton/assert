package assert

import (
	"fmt"
	"math"
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

func equal[T Comparable](actual, expected T) bool {
	return reflect.DeepEqual(actual, expected)
}

func equalDelta[T Numeric](actual, expected, delta T) bool {
	if delta < 0 {
		panic("delta must be positive")
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

	if math.IsNaN(actual) || math.IsNaN(expected) {
		return math.IsNaN(actual) && math.IsNaN(expected)
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
	var zero T

	switch reflect.TypeOf(zero).Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func orderable[T Numeric](actual, expected T) bool {
	return actual == actual && expected == expected
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

	return valueOfActual.Pointer() == valueOfExpected.Pointer(), valid
}

func isReference(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

func length[S Iterable[any]](object S) (int, reason) {
	valueOf := reflect.ValueOf(object)

	switch valueOf.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return valueOf.Len(), valid
	default:
		return 0, notIterable
	}
}

func contains[S Iterable[E], E Comparable](object S, element E) (bool, reason) {
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

func panics(fn func()) (panicked bool, value any) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			value = r
		}
	}()

	fn()

	return false, valid
}

func isNilError(err error) bool {
	if err == nil {
		return true
	}

	value := reflect.ValueOf(err)

	switch value.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface:
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
