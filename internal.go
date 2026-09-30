package assert

import (
	"cmp"
	"encoding/json"
	"errors"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strconv"
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
	invalidDelta     validity = "Should have non-negative delta"
	notFunction      validity = "Should not be function"
	elementType      validity = "Should have element of same type"
	notError         validity = "Should be error"
	invalidTarget    validity = "Should have pointer to error or interface target"
	nilFunction      validity = "Should be non-nil function"
)

func equal(actual, expected any) (bool, validity) {
	if isFunction(actual) || isFunction(expected) {
		return false, notFunction
	}

	return reflect.DeepEqual(actual, expected), valid
}

func isFunction(object any) bool {
	value := reflect.ValueOf(object)

	return value.Kind() == reflect.Func && !value.IsNil()
}

func equalDelta[T Numeric](actual, expected, delta T) (bool, validity) {
	switch {
	case delta < 0 || delta != delta:
		return false, invalidDelta
	case actual == expected || actual != actual && expected != expected:
		return true, valid
	case isFloat[T]():
		return floatDistance(actual, expected) <= delta, valid
	default:
		return integerDistance(actual, expected) <= uint64(delta), valid
	}
}

func distance[T Numeric](actual, expected T) any {
	if isFloat[T]() {
		return floatDistance(actual, expected)
	}

	return text(strconv.FormatUint(integerDistance(actual, expected), 10))
}

func floatDistance[T Numeric](actual, expected T) T {
	return max(actual, expected) - min(actual, expected)
}

func integerDistance[T Numeric](actual, expected T) uint64 {
	return uint64(max(actual, expected)) - uint64(min(actual, expected))
}

func isFloat[T Numeric]() bool {
	switch reflect.TypeFor[T]().Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
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
	value := reflect.ValueOf(object)

	switch value.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return value.Len(), valid
	case reflect.Pointer:
		if value.Type().Elem().Kind() == reflect.Array {
			return value.Len(), valid
		}
	}

	return 0, notIterable
}

func contains(object, element any) (bool, validity) {
	value := indirectArray(reflect.ValueOf(object))

	switch value.Kind() {
	case reflect.String:
		return containsSubstring(value, reflect.ValueOf(element))
	case reflect.Array, reflect.Slice, reflect.Map:
		return containsElement(value, element)
	default:
		return false, notIterable
	}
}

func containsSubstring(value, element reflect.Value) (bool, validity) {
	if element.Kind() != reflect.String {
		return false, elementType
	}

	return strings.Contains(value.String(), element.String()), valid
}

func containsElement(value reflect.Value, element any) (bool, validity) {
	if !isAssignable(reflect.TypeOf(element), value.Type().Elem()) {
		return false, elementType
	}

	return slices.ContainsFunc(elements(value), func(item reflect.Value) bool {
		return reflect.DeepEqual(item.Interface(), element)
	}), valid
}

func isAssignable(from, to reflect.Type) bool {
	if from == nil {
		return to.Kind() == reflect.Interface
	}

	return from.AssignableTo(to)
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

	if valueOfObject.Type() != valueOfAffix.Type() {
		return valueOfObject, valueOfAffix, typeMismatch
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

func unorderedDifference(actual, expected any) (extra, missing reflect.Value, reason validity) {
	valueOfActual := indirectArray(reflect.ValueOf(actual))
	valueOfExpected := indirectArray(reflect.ValueOf(expected))

	if !isArrayOrSlice(valueOfActual) || !isArrayOrSlice(valueOfExpected) {
		return extra, missing, notArrayOrSlice
	}

	if valueOfActual.Type() != valueOfExpected.Type() {
		return extra, missing, typeMismatch
	}

	var extraItems []reflect.Value
	missingItems := elements(valueOfExpected)

	for _, item := range elements(valueOfActual) {
		index := slices.IndexFunc(missingItems, func(candidate reflect.Value) bool {
			return reflect.DeepEqual(candidate.Interface(), item.Interface())
		})

		if index >= 0 {
			missingItems = slices.Delete(missingItems, index, index+1)
		} else {
			extraItems = append(extraItems, item)
		}
	}

	itemType := valueOfActual.Type().Elem()

	return sliceOf(itemType, extraItems), sliceOf(itemType, missingItems), valid
}

func indirectArray(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Pointer && value.Type().Elem().Kind() == reflect.Array {
		return value.Elem()
	}

	return value
}

func isArrayOrSlice(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		return true
	default:
		return false
	}
}

func elements(value reflect.Value) []reflect.Value {
	items := make([]reflect.Value, 0, value.Len())
	for _, item := range value.Seq2() {
		items = append(items, item)
	}

	return items
}

func sliceOf(itemType reflect.Type, items []reflect.Value) reflect.Value {
	return reflect.Append(reflect.MakeSlice(reflect.SliceOf(itemType), 0, len(items)), items...)
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
	if fn == nil {
		return false, nil, nilFunction
	}

	panicked, value = recovered(fn)
	if _, ok := value.(*runtime.PanicNilError); ok {
		value = nil
	}

	return panicked, value, valid
}

func recovered(fn func()) (panicked bool, value any) {
	defer func() {
		if panicked {
			value = recover()
		}
	}()

	panicked = true
	fn()

	return false, nil
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
	value := reflect.ValueOf(object)

	if zeroer, ok := object.(interface{ IsZero() bool }); ok && value.Kind() != reflect.Pointer {
		return zeroer.IsZero()
	}

	return object == nil || value.IsZero()
}
