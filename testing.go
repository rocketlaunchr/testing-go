package testing

import (
	"errors"
	"reflect"
	"testing"
)

type TestConfig struct {
	t     *testing.T
	C     Comparator
	Fatal bool
}

func NewTestConfig(t *testing.T) *TestConfig {
	return &TestConfig{
		t: t,
	}
}

func (tcfg TestConfig) Run1(name string, tc interface{}, f func(t *testing.T) interface{}) {
	tcfg.run2(name, tc, func(t *testing.T) (interface{}, error) {
		out := f(t)
		return out, nil
	}, 1)
}

func (tcfg TestConfig) RunErr(name string, tc interface{}, f func(t *testing.T) error) {
	tcfg.run2(name, tc, func(t *testing.T) (interface{}, error) {
		err := f(t)
		return nil, err
	}, 2)
}

func (tcfg TestConfig) Run2(name string, tc interface{}, f func(t *testing.T) (interface{}, error)) {
	tcfg.run2(name, tc, f, 0)
}

func (tcfg TestConfig) run2(name string, tc interface{}, f func(t *testing.T) (interface{}, error), mode int) {

	// Expected output value
	expOut, found := structVal(tc, "ExpOut")
	if !found {
		if mode != 2 {
			panic("ExpOut field not found in test case")
		}
	}

	if expOut == CustomTest {
		tcfg.t.Run(name, func(t *testing.T) { f(t) })
		return
	}

	// Expected Error
	expErr_, found := structVal(tc, "ExpErr")
	if !found {
		if mode != 1 {
			panic("ExpErr field not found in test case")
		}
	}
	expErr, _ := expErr_.(error)

	if expErr == CustomTest {
		tcfg.t.Run(name, func(t *testing.T) { f(t) })
		return
	}

	if expOut == PanicExpected {
		expOut = nil
		expErr = PanicExpected
	}

	comparator := tcfg.C
	if comparator == nil {
		comparator = reflect.DeepEqual
	}

	var errChecker errComparator
	errChecker = deepEqual

	var notVal bool
	if x, ok := expOut.(NotEqual); ok {
		// TODO: Make it recursive?
		notVal = true
		expOut = x.Val
		compCpy := comparator
		comparator = func(x, y interface{}) bool { return !compCpy(x, y) }
	}

	// Check for Is
	if errors.Is(expErr, Is{}) {
		// Bug: Assumes Is is immediate => recursively unwrap
		expErr = expErr.(Is).Err
		errChecker = func(err, target error) bool { return errors.Is(err, target) } // is
	}

	// Check for ErrAny
	var any bool
	if errors.Is(expErr, ErrAny) {
		any = true
	}

	// Check for NotEqual
	var notError bool
	if errors.Is(expErr, NotEqual{}) {
		if any {
			expErr = nil
			any = false
		} else {
			notError = true
			// Bug: Assumes NotEqual is immediate => recursively unwrap
			inner := expErr.(NotEqual).Val
			if inner == nil {
				expErr = nil
			} else {
				expErr = inner.(error)
			}
			errCheckerCpy := errChecker
			errChecker = func(err, target error) bool { return !errCheckerCpy(err, target) }
		}
	}

	tcfg.t.Run(name, func(t *testing.T) {
		gotVal, gotErr := func(t *testing.T) (gotVal interface{}, gotErr error) {
			defer func() {
				if recover() != nil {
					gotErr = PanicExpected
				}
			}()
			return f(t)
		}(t)

		if any {
			if !((expErr.(error) == nil && gotErr == nil) || (expErr.(error) != nil && gotErr != nil) && comparator(gotVal, expOut)) {
				t.Logf("got (%s, %s) ; want (%s, %s)", fmtVal(gotVal, false), fmtError(gotErr, false), fmtVal(expOut, notVal), fmtError(expErr, notError))
				if tcfg.Fatal {
					t.FailNow()
				} else {
					t.Fail()
				}
			}
			return
		}

		if !(errChecker(gotErr, expErr) && (gotErr == PanicExpected || comparator(gotVal, expOut))) {
			t.Logf("got (%s, %s) ; want (%s, %s)", fmtVal(gotVal, false), fmtError(gotErr, false), fmtVal(expOut, notVal), fmtError(expErr, notError))
			if tcfg.Fatal {
				t.FailNow()
			} else {
				t.Fail()
			}
		}
	})
}
