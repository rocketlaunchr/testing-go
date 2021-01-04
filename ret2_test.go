package testing_test

import (
	"errors"
	"math"
	"testing"

	. "github.com/rocketlaunchr/testing-go"
)

var errInvalidInput = errors.New("invalid input")

func sqrt(Panic bool, in float64) (float64, error) {
	if Panic {
		panic("panic")
	}

	if in < 0 {
		return 0, errInvalidInput
	}
	return math.Sqrt(in), nil
}

func TestRet2(t *testing.T) {
	testCases := []struct {
		shouldPanic bool
		in          float64
		ExpOut      interface{}
		ExpErr      error
	}{
		// Test panic
		{true, 0, NotEqual{"abcd"}, PanicExpected},
		{false, 0, 0.0, NotEqual{PanicExpected}},
		{false, 0, NotEqual{1.0}, NotEqual{PanicExpected}},

		{false, -1, NotEqual{"abcd"}, errInvalidInput},
		{false, 9, 3.0, NotEqual{errInvalidInput}},
	}

	tcfg := NewTestConfig(t)

	for idx, tc := range testCases {
		tcfg.Run2(Sprintf("[%d]: %v", idx, tc.in), tc, func(t *testing.T) (interface{}, error) {
			return sqrt(tc.shouldPanic, tc.in)
		})
	}
}
