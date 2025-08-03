package assert

import (
	"fmt"
	"testing"
)

func TestAssert(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	cases := []struct {
		actual bool
		result bool
	}{
		{true, true},
		{false, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Assert(%#v)", c.actual), func(t *testing.T) {
			if True(mockT, c.actual) != c.result {
				t.Errorf("True(%#v) should return return %#v", c.actual, c.result)
			}

			if False(mockT, c.actual) != !c.result {
				t.Errorf("False(%#v) should return return %#v", c.actual, !c.result)
			}
		})
	}

}

func TestEqual(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	type myType string

	cases := []struct {
		actual   any
		expected any
		result   bool
	}{
		{"Hello World", "Hello World", true},
		{123, 123, true},
		{123.5, 123.5, true},
		{123.5, 123, false},
		{[]byte("Hello World"), []byte("Hello World"), true},
		{nil, nil, true},
		{int32(123), int32(123), true},
		{uint64(123), uint64(123), true},
		{int32(123), int64(123), false},
		{10, uint(10), false},
		{&struct{}{}, &struct{}{}, true}, // pointer equality is based on equality of underlying value
		{myType("1"), myType("1"), true},
		{myType("1"), myType("2"), false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Equal(%#v,%#v)", c.actual, c.expected), func(t *testing.T) {
			if Equal(mockT, c.actual, c.expected) != c.result {
				t.Errorf("Equal(%#v,%#v) should return %#v", c.actual, c.expected, c.result)
			}

			if NotEqual(mockT, c.actual, c.expected) != !c.result {
				t.Errorf("NotEqual(%#v,%#v) should return %#v", c.actual, c.expected, !c.result)
			}
		})
	}
}

func TestSame(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	v1 := 1
	p1 := ptr(v1)
	p2 := p1

	cases := []struct {
		actual   any
		expected any
		result   bool
	}{
		{"Hello World", "Hello World", false},
		{123, 123, false},
		{nil, nil, false},
		{&struct{}{}, &struct{}{}, true}, // pointer equality is based on equality of underlying value
		{v1, v1, false},
		{v1, &v1, false},
		{&v1, &v1, true},
		{ptr(1), ptr(1), false},
		{p1, p2, true},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Equal(%#v,%#v)", c.actual, c.expected), func(t *testing.T) {
			if Same(mockT, c.actual, c.expected) != c.result {
				t.Errorf("Same(%#v,%#v) should return %#v", c.actual, c.expected, c.result)
			}

			if NotSame(mockT, c.actual, c.expected) != !c.result {
				t.Errorf("NotSame(%#v,%#v) should return %#v", c.actual, c.expected, !c.result)
			}
		})
	}
}

func ptr(i int) *int {
	return &i
}
