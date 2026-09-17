package assert

import "cmp"

// Comparable documents that a value is compared with reflect.DeepEqual.
// It does not restrict the type set; every type satisfies it.
type Comparable interface {
	any
}

// Reference documents that a value is expected to be a pointer, slice, map or channel.
// It does not restrict the type set; every type satisfies it and the kind is checked at runtime.
//
// Functions are deliberately not references: two function values have no comparable
// identity in Go, so [Same] and [NotSame] reject them even though [Nil] accepts them.
type Reference interface {
	any
}

// Ordered is the set of types that support the < operator: integers, floats and strings.
type Ordered interface {
	cmp.Ordered
}

// Numeric is the set of integer and floating-point types.
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

// Iterable documents that a value is expected to be a string, array, slice, map or channel.
// It does not restrict the type set; every type satisfies it and the kind is checked at runtime.
type Iterable interface {
	any
}
