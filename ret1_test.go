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
		{false, Not{PanicExpected}},
		{false, "abc"},
		{false, Not{"abcd"}},
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
		{false, Not{PanicExpected}},
		{false, errSample1},
		{false, ErrAny},
		{false, Not{errSample2}},
		{false, ErrContains{"1"}},
		{false, Not{ErrContains{"2"}}},
		{false, Not{nil}},

		{true, Is{PanicExpected}},
		{false, Is{Not{PanicExpected}}},
		{false, Is{errSample1}},
		{false, Is{ErrAny}},
		{false, Is{Not{errSample2}}},
		{false, Is{ErrContains{"1"}}},
		{false, Is{Not{ErrContains{"2"}}}},
		{false, Is{Not{nil}}},
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
		{false, Not{Not{nil}}},
		{false, Not{PanicExpected}},
		{false, Not{ErrAny}},
		{false, Not{errSample1}},
		{false, Not{ErrContains{"2"}}},

		{true, Is{PanicExpected}},
		{false, Is{nil}},
		{false, Is{Not{PanicExpected}}},
		{false, Is{Not{ErrAny}}},
		{false, Is{Not{errSample1}}},
		{false, Is{Not{ErrContains{"2"}}}},
	}

	tcfg := NewTestConfig(t)

	for idx, tc := range testCases {
		tcfg.RunErr(Sprintf("[%d]: %v", idx, tc.in), tc, func(t *testing.T) error {
			return retErrNil(tc.in)
		})
	}
}
