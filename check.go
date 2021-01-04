package testing

import (
	"errors"
	"reflect"
	"strings"
)

type Comparator func(x, y interface{}) bool

type errComparator func(err, target error) bool

func deepEqual(err, target error) bool {
	switch target := target.(type) {
	case ErrContains:
		if err == nil {
			return false
		}
		return strings.Contains(err.Error(), target.Str)
	}
	return reflect.DeepEqual(err, target)
}

func is(err, target error) bool {
	switch target := target.(type) {
	case ErrContains:
		if err == nil {
			return false
		}
		return strings.Contains(err.Error(), target.Str)
	}
	return errors.Is(err, target)
}
