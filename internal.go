package assert

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type validity string

const (
	valid            validity = ""
	notReference     validity = "Should be reference"
	notIterable      validity = "Should be iterable"
	notStringOrSlice validity = "Should be string or slice"
	notArrayOrSlice  validity = "Should be array or slice"
	elementType      validity = "Should have element of same type"
	typeMismatch     validity = "Should have same type"
	invalidDelta     validity = "Should have non-negative delta"
	invalidTarget    validity = "Should have pointer to error or interface target"
	notFunction      validity = "Should not be function"
)

const (
	formatLimit = 1024
	labelWidth  = 8
)

const jsonWhitespace = " \t\r\n"

type jsonNumber string

type field struct {
	label string
	value string
}

type entry struct {
	key   reflect.Value
	value reflect.Value
}

func fail(t Testing, messages []string, reason string, fields ...field) bool {
	t.Helper()

	var message strings.Builder
	message.WriteString(join(messages))
	message.WriteString(reason)

	for _, item := range fields {
		fmt.Fprintf(&message, "\n%*s: %s", labelWidth, item.label, truncate(item.value))
	}

	return Fail(t, message.String())
}

func errorFields(err error, fields ...field) []field {
	return append([]field{{"msg", fmt.Sprint(err)}, {"error", formatGoSyntax(err)}}, fields...)
}

func equal[T Comparable](actual, expected T) bool {
	return reflect.DeepEqual(actual, expected)
}

func deepEqual[T Comparable](actual, expected T) (bool, validity) {
	if isFunction(actual) || isFunction(expected) {
		return false, notFunction
	}

	return equal(actual, expected), valid
}

func isFunction(object any) bool {
	valueOf := reflect.ValueOf(object)

	return valueOf.Kind() == reflect.Func && !valueOf.IsNil()
}

func equalDelta[T Numeric](actual, expected, delta T) (bool, validity) {
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

func formatDistance[T Numeric](actual, expected T) string {
	if isFloat[T]() {
		return format(T(math.Abs(float64(expected) - float64(actual))))
	}

	distance := integerDistance(actual, expected)
	if typed := T(distance); typed >= 0 && uint64(typed) == distance {
		return format(typed)
	}

	return strconv.FormatUint(distance, 10)
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

func same[T Reference](actual, expected T) (bool, validity) {
	valueOfActual := reflect.ValueOf(actual)
	valueOfExpected := reflect.ValueOf(expected)

	if !isReference(valueOfActual) || !isReference(valueOfExpected) {
		return false, notReference
	}

	if !sameType(actual, expected) {
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

func sameType(actual, expected any) bool {
	return reflect.TypeOf(actual) == reflect.TypeOf(expected)
}

func isReference(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Chan:
		return true
	default:
		return false
	}
}

func length[S Iterable](object S) (int, validity) {
	valueOf := reflect.ValueOf(object)

	if !hasLength(valueOf) {
		return 0, notIterable
	}

	return valueOf.Len(), valid
}

func hasLength(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return true
	case reflect.Pointer:
		return value.Type().Elem().Kind() == reflect.Array
	default:
		return false
	}
}

func contains[S Iterable, E Comparable](object S, element E) (bool, validity) {
	valueOf := indirectArray(reflect.ValueOf(object))

	switch valueOf.Kind() {
	case reflect.String:
		return containsSubstring(valueOf, reflect.ValueOf(element))
	case reflect.Array, reflect.Slice, reflect.Map:
		return containsValue(valueOf, element)
	default:
		return false, notIterable
	}
}

func containsSubstring(object, element reflect.Value) (bool, validity) {
	if element.Kind() != reflect.String {
		return false, elementType
	}

	return strings.Contains(object.String(), element.String()), valid
}

func containsValue(object reflect.Value, element any) (bool, validity) {
	if !isAssignable(reflect.TypeOf(element), object.Type().Elem()) {
		return false, elementType
	}

	for _, item := range object.Seq2() {
		if equal(item.Interface(), element) {
			return true, valid
		}
	}

	return false, valid
}

func hasPrefix[S Iterable](object, prefix S) (bool, validity) {
	valueOfObject, valueOfPrefix, reason := affixValues(object, prefix)
	if reason != valid {
		return false, reason
	}

	return equalAt(valueOfObject, valueOfPrefix, 0), valid
}

func hasSuffix[S Iterable](object, suffix S) (bool, validity) {
	valueOfObject, valueOfSuffix, reason := affixValues(object, suffix)
	if reason != valid {
		return false, reason
	}

	return equalAt(valueOfObject, valueOfSuffix, valueOfObject.Len()-valueOfSuffix.Len()), valid
}

func affixValues[S Iterable](object, affix S) (valueOfObject, valueOfAffix reflect.Value, reason validity) {
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

	return equal(object.Slice(start, end).Interface(), part.Interface())
}

func unorderedDifference[S Iterable](actual, expected S) (extra, missing reflect.Value, reason validity) {
	valueOfActual := indirectArray(reflect.ValueOf(actual))
	valueOfExpected := indirectArray(reflect.ValueOf(expected))

	if !isArrayOrSlice(valueOfActual) || !isArrayOrSlice(valueOfExpected) {
		return extra, missing, notArrayOrSlice
	}

	if valueOfActual.Type() != valueOfExpected.Type() {
		return extra, missing, typeMismatch
	}

	extra, missing = difference(valueOfActual, valueOfExpected)

	return extra, missing, valid
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

func difference(actual, expected reflect.Value) (extra, missing reflect.Value) {
	sliceType := reflect.SliceOf(actual.Type().Elem())
	extra = reflect.MakeSlice(sliceType, 0, 0)
	missing = reflect.MakeSlice(sliceType, 0, 0)
	matched := make([]bool, expected.Len())

	for i := range actual.Len() {
		if j := indexUnmatched(expected, actual.Index(i), matched); j >= 0 {
			matched[j] = true
		} else {
			extra = reflect.Append(extra, actual.Index(i))
		}
	}

	for j, found := range matched {
		if !found {
			missing = reflect.Append(missing, expected.Index(j))
		}
	}

	return extra, missing
}

func indexUnmatched(list, element reflect.Value, matched []bool) int {
	for j := range list.Len() {
		if !matched[j] && equal(list.Index(j).Interface(), element.Interface()) {
			return j
		}
	}

	return -1
}

func isAssignable(from, to reflect.Type) bool {
	if from == nil {
		return to.Kind() == reflect.Interface
	}

	return from.AssignableTo(to)
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

	if rest := strings.TrimLeft(s[decoder.InputOffset():], jsonWhitespace); rest != "" {
		return nil, fmt.Errorf("unexpected %s after top-level value", format(rest))
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
	mantissa, exponent, _ := strings.Cut(strings.ToLower(number.String()), "e")
	unsigned, negative := strings.CutPrefix(mantissa, "-")
	integer, fraction, _ := strings.Cut(unsigned, ".")

	digits := strings.TrimLeft(integer+fraction, "0")
	if digits == "" {
		return "0"
	}

	significant := strings.TrimRight(digits, "0")

	scale, _ := new(big.Int).SetString(cmp.Or(exponent, "0"), 10)
	scale.Add(scale, big.NewInt(int64(len(digits)-len(significant)-len(fraction))))

	normalized := significant + "e" + scale.String()
	if negative {
		normalized = "-" + normalized
	}

	return jsonNumber(normalized)
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

func isZero(object any) bool {
	if isNil(object) {
		return true
	}

	valueOf := reflect.ValueOf(object)
	if valueOf.Kind() == reflect.Pointer {
		return false
	}

	if zeroer, ok := object.(interface{ IsZero() bool }); ok {
		return zeroer.IsZero()
	}

	return valueOf.IsZero()
}

func join(messages []string) string {
	return strings.Join(messages, "")
}

func format(object any) string {
	valueOf := reflect.ValueOf(object)

	switch valueOf.Kind() {
	case reflect.Pointer:
		if valueOf.IsNil() {
			return formatValue(object)
		}

		return fmt.Sprintf("[%p] %s", object, formatValue(valueOf.Elem().Interface()))
	case reflect.Slice, reflect.Map:
		return fmt.Sprintf("[%p] %s", object, formatValue(object))
	default:
		return formatValue(object)
	}
}

func formatValue(object any) string {
	valueOf := reflect.ValueOf(object)

	switch {
	case isNumber(valueOf):
		return fmt.Sprint(object)
	case isNil(object):
		return formatGoSyntax(object)
	case isBytes(valueOf):
		return fmt.Sprintf("%s(%q)", formatType(object), valueOf.Bytes())
	case isArrayOrSlice(valueOf):
		return formatComposite(valueOf.Type(), formatElements(valueOf))
	case valueOf.Kind() == reflect.Map:
		return formatComposite(valueOf.Type(), formatEntries(valueOf))
	default:
		return formatGoSyntax(object)
	}
}

func isNumber(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func isBytes(value reflect.Value) bool {
	return value.Kind() == reflect.Slice && value.Type().Elem().Kind() == reflect.Uint8
}

func formatElements(value reflect.Value) []string {
	items := make([]string, 0, value.Len())
	for _, item := range value.Seq2() {
		items = append(items, formatValue(item.Interface()))
	}

	return items
}

func formatEntries(value reflect.Value) []string {
	entries := make([]entry, 0, value.Len())
	for key, item := range value.Seq2() {
		entries = append(entries, entry{key, item})
	}

	slices.SortFunc(entries, func(a, b entry) int {
		return compareKeys(a.key, b.key)
	})

	items := make([]string, len(entries))
	for i, item := range entries {
		items[i] = formatValue(item.key.Interface()) + ":" + formatValue(item.value.Interface())
	}

	return items
}

func compareKeys(a, b reflect.Value) int {
	switch {
	case a.CanInt():
		return cmp.Compare(a.Int(), b.Int())
	case a.CanUint():
		return cmp.Compare(a.Uint(), b.Uint())
	case a.CanFloat():
		return cmp.Compare(a.Float(), b.Float())
	case a.Kind() == reflect.String:
		return strings.Compare(a.String(), b.String())
	default:
		return strings.Compare(formatValue(a.Interface()), formatValue(b.Interface()))
	}
}

func formatComposite(typ reflect.Type, items []string) string {
	return typ.String() + "{" + strings.Join(items, ", ") + "}"
}

func formatType(object any) string {
	if _, ok := object.([]byte); ok {
		return "[]byte"
	}

	return fmt.Sprintf("%T", object)
}

func formatGoSyntax(object any) string {
	return fmt.Sprintf("%#v", object)
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
