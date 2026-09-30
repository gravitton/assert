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
		fmt.Fprintf(&message, "\n%*s: %s", labelWidth, item.label, format(item.value))
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

func distinct(actual, expected any) (any, any) {
	if format(actual) != format(expected) || reflect.TypeOf(actual) == reflect.TypeOf(expected) {
		return actual, expected
	}

	return typed(actual), typed(expected)
}

func typed(object any) text {
	return text(fmt.Sprintf("%T(%s)", object, format(object)))
}

func typeOf(object any) text {
	return text(fmt.Sprintf("%T", object))
}

func identity(object any) text {
	return text(fmt.Sprintf("[%p] %#v", object, object))
}
