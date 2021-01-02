package testing

import "errors"

var (
	ErrPanic   = errors.New("panic expected")
	ErrAny     = errors.New("any error expected")
	CustomTest = errors.New("custom test")
)

type NotEqual struct {
	Val interface{}
}

func (NotEqual) Error() string { return "not equal" }
func (e NotEqual) Unwrap() error {
	x, _ := e.Val.(error)
	return x
}
func (e NotEqual) Is(target error) bool {
	_, ok := target.(NotEqual)
	return ok
}

type Is struct {
	Err error
}

func (Is) Error() string   { return "is" }
func (e Is) Unwrap() error { return e.Err }
func (e Is) Is(target error) bool {
	_, ok := target.(Is)
	return ok
}
