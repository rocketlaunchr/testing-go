package testing

import (
	"errors"
	"reflect"
	"strings"
)

// Comparator returns true if x and y are equal (as determined by the Comparator function).
// x and y are the returned and expected value respectively from a function being tested.
// The default Comparator uses reflect.DeepEqual. However, the github.com/google/go-cmp package
// may also be used. See: https://github.com/google/go-cmp.
type Comparator func(x, y interface{}) bool

type errComparator func(err, target error) bool

func deepEqual(err, target error) bool {
	switch target := target.(type) {
	case ErrContains:
		if err == nil {
			return false
		}
		return strings.Contains(err.Error(), target.Substr)
	}
	return reflect.DeepEqual(err, target)
}

func is(err, target error) bool {
	switch target := target.(type) {
	case ErrContains:
		if err == nil {
			return false
		}
		return strings.Contains(err.Error(), target.Substr)
	}
	return errors.Is(err, target)
}
