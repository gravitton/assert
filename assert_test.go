package assert

import (
	"errors"
	"fmt"
	"math"
	"testing"
	"time"
)

func TestAssert(t *testing.T) {
	testAssert(t, true, true)
	testAssert(t, false, false)
}

func TestEqual(t *testing.T) {
	testEqual(t, "Hello World", "Hello World", true)
	testEqual(t, "Hello World", "Hello World!", false)
	testEqual[testType](t, "A", "A", true)
	testEqual(t, []byte("Hello World"), []byte("Hello World"), true)

	testEqual(t, 123, 123, true)
	testEqual(t, 123.5, 123.5, true)
	testEqual(t, 123.5, 123.5000000001, false)
	testEqual(t, 123.5, 123, false)
	testEqual(t, int32(123), int32(123), true)
	testEqual(t, uint64(123), uint64(123), true)

	testEqual(t, testStruct{1, "a"}, testStruct{1, "a"}, true)
	testEqual(t, &testStruct{1, "a"}, &testStruct{1, "a"}, true)

	p := ptr(1)
	testEqual(t, p, p, true)
	testEqual(t, ptr(1), ptr(1), true)

	s := []int{1, 2}
	testEqual(t, s, s, true)
	testEqual(t, s, s[:], true)
	testEqual(t, s, s[:1], false)
	testEqual(t, &s, &s, true)

	testEqual(t, []int{1, 2, 3}, []int{1, 2, 3}, true)
	testEqual(t, []int{1, 2, 3}, []int{1, 2}, false)
	testEqual(t, &[]int{1, 2, 3}, &[]int{1, 2, 3}, true)
	testEqual(t, &[]int{1, 2, 3}, &[]int{1, 2}, false)

	m := map[string]int{"a": 1}
	testEqual(t, m, m, true)
	testEqual(t, map[string]int{"a": 1}, map[string]int{"a": 1}, true)
}

func TestEqualDelta(t *testing.T) {
	t.Parallel()

	testEqualDelta(t, 123, 123, 0, true)
	testEqualDelta(t, 123, 125, 2, true)
	testEqualDelta(t, 123, 125, 1, false)
	testEqualDelta(t, -10, -15, 5, true)
	testEqualDelta(t, -10, -15, 4, false)

	testEqualDelta(t, 123.0, 123.00001, 0.0001, true)
	testEqualDelta(t, 123.0, 123.00001, 0.000001, false)

	testEqualDelta(t, math.NaN(), math.NaN(), 1000, true)
	testEqualDelta(t, math.NaN(), 1, 1000, false)
	testEqualDelta(t, 2, math.NaN(), 1000, false)
	testEqualDelta(t, math.Inf(1), math.Inf(1), 1.0, true)
	testEqualDelta(t, math.Inf(1), math.Inf(-1), 1.0, false)

	testEqualDelta[uint32](t, 123, 125, 3, true)
	testEqualDelta(t, time.Millisecond*100, time.Millisecond*120, time.Millisecond*50, true)
}

func TestSame(t *testing.T) {
	testSame(t, "Hello World", "Hello World", false)
	testSame(t, 123, 123, false)

	v := 1
	p := &v
	testSame(t, v, v, false)
	testSame(t, &v, &v, true)
	testSame(t, p, &v, true)
	testSame(t, p, p, true)
	testSame(t, ptr(v), ptr(v), false)

	s := []int{1, 2}
	testSame(t, s, s, true)
	testSame(t, &s, &s, true)
	testSame(t, []int{1, 2}, []int{1, 2}, false)
	testSame(t, []byte("Hello World"), []byte("Hello World"), false)

	m := map[string]int{"a": 1}
	testSame(t, m, m, true)
	testSame(t, map[string]int{"a": 1}, map[string]int{"a": 1}, false)
}

func TestLength(t *testing.T) {
	testLength(t, []int{}, 0, true)
	testLength(t, []int{1, 2, 3}, 3, true)
	testLength(t, []int{1, 2, 3}, 2, false)
	testLength(t, "Hello", 5, true)
	testLength(t, map[string]bool{"a": true, "b": false}, 2, true)
}

func TestContains(t *testing.T) {
	testContains(t, []int{}, 0, false)
	testContains(t, []int{1, 2, 3}, 2, true)
	testContains(t, []int{1, 2, 3}, 4, false)
	testContains(t, "Hello", "e", true)
	testContains(t, "Hello", 2, false)
	testContains(t, map[string]bool{"a": true, "b": false}, true, true)
	testContains(t, map[string]bool{"a": true, "b": false}, "a", false)
}

func TestError(t *testing.T) {
	testError(t, nil, false)
	testError(t, errors.New("ooh"), true)
}

func TestErrorIs(t *testing.T) {
	err := errors.New("ooh")

	testErrorIs(t, nil, nil, true)
	testErrorIs(t, errors.New("ooh"), nil, false)
	testErrorIs(t, err, err, true)
	testErrorIs(t, errors.Join(errors.New("ooh1"), err), err, true)
}

func TestEqualJSON(t *testing.T) {
	testEqualJSON(t, "Hello World", "Hello World", false)
	testEqualJSON(t, "\"Hello World\"", "\"Hello World\"", true)
	testEqualJSON(t, "\"Hello World\"", "\"Hello World!\"", false)
	testEqualJSON(t, "123", "123", true)
	testEqualJSON(t, "123.0", "123", true)
	testEqualJSON(t, "123.3", "123", false)
	testEqualJSON(t, "false", "false", true)
	testEqualJSON(t, `{"x":10, "y":16}`, `{"x":10,"y":16.000}`, true)
}

type testType string

type testStruct struct {
	a int
	b string
}

type logger struct {
	LastError string
}

func (m *logger) Helper() {
}

func (m *logger) Errorf(format string, args ...any) {
	m.LastError = fmt.Sprintf(format, args...)
}

func (m *logger) Clear() {
	m.LastError = ""
}

var tt = &logger{}

func testAssert(t *testing.T, condition bool, result bool) {
	t.Helper()

	tt.Clear()
	if True(tt, condition) != result {
		t.Errorf("True(%#v) should return %#v: %s", condition, result, tt.LastError)
	}

	tt.Clear()
	if False(tt, condition) != !result {
		t.Errorf("False(%#v) should return %#v: %s", condition, !result, tt.LastError)
	}
}

func testEqual[T Comparable](t *testing.T, actual, expected T, result bool) {
	t.Helper()

	tt.Clear()
	if Equal(tt, actual, expected) != result {
		t.Errorf("Equal(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}

	tt.Clear()
	if NotEqual(tt, actual, expected) != !result {
		t.Errorf("NotEqual(%#v,%#v) should return %#v: %s", actual, expected, !result, tt.LastError)
	}
}

func testEqualDelta[T Numeric](t *testing.T, actual, expected, delta T, result bool) {
	t.Helper()

	tt.Clear()
	if EqualDelta(tt, actual, expected, delta) != result {
		t.Errorf("EqualDelta(%#v,%#v,%#v) should return %#v: %s", actual, expected, delta, result, tt.LastError)
	}
}

func testSame[T Reference](t *testing.T, actual, expected T, result bool) {
	t.Helper()

	tt.Clear()
	if Same(tt, actual, expected) != result {
		t.Errorf("Same(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}

	tt.Clear()
	if NotSame(tt, actual, expected) != !result {
		t.Errorf("NotSame(%#v,%#v) should return %#v: %s", actual, expected, !result, tt.LastError)
	}
}

func testLength[T any](t *testing.T, actual T, expected int, result bool) {
	t.Helper()

	tt.Clear()
	if Length(tt, actual, expected) != result {
		t.Errorf("Length(%#v,%#v) should return %#v: %s", actual, expected, !result, tt.LastError)
	}
}

func testContains[S Iterable[E], E Comparable](t *testing.T, object S, element E, result bool) {
	t.Helper()

	tt.Clear()
	if Contains(tt, object, element) != result {
		t.Errorf("Contains(%#v,%#v) should return %#v", object, element, result)
	}

	tt.Clear()
	if NotContains(tt, object, element) != !result {
		t.Errorf("NotContains(%#v,%#v) should return %#v", object, element, !result)
	}
}

func testError(t *testing.T, err error, result bool) {
	t.Helper()

	tt.Clear()
	if Error(tt, err) != result {
		t.Errorf("Error(%#v) should return %#v", err, result)
	}

	tt.Clear()
	if NoError(tt, err) != !result {
		t.Errorf("NoError(%#v) should return %#v", err, !result)
	}
}

func testErrorIs(t *testing.T, err, target error, result bool) {
	t.Helper()

	tt.Clear()
	if ErrorIs(tt, err, target) != result {
		t.Errorf("ErrorIs(%#v,%#v) should return %#v", err, target, result)
	}

	tt.Clear()
	if NotErrorIs(tt, err, target) != !result {
		t.Errorf("NotErrorIs(%#v,%#v) should return %#v", err, target, !result)
	}
}

func testEqualJSON(t *testing.T, actual, expected string, result bool) {
	t.Helper()

	tt.Clear()
	if EqualJSON(tt, actual, expected) != result {
		t.Errorf("EqualJSON(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}
}

func ptr(i int) *int {
	return &i
}
