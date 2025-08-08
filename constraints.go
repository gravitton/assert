package assert

type Comparable interface {
	any
}

type Reference interface {
	any
}

type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

type Iterable[T Comparable] interface {
	any
}
