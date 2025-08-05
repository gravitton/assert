package assert

import (
	"errors"
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
		t.Run(fmt.Sprintf("True(%#v)", c.actual), func(t *testing.T) {
			if True(mockT, c.actual) != c.result {
				t.Errorf("True(%#v) should return return %#v", c.actual, c.result)
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
		{[]int{1, 2}, []int{1, 2, 3}, false},
		{nil, nil, true},
		{int32(123), int32(123), true},
		{uint64(123), uint64(123), true},
		{int32(123), int64(123), false},
		{10, uint(10), false},
		{&struct{}{}, &struct{}{}, true}, // pointer equality is based on equality of underlying value
		{myType("1"), myType("1"), true},
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
		{[]byte("Hello World"), []byte("Hello World"), false},
		{[]int(nil), nil, false},
		{&struct{}{}, &struct{}{}, true}, // pointer equality is based on equality of underlying value
		{v1, v1, false},
		{v1, &v1, false},
		{&v1, &v1, true},
		{ptr(1), ptr(1), false},
		{p1, p2, true},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Same(%#v,%#v)", c.actual, c.expected), func(t *testing.T) {
			if Same(mockT, c.actual, c.expected) != c.result {
				t.Errorf("Same(%#v,%#v) should return %#v", c.actual, c.expected, c.result)
			}

			if NotSame(mockT, c.actual, c.expected) != !c.result {
				t.Errorf("NotSame(%#v,%#v) should return %#v", c.actual, c.expected, !c.result)
			}
		})
	}
}

func TestLength(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	cases := []struct {
		actual   any
		expected int
		result   bool
	}{
		{[]int{}, 0, true},
		{[]int{1, 2, 3}, 3, true},
		{[]int{1, 2, 3}, 2, false},
		{"Hello", 5, true},
		{map[string]bool{"a": true, "b": false}, 2, true},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Length(%#v,%#v)", c.actual, c.expected), func(t *testing.T) {
			if Length(mockT, c.actual, c.expected) != c.result {
				t.Errorf("Length(%#v,%#v) should return %#v", c.actual, c.expected, c.result)
			}
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	cases := []struct {
		object  any
		element any
		result  bool
	}{
		{[]int{}, 0, false},
		{[]int{1, 2, 3}, 2, true},
		{[]int{1, 2, 3}, 4, false},
		{"Hello", "e", true},
		{"Hello", 2, false},
		{map[string]bool{"a": true, "b": false}, true, true},
		{map[string]bool{"a": true, "b": false}, "a", false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Contains(%#v,%#v)", c.object, c.element), func(t *testing.T) {
			if Contains(mockT, c.object, c.element) != c.result {
				t.Errorf("Contains(%#v,%#v) should return %#v", c.object, c.element, c.result)
			}

			if NotContains(mockT, c.object, c.element) != !c.result {
				t.Errorf("NotContains(%#v,%#v) should return %#v", c.object, c.element, c.result)
			}
		})
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	cases := []struct {
		err    error
		result bool
	}{
		{nil, false},
		{errors.New("ooh"), true},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Error(%#v)", c.err), func(t *testing.T) {
			if Error(mockT, c.err) != c.result {
				t.Errorf("Error(%#v) should return %#v", c.err, c.result)
			}

			if NoError(mockT, c.err) != !c.result {
				t.Errorf("NoError(%#v) should return %#v", c.err, !c.result)
			}
		})
	}
}

func TestErrorIs(t *testing.T) {
	t.Parallel()

	mockT := new(testing.T)

	err := errors.New("ooh")

	cases := []struct {
		err    error
		target error
		result bool
	}{
		{nil, nil, true},
		{errors.New("ooh"), nil, false},
		{err, err, true},
		{errors.Join(errors.New("ooh1"), err), err, true},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("ErrorIs(%#v,%#v)", c.err, c.target), func(t *testing.T) {
			if ErrorIs(mockT, c.err, c.target) != c.result {
				t.Errorf("ErrorIs(%#v,%#v) should return %#v", c.err, c.target, c.result)
			}

			if NotErrorIs(mockT, c.err, c.target) != !c.result {
				t.Errorf("NotErrorIs(%#v,%#v) should return %#v", c.err, c.target, !c.result)
			}
		})
	}
}

func ptr(i int) *int {
	return &i
}
