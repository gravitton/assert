package assert

type Testing interface {
	Errorf(format string, args ...any)
}

type Helper interface {
	Helper()
}
