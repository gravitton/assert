package assert

import (
	"fmt"
	"reflect"
	"strings"
)

const labelWidth = 8

type field struct {
	label string
	value any
}

type text string

func (t text) GoString() string {
	return string(t)
}

func fail(t Testing, messages []string, reason string, fields ...field) bool {
	t.Helper()

	var message strings.Builder
	message.WriteString(strings.Join(messages, ""))
	message.WriteString(reason)

	for _, item := range fields {
		fmt.Fprintf(&message, "\n%*s: %s", labelWidth, item.label, indent(format(item.value)))
	}

	return Fail(t, message.String())
}

func format(object any) string {
	value := reflect.ValueOf(object)

	switch {
	case value.CanInt() || value.CanUint() || value.CanFloat():
		return fmt.Sprint(object)
	case value.Kind() == reflect.Pointer && !value.IsNil():
		return "&" + format(value.Elem().Interface())
	default:
		return fmt.Sprintf("%#v", object)
	}
}

func indent(s string) string {
	return strings.ReplaceAll(s, "\n", "\n"+strings.Repeat(" ", labelWidth+2))
}

func errorFields(err error, fields ...field) []field {
	return append([]field{{"msg", errorText(err)}, {"error", err}}, fields...)
}

func errorText(err error) text {
	return text(fmt.Sprint(err))
}

func panicFields(value any, fields ...field) []field {
	if err, ok := value.(error); ok {
		return append([]field{{"msg", errorText(err)}}, fields...)
	}

	return fields
}

func comparisonFields(actual, expected any) []field {
	switch {
	case format(actual) != format(expected):
		return []field{{"actual", actual}, {"expected", expected}}
	case reflect.TypeOf(actual) != reflect.TypeOf(expected):
		return []field{{"actual", typed(actual)}, {"expected", typed(expected)}}
	default:
		return []field{{"actual", actual}, {"expected", expected}, {"hint", text("values print the same but are not deeply equal")}}
	}
}

func typed(object any) text {
	value, name := format(object), fmt.Sprintf("%T", object)
	if strings.HasPrefix(value, name) {
		return text(value)
	}

	return text(name + "(" + value + ")")
}

func typeName(object any) text {
	return text(fmt.Sprintf("%T", object))
}

func identity(object any) text {
	return text(fmt.Sprintf("[%p] %#v", object, object))
}
