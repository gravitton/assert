package assert

import (
	"cmp"
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

type validity string

const (
	valid            validity = ""
	notReference     validity = "Should be reference"
	notIterable      validity = "Should be iterable"
	notStringOrSlice validity = "Should be string or slice"
	notArrayOrSlice  validity = "Should be array or slice"
	typeMismatch     validity = "Should have same type"
	notNumber        validity = "Should not be NaN"
	notError         validity = "Should be error"
	invalidTarget    validity = "Should have pointer to error or interface target"
)

func equal(actual, expected any) (bool, validity) {
	return reflect.DeepEqual(actual, expected), valid
}

func equalDelta[T Numeric](actual, expected, delta T) (bool, validity) {
	return distance(actual, expected) <= delta, valid
}

func distance[T Numeric](actual, expected T) T {
	return max(actual, expected) - min(actual, expected)
}

func compare[T Ordered](actual, bound T) (int, validity) {
	if actual != actual || bound != bound {
		return 0, notNumber
	}

	return cmp.Compare(actual, bound), valid
}

func same(actual, expected any) (bool, validity) {
	valueOfActual := reflect.ValueOf(actual)
	valueOfExpected := reflect.ValueOf(expected)

	if !isReference(valueOfActual) || !isReference(valueOfExpected) {
		return false, notReference
	}

	if valueOfActual.Type() != valueOfExpected.Type() {
		return false, typeMismatch
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

func length(object any) (int, validity) {
	value := reflect.Indirect(reflect.ValueOf(object))

	switch value.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return value.Len(), valid
	default:
		return 0, notIterable
	}
}

func contains(object, element any) (bool, validity) {
	value := reflect.Indirect(reflect.ValueOf(object))

	switch value.Kind() {
	case reflect.String:
		return strings.Contains(value.String(), reflect.ValueOf(element).String()), valid
	case reflect.Array, reflect.Slice, reflect.Map:
		return slices.ContainsFunc(elements(value), func(item any) bool {
			return reflect.DeepEqual(item, element)
		}), valid
	default:
		return false, notIterable
	}
}

func hasPrefix(object, prefix any) (bool, validity) {
	valueOfObject, valueOfPrefix, reason := affixValues(object, prefix)
	if reason != valid {
		return false, reason
	}

	return equalAt(valueOfObject, valueOfPrefix, 0), valid
}

func hasSuffix(object, suffix any) (bool, validity) {
	valueOfObject, valueOfSuffix, reason := affixValues(object, suffix)
	if reason != valid {
		return false, reason
	}

	return equalAt(valueOfObject, valueOfSuffix, valueOfObject.Len()-valueOfSuffix.Len()), valid
}

func affixValues(object, affix any) (valueOfObject, valueOfAffix reflect.Value, reason validity) {
	valueOfObject = reflect.ValueOf(object)
	valueOfAffix = reflect.ValueOf(affix)

	if !isStringOrSlice(valueOfObject) || !isStringOrSlice(valueOfAffix) {
		return valueOfObject, valueOfAffix, notStringOrSlice
	}

	return valueOfObject, valueOfAffix, valid
}

func isStringOrSlice(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Slice:
		return true
	default:
		return false
	}
}

func equalAt(object, part reflect.Value, start int) bool {
	end := start + part.Len()
	if start < 0 || end > object.Len() {
		return false
	}

	if part.Len() == 0 {
		return true
	}

	return reflect.DeepEqual(object.Slice(start, end).Interface(), part.Interface())
}

func unorderedDifference(actual, expected any) (extra, missing []any, reason validity) {
	valueOfActual := reflect.Indirect(reflect.ValueOf(actual))
	valueOfExpected := reflect.Indirect(reflect.ValueOf(expected))

	if !isArrayOrSlice(valueOfActual) || !isArrayOrSlice(valueOfExpected) {
		return nil, nil, notArrayOrSlice
	}

	missing = elements(valueOfExpected)

	for _, item := range elements(valueOfActual) {
		index := slices.IndexFunc(missing, func(candidate any) bool {
			return reflect.DeepEqual(candidate, item)
		})

		if index >= 0 {
			missing = slices.Delete(missing, index, index+1)
		} else {
			extra = append(extra, item)
		}
	}

	return extra, missing, valid
}

func isArrayOrSlice(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func elements(value reflect.Value) []any {
	items := make([]any, 0, value.Len())
	for _, item := range value.Seq2() {
		items = append(items, item.Interface())
	}

	return items
}

func errorIs(err, target error) bool {
	return errors.Is(err, target)
}

func errorContains(err error, substr string) (bool, validity) {
	if err == nil {
		return false, notError
	}

	return strings.Contains(err.Error(), substr), valid
}

func errorAs(err error, target any) (bool, validity) {
	if !isErrorTarget(target) {
		return false, invalidTarget
	}

	return errors.As(err, target), valid
}

func isErrorTarget(target any) bool {
	valueOf := reflect.ValueOf(target)
	if valueOf.Kind() != reflect.Pointer || valueOf.IsNil() {
		return false
	}

	elem := valueOf.Type().Elem()

	return elem.Kind() == reflect.Interface || elem.Implements(reflect.TypeFor[error]())
}

func matches(actual, pattern string) (bool, error) {
	return regexp.MatchString(pattern, actual)
}

func decodeJSON(s string) (any, error) {
	var value any
	err := json.Unmarshal([]byte(s), &value)

	return value, err
}

func panics(fn func()) (panicked bool, value any, reason validity) {
	value = recovered(fn)

	return value != nil, value, valid
}

func recovered(fn func()) (value any) {
	defer func() {
		value = recover()
	}()

	fn()

	return nil
}

func panicsWith(value, expected any) bool {
	if target, ok := expected.(error); ok {
		if err, ok := value.(error); ok && errors.Is(err, target) {
			return true
		}
	}

	return reflect.DeepEqual(value, expected)
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

func isZero(object any) bool {
	return object == nil || reflect.ValueOf(object).IsZero()
}
