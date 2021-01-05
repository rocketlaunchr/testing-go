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

// NotEqual means the expected value/error is expected to be not equal.
//
// Example:
//
//  testCases := []struct {
//      in     bool
//      ExpErr error
//  }{
//      {false, NotEqual{ErrContains{"database error"}}},
//  }
type NotEqual struct{ Val interface{} }

// Error ...
func (NotEqual) Error() string { return "not equal" }

// Unwrap ...
func (e NotEqual) Unwrap() error {
	x, _ := e.Val.(error)
	return x
}

// Is ...
func (e NotEqual) Is(target error) bool {
	_, ok := target.(NotEqual)
	return ok
}

// Is indicates that errors.Is() be used to test if an expected error matches the observed error.
// When not set, a simple equality check is performed using reflect.DeepEqual.
//
// See: https://pkg.go.dev/errors#Is
type Is struct{ Err error }

// Error ...
func (Is) Error() string { return "is" }

// Unwrap ...
func (e Is) Unwrap() error { return e.Err }

// Is ...
func (e Is) Is(target error) bool {
	_, ok := target.(Is)
	return ok
}

// ErrContains is used to check if the observed error contains a particular substring.
// It uses strings.Contain.
//
// Example:
//
//  testCases := []struct {
//      in     bool
//      ExpErr error
//  }{
//      {false, ErrContains{"database error"}},
//  }
type ErrContains struct{ Str string }

// Error ...
func (e ErrContains) Error() string { return "error contains: " + e.Str }

// Is ...
func (e ErrContains) Is(target error) bool {
	return strings.Contains(target.Error(), e.Str)
}
