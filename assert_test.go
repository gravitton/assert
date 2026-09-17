package assert

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode/utf8"
	"unsafe"
)

func TestFail(t *testing.T) {
	tt := newLogger()
	if Fail(tt, "custom failure") != false {
		t.Errorf("Fail should return false")
	}
	if tt.LastError != "custom failure" {
		t.Errorf("Fail should report message, got: %s", tt.LastError)
	}
}

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

	testEqual(t, (*int)(nil), ptr(1), false)
	testEqual(t, (*int)(nil), (*int)(nil), true)
}

func TestEqualDelta(t *testing.T) {
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
	testEqualDelta[uint64](t, 1<<60, 1<<60+1, 0, false)
	testEqualDelta[uint64](t, 1<<60, 1<<60+1, 1, true)
	testEqualDelta[int64](t, math.MaxInt64, math.MinInt64, math.MaxInt64, false)
	testEqualDelta[int8](t, 127, -128, 127, false)

	testEqualDeltaInvalid(t, 1, 2, -1)
	testEqualDeltaInvalid(t, 1.0, 2.0, math.NaN())
	testEqualDelta[uintptr](t, 10, 12, 2, true)
}

func TestNil(t *testing.T) {
	var p *int
	var s []int
	var m map[string]int
	var c chan int
	var f func()
	var e error

	testNil(t, nil, true)
	testNil(t, p, true)
	testNil(t, s, true)
	testNil(t, m, true)
	testNil(t, c, true)
	testNil(t, f, true)
	testNil(t, e, true)
	testNil(t, unsafe.Pointer(nil), true)
	testNil(t, unsafe.Pointer(ptr(1)), false)
	testNil(t, ptr(1), false)
	testNil(t, []int{}, false)
	testNil(t, map[string]int{}, false)
	testNil(t, 0, false)
	testNil(t, "", false)
	testNil(t, testStruct{}, false)
}

func TestSame(t *testing.T) {
	testSameInvalid(t, "Hello World", "Hello World")
	testSameInvalid(t, 123, 123)
	testSameInvalid[any](t, nil, nil)

	v := 1
	p := &v
	testSameInvalid(t, v, v)
	testSameInvalid[any](t, &v, 5)
	testSameInvalid[any](t, 5, &v)
	testSame(t, &v, &v, true)
	testSame(t, p, &v, true)
	testSame(t, p, p, true)
	testSame(t, ptr(v), ptr(v), false)

	s := []int{1, 2}
	testSame(t, s, s, true)
	testSame(t, s, s[:], true)
	testSame(t, s, s[:0], false)
	testSame(t, s, s[:1], false)
	testSame(t, s[:1], s[:1], true)
	testSame(t, s, s[1:], false)
	testSame(t, &s, &s, true)
	testSame(t, []int{1, 2}, []int{1, 2}, false)
	testSame(t, []byte("Hello World"), []byte("Hello World"), false)

	m := map[string]int{"a": 1}
	testSame(t, m, m, true)
	testSame(t, map[string]int{"a": 1}, map[string]int{"a": 1}, false)

	var pair struct{ A int }
	testSame[any](t, &pair, &pair.A, false)
	testSame[any](t, &pair, &pair, true)
	testSame(t, (*int)(nil), (*int)(nil), true)
}

func TestGreater(t *testing.T) {
	testOrder(t, Greater, "Greater", 2, 1, true)
	testOrder(t, Greater, "Greater", 1, 1, false)
	testOrder(t, Greater, "Greater", 0, 1, false)
	testOrder(t, Greater, "Greater", -1, -2, true)
	testOrder[float64](t, Greater, "Greater", 1.1, 1.0, true)
	testOrder[float64](t, Greater, "Greater", 1.0, 1.0, false)
	testOrder[uint32](t, Greater, "Greater", 5, 3, true)
	testOrder(t, Greater, "Greater", math.NaN(), 1.0, false)
	testOrder(t, Greater, "Greater", 1.0, math.NaN(), false)
	testOrder(t, Greater, "Greater", "b", "a", true)
	testOrder(t, Greater, "Greater", "a", "b", false)
}

func TestGreaterOrEqual(t *testing.T) {
	testOrder(t, GreaterOrEqual, "GreaterOrEqual", 2, 1, true)
	testOrder(t, GreaterOrEqual, "GreaterOrEqual", 1, 1, true)
	testOrder(t, GreaterOrEqual, "GreaterOrEqual", 0, 1, false)
	testOrder(t, GreaterOrEqual, "GreaterOrEqual", -1, -2, true)
	testOrder[float64](t, GreaterOrEqual, "GreaterOrEqual", 1.0, 1.0, true)
	testOrder[float64](t, GreaterOrEqual, "GreaterOrEqual", 0.9, 1.0, false)
	testOrder(t, GreaterOrEqual, "GreaterOrEqual", math.NaN(), math.NaN(), false)
}

func TestLess(t *testing.T) {
	testOrder(t, Less, "Less", 1, 2, true)
	testOrder(t, Less, "Less", 1, 1, false)
	testOrder(t, Less, "Less", 2, 1, false)
	testOrder(t, Less, "Less", -2, -1, true)
	testOrder[float64](t, Less, "Less", 1.0, 1.1, true)
	testOrder[float64](t, Less, "Less", 1.0, 1.0, false)
	testOrder[uint32](t, Less, "Less", 3, 5, true)
	testOrder(t, Less, "Less", math.NaN(), 1.0, false)
	testOrder(t, Less, "Less", "a", "b", true)
	testOrder[testType](t, Less, "Less", "b", "a", false)
}

func TestLessOrEqual(t *testing.T) {
	testOrder(t, LessOrEqual, "LessOrEqual", 1, 2, true)
	testOrder(t, LessOrEqual, "LessOrEqual", 1, 1, true)
	testOrder(t, LessOrEqual, "LessOrEqual", 2, 1, false)
	testOrder(t, LessOrEqual, "LessOrEqual", -2, -1, true)
	testOrder[float64](t, LessOrEqual, "LessOrEqual", 1.0, 1.0, true)
	testOrder[float64](t, LessOrEqual, "LessOrEqual", 1.1, 1.0, false)
	testOrder(t, LessOrEqual, "LessOrEqual", 1.0, math.NaN(), false)
}

func TestLength(t *testing.T) {
	testLength(t, []int{}, 0, true)
	testLength(t, []int{1, 2, 3}, 3, true)
	testLength(t, []int{1, 2, 3}, 2, false)
	testLength(t, "Hello", 5, true)
	testLength(t, map[string]bool{"a": true, "b": false}, 2, true)
	testLength(t, [2]int{1, 2}, 2, true)
	testLength(t, bufferedChan(3), 3, true)
	testLength(t, []int(nil), 0, true)
	testLength(t, map[string]int(nil), 0, true)
	testLength(t, (chan int)(nil), 0, true)
	testLength(t, 5, 1, false)
	testLength[any](t, nil, 0, false)
}

func TestEmpty(t *testing.T) {
	testEmpty(t, []int{}, true)
	testEmpty(t, []int{1}, false)
	testEmpty(t, "", true)
	testEmpty(t, "a", false)
	testEmpty(t, map[string]bool{}, true)
	testEmpty(t, map[string]bool{"a": true}, false)
	testEmpty(t, []int(nil), true)
	testEmpty(t, map[string]int(nil), true)
	testEmpty(t, (chan int)(nil), true)
	testEmpty(t, bufferedChan(0), true)
	testEmpty(t, bufferedChan(1), false)

	tt := newLogger()
	if Empty(tt, 5) != false {
		t.Errorf("Empty(5) should return false: %s", tt.LastError)
	}

	tt = newLogger()
	if NotEmpty(tt, 5) != false {
		t.Errorf("NotEmpty(5) should return false: %s", tt.LastError)
	}
}

func TestContains(t *testing.T) {
	testContains(t, []int{}, 0, false)
	testContains(t, []int{1, 2, 3}, 2, true)
	testContains(t, []int{1, 2, 3}, 4, false)
	testContains(t, "Hello", "e", true)
	testContains(t, map[string]bool{"a": true, "b": false}, true, true)
	testContains(t, map[string]bool{"a": true, "b": true}, false, false)

	testContains(t, [3]int{1, 2, 3}, 2, true)
	testContains(t, []testType{"a"}, testType("a"), true)
	testContainsInvalid(t, []testType{"a"}, "a")
	testContainsInvalid(t, []string{"a"}, testType("a"))
	testContains(t, testType("Hello"), testType("e"), true)
	testContains[[]any, any](t, []any{1, nil}, nil, true)
	testContains[[]any, any](t, []any{1, 2}, nil, false)

	p := ptr(1)
	testContains(t, []*int{p}, p, true)
	testContains(t, []*int{p}, ptr(1), true)
	testContains(t, []*int{p}, ptr(2), false)
	testContains[[]*int, *int](t, []*int{p}, nil, false)
	testContains[[]*int, *int](t, []*int{p, nil}, nil, true)
	testContains(t, []any{1, "two", 3}, 3, true)
	testContains(t, []any{1, "two", 3}, 4, false)
	testContains(t, []any{1, "two", 3}, "two", true)

	testContainsInvalid(t, "Hello", 2)
	testContainsInvalid(t, "<int Value>", 2)
	testContainsInvalid(t, "true", true)
	testContainsInvalid(t, map[string]bool{"a": true, "b": false}, "a")
	testContainsInvalid(t, []int{1}, "x")
	testContainsInvalid[[]int, any](t, []int{1}, nil)
	testContainsInvalid[any, int](t, nil, 5)
	testContainsInvalid(t, 5, 5)
	testContainsInvalid(t, bufferedChan(1), 1)
}

func TestError(t *testing.T) {
	var err *testErr = nil

	testError(t, nil, false)
	testError(t, errors.New("ooh"), true)
	testError(t, err, false)
	testError(t, testErr{}, true)
	testError(t, testSliceErr(nil), false)
	testError(t, testSliceErr{"a"}, true)
}

func TestErrorIs(t *testing.T) {
	err := errors.New("ooh")

	testErrorIs(t, nil, nil, true)
	testErrorIs(t, errors.New("ooh"), nil, false)
	testErrorIs(t, err, err, true)
	testErrorIs(t, errors.Join(errors.New("ooh1"), err), err, true)
}

func TestErrorAs(t *testing.T) {
	var target *testPtrErr
	var iface interface{ Error() string }

	testErrorAs(t, &testPtrErr{"a"}, &target, true)
	testErrorAs(t, fmt.Errorf("wrapped: %w", &testPtrErr{"a"}), &target, true)
	testErrorAs(t, errors.New("plain"), &target, false)
	testErrorAs(t, nil, &target, false)
	testErrorAs(t, errors.New("plain"), &iface, true)
}

func TestMatches(t *testing.T) {
	testMatches(t, "Hello World", `^Hello`, true)
	testMatches(t, "Hello World", `World$`, true)
	testMatches(t, "Hello World", `\d+`, false)
	testMatches(t, "abc123", `\d+`, true)
	testMatches(t, "Hello World", `[`, false) // invalid regexp

	tt := newLogger()
	if NotMatches(tt, "Hello World", `[`) != false {
		t.Errorf("NotMatches with invalid pattern should return false: %s", tt.LastError)
	}
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
	testEqualJSON(t, `{"id":9007199254740993}`, `{"id":9007199254740992}`, false)
	testEqualJSON(t, `{"id":9007199254740993}`, `{"id":9007199254740993}`, true)
	testEqualJSON(t, `[1e2, 0.5]`, `[100, 5e-1]`, true)
	testEqualJSON(t, `[1, [2, {"a": 3.0}]]`, `[1, [2, {"a": 3}]]`, true)
	testEqualJSON(t, `1e9999999`, `1e9999999`, true)
	testEqualJSON(t, `1e9999999`, `1e9999998`, false)
	testEqualJSON(t, `1 2`, `1`, false)
	testEqualJSON(t, `1`, `1 }`, false)
}

func TestJSON(t *testing.T) {
	testJSON(t, map[string]any{"a": 1, "b": true, "c": "Hello"}, `{"a":1,"b":true,"c":"Hello"}`, true)
	testJSON(t, "Hello", "", false)
	testJSON(t, struct {
		A int `json:"a"`
		B string
	}{
		A: 1,
		B: "Hello",
	}, `{"a":1,"B":"Hello"}`, true)
}

func TestPanics(t *testing.T) {
	testPanics(t, func() { panic("boom") }, true)
	testPanics(t, func() {}, false)
	testPanics(t, func() { panic(nil) }, true)
}

func TestPanicsWith(t *testing.T) {
	testPanicsWith(t, func() {}, "boom", false)
	testPanicsWith(t, func() { panic("boom") }, "boom", true)
	testPanicsWith(t, func() { panic("boom") }, "other", false)
	testPanicsWith(t, func() { panic("boom") }, testType("boom"), false)

	err := errors.New("oops")
	testPanicsWith(t, func() { panic(err) }, err, true)
	testPanicsWith(t, func() { panic(fmt.Errorf("wrapped: %w", err)) }, err, true)
	testPanicsWith(t, func() { panic(errors.New("other")) }, err, false)
	testPanicsWith(t, func() { panic("not-an-error") }, err, false)

	testPanicsWith(t, func() { panic(testSliceErr{"a"}) }, testSliceErr{"a"}, true)
	testPanicsWith(t, func() { panic(testSliceErr{"a"}) }, testSliceErr{"b"}, false)

	testPanicsWith(t, func() { panic(nil) }, nil, true)
	testPanicsWith(t, func() { panic("boom") }, nil, false)
}

func TestNotPanics(t *testing.T) {
	testNotPanics(t, func() { panic("boom") }, false)
	testNotPanics(t, func() {}, true)
}

func TestMessages(t *testing.T) {
	err := errors.New("oops")
	target := errors.New("target")
	x := 1
	var ptrErr *testPtrErr

	cases := []struct {
		name     string
		assert   func(t Testing) bool
		expected string
	}{
		{"Failf", func(t Testing) bool { return Failf(t, "x=%d", 1) }, "x=1"},
		{"True", func(t Testing) bool { return True(t, false, "ctx: ") }, "ctx: Should be true"},
		{"False", func(t Testing) bool { return False(t, true, "ctx: ") }, "ctx: Should be false"},
		{"Same", func(t Testing) bool { return Same(t, ptr(1), ptr(1), "ctx: ") }, "ctx: Should be same\n"},
		{"Same invalid", func(t Testing) bool { return Same(t, 1, 1, "ctx: ") }, "ctx: Should be reference\n  actual: 1\nexpected: 1"},
		{"NotSame", func(t Testing) bool { return NotSame(t, &x, &x, "ctx: ") }, "ctx: Should not be same\n"},
		{"NotSame invalid", func(t Testing) bool { return NotSame(t, 1, 1, "ctx: ") }, "ctx: Should be reference\n  actual: 1\nexpected: 1"},
		{"Nil", func(t Testing) bool { return Nil(t, &x, "ctx: ") }, "ctx: Should be nil\n  actual: ["},
		{"NotNil", func(t Testing) bool { return NotNil(t, (*int)(nil), "ctx: ") }, "ctx: Should not be nil\n  actual: (*int)(nil)"},
		{"Equal", func(t Testing) bool { return Equal(t, 1, 2, "ctx: ") }, "ctx: Should be equal\n  actual: 1\nexpected: 2"},
		{"Equal nil pointer", func(t Testing) bool { return Equal(t, (*int)(nil), &x, "ctx: ") }, "ctx: Should be equal\n  actual: (*int)(nil)\nexpected: ["},
		{"NotEqual", func(t Testing) bool { return NotEqual(t, 1, 1, "ctx: ") }, "ctx: Should not be equal\n  actual: 1"},
		{"EqualDelta", func(t Testing) bool { return EqualDelta(t, 1, 3, 1, "ctx: ") }, "ctx: Should be equal in delta\n  actual: 1\nexpected: 3"},
		{"NotEqualDelta", func(t Testing) bool { return NotEqualDelta(t, 1, 2, 1, "ctx: ") }, "ctx: Should not be equal in delta\n  actual: 1\nexpected: 2"},
		{"EqualDelta invalid", func(t Testing) bool { return EqualDelta(t, 1, 2, -1, "ctx: ") }, "ctx: Should have non-negative delta\n   delta: -1"},
		{"NotEqualDelta invalid", func(t Testing) bool { return NotEqualDelta(t, 1, 2, -1, "ctx: ") }, "ctx: Should have non-negative delta\n   delta: -1"},
		{"Greater", func(t Testing) bool { return Greater(t, 1, 2, "ctx: ") }, "ctx: Should be greater\n  actual: 1\nexpected: 2"},
		{"GreaterOrEqual", func(t Testing) bool { return GreaterOrEqual(t, 1, 2, "ctx: ") }, "ctx: Should be greater or equal\n  actual: 1\nexpected: 2"},
		{"Less", func(t Testing) bool { return Less(t, 2, 1, "ctx: ") }, "ctx: Should be less\n  actual: 2\nexpected: 1"},
		{"LessOrEqual", func(t Testing) bool { return LessOrEqual(t, 2, 1, "ctx: ") }, "ctx: Should be less or equal\n  actual: 2\nexpected: 1"},
		{"Length", func(t Testing) bool { return Length(t, []int{1}, 2, "ctx: ") }, "ctx: Should have length\n  object: ["},
		{"Length invalid", func(t Testing) bool { return Length(t, 5, 2, "ctx: ") }, "ctx: Should be iterable\n  object: 5"},
		{"Empty", func(t Testing) bool { return Empty(t, []int{1}, "ctx: ") }, "ctx: Should be empty\n  object: ["},
		{"Empty invalid", func(t Testing) bool { return Empty(t, 5, "ctx: ") }, "ctx: Should be iterable\n  object: 5"},
		{"NotEmpty", func(t Testing) bool { return NotEmpty(t, []int{}, "ctx: ") }, "ctx: Should not be empty\n  object: ["},
		{"NotEmpty invalid", func(t Testing) bool { return NotEmpty(t, 5, "ctx: ") }, "ctx: Should be iterable\n  object: 5"},
		{"Contains", func(t Testing) bool { return Contains(t, []int{1}, 2, "ctx: ") }, "ctx: Should contain element\n  object: ["},
		{"Contains invalid", func(t Testing) bool { return Contains(t, 5, 2, "ctx: ") }, "ctx: Should be iterable\n  object: 5\n element: 2"},
		{"Contains element type", func(t Testing) bool { return Contains(t, []int{1}, "a", "ctx: ") }, "ctx: Should have element of same type\n  object: ["},
		{"NotContains", func(t Testing) bool { return NotContains(t, []int{1}, 1, "ctx: ") }, "ctx: Should not contain element\n  object: ["},
		{"NotContains invalid", func(t Testing) bool { return NotContains(t, 5, 2, "ctx: ") }, "ctx: Should be iterable\n  object: 5\n element: 2"},
		{"Error", func(t Testing) bool { return Error(t, nil, "ctx: ") }, "ctx: Should be error"},
		{"NoError", func(t Testing) bool { return NoError(t, err, "ctx: ") }, "ctx: Should not be error\n     msg: oops\n   error: &errors.errorString{s:\"oops\"}"},
		{"ErrorIs", func(t Testing) bool { return ErrorIs(t, err, target, "ctx: ") }, "ctx: Should match error\n     msg: oops\n   error: &errors.errorString{s:\"oops\"}\n  target: &errors.errorString{s:\"target\"}"},
		{"NotErrorIs", func(t Testing) bool { return NotErrorIs(t, err, err, "ctx: ") }, "ctx: Should not match error\n"},
		{"ErrorAs", func(t Testing) bool { return ErrorAs(t, err, &ptrErr, "ctx: ") }, "ctx: Should be assignable to target\n     msg: oops\n   error: &errors.errorString{s:\"oops\"}\n  target: **assert.testPtrErr"},
		{"Matches", func(t Testing) bool { return Matches(t, "a", "b", "ctx: ") }, "ctx: Should match regexp\n  actual: a\n pattern: b"},
		{"Matches invalid", func(t Testing) bool { return Matches(t, "a", "[", "ctx: ") }, "ctx: Should be valid regexp\n pattern: [\n     err: "},
		{"NotMatches", func(t Testing) bool { return NotMatches(t, "a", "a", "ctx: ") }, "ctx: Should not match regexp\n  actual: a\n pattern: a"},
		{"NotMatches invalid", func(t Testing) bool { return NotMatches(t, "a", "[", "ctx: ") }, "ctx: Should be valid regexp\n pattern: [\n     err: "},
		{"EqualJSON", func(t Testing) bool { return EqualJSON(t, "1", "2", "ctx: ") }, "ctx: Should be equal JSON\n  actual: 1\nexpected: 2"},
		{"EqualJSON trailing", func(t Testing) bool { return EqualJSON(t, "1 x", "1", "ctx: ") }, "ctx: Should be valid JSON\n  actual: 1 x\n     err: invalid character 'x' after top-level value"},
		{"EqualJSON invalid actual", func(t Testing) bool { return EqualJSON(t, "x", "2", "ctx: ") }, "ctx: Should be valid JSON\n  actual: x\n     err: "},
		{"EqualJSON invalid expected", func(t Testing) bool { return EqualJSON(t, "1", "x", "ctx: ") }, "ctx: Should be valid JSON\nexpected: x\n     err: "},
		{"JSON", func(t Testing) bool { return JSON(t, 1, "2", "ctx: ") }, "ctx: Should be equal JSON\n  actual: 1\nexpected: 2"},
		{"JSON unmarshalable", func(t Testing) bool { return JSON(t, make(chan int), "2", "ctx: ") }, "ctx: Should be marshalable\n  actual: [0x"},
		{"Panics", func(t Testing) bool { return Panics(t, func() {}, "ctx: ") }, "ctx: Should panic"},
		{"PanicsWith", func(t Testing) bool { return PanicsWith(t, func() {}, "b", "ctx: ") }, "ctx: Should panic\nexpected: \"b\""},
		{"PanicsWith value", func(t Testing) bool { return PanicsWith(t, func() { panic("a") }, "b", "ctx: ") }, "ctx: Should panic with value\n  actual: \"a\"\nexpected: \"b\""},
		{"PanicsWith error", func(t Testing) bool { return PanicsWith(t, func() { panic("a") }, err, "ctx: ") }, "ctx: Should panic with value\n  actual: \"a\"\nexpected: ["},
		{"PanicsWith wrong error", func(t Testing) bool { return PanicsWith(t, func() { panic(err) }, target, "ctx: ") }, "ctx: Should panic with value\n  actual: ["},
		{"NotPanics nil", func(t Testing) bool { return NotPanics(t, func() { panic(nil) }, "ctx: ") }, "ctx: Should not panic\n  value: <nil>"},
		{"NotPanics", func(t Testing) bool { return NotPanics(t, func() { panic("a") }, "ctx: ") }, "ctx: Should not panic\n  value: \"a\""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tt := newLogger()
			if c.assert(tt) {
				t.Errorf("should return false")
			}

			if !strings.HasPrefix(tt.LastError, c.expected) {
				t.Errorf("message should start with %q, got %q", c.expected, tt.LastError)
			}
		})
	}
}

type testType string

type testStruct struct {
	a int
	b string
}

type testErr struct{}

func (t testErr) Error() string {
	return "Custom error"
}

type testPtrErr struct {
	msg string
}

func (t *testPtrErr) Error() string {
	return t.msg
}

type testSliceErr []string

func (t testSliceErr) Error() string {
	return "Slice error"
}

type logger struct {
	LastError string
}

func newLogger() *logger {
	return &logger{}
}

func (m *logger) Helper() {
}

func (m *logger) Errorf(format string, args ...any) {
	m.LastError = fmt.Sprintf(format, args...)
}

func testAssert(t *testing.T, condition bool, result bool) {
	t.Helper()

	tt := newLogger()
	if True(tt, condition) != result {
		t.Errorf("True(%#v) should return %#v: %s", condition, result, tt.LastError)
	}

	tt = newLogger()
	if False(tt, condition) != !result {
		t.Errorf("False(%#v) should return %#v: %s", condition, !result, tt.LastError)
	}
}

func testEqual[T Comparable](t *testing.T, actual, expected T, result bool) {
	t.Helper()

	tt := newLogger()
	if Equal(tt, actual, expected) != result {
		t.Errorf("Equal(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}

	tt = newLogger()
	if NotEqual(tt, actual, expected) != !result {
		t.Errorf("NotEqual(%#v,%#v) should return %#v: %s", actual, expected, !result, tt.LastError)
	}
}

func testEqualDelta[T Numeric](t *testing.T, actual, expected, delta T, result bool) {
	t.Helper()

	tt := newLogger()
	if EqualDelta(tt, actual, expected, delta) != result {
		t.Errorf("EqualDelta(%#v,%#v,%#v) should return %#v: %s", actual, expected, delta, result, tt.LastError)
	}

	tt = newLogger()
	if NotEqualDelta(tt, actual, expected, delta) != !result {
		t.Errorf("NotEqualDelta(%#v,%#v,%#v) should return %#v: %s", actual, expected, delta, !result, tt.LastError)
	}
}

func testSame[T Reference](t *testing.T, actual, expected T, result bool) {
	t.Helper()

	tt := newLogger()
	if Same(tt, actual, expected) != result {
		t.Errorf("Same(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}

	tt = newLogger()
	if NotSame(tt, actual, expected) != !result {
		t.Errorf("NotSame(%#v,%#v) should return %#v: %s", actual, expected, !result, tt.LastError)
	}
}

func testSameInvalid[T Reference](t *testing.T, actual, expected T) {
	t.Helper()

	tt := newLogger()
	if Same(tt, actual, expected) != false {
		t.Errorf("Same(%#v,%#v) should return false: %s", actual, expected, tt.LastError)
	}

	tt = newLogger()
	if NotSame(tt, actual, expected) != false {
		t.Errorf("NotSame(%#v,%#v) should return false: %s", actual, expected, tt.LastError)
	}
}

func testOrder[T Ordered](t *testing.T, assertion func(Testing, T, T, ...string) bool, name string, actual, expected T, result bool) {
	t.Helper()

	tt := newLogger()
	if assertion(tt, actual, expected) != result {
		t.Errorf("%s(%#v,%#v) should return %#v: %s", name, actual, expected, result, tt.LastError)
	}
}

func testLength[T any](t *testing.T, actual T, expected int, result bool) {
	t.Helper()

	tt := newLogger()
	if Length(tt, actual, expected) != result {
		t.Errorf("Length(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}
}

func testEmpty[T any](t *testing.T, object T, result bool) {
	t.Helper()

	tt := newLogger()
	if Empty(tt, object) != result {
		t.Errorf("Empty(%#v) should return %#v: %s", object, result, tt.LastError)
	}

	tt = newLogger()
	if NotEmpty(tt, object) != !result {
		t.Errorf("NotEmpty(%#v) should return %#v: %s", object, !result, tt.LastError)
	}
}

func testContains[S Iterable, E Comparable](t *testing.T, object S, element E, result bool) {
	t.Helper()

	tt := newLogger()
	if Contains(tt, object, element) != result {
		t.Errorf("Contains(%#v,%#v) should return %#v: %s", object, element, result, tt.LastError)
	}

	tt = newLogger()
	if NotContains(tt, object, element) != !result {
		t.Errorf("NotContains(%#v,%#v) should return %#v: %s", object, element, !result, tt.LastError)
	}
}

func testError(t *testing.T, err error, result bool) {
	t.Helper()

	tt := newLogger()
	if Error(tt, err) != result {
		t.Errorf("Error(%#v) should return %#v", err, result)
	}

	tt = newLogger()
	if NoError(tt, err) != !result {
		t.Errorf("NoError(%#v) should return %#v", err, !result)
	}
}

func testErrorIs(t *testing.T, err, target error, result bool) {
	t.Helper()

	tt := newLogger()
	if ErrorIs(tt, err, target) != result {
		t.Errorf("ErrorIs(%#v,%#v) should return %#v", err, target, result)
	}

	tt = newLogger()
	if NotErrorIs(tt, err, target) != !result {
		t.Errorf("NotErrorIs(%#v,%#v) should return %#v", err, target, !result)
	}
}

func testMatches(t *testing.T, actual, pattern string, result bool) {
	t.Helper()

	tt := newLogger()
	if Matches(tt, actual, pattern) != result {
		t.Errorf("Matches(%#v,%#v) should return %#v: %s", actual, pattern, result, tt.LastError)
	}

	// NotMatches inverts the result only for valid patterns
	if _, err := regexp.Compile(pattern); err == nil {
		tt = newLogger()
		if NotMatches(tt, actual, pattern) != !result {
			t.Errorf("NotMatches(%#v,%#v) should return %#v: %s", actual, pattern, !result, tt.LastError)
		}
	}
}

func testEqualJSON(t *testing.T, actual, expected string, result bool) {
	t.Helper()

	tt := newLogger()
	if EqualJSON(tt, actual, expected) != result {
		t.Errorf("EqualJSON(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}
}

func testPanics(t *testing.T, fn func(), result bool) {
	t.Helper()

	tt := newLogger()
	if Panics(tt, fn) != result {
		t.Errorf("Panics() should return %#v: %s", result, tt.LastError)
	}
}

func testPanicsWith(t *testing.T, fn func(), expected any, result bool) {
	t.Helper()

	tt := newLogger()
	if PanicsWith(tt, fn, expected) != result {
		t.Errorf("PanicsWith(%#v) should return %#v: %s", expected, result, tt.LastError)
	}
}

func testEqualDeltaInvalid[T Numeric](t *testing.T, actual, expected, delta T) {
	t.Helper()

	tt := newLogger()
	if EqualDelta(tt, actual, expected, delta) {
		t.Errorf("EqualDelta(%v,%v,%v) should return false", actual, expected, delta)
	}

	tt = newLogger()
	if NotEqualDelta(tt, actual, expected, delta) {
		t.Errorf("NotEqualDelta(%v,%v,%v) should return false", actual, expected, delta)
	}
}

func TestFormatTruncate(t *testing.T) {
	long := strings.Repeat("é", formatLimit)
	formatted := format(long)

	if !strings.HasSuffix(formatted, "…") {
		t.Errorf("long value should be truncated, got %d bytes", len(formatted))
	}

	if !utf8.ValidString(formatted) {
		t.Errorf("truncated value should stay valid UTF-8")
	}

	if format("short") != `"short"` {
		t.Errorf("short value should not be truncated")
	}
}

func testNil(t *testing.T, object any, result bool) {
	t.Helper()

	tt := newLogger()
	if Nil(tt, object) != result {
		t.Errorf("Nil(%#v) should return %#v: %s", object, result, tt.LastError)
	}

	tt = newLogger()
	if NotNil(tt, object) != !result {
		t.Errorf("NotNil(%#v) should return %#v: %s", object, !result, tt.LastError)
	}
}

func testErrorAs(t *testing.T, err error, target any, result bool) {
	t.Helper()

	tt := newLogger()
	if ErrorAs(tt, err, target) != result {
		t.Errorf("ErrorAs(%#v,%T) should return %#v: %s", err, target, result, tt.LastError)
	}
}

func testNotPanics(t *testing.T, fn func(), result bool) {
	t.Helper()

	tt := newLogger()
	if NotPanics(tt, fn) != result {
		t.Errorf("NotPanics() should return %#v: %s", result, tt.LastError)
	}
}

func testJSON(t *testing.T, actual any, expected string, result bool) {
	t.Helper()

	tt := newLogger()
	if JSON(tt, actual, expected) != result {
		t.Errorf("JSON(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}
}

func testContainsInvalid[S Iterable, E Comparable](t *testing.T, object S, element E) {
	t.Helper()

	tt := newLogger()
	if Contains(tt, object, element) != false {
		t.Errorf("Contains(%#v,%#v) should return false: %s", object, element, tt.LastError)
	}

	tt = newLogger()
	if NotContains(tt, object, element) != false {
		t.Errorf("NotContains(%#v,%#v) should return false: %s", object, element, tt.LastError)
	}
}

func ptr(i int) *int {
	return &i
}

func bufferedChan(n int) chan int {
	c := make(chan int, n)
	for i := 0; i < n; i++ {
		c <- i
	}

	return c
}
