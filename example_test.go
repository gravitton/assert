package assert_test

import (
	"errors"
	"fmt"
	"io"
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

func ExampleTrue() {
	assert.True(t, 1 > 2)

	// Output:
	// Should be true
}

func ExampleNil() {
	x := 1
	assert.Nil(t, &x)

	// Output:
	// Should be nil
	//   actual: &1
}

func ExampleZero() {
	assert.Zero(t, 5*time.Second)

	// Output:
	// Should be zero
	//   actual: 5s
}

func ExampleSame() {
	assert.Same[any](t, new(int), new(int64))

	// Output:
	// Should have same type
	//   actual: *int
	// expected: *int64
}

func ExampleEqual() {
	assert.Equal(t, []string{"a", "b"}, []string{"a", "c"})

	// Output:
	// Should be equal
	//   actual: []string{"a", "b"}
	// expected: []string{"a", "c"}
}

func ExampleEqual_types() {
	assert.Equal[any](t, 1, int64(1))

	// Output:
	// Should be equal
	//   actual: int(1)
	// expected: int64(1)
}

func ExampleEqual_pointers() {
	a, b := 1, 2
	assert.Equal(t, &a, &b)

	// Output:
	// Should be equal
	//   actual: &1
	// expected: &2
}

func ExampleEqual_messages() {
	assert.Equal(t, "Ann", "Bob", "user.Name: ")

	// Output:
	// user.Name: Should be equal
	//   actual: "Ann"
	// expected: "Bob"
}

func ExampleEqualDelta() {
	assert.EqualDelta(t, time.Second, 2*time.Second, 100*time.Millisecond)

	// Output:
	// Should be equal in delta
	//   actual: 1s
	// expected: 2s
	//    delta: 100ms
	//     diff: 1s
}

func ExampleGreater() {
	assert.Greater(t, uint(1), 2)

	// Output:
	// Should be greater
	//   actual: 1
	//    bound: 2
}

func ExampleLength() {
	assert.Length(t, []string{"a"}, 2)

	// Output:
	// Should have length
	//   actual: []string{"a"}
	//   length: 1
	// expected: 2
}

func ExampleContains() {
	assert.Contains(t, []int{1, 2}, 3)

	// Output:
	// Should contain element
	//   actual: []int{1, 2}
	//  element: 3
}

func ExampleContains_elementType() {
	assert.Contains(t, []int64{1, 2}, 1)

	// Output:
	// Should have element of same type
	//   actual: []int64{1, 2}
	//  element: 1
}

func ExampleEqualUnordered() {
	assert.EqualUnordered(t, []int{1, 2, 2}, []int{2, 1, 3})

	// Output:
	// Should be equal in any order
	//   actual: []int{1, 2, 2}
	// expected: []int{2, 1, 3}
	//    extra: []int{2}
	//  missing: []int{3}
}

func ExampleHasPrefix() {
	assert.HasPrefix(t, "/api/users", "/v1/")

	// Output:
	// Should have prefix
	//   actual: "/api/users"
	//   prefix: "/v1/"
}

func ExampleNoError() {
	assert.NoError(t, errors.New("boom"))

	// Output:
	// Should not be error
	//    error: &errors.errorString{s:"boom"}
}

func ExampleErrorIs() {
	assert.ErrorIs(t, errors.New("boom"), io.EOF)

	// Output:
	// Should match error
	//    error: &errors.errorString{s:"boom"}
	//   target: &errors.errorString{s:"EOF"}
}

func ExampleErrorContains() {
	assert.ErrorContains(t, errors.New("boom"), "bang")

	// Output:
	// Should contain substring
	//    error: &errors.errorString{s:"boom"}
	//   substr: "bang"
}

func ExampleMatches() {
	assert.Matches(t, "abc", `^\d+$`)

	// Output:
	// Should match regexp
	//   actual: "abc"
	//  pattern: "^\\d+$"
}

func ExampleEqualJSON() {
	assert.EqualJSON(t, `{"id":1}`, `{"id":2}`)

	// Output:
	// Should be equal JSON
	//   actual: "{\"id\":1}"
	// expected: "{\"id\":2}"
}

func ExamplePanics() {
	assert.Panics(t, func() {
	})

	// Output:
	// Should panic
}

func ExamplePanicsWith() {
	assert.PanicsWith(t, func() {
		panic("boom")
	}, "bang")

	// Output:
	// Should panic with value
	//   actual: "boom"
	// expected: "bang"
}

func ExampleNotPanics() {
	assert.NotPanics(t, func() {
		panic("boom")
	})

	// Output:
	// Should not panic
	//    value: "boom"
}
