package assert

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

func equal[T Comparable](actual, expected T) bool {
	return reflect.DeepEqual(actual, expected)
}

func equalDelta[T Numeric](actual, expected, delta T) bool {
	if delta < 0 {
		panic("delta must be positive")
	}

	if expected == actual {
		return true
	}

	actualFloat := float64(actual)
	expectedFloat := float64(expected)

	fmt.Println(actual, expected)
	if math.IsNaN(actualFloat) && math.IsNaN(expectedFloat) {
		return true
	} else if math.IsNaN(actualFloat) || math.IsNaN(expectedFloat) {
		return false
	}

	diff := expectedFloat - actualFloat
	d := float64(delta)
	fmt.Println(actual, expected, diff, d)
	if diff < -d || diff > d {
		return false
	}

	return true
}

func same[T Reference](actual, expected T) (valid bool, ok bool) {
	valueOfActual := reflect.ValueOf(actual)
	valueOfExpected := reflect.ValueOf(expected)

	switch valueOfActual.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return valueOfActual.Pointer() == valueOfExpected.Pointer(), true
	default:
		return false, false
	}
}

func length[S Iterable[any]](object S) int {
	return reflect.ValueOf(object).Len()
}

func contains[S Iterable[E], E Comparable](object S, element E) (found bool, ok bool) {
	valueOf := reflect.ValueOf(object)
	typeOf := reflect.TypeOf(object)
	if typeOf == nil {
		return false, false
	}

	kind := typeOf.Kind()
	if kind == reflect.String {
		elementValue := reflect.ValueOf(element)
		return strings.Contains(valueOf.String(), elementValue.String()), true
	}

	if kind == reflect.Map {
		for _, key := range valueOf.MapKeys() {
			value, valid := valueOf.MapIndex(key).Interface().(E)
			if !valid {
				return false, false
			}
			if equal(value, element) {
				return true, true
			}
		}
		return false, true
	}

	for i := 0; i < valueOf.Len(); i++ {
		value, valid := valueOf.Index(i).Interface().(E)
		if !valid {
			return false, false
		}

		if equal(value, element) {
			return true, true
		}
	}

	return false, true
}

func print(object any) string {
	valueOf := reflect.ValueOf(object)
	switch valueOf.Kind() {
	case reflect.Pointer:
		return fmt.Sprintf("[%p] %#v", object, valueOf.Elem().Interface())
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return fmt.Sprintf("[%[1]p] %#[1]v", object)
	default:
		return fmt.Sprintf("%#v", object)
	}
}
