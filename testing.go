package testing

import (
	"reflect"
	"testing"
)

// TestConfig ...
//
// See ret1_test.go file for usage example.
type TestConfig struct {
	t     *testing.T
	C     Comparator
	Fatal bool
}

// NewTestConfig creates a new TestConfig.
func NewTestConfig(t *testing.T) *TestConfig {
	return &TestConfig{
		t: t,
	}
}

// Run1 is used when testing a function that returns a single (non-error) value.
//
// See ret1_test.go file for usage example.
func (tcfg TestConfig) Run1(name string, tc interface{}, f func(t *testing.T) interface{}) {
	tcfg.run2(name, tc, func(t *testing.T) (interface{}, error) {
		out := f(t)
		return out, nil
	}, 1)
}

// RunErr is used when testing a function that returns a single error value.
//
// See ret1_test.go file.
func (tcfg TestConfig) RunErr(name string, tc interface{}, f func(t *testing.T) error) {
	tcfg.run2(name, tc, func(t *testing.T) (interface{}, error) {
		err := f(t)
		return nil, err
	}, 2)
}

// Run2 is used when testing a function that returns a value and an error.
//
// See ret2_test.go file for usage example.
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

	var neq = NotEqual{PanicExpected}
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
	if x, ok := expOut.(NotEqual); ok {
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
			case NotEqual:
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
					// ErrAny & NotEqual together is equivalent to checking for nil
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
