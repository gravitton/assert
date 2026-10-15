package assert

import "cmp"

// Comparable documents an argument compared with reflect.DeepEqual.
// It does not restrict the type set.
type Comparable interface {
	any
}

// Reference documents an argument that must be a pointer, slice, map or channel.
// It does not restrict the type set; the kind is checked at runtime. A function is not a
// reference, since Go gives function values no identity.
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

// Iterable documents an argument that must be a string, array, array pointer, slice, map or
// channel. It does not restrict the type set; the kind is checked at runtime. An array pointer
// stands for the array it points to.
type Iterable interface {
	any
}
