package testing

import (
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"strings"
)

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

type ErrContains struct {
	Str string
}

func (e ErrContains) Error() string { return "error contains: " + e.Str }

func (e ErrContains) Is(target error) bool {
	s := spew.NewDefaultConfig()
	s.DisableMethods = false
	s.DisablePointerMethods = false

	fmt.Printf("target: %s", s.Sdump(target))
	fmt.Printf("e: %s", s.Sdump(e))
	fmt.Println(target.Error(), e.Str)
	fmt.Printf("return: %v\n", strings.Contains(target.Error(), e.Str))
	fmt.Println("-----------")
	return strings.Contains(target.Error(), e.Str)
}

// func Contains(s, substr string) bool

// https://golang.org/src/errors/wrap.go?s=1170:1201#L29
