package testing_test

import (
	"errors"
	"testing"

	. "github.com/rocketlaunchr/testing-go"
)

func ret1Val(Panic bool) string {
	if Panic {
		panic("panic")
	}
	return "abc"
}

func TestRet1(t *testing.T) {
	testCases := []struct {
		in     bool
		ExpOut interface{}
	}{
		{true, PanicExpected},
		{false, NotEqual{PanicExpected}},
		{false, "abc"},
		{false, NotEqual{"abcd"}},
	}

	tcfg := NewTestConfig(t)

	for idx, tc := range testCases {
		tcfg.Run1(Sprintf("[%d]: %v", idx, tc.in), tc, func(t *testing.T) interface{} {
			return ret1Val(tc.in)
		})
	}
}

var (
	errSample1 = errors.New("sample error 1")
	errSample2 = errors.New("sample error 2")
)

func retErr(Panic bool) error {
	if Panic {
		panic("panic")
	}
	return errSample1
}

func TestRetErr(t *testing.T) {
	testCases := []struct {
		in     bool
		ExpErr error
	}{
		{true, PanicExpected},
		{false, NotEqual{PanicExpected}},
		{false, errSample1},
		{false, ErrAny},
		{false, NotEqual{errSample2}},
		{false, ErrContains{"1"}},
		{false, NotEqual{ErrContains{"2"}}},
		{false, NotEqual{nil}},

		{true, Is{PanicExpected}},
		{false, Is{NotEqual{PanicExpected}}},
		{false, Is{errSample1}},
		{false, Is{ErrAny}},
		{false, Is{NotEqual{errSample2}}},
		{false, Is{ErrContains{"1"}}},
		{false, Is{NotEqual{ErrContains{"2"}}}},
		{false, Is{NotEqual{nil}}},
	}

	tcfg := NewTestConfig(t)

	for idx, tc := range testCases {
		tcfg.RunErr(Sprintf("[%d]: %v", idx, tc.in), tc, func(t *testing.T) error {
			return retErr(tc.in)
		})
	}
}

func retErrNil(Panic bool) error {
	if Panic {
		panic("panic")
	}
	return nil
}

func TestRetErrNil(t *testing.T) {
	testCases := []struct {
		in     bool
		ExpErr error
	}{
		{true, PanicExpected},
		{false, nil},
		{false, NotEqual{NotEqual{nil}}},
		{false, NotEqual{PanicExpected}},
		{false, NotEqual{ErrAny}},
		{false, NotEqual{errSample1}},
		{false, NotEqual{ErrContains{"2"}}},

		{true, Is{PanicExpected}},
		{false, Is{nil}},
		{false, Is{NotEqual{PanicExpected}}},
		{false, Is{NotEqual{ErrAny}}},
		{false, Is{NotEqual{errSample1}}},
		{false, Is{NotEqual{ErrContains{"2"}}}},
	}

	tcfg := NewTestConfig(t)

	for idx, tc := range testCases {
		tcfg.RunErr(Sprintf("[%d]: %v", idx, tc.in), tc, func(t *testing.T) error {
			return retErrNil(tc.in)
		})
	}
}
