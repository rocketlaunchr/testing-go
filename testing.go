package testing

import (
	"errors"
	// "github.com/davecgh/go-spew/spew"
	"fmt"
	"reflect"
	"testing"
)

type TestConfig struct {
	t *testing.T
	C Comparator
	// Fmt
	// Fail
	// string contains
}

func NewTestConfig(t *testing.T) *TestConfig {
	return &TestConfig{
		t: t,
	}
}

func (tcfg TestConfig) Run1(name string, tc interface{}, f func(t *testing.T) interface{}) {

	tcfg.t.Run(name, func(t *testing.T) {
		f(t)
	})
}

func (tcfg TestConfig) Run2(name string, tc interface{}, f func(t *testing.T) (interface{}, error)) {

	comparator := tcfg.C
	if comparator == nil {
		comparator = reflect.DeepEqual
	}

	var errChecker errComparator
	errChecker = func(err, target error) bool { return reflect.DeepEqual(err, target) }

	// Expected output value
	expOut, found := structVal(tc, "ExpOut")
	if !found {
		panic("ExpOut field not found in test case")
	}
	var notVal bool
	if x, ok := expOut.(NotEqual); ok {
		// TODO: Make it recursive
		notVal = true
		expOut = x.Val
		compCpy := comparator
		comparator = func(x, y interface{}) bool { return !compCpy(x, y) }
	}

	// Expected Error
	expErr_, found := structVal(tc, "ExpErr")
	if !found {
		panic("ExpErr field not found in test case")
	}
	expErr, _ := expErr_.(error)

	// Check for Is
	if errors.Is(expErr, Is{}) {
		errChecker = func(err, target error) bool { return errors.Is(err, target) }
	}

	// Check for NotEqual
	var notError bool
	if errors.Is(expErr, NotEqual{}) {
		notError = true
		inner := expErr.(NotEqual).Val
		if inner == nil {
			expErr = nil
		} else {
			expErr = inner.(error)
		}
		errCheckerCpy := errChecker
		errChecker = func(err, target error) bool { return !errCheckerCpy(err, target) }
	}

	var any bool
	if errors.Is(expErr, ErrAny) {
		// any = true
		if notError {
			fmt.Println("A----")
			// panic("not ErrAny")
			expErr = nil
		} else {
			fmt.Println("B----")
			any = true
		}
	}

	tcfg.t.Run(name, func(t *testing.T) {

		gotVal, gotErr := func(t *testing.T) (gotVal interface{}, gotErr error) {
			defer func() {
				if recover() != nil {
					gotErr = ErrPanic
				}
			}()
			return f(t)
		}(t)

		// Test myself
		if expOut == CustomTest || expErr == CustomTest {
			return
		}

		if any {
			if !((expErr.(error) == nil && gotErr == nil) || (expErr.(error) != nil && gotErr != nil) && comparator(gotVal, expOut)) {
				t.Errorf("got (%+#v, %s) ; want (%+#v, %s)", gotVal, fmtError(gotErr, false), expOut, fmtError(expErr, notError))
			}
			return
		}

		if !(errChecker(gotErr, expErr) && (gotErr == ErrPanic || comparator(gotVal, expOut))) {
			if notVal {
				t.Errorf("got (%+#v, %s) ; want (NOT %+#v, %s)", gotVal, fmtError(gotErr, false), expOut, fmtError(expErr, notError))
			} else {
				t.Errorf("got (%+#v, %s) ; want (%+#v, %s)", gotVal, fmtError(gotErr, false), expOut, fmtError(expErr, notError))
			}
		}
	})
}
