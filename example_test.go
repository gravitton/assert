package assert_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"math"
	"time"

	"github.com/gravitton/assert"
)

type stdout struct{}

func (stdout) Helper() {
}

func (stdout) Errorf(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}

var t stdout

func Example_messages() {
	assert.Equal(t, "Ann", "Bob", "user.Name: ")

	// Output:
	// user.Name: Should be equal
	//   actual: "Ann"
	// expected: "Bob"
}

func ExampleFail() {
	assert.Fail(t, "custom failure")

	// Output:
	// custom failure
}

func ExampleFailf() {
	assert.Failf(t, "x=%d", 1)

	// Output:
	// x=1
}

func ExampleTrue() {
	assert.True(t, 1 > 2)

	// Output:
	// Should be true
}

func ExampleFalse() {
	assert.False(t, 2 > 1)

	// Output:
	// Should be false
}

func ExampleNil() {
	x := 1
	assert.Nil(t, &x)

	// Output:
	// Should be nil
	//   actual: &1
}

func ExampleNotNil() {
	assert.NotNil(t, (*int)(nil))

	// Output:
	// Should not be nil
	//   actual: (*int)(nil)
}

func ExampleZero() {
	assert.Zero(t, 5*time.Second)

	// Output:
	// Should be zero
	//   actual: 5s
}

func ExampleNotZero() {
	assert.NotZero(t, "")

	// Output:
	// Should not be zero
	//   actual: ""
}

func ExampleSame() {
	assert.Same[any](t, new(int), new(int64))
	assert.Same(t, 1, 1)

	// Output:
	// Should have same type
	//   actual: *int(&0)
	// expected: *int64(&0)
	// Should be reference
	//   actual: int(1)
	// expected: int(1)
}

func ExampleNotSame() {
	assert.NotSame(t, 1, 1)

	// Output:
	// Should be reference
	//   actual: int(1)
	// expected: int(1)
}

func ExampleEqual() {
	x := 1
	assert.Equal(t, []string{"a", "b"}, []string{"a", "c"})
	assert.Equal(t, (*int)(nil), &x)
	assert.Equal[any](t, 1, int64(1))
	assert.Equal(t, uint(5), uint(6))
	assert.Equal(t, math.NaN(), math.NaN())

	// Output:
	// Should be equal
	//   actual: []string{"a", "b"}
	// expected: []string{"a", "c"}
	// Should be equal
	//   actual: (*int)(nil)
	// expected: &1
	// Should be equal
	//   actual: int(1)
	// expected: int64(1)
	// Should be equal
	//   actual: 5
	// expected: 6
	// Should be equal
	//   actual: NaN
	// expected: NaN
}

func ExampleNotEqual() {
	assert.NotEqual(t, 1, 1)

	// Output:
	// Should not be equal
	//   actual: 1
}

func ExampleEqualDelta() {
	assert.EqualDelta(t, 1, 3, 1)
	assert.EqualDelta(t, 1.5, 1.0, 0.25)
	assert.EqualDelta(t, time.Second, 2*time.Second, 100*time.Millisecond)
	assert.EqualDelta[int8](t, 127, -128, 1)
	assert.EqualDelta(t, 1, 2, -1)

	// Output:
	// Should be equal in delta
	//   actual: 1
	// expected: 3
	//    delta: 1
	//     diff: 2
	// Should be equal in delta
	//   actual: 1.5
	// expected: 1
	//    delta: 0.25
	//     diff: 0.5
	// Should be equal in delta
	//   actual: 1s
	// expected: 2s
	//    delta: 100ms
	//     diff: 1s
	// Should be equal in delta
	//   actual: 127
	// expected: -128
	//    delta: 1
	//     diff: 255
	// Should have non-negative delta
	//    delta: -1
}

func ExampleNotEqualDelta() {
	assert.NotEqualDelta(t, 1, 2, 1)

	// Output:
	// Should not be equal in delta
	//   actual: 1
	// expected: 2
	//    delta: 1
	//     diff: 1
}

func ExampleGreater() {
	assert.Greater(t, 1, 2)
	assert.Greater(t, time.Second, time.Minute)
	assert.Greater(t, math.NaN(), 1)

	// Output:
	// Should be greater
	//   actual: 1
	//    bound: 2
	// Should be greater
	//   actual: 1s
	//    bound: 1m0s
	// Should not be NaN
	//   actual: NaN
	//    bound: 1
}

func ExampleGreaterOrEqual() {
	assert.GreaterOrEqual(t, 1, 2)

	// Output:
	// Should be greater or equal
	//   actual: 1
	//    bound: 2
}

func ExampleLess() {
	assert.Less(t, "b", "a")

	// Output:
	// Should be less
	//   actual: "b"
	//    bound: "a"
}

func ExampleLessOrEqual() {
	assert.LessOrEqual(t, 2, 1)

	// Output:
	// Should be less or equal
	//   actual: 2
	//    bound: 1
}

func ExampleLength() {
	assert.Length(t, []string{"a"}, 2)
	assert.Length(t, 5, 2)

	// Output:
	// Should have length
	//   actual: []string{"a"}
	//   length: 1
	// expected: 2
	// Should be iterable
	//   actual: 5
}

func ExampleEmpty() {
	assert.Empty(t, []int{1})

	// Output:
	// Should be empty
	//   actual: []int{1}
}

func ExampleNotEmpty() {
	assert.NotEmpty(t, "")

	// Output:
	// Should not be empty
	//   actual: ""
}

func ExampleContains() {
	assert.Contains(t, []int{1, 2}, 3)
	assert.Contains(t, "abc", "x")
	assert.Contains(t, []int64{1, 2}, 1)
	assert.Contains(t, 5, 2)

	// Output:
	// Should contain element
	//   actual: []int{1, 2}
	//  element: 3
	// Should contain element
	//   actual: "abc"
	//  element: "x"
	// Should have element of same type
	//   actual: []int64{1, 2}
	//  element: int(1)
	// Should be iterable
	//   actual: 5
	//  element: 2
}

func ExampleNotContains() {
	assert.NotContains(t, map[string]int{"a": 1}, 1)

	// Output:
	// Should not contain element
	//   actual: map[string]int{"a":1}
	//  element: 1
}

func ExampleEqualUnordered() {
	assert.EqualUnordered(t, []int{1, 2, 2}, []int{2, 1, 3})
	assert.EqualUnordered(t, [1]time.Duration{time.Second}, [1]time.Duration{time.Minute})
	assert.EqualUnordered[any](t, []int{}, []uint{})
	assert.EqualUnordered(t, "a", "a")

	// Output:
	// Should be equal in any order
	//   actual: []int{1, 2, 2}
	// expected: []int{2, 1, 3}
	//    extra: []int{2}
	//  missing: []int{3}
	// Should be equal in any order
	//   actual: [1]time.Duration{1000000000}
	// expected: [1]time.Duration{60000000000}
	//    extra: []time.Duration{1000000000}
	//  missing: []time.Duration{60000000000}
	// Should have same type
	//   actual: []int{}
	// expected: []uint{}
	// Should be array or slice
	//   actual: "a"
	// expected: "a"
}

func ExampleHasPrefix() {
	assert.HasPrefix(t, "/api/users", "/v1/")
	assert.HasPrefix[any](t, "1", json.Number("1"))
	assert.HasPrefix(t, 1, 1)

	// Output:
	// Should have prefix
	//   actual: "/api/users"
	//   prefix: "/v1/"
	// Should have same type
	//   actual: string("1")
	//   prefix: json.Number("1")
	// Should be string or slice
	//   actual: 1
	//   prefix: 1
}

func ExampleHasSuffix() {
	assert.HasSuffix(t, []string{"a", "b"}, []string{"a"})

	// Output:
	// Should have suffix
	//   actual: []string{"a", "b"}
	//   suffix: []string{"a"}
}

func ExampleError() {
	assert.Error(t, nil)

	// Output:
	// Should be error
}

func ExampleNoError() {
	var pathErr *fs.PathError
	assert.NoError(t, errors.New("boom"))
	assert.NoError(t, errors.New("line 1\nline 2"))
	assert.NoError(t, pathErr)

	// Output:
	// Should not be error
	//      msg: boom
	//    error: &errors.errorString{s:"boom"}
	// Should not be error
	//      msg: line 1
	//           line 2
	//    error: &errors.errorString{s:"line 1\nline 2"}
	// Should not be error
	//      msg: <nil>
	//    error: (*fs.PathError)(nil)
}

func ExampleErrorIs() {
	assert.ErrorIs(t, errors.New("boom"), io.EOF)

	// Output:
	// Should match error
	//      msg: boom
	//    error: &errors.errorString{s:"boom"}
	//   target: &errors.errorString{s:"EOF"}
}

func ExampleNotErrorIs() {
	assert.NotErrorIs(t, io.EOF, io.EOF)

	// Output:
	// Should not match error
	//      msg: EOF
	//    error: &errors.errorString{s:"EOF"}
	//   target: &errors.errorString{s:"EOF"}
}

func ExampleErrorAs() {
	var pathErr *fs.PathError
	assert.ErrorAs(t, errors.New("boom"), &pathErr)
	assert.ErrorAs(t, errors.New("boom"), pathErr)

	// Output:
	// Should be assignable to target
	//      msg: boom
	//    error: &errors.errorString{s:"boom"}
	//   target: **fs.PathError
	// Should have pointer to error or interface target
	//   target: *fs.PathError
}

func ExampleErrorContains() {
	assert.ErrorContains(t, errors.New("boom"), "bang")
	assert.ErrorContains(t, nil, "bang")

	// Output:
	// Should contain substring
	//      msg: boom
	//    error: &errors.errorString{s:"boom"}
	//   substr: "bang"
	// Should be error
	//   substr: "bang"
}

func ExampleMatches() {
	assert.Matches(t, "abc", `^\d+$`)
	assert.Matches(t, "abc", `[`)

	// Output:
	// Should match regexp
	//   actual: "abc"
	//  pattern: ^\d+$
	// Should be valid regexp
	//  pattern: [
	//      err: error parsing regexp: missing closing ]: `[`
}

func ExampleNotMatches() {
	assert.NotMatches(t, "123", `^\d+$`)

	// Output:
	// Should not match regexp
	//   actual: "123"
	//  pattern: ^\d+$
}

func ExampleEqualJSON() {
	assert.EqualJSON(t, `{"id":1}`, `{"id":2}`)
	assert.EqualJSON(t, "[\n  1\n]", `[2]`)
	assert.EqualJSON(t, `1 x`, `1`)
	assert.EqualJSON(t, `1`, `x`)

	// Output:
	// Should be equal JSON
	//   actual: {"id":1}
	// expected: {"id":2}
	// Should be equal JSON
	//   actual: [
	//             1
	//           ]
	// expected: [2]
	// Should be valid JSON
	//   actual: 1 x
	//      err: invalid character 'x' after top-level value
	// Should be valid JSON
	// expected: x
	//      err: invalid character 'x' looking for beginning of value
}

func ExampleJSON() {
	assert.JSON(t, map[string]int{"id": 1}, `{"id":2}`)

	// Output:
	// Should be equal JSON
	//   actual: {"id":1}
	// expected: {"id":2}
}

func ExamplePanics() {
	assert.Panics(t, func() {
	})
	assert.Panics(t, nil)

	// Output:
	// Should panic
	// Should be non-nil function
}

func ExamplePanicsWith() {
	assert.PanicsWith(t, func() {
		panic("boom")
	}, "bang")
	assert.PanicsWith(t, func() {
	}, "bang")
	assert.PanicsWith(t, func() {
		panic(1)
	}, int64(1))
	assert.PanicsWith(t, func() {
		panic(io.EOF)
	}, io.ErrUnexpectedEOF)

	// Output:
	// Should panic with value
	//   actual: "boom"
	// expected: "bang"
	// Should panic
	// expected: "bang"
	// Should panic with value
	//   actual: int(1)
	// expected: int64(1)
	// Should panic with value
	//   actual: &errors.errorString{s:"EOF"}
	// expected: &errors.errorString{s:"unexpected EOF"}
}

func ExampleNotPanics() {
	assert.NotPanics(t, func() {
		panic("boom")
	})
	assert.NotPanics(t, func() {
		panic(nil)
	})

	// Output:
	// Should not panic
	//    value: "boom"
	// Should not panic
	//    value: <nil>
}
