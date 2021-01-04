package testing

import "fmt"

// Sprintf convenience function.
//
// See: https://pkg.go.dev/fmt#Sprintf
func Sprintf(format string, a ...interface{}) string { return fmt.Sprintf(format, a...) }

// Errorf convenience function.
//
// See: https://pkg.go.dev/fmt#Errorf
func Errorf(format string, a ...interface{}) error { return fmt.Errorf(format, a...) }

func fmtError(err error, not bool) (rstr string) {
	defer func() {
		if not {
			if err == ErrAny {
				rstr = "<nil>"
			} else {
				rstr = "NOT " + rstr
			}
		}
	}()

	if err == nil {
		return "<nil>"
	} else if err == PanicExpected {
		return "<panic>"
	} else if err == ErrAny {
		return "<any error>"
	}

	switch err := err.(type) {
	case ErrContains:
		return "<contains: \"" + err.Str + "\">"
	}
	return fmt.Sprintf("%+#v", err)
}

func fmtVal(val interface{}, not bool) string {
	if not {
		return fmt.Sprintf("NOT %+#v", val)
	}
	return fmt.Sprintf("%+#v", val)
}
