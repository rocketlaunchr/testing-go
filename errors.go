package testing

import (
	"errors"
	"strings"
)

var (
	// PanicExpected indicates that the function being tested is expected to panic.
	PanicExpected = errors.New("panic expected")

	// ErrAny indicates that the function being tested is expected to return an (unspecified) error.
	ErrAny = errors.New("any error expected")

	// CustomTest indicates that you want to perform your own test inside the Run1, RunErr or Run2 function.
	// You will be expected to signify when a test failed yourself i.e. using t.Errorf(...).
	CustomTest = errors.New("custom test")
)

// Not means the expected value/error is expected to be not equal.
//
// Example:
//
//  testCases := []struct {
//      in     bool
//      ExpErr error
//  }{
//      {false, Not{ErrContains{"database error"}}},
//  }
type Not struct{ Val interface{} }

func (Not) Error() string { return "not equal" }
func (e Not) Unwrap() error {
	x, _ := e.Val.(error)
	return x
}
func (e Not) Is(target error) bool {
	_, ok := target.(Not)
	return ok
}

// Is indicates that errors.Is() be used to test if an expected error matches the observed error.
// When not set, a simple equality check is performed using reflect.DeepEqual.
//
// See: https://pkg.go.dev/errors#Is
type Is struct{ Err error }

func (Is) Error() string   { return "is" }
func (e Is) Unwrap() error { return e.Err }
func (e Is) Is(target error) bool {
	_, ok := target.(Is)
	return ok
}

// ErrContains is used to check if the observed error contains a particular substring.
// It uses strings.Contains().
//
// Example:
//
//  testCases := []struct {
//      in     bool
//      ExpErr error
//  }{
//      {false, ErrContains{"database error"}},
//  }
type ErrContains struct{ Substr string }

func (e ErrContains) Error() string { return "error contains: " + e.Substr }
func (e ErrContains) Is(target error) bool {
	return strings.Contains(target.Error(), e.Substr)
}
