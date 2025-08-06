package assert

import (
	"reflect"
	"strings"
)

func equal(actual, expected any) bool {
	return reflect.DeepEqual(actual, expected)
}

func same(actual, expected any) bool {
	if reflect.ValueOf(actual).Kind() != reflect.Ptr || reflect.ValueOf(expected).Kind() != reflect.Ptr {
		return false
	}

	return actual == expected
}

func length(object any) int {
	return reflect.ValueOf(object).Len()
}

func contains(object any, element any) (found bool, ok bool) {
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
			if equal(valueOf.MapIndex(key).Interface(), element) {
				return true, true
			}
		}
		return false, true
	}

	for i := 0; i < valueOf.Len(); i++ {
		if equal(valueOf.Index(i).Interface(), element) {
			return true, true
		}
	}

	return false, true
}
