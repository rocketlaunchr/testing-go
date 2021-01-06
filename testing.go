package testing

import (
	"reflect"
	"testing"
)

// TestConfig ...
//
// See ret1_test.go file for usage example.
type TestConfig struct {
	t *testing.T

	// C sets the function used to check if the observed return value (non error) matches
	// the expected return value. By default, reflect.DeepEqual is used. However, the github.com/google/go-cmp
	// package  may also be used. See: https://github.com/google/go-cmp.
	C Comparator

	// Fatal marks the test as having failed and stops its execution.
	//
	// See: https://golang.org/pkg/testing/#T.FailNow
	Fatal bool
}

// NewTestConfig creates a new TestConfig.
func NewTestConfig(t *testing.T) *TestConfig {
	return &TestConfig{
		t: t,
	}
}

// Run1 is used when testing a function that returns a single (non-error) value. tc is a struct that is
// expected to have a field named: ExpOut. It can be of the relevant type, but if you want to test
// for panics, then it must be interface{} so you can use PanicExpected (which is an error type).
//
// See ret1_test.go file for usage example.
//
// NOTE: Be wary of how Go interprets in-place constants. If you are expecting a float64, then don't type
// 1. Instead type 1.0. See: https://blog.golang.org/constants.
func (tcfg TestConfig) Run1(name string, tc interface{}, f func(t *testing.T) interface{}) {
	tcfg.run2(name, tc, func(t *testing.T) (interface{}, error) {
		out := f(t)
		return out, nil
	}, 1)
}

// RunErr is used when testing a function that returns a single error value. tc is a struct that is
// expected to have a field named: ExpErr of type error.
//
// See ret1_test.go file for usage example.
func (tcfg TestConfig) RunErr(name string, tc interface{}, f func(t *testing.T) error) {
	tcfg.run2(name, tc, func(t *testing.T) (interface{}, error) {
		err := f(t)
		return nil, err
	}, 2)
}

// Run2 is used when testing a function that returns a value and an error. tc is a struct that is
// expected to have 2 fields named: ExpOut and ExpErr (of type error). ExpOut can be of the relevant
// type, but if you want to test for panics, then it must be interface{} so you can use
// PanicExpected (which is an error type).
//
// See ret2_test.go file for usage example.
//
// NOTE: Be wary of how Go interprets in-place constants. If you are expecting a float64, then don't type
// 1. Instead type 1.0. See: https://blog.golang.org/constants.
func (tcfg TestConfig) Run2(name string, tc interface{}, f func(t *testing.T) (interface{}, error)) {
	tcfg.run2(name, tc, f, 0)
}

func (tcfg TestConfig) run2(name string, tc interface{}, f func(t *testing.T) (interface{}, error), mode int) {

	comparator := tcfg.C
	if comparator == nil {
		comparator = reflect.DeepEqual
	}

	var errChecker errComparator
	errChecker = deepEqual

	// Expected Error
	_expErr, found := structVal(tc, "ExpErr")
	if !found {
		if mode != 1 {
			panic("ExpErr field not found in test case")
		}
	}
	expErr, _ := _expErr.(error)

	if expErr == CustomTest {
		tcfg.t.Run(name, func(t *testing.T) { f(t) })
		return
	}

	// Expected output value
	expOut, found := structVal(tc, "ExpOut")
	if !found {
		if mode != 2 {
			panic("ExpOut field not found in test case")
		}
	}

	//
	// Expected Value
	//

	var neq = Not{PanicExpected}
	if expOut == CustomTest {
		tcfg.t.Run(name, func(t *testing.T) { f(t) })
		return
	} else if expOut == neq {
		// expErr = neq
		// expOut = nil
	} else if expOut == PanicExpected {
		expErr = PanicExpected
		expOut = nil
	}

	var notVal bool
	if x, ok := expOut.(Not); ok {
		notVal = true
		expOut = x.Val
		compCpy := comparator
		comparator = func(x, y interface{}) bool { return !compCpy(x, y) }
	}

	//
	// Expected Error
	//

	var (
		notError bool
		any      bool
	)

	if expErr != nil {
		err := expErr
	FOR:
		for {
			// Do we recognise this error as an internal type?
			switch x := err.(type) {
			case Not:
				notError = !notError
				errCheckerCpy := errChecker
				errChecker = func(err, target error) bool { return !errCheckerCpy(err, target) }
				if x.Val == nil || x.Val.(error) == nil {
					err = nil
					break FOR
				} else {
					err = x.Val.(error)
					continue FOR
				}
			case Is:
				errChecker = is
				if x.Err == nil {
					err = nil
					break FOR
				} else {
					err = x.Err
					continue FOR
				}
			case ErrContains:
				err = x
				break FOR
			default:
				switch err {
				case PanicExpected:
					break FOR
				case ErrAny:
					// ErrAny & Not together is equivalent to checking for nil
					if notError {
						notError = !notError
						errCheckerCpy := errChecker
						errChecker = func(err, target error) bool { return !errCheckerCpy(err, target) }
						err = nil
					} else {
						any = true
					}
					break FOR
				default:
					err = x
					break FOR
				}
			}
		}
		expErr = err
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
				if mode == 0 {
					t.Logf("got (%s, %s) ; want (%s, %s)", fmtVal(gotVal, false), fmtError(gotErr, false), fmtVal(expOut, notVal), fmtError(expErr, notError))
				} else {
					if expErr == PanicExpected || expErr == neq {
						t.Logf("got %s ; want %s", fmtVal(gotVal, false), fmtError(expErr, notError))
					} else {
						t.Logf("got %s ; want %s", fmtVal(gotVal, false), fmtVal(expOut, notVal))
					}
				}
				if tcfg.Fatal {
					t.FailNow()
				} else {
					t.Fail()
				}
			}
			return
		}

		if !(errChecker(gotErr, expErr) && (gotErr == PanicExpected || comparator(gotVal, expOut))) {
			if mode == 0 {
				t.Logf("got (%s, %s) ; want (%s, %s)", fmtVal(gotVal, false), fmtError(gotErr, false), fmtVal(expOut, notVal), fmtError(expErr, notError))
			} else {
				if expErr == PanicExpected || expErr == neq {
					t.Logf("got %s ; want %s", fmtVal(gotVal, false), fmtError(expErr, notError))
				} else {
					t.Logf("got %s ; want %s", fmtVal(gotVal, false), fmtVal(expOut, notVal))
				}
			}
			if tcfg.Fatal {
				t.FailNow()
			} else {
				t.Fail()
			}
		}
	})
}
