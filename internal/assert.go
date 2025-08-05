package internal

import (
	"errors"
	"reflect"
	"strings"
)

func Assert(condition bool) bool {
	return condition
}

func Equal(actual, expected any) bool {
	return reflect.DeepEqual(actual, expected)
}

func Same(actual, expected any) bool {
	if reflect.ValueOf(actual).Kind() != reflect.Ptr || reflect.ValueOf(expected).Kind() != reflect.Ptr {
		return false
	}

	return actual == expected
}

func Length(object any) int {
	return reflect.ValueOf(object).Len()
}

func Contains(object any, element any) (found bool, ok bool) {
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
			if Equal(valueOf.MapIndex(key).Interface(), element) {
				return true, true
			}
		}
		return false, true
	}

	for i := 0; i < valueOf.Len(); i++ {
		if Equal(valueOf.Index(i).Interface(), element) {
			return true, true
		}
	}

	return false, true
}

func Error(err error) bool {
	return err != nil
}

func ErrorIs(err, target error) bool {
	return errors.Is(err, target)
}
