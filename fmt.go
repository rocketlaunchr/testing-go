package testing

import "fmt"

func Sprintf(format string, a ...interface{}) string { return fmt.Sprintf(format, a...) }
func Errorf(format string, a ...interface{}) error   { return fmt.Errorf(format, a...) }

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
	} else if err == ErrPanic {
		return "<panic>"
	} else if err == ErrAny {
		return "<any error>"
	}
	return fmt.Sprintf("%+#v", err)
}
