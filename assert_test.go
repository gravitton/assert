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
		testEqualInvalid(t, f, f)
		testEqualInvalid(t, f, nil)
		testEqualInvalid[any](t, 1, f)
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
	t.Run("large numbers", func(t *testing.T) {
		testEqualJSON(t, `{"id":9007199254740993}`, `{"id":9007199254740992}`, false)
		testEqualJSON(t, `{"id":9007199254740993}`, `{"id":9007199254740993}`, true)
		testEqualJSON(t, `1e9999999`, `1e9999999`, true)
		testEqualJSON(t, `1e9999999`, `1e9999998`, false)
		testEqualJSON(t, `1e1000001`, `10e1000000`, true)
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
		testEqualJSON(t, `{"a":1,"a":2}`, `{"a":2}`, false)
		testEqualJSON(t, `[{"a":1}, {"a":1}]`, `[{"a":1}, {"a":1}]`, true)
		testEqualJSON(t, `{"a":{"b":1,"b":1}}`, `{"a":{"b":1}}`, false)
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

func TestMessages(t *testing.T) {
	err := errors.New("oops")
	target := errors.New("target")
	x := 1
	var ptrErr *testPtrErr

	t.Run("conditions", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Failf", func(t Testing) bool { return Failf(t, "x=%d", 1) }, "x=1"},
			{"True", func(t Testing) bool { return True(t, false, "ctx: ") }, "ctx: Should be true"},
			{"False", func(t Testing) bool { return False(t, true, "ctx: ") }, "ctx: Should be false"},
		})
	})
	t.Run("nil and zero", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Nil", func(t Testing) bool { return Nil(t, &x, "ctx: ") }, "ctx: Should be nil\n  actual: &1"},
			{"NotNil", func(t Testing) bool { return NotNil(t, (*int)(nil), "ctx: ") }, "ctx: Should not be nil\n  actual: (*int)(nil)"},
			{"Zero", func(t Testing) bool { return Zero(t, 1, "ctx: ") }, "ctx: Should be zero\n  actual: 1"},
			{"NotZero", func(t Testing) bool { return NotZero(t, "", "ctx: ") }, "ctx: Should not be zero\n  actual: \"\""},
		})
	})
	t.Run("identity", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Same", func(t Testing) bool { return Same(t, ptr(1), ptr(1), "ctx: ") }, "ctx: Should be same\n"},
			{"Same type", func(t Testing) bool { return Same[any](t, &x, new(int64), "ctx: ") }, "ctx: Should have same type\n  actual: *int\nexpected: *int64"},
			{"Same invalid", func(t Testing) bool { return Same(t, 1, 1, "ctx: ") }, "ctx: Should be reference\n  actual: 1\nexpected: 1"},
			{"NotSame", func(t Testing) bool { return NotSame(t, &x, &x, "ctx: ") }, "ctx: Should not be same\n"},
			{"NotSame invalid", func(t Testing) bool { return NotSame(t, 1, 1, "ctx: ") }, "ctx: Should be reference\n  actual: 1\nexpected: 1"},
		})
	})
	t.Run("equality", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Equal", func(t Testing) bool { return Equal(t, 1, 2, "ctx: ") }, "ctx: Should be equal\n  actual: 1\nexpected: 2"},
			{"Equal nil pointer", func(t Testing) bool { return Equal(t, (*int)(nil), &x, "ctx: ") }, "ctx: Should be equal\n  actual: (*int)(nil)\nexpected: &1"},
			{"Equal types", func(t Testing) bool { return Equal[any](t, 1, int64(1), "ctx: ") }, "ctx: Should be equal\n  actual: int(1)\nexpected: int64(1)"},
			{"Equal NaN", func(t Testing) bool { return Equal(t, math.NaN(), math.NaN(), "ctx: ") }, "ctx: Should be equal\n  actual: NaN\nexpected: NaN"},
			{"Equal unsigned", func(t Testing) bool { return Equal(t, uint(5), uint(6), "ctx: ") }, "ctx: Should be equal\n  actual: 5\nexpected: 6"},
			{"Equal function", func(t Testing) bool { return Equal(t, TestMessages, TestMessages, "ctx: ") }, "ctx: Should not be function\n  actual: (func(*testing.T))(0x"},
			{"NotEqual", func(t Testing) bool { return NotEqual(t, 1, 1, "ctx: ") }, "ctx: Should not be equal\n  actual: 1"},
			{"NotEqual function", func(t Testing) bool { return NotEqual(t, TestMessages, nil, "ctx: ") }, "ctx: Should not be function\n  actual: (func(*testing.T))(0x"},
		})
	})
	t.Run("delta", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"EqualDelta", func(t Testing) bool { return EqualDelta(t, 1, 3, 1, "ctx: ") }, "ctx: Should be equal in delta\n  actual: 1\nexpected: 3\n   delta: 1\n    diff: 2"},
			{"NotEqualDelta", func(t Testing) bool { return NotEqualDelta(t, 1, 2, 1, "ctx: ") }, "ctx: Should not be equal in delta\n  actual: 1\nexpected: 2\n   delta: 1\n    diff: 1"},
			{"EqualDelta float", func(t Testing) bool { return EqualDelta(t, 1.5, 1.0, 0.25, "ctx: ") }, "ctx: Should be equal in delta\n  actual: 1.5\nexpected: 1\n   delta: 0.25\n    diff: 0.5"},
			{"EqualDelta duration", func(t Testing) bool { return EqualDelta(t, time.Second, 2*time.Second, time.Second/2, "ctx: ") }, "ctx: Should be equal in delta\n  actual: 1s\nexpected: 2s\n   delta: 500ms\n    diff: 1s"},
			{"EqualDelta overflow", func(t Testing) bool { return EqualDelta[int8](t, 127, -128, 1, "ctx: ") }, "ctx: Should be equal in delta\n  actual: 127\nexpected: -128\n   delta: 1\n    diff: 255"},
			{"EqualDelta invalid", func(t Testing) bool { return EqualDelta(t, 1, 2, -1, "ctx: ") }, "ctx: Should have non-negative delta\n   delta: -1"},
			{"NotEqualDelta invalid", func(t Testing) bool { return NotEqualDelta(t, 1, 2, -1, "ctx: ") }, "ctx: Should have non-negative delta\n   delta: -1"},
		})
	})
	t.Run("ordering", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Greater", func(t Testing) bool { return Greater(t, 1, 2, "ctx: ") }, "ctx: Should be greater\n  actual: 1\n   bound: 2"},
			{"Greater duration", func(t Testing) bool { return Greater(t, time.Second, time.Minute, "ctx: ") }, "ctx: Should be greater\n  actual: 1s\n   bound: 1m0s"},
			{"GreaterOrEqual", func(t Testing) bool { return GreaterOrEqual(t, 1, 2, "ctx: ") }, "ctx: Should be greater or equal\n  actual: 1\n   bound: 2"},
			{"Less", func(t Testing) bool { return Less(t, 2, 1, "ctx: ") }, "ctx: Should be less\n  actual: 2\n   bound: 1"},
			{"LessOrEqual", func(t Testing) bool { return LessOrEqual(t, 2, 1, "ctx: ") }, "ctx: Should be less or equal\n  actual: 2\n   bound: 1"},
		})
	})
	t.Run("length", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Length", func(t Testing) bool { return Length(t, []int{1}, 2, "ctx: ") }, "ctx: Should have length\n  actual: ["},
			{"Length array", func(t Testing) bool { return Length(t, [1]int{1}, 2, "ctx: ") }, "ctx: Should have length\n  actual: [1]int{1}\n  length: 1\nexpected: 2"},
			{"Length invalid", func(t Testing) bool { return Length(t, 5, 2, "ctx: ") }, "ctx: Should be iterable\n  actual: 5"},
			{"Empty", func(t Testing) bool { return Empty(t, []int{1}, "ctx: ") }, "ctx: Should be empty\n  actual: ["},
			{"Empty invalid", func(t Testing) bool { return Empty(t, 5, "ctx: ") }, "ctx: Should be iterable\n  actual: 5"},
			{"NotEmpty", func(t Testing) bool { return NotEmpty(t, []int{}, "ctx: ") }, "ctx: Should not be empty\n  actual: ["},
			{"NotEmpty invalid", func(t Testing) bool { return NotEmpty(t, 5, "ctx: ") }, "ctx: Should be iterable\n  actual: 5"},
		})
	})
	t.Run("contents", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Contains", func(t Testing) bool { return Contains(t, []int{1}, 2, "ctx: ") }, "ctx: Should contain element\n  actual: ["},
			{"Contains invalid", func(t Testing) bool { return Contains(t, 5, 2, "ctx: ") }, "ctx: Should be iterable\n  actual: 5\n element: 2"},
			{"Contains element type", func(t Testing) bool { return Contains(t, []int{1}, "a", "ctx: ") }, "ctx: Should have element of same type\n  actual: ["},
			{"NotContains", func(t Testing) bool { return NotContains(t, []int{1}, 1, "ctx: ") }, "ctx: Should not contain element\n  actual: ["},
			{"NotContains invalid", func(t Testing) bool { return NotContains(t, 5, 2, "ctx: ") }, "ctx: Should be iterable\n  actual: 5\n element: 2"},
			{"EqualUnordered", func(t Testing) bool { return EqualUnordered(t, []int{1, 2, 2}, []int{3, 2, 1}, "ctx: ") }, "ctx: Should be equal in any order\n  actual: ["},
			{"EqualUnordered difference", func(t Testing) bool { return EqualUnordered(t, [3]int{1, 2, 2}, [3]int{3, 2, 1}, "ctx: ") }, "ctx: Should be equal in any order\n  actual: [3]int{1, 2, 2}\nexpected: [3]int{3, 2, 1}\n   extra: []int{2}\n missing: []int{3}"},
			{"EqualUnordered durations", func(t Testing) bool { return EqualUnordered(t, [1]time.Duration{1e9}, [1]time.Duration{6e10}, "ctx: ") }, "ctx: Should be equal in any order\n  actual: [1]time.Duration{1000000000}\nexpected: [1]time.Duration{60000000000}\n   extra: []time.Duration{1000000000}\n missing: []time.Duration{60000000000}"},
			{"EqualUnordered invalid", func(t Testing) bool { return EqualUnordered(t, "a", "a", "ctx: ") }, "ctx: Should be array or slice\n  actual: \"a\"\nexpected: \"a\""},
			{"EqualUnordered type", func(t Testing) bool { return EqualUnordered[any](t, []int{}, []uint{}, "ctx: ") }, "ctx: Should have same type\n  actual: ["},
			{"HasPrefix", func(t Testing) bool { return HasPrefix(t, "ab", "b", "ctx: ") }, "ctx: Should have prefix\n  actual: \"ab\"\n  prefix: \"b\""},
			{"HasPrefix invalid", func(t Testing) bool { return HasPrefix(t, 1, 1, "ctx: ") }, "ctx: Should be string or slice\n  actual: 1\n  prefix: 1"},
			{"HasPrefix type", func(t Testing) bool { return HasPrefix[any](t, "a", testType("a"), "ctx: ") }, "ctx: Should have same type\n  actual: string(\"a\")\n  prefix: assert.testType(\"a\")"},
			{"HasSuffix", func(t Testing) bool { return HasSuffix(t, "ab", "a", "ctx: ") }, "ctx: Should have suffix\n  actual: \"ab\"\n  suffix: \"a\""},
			{"HasSuffix type", func(t Testing) bool { return HasSuffix[any](t, []int{1}, []uint{1}, "ctx: ") }, "ctx: Should have same type\n  actual: []int{1}\n  suffix: []uint{0x1}"},
			{"HasSuffix invalid", func(t Testing) bool { return HasSuffix(t, 1, 1, "ctx: ") }, "ctx: Should be string or slice\n  actual: 1\n  suffix: 1"},
		})
	})
	t.Run("errors", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Error", func(t Testing) bool { return Error(t, nil, "ctx: ") }, "ctx: Should be error"},
			{"NoError", func(t Testing) bool { return NoError(t, err, "ctx: ") }, "ctx: Should not be error\n     msg: oops\n   error: &errors.errorString{s:\"oops\"}"},
			{"NoError multi-line", func(t Testing) bool { return NoError(t, errors.New("a\nb"), "ctx: ") }, "ctx: Should not be error\n     msg: a\n          b\n   error: &errors.errorString{s:\"a\\nb\"}"},
			{"NoError typed nil", func(t Testing) bool { return NoError(t, ptrErr, "ctx: ") }, "ctx: Should not be error\n     msg: <nil>\n   error: (*assert.testPtrErr)(nil)"},
			{"ErrorIs", func(t Testing) bool { return ErrorIs(t, err, target, "ctx: ") }, "ctx: Should match error\n     msg: oops\n   error: &errors.errorString{s:\"oops\"}\n  target: &errors.errorString{s:\"target\"}"},
			{"NotErrorIs", func(t Testing) bool { return NotErrorIs(t, err, err, "ctx: ") }, "ctx: Should not match error\n"},
			{"ErrorAs", func(t Testing) bool { return ErrorAs(t, err, &ptrErr, "ctx: ") }, "ctx: Should be assignable to target\n     msg: oops\n   error: &errors.errorString{s:\"oops\"}\n  target: **assert.testPtrErr"},
			{"ErrorAs invalid", func(t Testing) bool { return ErrorAs(t, err, ptrErr, "ctx: ") }, "ctx: Should have pointer to error or interface target\n  target: *assert.testPtrErr"},
			{"ErrorContains", func(t Testing) bool { return ErrorContains(t, err, "x", "ctx: ") }, "ctx: Should contain substring\n     msg: oops\n   error: &errors.errorString{s:\"oops\"}\n  substr: \"x\""},
			{"ErrorContains nil", func(t Testing) bool { return ErrorContains(t, nil, "x", "ctx: ") }, "ctx: Should be error\n  substr: \"x\""},
		})
	})
	t.Run("regexp", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Matches", func(t Testing) bool { return Matches(t, "a", "b", "ctx: ") }, "ctx: Should match regexp\n  actual: \"a\"\n pattern: b"},
			{"Matches invalid", func(t Testing) bool { return Matches(t, "a", "[", "ctx: ") }, "ctx: Should be valid regexp\n pattern: [\n     err: "},
			{"NotMatches", func(t Testing) bool { return NotMatches(t, "a", "a", "ctx: ") }, "ctx: Should not match regexp\n  actual: \"a\"\n pattern: a"},
			{"NotMatches invalid", func(t Testing) bool { return NotMatches(t, "a", "[", "ctx: ") }, "ctx: Should be valid regexp\n pattern: [\n     err: "},
		})
	})
	t.Run("JSON", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"EqualJSON", func(t Testing) bool { return EqualJSON(t, "1", "2", "ctx: ") }, "ctx: Should be equal JSON\n  actual: 1\nexpected: 2"},
			{"EqualJSON multi-line", func(t Testing) bool { return EqualJSON(t, "[\n1\n]", "2", "ctx: ") }, "ctx: Should be equal JSON\n  actual: [\n          1\n          ]\nexpected: 2"},
			{"EqualJSON duplicate", func(t Testing) bool { return EqualJSON(t, `{"a":1,"a":1}`, `{"a":1}`, "ctx: ") }, "ctx: Should be valid JSON\n  actual: {\"a\":1,\"a\":1}\n     err: duplicate key \"a\""},
			{"EqualJSON trailing", func(t Testing) bool { return EqualJSON(t, "1 x", "1", "ctx: ") }, "ctx: Should be valid JSON\n  actual: 1 x\n     err: unexpected \"x\" after top-level value"},
			{"EqualJSON invalid actual", func(t Testing) bool { return EqualJSON(t, "x", "2", "ctx: ") }, "ctx: Should be valid JSON\n  actual: x\n     err: "},
			{"EqualJSON invalid expected", func(t Testing) bool { return EqualJSON(t, "1", "x", "ctx: ") }, "ctx: Should be valid JSON\nexpected: x\n     err: "},
			{"JSON", func(t Testing) bool { return JSON(t, 1, "2", "ctx: ") }, "ctx: Should be equal JSON\n  actual: 1\nexpected: 2"},
			{"JSON unmarshalable", func(t Testing) bool { return JSON(t, make(chan int), "2", "ctx: ") }, "ctx: Should be marshalable\n  actual: (chan int)(0x"},
		})
	})
	t.Run("panics", func(t *testing.T) {
		testMessages(t, []messageCase{
			{"Panics", func(t Testing) bool { return Panics(t, func() {}, "ctx: ") }, "ctx: Should panic"},
			{"Panics nil", func(t Testing) bool { return Panics(t, nil, "ctx: ") }, "ctx: Should be non-nil function"},
			{"PanicsWith nil", func(t Testing) bool { return PanicsWith(t, nil, 1, "ctx: ") }, "ctx: Should be non-nil function"},
			{"NotPanics nil function", func(t Testing) bool { return NotPanics(t, nil, "ctx: ") }, "ctx: Should be non-nil function"},
			{"PanicsWith types", func(t Testing) bool { return PanicsWith(t, func() { panic(1) }, int64(1), "ctx: ") }, "ctx: Should panic with value\n  actual: int(1)\nexpected: int64(1)"},
			{"PanicsWith", func(t Testing) bool { return PanicsWith(t, func() {}, "b", "ctx: ") }, "ctx: Should panic\nexpected: \"b\""},
			{"PanicsWith value", func(t Testing) bool { return PanicsWith(t, func() { panic("a") }, "b", "ctx: ") }, "ctx: Should panic with value\n  actual: \"a\"\nexpected: \"b\""},
			{"PanicsWith error", func(t Testing) bool { return PanicsWith(t, func() { panic("a") }, err, "ctx: ") }, "ctx: Should panic with value\n  actual: \"a\"\nexpected: &errors.errorString{s:\"oops\"}"},
			{"PanicsWith wrong error", func(t Testing) bool { return PanicsWith(t, func() { panic(err) }, target, "ctx: ") }, "ctx: Should panic with value\n  actual: &errors.errorString{s:\"oops\"}\nexpected: &errors.errorString{s:\"target\"}"},
			{"NotPanics nil", func(t Testing) bool { return NotPanics(t, func() { panic(nil) }, "ctx: ") }, "ctx: Should not panic\n   value: <nil>"},
			{"NotPanics", func(t Testing) bool { return NotPanics(t, func() { panic("a") }, "ctx: ") }, "ctx: Should not panic\n   value: \"a\""},
		})
	})
}

type messageCase struct {
	name     string
	assert   func(t Testing) bool
	expected string
}

type testType string

type testBytes []byte

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

func testEqualInvalid[T Comparable](t *testing.T, actual, expected T) {
	t.Helper()

	tt := newLogger()
	if Equal(tt, actual, expected) {
		t.Errorf("Equal(%#v,%#v) should return false: %s", actual, expected, tt.LastError)
	}

	tt = newLogger()
	if NotEqual(tt, actual, expected) {
		t.Errorf("NotEqual(%#v,%#v) should return false: %s", actual, expected, tt.LastError)
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

func TestFormat(t *testing.T) {
	t.Run("values", func(t *testing.T) {
		u := uint(5)
		sl := []time.Duration{time.Second}

		cases := []struct {
			name     string
			object   any
			expected string
		}{
			{"nil", nil, "<nil>"},
			{"int", -3, "-3"},
			{"unsigned", uint8(5), "5"},
			{"float", float32(1.5), "1.5"},
			{"duration", 1500 * time.Millisecond, "1.5s"},
			{"uintptr", uintptr(10), "0xa"},
			{"string", testType("a"), `"a"`},
			{"unsigned pointer", &u, "&5"},
			{"nil pointer", (*uint)(nil), "(*uint)(nil)"},
			{"pointer to slice", &sl, "&[]time.Duration{1000000000}"},
			{"pointer to struct", &testStruct{1, "a"}, `&assert.testStruct{a:1, b:"a"}`},
			{"bytes", []byte("a\xff"), `[]byte("a\xff")`},
			{"named bytes", testBytes("{}"), `assert.testBytes("{}")`},
			{"nil bytes", []byte(nil), "[]byte(nil)"},
			{"slice", []int{1}, "[]int{1}"},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				if actual := format(c.object); actual != c.expected {
					t.Errorf("format(%#v) should return %q, got %q", c.object, c.expected, actual)
				}
			})
		}
	})
	t.Run("references", func(t *testing.T) {
		x := 1
		s := []int{1}

		cases := []struct {
			name     string
			object   any
			expected string
		}{
			{"pointer", &x, fmt.Sprintf("[%p] &1", &x)},
			{"slice", s, fmt.Sprintf("[%p] []int{1}", s)},
			{"nil slice", []int(nil), "[0x0] []int(nil)"},
			{"value", 1, "1"},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				if actual := formatReference(c.object); actual != c.expected {
					t.Errorf("formatReference(%#v) should return %q, got %q", c.object, c.expected, actual)
				}
			})
		}
	})
	t.Run("window", func(t *testing.T) {
		long := "a" + strings.Repeat("é", formatLimit)
		cut := window(long, 0)

		if !strings.HasSuffix(cut, ellipsis) || strings.HasPrefix(cut, ellipsis) {
			t.Errorf("long value should be cut at the end, got %q", cut)
		}

		if len(cut) > formatLimit || !utf8.ValidString(cut) {
			t.Errorf("cut value should fit the limit and stay valid UTF-8, got %d bytes", len(cut))
		}

		if window(cut, 0) != cut {
			t.Errorf("cut value should not be cut again")
		}

		if window("short", 3) != "short" {
			t.Errorf("short value should not be cut")
		}

		middle := window(long, len(long)/2)
		if !strings.HasPrefix(middle, ellipsis) || !strings.HasSuffix(middle, ellipsis) || !utf8.ValidString(middle) {
			t.Errorf("value should be cut around the position, got %q", middle)
		}

		end := window(long, len(long))
		if !strings.HasPrefix(end, ellipsis) || strings.HasSuffix(end, ellipsis) || len(end) < formatLimit-4 {
			t.Errorf("value should be cut to its full tail, got %d bytes", len(end))
		}
	})
	t.Run("pairs", func(t *testing.T) {
		common := strings.Repeat("a", 2*formatLimit)

		actual, expected := formatPair(common+"x"+common, common+"y"+common)
		if !strings.Contains(actual, "ax") || !strings.Contains(expected, "ay") {
			t.Errorf("pair should show the first difference, got %q", actual)
		}

		actual, expected = formatSuffixPair(common+"x"+common, "y"+common)
		if !strings.Contains(actual, "xa") || !strings.Contains(expected, "ya") {
			t.Errorf("suffix pair should show the last difference, got %q", actual)
		}

		actual, expected = formatPair(1, int64(1))
		if actual != "int(1)" || expected != "int64(1)" {
			t.Errorf("pair should show types of equally printed values, got %q and %q", actual, expected)
		}

		actual, expected = formatPair(math.NaN(), math.NaN())
		if actual != "NaN" || expected != "NaN" {
			t.Errorf("pair should not show types of the same type, got %q and %q", actual, expected)
		}
	})
}

func TestMessageTruncate(t *testing.T) {
	long := strings.Repeat("a", 2*formatLimit)
	limit := 3*formatLimit + 100

	cases := []struct {
		name   string
		assert func(t Testing) bool
	}{
		{"EqualJSON", func(t Testing) bool { return EqualJSON(t, `"`+long+`"`, `"`+long+`b"`) }},
		{"EqualJSON invalid", func(t Testing) bool { return EqualJSON(t, long, "1") }},
		{"Matches", func(t Testing) bool { return Matches(t, long, long+"b") }},
		{"NoError", func(t Testing) bool { return NoError(t, errors.New(long)) }},
		{"Equal", func(t Testing) bool { return Equal(t, long, long+"b") }},
		{"ErrorIs", func(t Testing) bool { return ErrorIs(t, errors.New(long), errors.New(long)) }},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tt := newLogger()
			c.assert(tt)

			if len(tt.LastError) > limit {
				t.Errorf("message should be truncated, got %d bytes", len(tt.LastError))
			}
		})
	}
}

func TestMessageWindow(t *testing.T) {
	long := strings.Repeat("a", 2*formatLimit)

	cases := []struct {
		name     string
		assert   func(t Testing) bool
		expected []string
	}{
		{"Equal", func(t Testing) bool { return Equal(t, long+"x"+long, long+"y"+long) }, []string{"\n  actual: …a", "ax", "\nexpected: …a", "ay"}},
		{"Equal slices", func(t Testing) bool { return Equal(t, []string{long, "x"}, []string{long, "y"}) }, []string{`", "x"}`, `", "y"}`}},
		{"HasPrefix", func(t Testing) bool { return HasPrefix(t, long+"x"+long, long+"y") }, []string{"ax", `ay"`}},
		{"HasSuffix", func(t Testing) bool { return HasSuffix(t, long+"x"+long, "y"+long) }, []string{"\n  actual: …a", "xa", "\n  suffix: \"ya"}},
		{"EqualJSON", func(t Testing) bool { return EqualJSON(t, `"`+long+`x"`, `"`+long+`y"`) }, []string{`ax"`, `ay"`}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tt := newLogger()
			c.assert(tt)

			for _, expected := range c.expected {
				if !strings.Contains(tt.LastError, expected) {
					t.Errorf("message should contain %q, got %q", expected, tt.LastError)
				}
			}
		})
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
