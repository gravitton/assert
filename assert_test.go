package assert

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"testing"
	"time"
	"unsafe"
)

func TestFail(t *testing.T) {
	tt := newLogger()
	if Fail(tt, "custom failure") {
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
	t.Run("strings", func(t *testing.T) {
		testEqual(t, "Hello World", "Hello World", true)
		testEqual(t, "Hello World", "Hello World!", false)
		testEqual[testType](t, "A", "A", true)
		testEqual(t, []byte("Hello World"), []byte("Hello World"), true)
	})
	t.Run("numbers", func(t *testing.T) {
		testEqual(t, 123, 123, true)
		testEqual(t, 123.5, 123.5, true)
		testEqual(t, 123.5, 123.5000000001, false)
		testEqual(t, 123.5, 123, false)
		testEqual(t, int32(123), int32(123), true)
		testEqual(t, uint64(123), uint64(123), true)
	})
	t.Run("structs", func(t *testing.T) {
		testEqual(t, testStruct{1, "a"}, testStruct{1, "a"}, true)
		testEqual(t, &testStruct{1, "a"}, &testStruct{1, "a"}, true)
	})
	t.Run("pointers", func(t *testing.T) {
		p := ptr(1)
		testEqual(t, p, p, true)
		testEqual(t, ptr(1), ptr(1), true)
		testEqual(t, (*int)(nil), ptr(1), false)
		testEqual(t, (*int)(nil), (*int)(nil), true)
	})
	t.Run("slices", func(t *testing.T) {
		s := []int{1, 2}
		testEqual(t, s, s, true)
		testEqual(t, s, s[:], true)
		testEqual(t, s, s[:1], false)
		testEqual(t, &s, &s, true)

		testEqual(t, []int{1, 2, 3}, []int{1, 2, 3}, true)
		testEqual(t, []int{1, 2, 3}, []int{1, 2}, false)
		testEqual(t, &[]int{1, 2, 3}, &[]int{1, 2, 3}, true)
		testEqual(t, &[]int{1, 2, 3}, &[]int{1, 2}, false)
	})
	t.Run("maps", func(t *testing.T) {
		m := map[string]int{"a": 1}
		testEqual(t, m, m, true)
		testEqual(t, map[string]int{"a": 1}, map[string]int{"a": 1}, true)
	})
	t.Run("functions", func(t *testing.T) {
		f := func() {}

		testEqual(t, (func())(nil), (func())(nil), true)
		testEqual(t, f, f, false)
		testEqual(t, f, nil, false)
		testEqual[any](t, 1, f, false)
		testEqual(t, testFuncStruct{f}, testFuncStruct{f}, false)
		testEqual(t, testFuncStruct{}, testFuncStruct{}, true)
	})
}

func TestEqualDelta(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testEqualDelta(t, 123, 123, 0, true)
		testEqualDelta(t, 123, 125, 2, true)
		testEqualDelta(t, 123, 125, 1, false)
		testEqualDelta(t, -10, -15, 5, true)
		testEqualDelta(t, -10, -15, 4, false)
		testEqualDelta[uint32](t, 123, 125, 3, true)
		testEqualDelta[uintptr](t, 10, 12, 2, true)
	})
	t.Run("integer bounds", func(t *testing.T) {
		testEqualDelta[uint64](t, 1<<60, 1<<60+1, 0, false)
		testEqualDelta[uint64](t, 1<<60, 1<<60+1, 1, true)
		testEqualDelta[int64](t, math.MaxInt64, math.MinInt64, math.MaxInt64, false)
		testEqualDelta[int8](t, 127, -128, 127, false)
	})
	t.Run("floats", func(t *testing.T) {
		testEqualDelta(t, 123.0, 123.00001, 0.0001, true)
		testEqualDelta(t, 123.0, 123.00001, 0.000001, false)
	})
	t.Run("NaN and infinity", func(t *testing.T) {
		testEqualDelta(t, math.NaN(), math.NaN(), 1000, true)
		testEqualDelta(t, math.NaN(), 1, 1000, false)
		testEqualDelta(t, 2, math.NaN(), 1000, false)
		testEqualDelta(t, math.Inf(1), math.Inf(1), 1.0, true)
		testEqualDelta(t, math.Inf(1), math.Inf(-1), 1.0, false)
	})
	t.Run("durations", func(t *testing.T) {
		testEqualDelta(t, time.Millisecond*100, time.Millisecond*120, time.Millisecond*50, true)
	})
	t.Run("invalid delta", func(t *testing.T) {
		testEqualDeltaInvalid(t, 1, 2, -1)
		testEqualDeltaInvalid(t, 1.0, 2.0, math.NaN())
	})
}

func TestNil(t *testing.T) {
	t.Run("untyped nil", func(t *testing.T) {
		testNil(t, nil, true)
	})
	t.Run("typed nil", func(t *testing.T) {
		var p *int
		var s []int
		var m map[string]int
		var c chan int
		var f func()
		var e error

		testNil(t, p, true)
		testNil(t, s, true)
		testNil(t, m, true)
		testNil(t, c, true)
		testNil(t, f, true)
		testNil(t, e, true)
		testNil(t, unsafe.Pointer(nil), true)
	})
	t.Run("non-nil", func(t *testing.T) {
		testNil(t, unsafe.Pointer(ptr(1)), false)
		testNil(t, ptr(1), false)
		testNil(t, []int{}, false)
		testNil(t, map[string]int{}, false)
		testNil(t, 0, false)
		testNil(t, "", false)
		testNil(t, testStruct{}, false)
	})
}

func TestZero(t *testing.T) {
	t.Run("basic types", func(t *testing.T) {
		testZero(t, nil, true)
		testZero(t, 0, true)
		testZero(t, 1, false)
		testZero(t, math.Copysign(0, -1), true)
		testZero(t, "", true)
		testZero(t, "a", false)
		testZero(t, false, true)
	})
	t.Run("composite types", func(t *testing.T) {
		testZero(t, testStruct{}, true)
		testZero(t, testStruct{b: "a"}, false)
		testZero(t, [2]int{}, true)
		testZero(t, [2]int{0, 1}, false)
	})
	t.Run("references", func(t *testing.T) {
		var e error = (*testPtrErr)(nil)

		testZero(t, (*int)(nil), true)
		testZero(t, ptr(0), false)
		testZero(t, []int(nil), true)
		testZero(t, []int{}, false)
		testZero(t, map[string]int(nil), true)
		testZero(t, map[string]int{}, false)
		testZero(t, e, true)
	})
	t.Run("IsZero method", func(t *testing.T) {
		testZero(t, testQuantity{}, true)
		testZero(t, testQuantity{unit: "kg"}, true)
		testZero(t, testQuantity{value: 1}, false)
		testZero(t, &testQuantity{unit: "kg"}, false)
		testZero(t, &testQuantity{value: 1}, false)
		testZero(t, (*testQuantity)(nil), true)
		testZero(t, &time.Time{}, false)
		testZero(t, time.Time{}, true)
		testZero(t, time.Unix(0, 0), false)
		testZero(t, time.Time{}.In(time.FixedZone("CET", 3600)), true)
		testZero(t, struct {
			time.Time
			X int
		}{X: 1}, true)
	})
}

func TestSame(t *testing.T) {
	t.Run("invalid", func(t *testing.T) {
		v := 1

		testSameInvalid(t, "Hello World", "Hello World")
		testSameInvalid(t, 123, 123)
		testSameInvalid[any](t, nil, nil)
		testSameInvalid(t, v, v)
		testSameInvalid[any](t, &v, 5)
		testSameInvalid[any](t, 5, &v)
	})
	t.Run("pointers", func(t *testing.T) {
		v := 1
		p := &v

		testSame(t, &v, &v, true)
		testSame(t, p, &v, true)
		testSame(t, p, p, true)
		testSame(t, ptr(v), ptr(v), false)
		testSame(t, (*int)(nil), (*int)(nil), true)
		testSame[any](t, (*int)(nil), (*int64)(nil), false)
	})
	t.Run("slices", func(t *testing.T) {
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
		testSame(t, []int(nil), []int(nil), true)
	})
	t.Run("maps", func(t *testing.T) {
		m := map[string]int{"a": 1}

		testSame(t, m, m, true)
		testSame(t, map[string]int{"a": 1}, map[string]int{"a": 1}, false)
		testSame(t, map[string]int(nil), map[string]int(nil), true)
	})
	t.Run("types", func(t *testing.T) {
		var pair struct{ A int }

		testSame[any](t, &pair, &pair.A, false)
		testSame[any](t, &pair, &pair, true)
	})
}

func TestGreater(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testOrder(t, Greater, "Greater", 2, 1, true)
		testOrder(t, Greater, "Greater", 1, 1, false)
		testOrder(t, Greater, "Greater", 0, 1, false)
		testOrder(t, Greater, "Greater", -1, -2, true)
		testOrder[uint32](t, Greater, "Greater", 5, 3, true)
	})
	t.Run("floats", func(t *testing.T) {
		testOrder[float64](t, Greater, "Greater", 1.1, 1.0, true)
		testOrder[float64](t, Greater, "Greater", 1.0, 1.0, false)
		testOrder(t, Greater, "Greater", math.NaN(), 1.0, false)
		testOrder(t, Greater, "Greater", 1.0, math.NaN(), false)
	})
	t.Run("strings", func(t *testing.T) {
		testOrder(t, Greater, "Greater", "b", "a", true)
		testOrder(t, Greater, "Greater", "a", "b", false)
	})
}

func TestGreaterOrEqual(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testOrder(t, GreaterOrEqual, "GreaterOrEqual", 2, 1, true)
		testOrder(t, GreaterOrEqual, "GreaterOrEqual", 1, 1, true)
		testOrder(t, GreaterOrEqual, "GreaterOrEqual", 0, 1, false)
		testOrder(t, GreaterOrEqual, "GreaterOrEqual", -1, -2, true)
	})
	t.Run("floats", func(t *testing.T) {
		testOrder[float64](t, GreaterOrEqual, "GreaterOrEqual", 1.0, 1.0, true)
		testOrder[float64](t, GreaterOrEqual, "GreaterOrEqual", 0.9, 1.0, false)
		testOrder(t, GreaterOrEqual, "GreaterOrEqual", math.NaN(), math.NaN(), false)
	})
}

func TestLess(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testOrder(t, Less, "Less", 1, 2, true)
		testOrder(t, Less, "Less", 1, 1, false)
		testOrder(t, Less, "Less", 2, 1, false)
		testOrder(t, Less, "Less", -2, -1, true)
		testOrder[uint32](t, Less, "Less", 3, 5, true)
	})
	t.Run("floats", func(t *testing.T) {
		testOrder[float64](t, Less, "Less", 1.0, 1.1, true)
		testOrder[float64](t, Less, "Less", 1.0, 1.0, false)
		testOrder(t, Less, "Less", math.NaN(), 1.0, false)
	})
	t.Run("strings", func(t *testing.T) {
		testOrder(t, Less, "Less", "a", "b", true)
		testOrder[testType](t, Less, "Less", "b", "a", false)
	})
}

func TestLessOrEqual(t *testing.T) {
	t.Run("integers", func(t *testing.T) {
		testOrder(t, LessOrEqual, "LessOrEqual", 1, 2, true)
		testOrder(t, LessOrEqual, "LessOrEqual", 1, 1, true)
		testOrder(t, LessOrEqual, "LessOrEqual", 2, 1, false)
		testOrder(t, LessOrEqual, "LessOrEqual", -2, -1, true)
	})
	t.Run("floats", func(t *testing.T) {
		testOrder[float64](t, LessOrEqual, "LessOrEqual", 1.0, 1.0, true)
		testOrder[float64](t, LessOrEqual, "LessOrEqual", 1.1, 1.0, false)
		testOrder(t, LessOrEqual, "LessOrEqual", 1.0, math.NaN(), false)
	})
}

func TestLength(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		testLength(t, "Hello", 5, true)
	})
	t.Run("slices", func(t *testing.T) {
		testLength(t, []int{}, 0, true)
		testLength(t, []int{1, 2, 3}, 3, true)
		testLength(t, []int{1, 2, 3}, 2, false)
		testLength(t, []int(nil), 0, true)
	})
	t.Run("arrays", func(t *testing.T) {
		testLength(t, [2]int{1, 2}, 2, true)
		testLength(t, &[3]int{1, 2, 3}, 3, true)
		testLength(t, (*[2]int)(nil), 2, true)
	})
	t.Run("maps", func(t *testing.T) {
		testLength(t, map[string]bool{"a": true, "b": false}, 2, true)
		testLength(t, map[string]int(nil), 0, true)
	})
	t.Run("channels", func(t *testing.T) {
		testLength(t, bufferedChan(3), 3, true)
		testLength(t, (chan int)(nil), 0, true)
	})
	t.Run("invalid", func(t *testing.T) {
		testLength(t, ptr(1), 1, false)
		testLength(t, 5, 1, false)
		testLength[any](t, nil, 0, false)
	})
}

func TestEmpty(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		testEmpty(t, "", true)
		testEmpty(t, "a", false)
	})
	t.Run("slices", func(t *testing.T) {
		testEmpty(t, []int{}, true)
		testEmpty(t, []int{1}, false)
		testEmpty(t, []int(nil), true)
	})
	t.Run("arrays", func(t *testing.T) {
		testEmpty(t, &[0]int{}, true)
		testEmpty(t, &[1]int{}, false)
	})
	t.Run("maps", func(t *testing.T) {
		testEmpty(t, map[string]bool{}, true)
		testEmpty(t, map[string]bool{"a": true}, false)
		testEmpty(t, map[string]int(nil), true)
	})
	t.Run("channels", func(t *testing.T) {
		testEmpty(t, (chan int)(nil), true)
		testEmpty(t, bufferedChan(0), true)
		testEmpty(t, bufferedChan(1), false)
	})
	t.Run("invalid", func(t *testing.T) {
		tt := newLogger()
		if Empty(tt, 5) {
			t.Errorf("Empty(5) should return false: %s", tt.LastError)
		}

		tt = newLogger()
		if NotEmpty(tt, 5) {
			t.Errorf("NotEmpty(5) should return false: %s", tt.LastError)
		}
	})
}

func TestContains(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		testContains(t, "Hello", "e", true)
		testContains(t, testType("Hello"), testType("e"), true)
	})
	t.Run("slices", func(t *testing.T) {
		testContains(t, []int{}, 0, false)
		testContains(t, []int{1, 2, 3}, 2, true)
		testContains(t, []int{1, 2, 3}, 4, false)
		testContains(t, []testType{"a"}, testType("a"), true)
	})
	t.Run("arrays", func(t *testing.T) {
		testContains(t, [3]int{1, 2, 3}, 2, true)
		testContains(t, &[3]int{1, 2, 3}, 2, true)
		testContains(t, &[3]int{1, 2, 3}, 4, false)
	})
	t.Run("map values", func(t *testing.T) {
		testContains(t, map[string]bool{"a": true, "b": false}, true, true)
		testContains(t, map[string]bool{"a": true, "b": true}, false, false)
	})
	t.Run("interfaces", func(t *testing.T) {
		testContains[[]any, any](t, []any{1, nil}, nil, true)
		testContains[[]any, any](t, []any{1, 2}, nil, false)
		testContains(t, []any{1, "two", 3}, 3, true)
		testContains(t, []any{1, "two", 3}, 4, false)
		testContains(t, []any{1, "two", 3}, "two", true)
		testContains(t, []error{testSliceErr{"a"}}, testSliceErr{"a"}, true)
	})
	t.Run("pointers", func(t *testing.T) {
		p := ptr(1)

		testContains(t, []*int{p}, p, true)
		testContains(t, []*int{p}, ptr(1), true)
		testContains(t, []*int{p}, ptr(2), false)
		testContains[[]*int, *int](t, []*int{p}, nil, false)
		testContains[[]*int, *int](t, []*int{p, nil}, nil, true)
	})
	t.Run("element type", func(t *testing.T) {
		testContainsInvalid(t, []testType{"a"}, "a")
		testContainsInvalid(t, []string{"a"}, testType("a"))
		testContainsInvalid(t, []testSliceErr{{"a"}}, []string{"a"})
		testContainsInvalid(t, []<-chan int{bufferedChan(0)}, bufferedChan(0))
		testContainsInvalid(t, testType("Hello"), "e")
		testContainsInvalid(t, "Hello", testType("e"))
		testContainsInvalid[string, any](t, "Hello", nil)
		testContainsInvalid(t, "Hello", 2)
		testContainsInvalid(t, "<int Value>", 2)
		testContainsInvalid(t, "true", true)
		testContainsInvalid(t, map[string]bool{"a": true, "b": false}, "a")
		testContainsInvalid(t, []int{1}, "x")
		testContainsInvalid[[]int, any](t, []int{1}, nil)
	})
	t.Run("invalid", func(t *testing.T) {
		testContainsInvalid[any, int](t, nil, 5)
		testContainsInvalid(t, 5, 5)
		testContainsInvalid(t, bufferedChan(1), 1)
		testContainsInvalid(t, (*[3]int)(nil), 1)
		testContainsInvalid(t, ptr(1), 1)
	})
}

func TestEqualUnordered(t *testing.T) {
	t.Run("slices", func(t *testing.T) {
		testEqualUnordered(t, []int{1, 2, 3}, []int{3, 1, 2}, true)
		testEqualUnordered(t, []int{1, 2, 3}, []int{1, 2, 3}, true)
		testEqualUnordered(t, []int{1, 2}, []int{1, 2, 3}, false)
		testEqualUnordered(t, []int{1, 2, 3}, []int{1, 2}, false)
	})
	t.Run("duplicates", func(t *testing.T) {
		testEqualUnordered(t, []int{1, 2, 2}, []int{2, 1, 2}, true)
		testEqualUnordered(t, []int{1, 2, 2}, []int{1, 1, 2}, false)
		testEqualUnordered(t, [2]int{1, 1}, [2]int{1, 2}, false)
	})
	t.Run("nil and empty", func(t *testing.T) {
		testEqualUnordered(t, []int(nil), []int{}, true)
		testEqualUnordered(t, []int(nil), []int(nil), true)
	})
	t.Run("arrays", func(t *testing.T) {
		testEqualUnordered(t, [3]int{1, 2, 3}, [3]int{2, 3, 1}, true)
		testEqualUnordered(t, &[3]int{1, 2, 3}, &[3]int{2, 3, 1}, true)
		testEqualUnordered(t, &[3]int{1, 2, 3}, &[3]int{2, 3, 2}, false)
	})
	t.Run("element types", func(t *testing.T) {
		testEqualUnordered(t, testSliceErr{"a", "b"}, testSliceErr{"b", "a"}, true)
		testEqualUnordered(t, []*int{ptr(1), ptr(2)}, []*int{ptr(2), ptr(1)}, true)
		testEqualUnordered(t, []any{1, "a", nil}, []any{nil, 1, "a"}, true)
		testEqualUnordered(t, []testStruct{{1, "a"}, {2, "b"}}, []testStruct{{2, "b"}, {1, "a"}}, true)
	})
	t.Run("invalid", func(t *testing.T) {
		testEqualUnordered(t, "ab", "ba", false)
		testEqualUnordered(t, 1, 1, false)
		testEqualUnordered(t, map[int]int{1: 1}, map[int]int{1: 1}, false)
		testEqualUnordered[any](t, nil, []int{}, false)
		testEqualUnordered[any](t, []int{1}, []int64{1}, false)
		testEqualUnordered[any](t, []int{1}, [1]int{1}, false)
		testEqualUnordered(t, (*[1]int)(nil), (*[1]int)(nil), false)
	})
}

func TestHasPrefix(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		testAffix(t, HasPrefix, "HasPrefix", "Hello World", "Hello", true)
		testAffix(t, HasPrefix, "HasPrefix", "Hello World", "World", false)
		testAffix(t, HasPrefix, "HasPrefix", "Hello", "Hello World", false)
		testAffix(t, HasPrefix, "HasPrefix", "Hello", "", true)
		testAffix(t, HasPrefix, "HasPrefix", "", "", true)
		testAffix[testType](t, HasPrefix, "HasPrefix", "Hello", "He", true)
	})
	t.Run("slices", func(t *testing.T) {
		testAffix(t, HasPrefix, "HasPrefix", []int{1, 2, 3}, []int{1, 2}, true)
		testAffix(t, HasPrefix, "HasPrefix", []int{1, 2, 3}, []int{2, 3}, false)
		testAffix(t, HasPrefix, "HasPrefix", []int{1, 2, 3}, []int{1, 2, 3}, true)
		testAffix(t, HasPrefix, "HasPrefix", []int{1, 2}, []int{1, 2, 3}, false)
		testAffix(t, HasPrefix, "HasPrefix", []int{1}, nil, true)
		testAffix(t, HasPrefix, "HasPrefix", []int(nil), []int{}, true)
		testAffix(t, HasPrefix, "HasPrefix", []*int{ptr(1), ptr(2)}, []*int{ptr(1)}, true)
		testAffix(t, HasPrefix, "HasPrefix", testSliceErr{"a", "b"}, testSliceErr{"a"}, true)
	})
	t.Run("invalid", func(t *testing.T) {
		testAffix(t, HasPrefix, "HasPrefix", 5, 5, false)
		testAffix(t, HasPrefix, "HasPrefix", [2]int{1, 2}, [2]int{1, 2}, false)
		testAffix(t, HasPrefix, "HasPrefix", map[int]int{}, map[int]int{}, false)
		testAffix[any](t, HasPrefix, "HasPrefix", nil, "a", false)
		testAffix[any](t, HasPrefix, "HasPrefix", "a", nil, false)
		testAffix[any](t, HasPrefix, "HasPrefix", "a", testType("a"), false)
		testAffix[any](t, HasPrefix, "HasPrefix", []int{1}, []int64{1}, false)
	})
}

func TestHasSuffix(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		testAffix(t, HasSuffix, "HasSuffix", "Hello World", "World", true)
		testAffix(t, HasSuffix, "HasSuffix", "Hello World", "Hello", false)
		testAffix(t, HasSuffix, "HasSuffix", "World", "Hello World", false)
		testAffix(t, HasSuffix, "HasSuffix", "Hello", "", true)
		testAffix(t, HasSuffix, "HasSuffix", "", "", true)
		testAffix[testType](t, HasSuffix, "HasSuffix", "Hello", "lo", true)
	})
	t.Run("slices", func(t *testing.T) {
		testAffix(t, HasSuffix, "HasSuffix", []int{1, 2, 3}, []int{2, 3}, true)
		testAffix(t, HasSuffix, "HasSuffix", []int{1, 2, 3}, []int{1, 2}, false)
		testAffix(t, HasSuffix, "HasSuffix", []int{1, 2, 3}, []int{1, 2, 3}, true)
		testAffix(t, HasSuffix, "HasSuffix", []int{2, 3}, []int{1, 2, 3}, false)
		testAffix(t, HasSuffix, "HasSuffix", []int{1}, nil, true)
		testAffix(t, HasSuffix, "HasSuffix", []int(nil), []int{}, true)
		testAffix(t, HasSuffix, "HasSuffix", []*int{ptr(1), ptr(2)}, []*int{ptr(2)}, true)
		testAffix(t, HasSuffix, "HasSuffix", testSliceErr{"a", "b"}, testSliceErr{"b"}, true)
	})
	t.Run("invalid", func(t *testing.T) {
		testAffix(t, HasSuffix, "HasSuffix", 5, 5, false)
		testAffix(t, HasSuffix, "HasSuffix", [2]int{1, 2}, [2]int{1, 2}, false)
		testAffix(t, HasSuffix, "HasSuffix", map[int]int{}, map[int]int{}, false)
		testAffix[any](t, HasSuffix, "HasSuffix", nil, "a", false)
		testAffix[any](t, HasSuffix, "HasSuffix", "a", nil, false)
		testAffix[any](t, HasSuffix, "HasSuffix", "a", testType("a"), false)
		testAffix[any](t, HasSuffix, "HasSuffix", []int{1}, []int64{1}, false)
	})
}

func TestError(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		testError(t, nil, false)
	})
	t.Run("non-nil", func(t *testing.T) {
		testError(t, errors.New("ooh"), true)
		testError(t, testErr{}, true)
		testError(t, testSliceErr{"a"}, true)
	})
	t.Run("typed nil", func(t *testing.T) {
		var err *testErr

		testError(t, err, true)
		testError(t, testSliceErr(nil), true)
	})
}

func TestErrorIs(t *testing.T) {
	err := errors.New("ooh")

	t.Run("nil", func(t *testing.T) {
		testErrorIs(t, nil, nil, true)
		testErrorIs(t, errors.New("ooh"), nil, false)
	})
	t.Run("direct", func(t *testing.T) {
		testErrorIs(t, err, err, true)
	})
	t.Run("wrapped", func(t *testing.T) {
		testErrorIs(t, errors.Join(errors.New("ooh1"), err), err, true)
	})
}

func TestErrorAs(t *testing.T) {
	var target *testPtrErr
	var iface interface{ Error() string }

	t.Run("assignable", func(t *testing.T) {
		testErrorAs(t, &testPtrErr{"a"}, &target, true)
		testErrorAs(t, fmt.Errorf("wrapped: %w", &testPtrErr{"a"}), &target, true)
		testErrorAs(t, errors.New("plain"), &iface, true)
	})
	t.Run("not assignable", func(t *testing.T) {
		testErrorAs(t, errors.New("plain"), &target, false)
		testErrorAs(t, nil, &target, false)
	})
	t.Run("invalid target", func(t *testing.T) {
		testErrorAs(t, errors.New("plain"), nil, false)
		testErrorAs(t, errors.New("plain"), target, false)
		testErrorAs(t, errors.New("plain"), (*error)(nil), false)
		testErrorAs(t, errors.New("plain"), ptr(1), false)
	})
}

func TestErrorContains(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		testErrorContains(t, errors.New("file not found"), "not found", true)
		testErrorContains(t, errors.New("file not found"), "", true)
		testErrorContains(t, fmt.Errorf("open: %w", errors.New("not found")), "open: not", true)
		testErrorContains(t, testSliceErr(nil), "Slice", true)
		testErrorContains(t, (*testPtrErr)(nil), "<nil>", true)
	})
	t.Run("not contains", func(t *testing.T) {
		testErrorContains(t, errors.New("file not found"), "denied", false)
	})
	t.Run("nil", func(t *testing.T) {
		testErrorContains(t, nil, "", false)
	})
}

func TestMatches(t *testing.T) {
	t.Run("matching", func(t *testing.T) {
		testMatches(t, "Hello World", `^Hello`, true)
		testMatches(t, "Hello World", `World$`, true)
		testMatches(t, "abc123", `\d+`, true)
	})
	t.Run("not matching", func(t *testing.T) {
		testMatches(t, "Hello World", `\d+`, false)
	})
	t.Run("invalid pattern", func(t *testing.T) {
		testMatches(t, "Hello World", `[`, false)

		tt := newLogger()
		if NotMatches(tt, "Hello World", `[`) {
			t.Errorf("NotMatches with invalid pattern should return false: %s", tt.LastError)
		}
	})
}

func TestEqualJSON(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		testEqualJSON(t, "\"Hello World\"", "\"Hello World\"", true)
		testEqualJSON(t, "\"Hello World\"", "\"Hello World!\"", false)
		testEqualJSON(t, "false", "false", true)
	})
	t.Run("structures", func(t *testing.T) {
		testEqualJSON(t, `{"x":10, "y":16}`, `{"x":10,"y":16.000}`, true)
		testEqualJSON(t, `[1, [2, {"a": 3.0}]]`, `[1, [2, {"a": 3}]]`, true)
	})
	t.Run("numbers", func(t *testing.T) {
		testEqualJSON(t, "123", "123", true)
		testEqualJSON(t, "123.0", "123", true)
		testEqualJSON(t, "123.3", "123", false)
		testEqualJSON(t, `[1e2, 0.5]`, `[100, 5e-1]`, true)
		testEqualJSON(t, `1E+2`, `100`, true)
		testEqualJSON(t, `1.50`, `15e-1`, true)
		testEqualJSON(t, `0.00`, `-0e5`, true)
		testEqualJSON(t, `-1`, `1`, false)
		testEqualJSON(t, `0.001`, `1e-3`, true)
		testEqualJSON(t, `0.001`, `1e-2`, false)
	})
	t.Run("float64 numbers", func(t *testing.T) {
		testEqualJSON(t, `{"id":9007199254740993}`, `{"id":9007199254740992}`, true)
		testEqualJSON(t, `{"id":9007199254740994}`, `{"id":9007199254740992}`, false)
		testEqualJSON(t, `1.0000000000000001`, `1`, true)
		testEqualJSON(t, `1e308`, `10e307`, true)
		testEqualJSON(t, `1e9999999`, `1e9999999`, false)
	})
	t.Run("invalid", func(t *testing.T) {
		testEqualJSON(t, "Hello World", "Hello World", false)
		testEqualJSON(t, `1 2`, `1`, false)
		testEqualJSON(t, `1`, `1 }`, false)
		testEqualJSON(t, "1\f", `1`, false)
		testEqualJSON(t, "1\u00a0", `1`, false)
		testEqualJSON(t, `{1:2}`, `{1:2}`, false)
		testEqualJSON(t, `{"a" 1}`, `{"a" 1}`, false)
		testEqualJSON(t, `{"a":1 "b":2}`, `{"a":1 "b":2}`, false)
		testEqualJSON(t, `{"a":1,}`, `{"a":1,}`, false)
		testEqualJSON(t, `[1 2]`, `[1 2]`, false)
		testEqualJSON(t, `[1,]`, `[1,]`, false)
		testEqualJSON(t, `[1`, `[1`, false)
		testEqualJSON(t, ``, ``, false)
	})
	t.Run("duplicate keys", func(t *testing.T) {
		testEqualJSON(t, `{"a":1,"a":2}`, `{"a":2}`, true)
		testEqualJSON(t, `{"a":1,"a":2}`, `{"a":1}`, false)
		testEqualJSON(t, `[{"a":1}, {"a":1}]`, `[{"a":1}, {"a":1}]`, true)
		testEqualJSON(t, `{"a":{"b":1,"b":1}}`, `{"a":{"b":1}}`, true)
	})
	t.Run("empty", func(t *testing.T) {
		testEqualJSON(t, `{}`, ` { } `, true)
		testEqualJSON(t, `[]`, `[ ]`, true)
		testEqualJSON(t, `[]`, `{}`, false)
		testEqualJSON(t, `null`, `null`, true)
	})
	t.Run("whitespace", func(t *testing.T) {
		testEqualJSON(t, " \t\r\n1 \t\r\n", `1`, true)
	})
}

func TestJSON(t *testing.T) {
	t.Run("maps", func(t *testing.T) {
		testJSON(t, map[string]any{"a": 1, "b": true, "c": "Hello"}, `{"a":1,"b":true,"c":"Hello"}`, true)
	})
	t.Run("structs", func(t *testing.T) {
		testJSON(t, struct {
			A int `json:"a"`
			B string
		}{
			A: 1,
			B: "Hello",
		}, `{"a":1,"B":"Hello"}`, true)
	})
	t.Run("invalid", func(t *testing.T) {
		testJSON(t, "Hello", "", false)
	})
}

func TestPanics(t *testing.T) {
	testPanics(t, func() { panic("boom") }, true)
	testPanics(t, func() {}, false)
	testPanics(t, func() { panic(nil) }, true)
	testPanics(t, nil, false)
}

func TestPanicsWith(t *testing.T) {
	err := errors.New("oops")

	t.Run("no panic", func(t *testing.T) {
		testPanicsWith(t, func() {}, "boom", false)
	})
	t.Run("values", func(t *testing.T) {
		testPanicsWith(t, func() { panic("boom") }, "boom", true)
		testPanicsWith(t, func() { panic("boom") }, "other", false)
		testPanicsWith(t, func() { panic("boom") }, testType("boom"), false)
	})
	t.Run("errors", func(t *testing.T) {
		testPanicsWith(t, func() { panic(err) }, err, true)
		testPanicsWith(t, func() { panic(fmt.Errorf("wrapped: %w", err)) }, err, true)
		testPanicsWith(t, func() { panic(errors.New("other")) }, err, false)
		testPanicsWith(t, func() { panic("not-an-error") }, err, false)
		testPanicsWith(t, func() { panic(testSliceErr{"a"}) }, testSliceErr{"a"}, true)
		testPanicsWith(t, func() { panic(testSliceErr{"a"}) }, testSliceErr{"b"}, false)
	})
	t.Run("nil", func(t *testing.T) {
		testPanicsWith(t, func() { panic(nil) }, nil, true)
		testPanicsWith(t, func() { panic("boom") }, nil, false)
		testPanicsWith(t, nil, nil, false)
	})
}

func TestNotPanics(t *testing.T) {
	testNotPanics(t, func() { panic("boom") }, false)
	testNotPanics(t, func() {}, true)
	testNotPanics(t, nil, false)
}

func TestMessagesWithAddresses(t *testing.T) {
	x := 1
	p := &x
	s := []int{1, 2}
	c := make(chan int)

	var cycle selfPointer
	cycle = &cycle

	testMessages(t, []messageCase{
		{"Same", func(t Testing) bool { return Same(t, s, s[:1]) }, fmt.Sprintf("Should be same\n  actual: [%p] []int{1, 2}\nexpected: [%p] []int{1}", s, s)},
		{"NotSame", func(t Testing) bool { return NotSame(t, &x, &x) }, fmt.Sprintf("Should not be same\n  actual: [%p] (*int)(%p)", &x, &x)},
		{"Equal function", func(t Testing) bool { return Equal(t, ptr, ptr) }, fmt.Sprintf("Should be equal\n  actual: (func(int) *int)(%p)\nexpected: (func(int) *int)(%p)\n    hint: values print the same but are not deeply equal", ptr, ptr)},
		{"JSON", func(t Testing) bool { return JSON(t, c, "1") }, fmt.Sprintf("Should be marshalable\n  actual: (chan int)(%p)", c)},
		{"nested pointer", func(t Testing) bool { return Equal(t, &p, nil) }, fmt.Sprintf("Should be equal\n  actual: &(*int)(%p)\nexpected: (**int)(nil)", p)},
		{"pointer cycle", func(t Testing) bool { return Nil(t, cycle) }, fmt.Sprintf("Should be nil\n  actual: &(assert.selfPointer)(%p)", cycle)},
	})
}

func TestMessages(t *testing.T) {
	testMessages(t, []messageCase{
		{"bytes", func(t Testing) bool { return Contains(t, []byte("abc"), 'a') }, "Should have element of same type\n  actual: []byte{0x61, 0x62, 0x63}\n element: int32(97)"},
		{"Contains nil array pointer", func(t Testing) bool { return Contains(t, (*[3]int)(nil), 1) }, "Should be non-nil array pointer\n  actual: (*[3]int)(nil)\n element: 1"},
		{"EqualUnordered nil array pointer", func(t Testing) bool { return EqualUnordered(t, &[1]int{1}, nil) }, "Should be non-nil array pointer\n  actual: &[1]int{1}\nexpected: (*[1]int)(nil)"},
	})
}

type messageCase struct {
	name     string
	assert   func(t Testing) bool
	expected string
}

type testType string

type selfPointer *selfPointer

type testFuncStruct struct {
	f func()
}

type testStruct struct {
	a int
	b string
}

type testQuantity struct {
	value int
	unit  string
}

func (z testQuantity) IsZero() bool {
	return z.value == 0
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

func testMessages(t *testing.T, cases []messageCase) {
	t.Helper()

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
	if Same(tt, actual, expected) {
		t.Errorf("Same(%#v,%#v) should return false: %s", actual, expected, tt.LastError)
	}

	tt = newLogger()
	if NotSame(tt, actual, expected) {
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

func testEmpty[T any](t *testing.T, actual T, result bool) {
	t.Helper()

	tt := newLogger()
	if Empty(tt, actual) != result {
		t.Errorf("Empty(%#v) should return %#v: %s", actual, result, tt.LastError)
	}

	tt = newLogger()
	if NotEmpty(tt, actual) != !result {
		t.Errorf("NotEmpty(%#v) should return %#v: %s", actual, !result, tt.LastError)
	}
}

func testContains[S Iterable, E Comparable](t *testing.T, actual S, element E, result bool) {
	t.Helper()

	tt := newLogger()
	if Contains(tt, actual, element) != result {
		t.Errorf("Contains(%#v,%#v) should return %#v: %s", actual, element, result, tt.LastError)
	}

	tt = newLogger()
	if NotContains(tt, actual, element) != !result {
		t.Errorf("NotContains(%#v,%#v) should return %#v: %s", actual, element, !result, tt.LastError)
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

func testNil(t *testing.T, actual any, result bool) {
	t.Helper()

	tt := newLogger()
	if Nil(tt, actual) != result {
		t.Errorf("Nil(%#v) should return %#v: %s", actual, result, tt.LastError)
	}

	tt = newLogger()
	if NotNil(tt, actual) != !result {
		t.Errorf("NotNil(%#v) should return %#v: %s", actual, !result, tt.LastError)
	}
}

func testZero(t *testing.T, actual any, result bool) {
	t.Helper()

	tt := newLogger()
	if Zero(tt, actual) != result {
		t.Errorf("Zero(%#v) should return %#v: %s", actual, result, tt.LastError)
	}

	tt = newLogger()
	if NotZero(tt, actual) != !result {
		t.Errorf("NotZero(%#v) should return %#v: %s", actual, !result, tt.LastError)
	}
}

func testEqualUnordered[S Iterable](t *testing.T, actual, expected S, result bool) {
	t.Helper()

	tt := newLogger()
	if EqualUnordered(tt, actual, expected) != result {
		t.Errorf("EqualUnordered(%#v,%#v) should return %#v: %s", actual, expected, result, tt.LastError)
	}
}

func testAffix[S Iterable](t *testing.T, assertion func(Testing, S, S, ...string) bool, name string, actual, affix S, result bool) {
	t.Helper()

	tt := newLogger()
	if assertion(tt, actual, affix) != result {
		t.Errorf("%s(%#v,%#v) should return %#v: %s", name, actual, affix, result, tt.LastError)
	}
}

func testErrorContains(t *testing.T, err error, substr string, result bool) {
	t.Helper()

	tt := newLogger()
	if ErrorContains(tt, err, substr) != result {
		t.Errorf("ErrorContains(%#v,%#v) should return %#v: %s", err, substr, result, tt.LastError)
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

func testContainsInvalid[S Iterable, E Comparable](t *testing.T, actual S, element E) {
	t.Helper()

	tt := newLogger()
	if Contains(tt, actual, element) {
		t.Errorf("Contains(%#v,%#v) should return false: %s", actual, element, tt.LastError)
	}

	tt = newLogger()
	if NotContains(tt, actual, element) {
		t.Errorf("NotContains(%#v,%#v) should return false: %s", actual, element, tt.LastError)
	}
}

func ptr(i int) *int {
	return &i
}

func bufferedChan(n int) chan int {
	c := make(chan int, n)
	for i := range n {
		c <- i
	}

	return c
}
