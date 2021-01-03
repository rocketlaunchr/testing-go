package testing

import (
	"errors"
	"reflect"
	"strings"

	"fmt"
	"github.com/davecgh/go-spew/spew"
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

	s := spew.NewDefaultConfig()
	s.DisableMethods = false
	s.DisablePointerMethods = false

	s.Dump(err)
	s.Dump(target)
	fmt.Println("-------")

	switch target := target.(type) {
	case ErrContains:
		return strings.Contains(err.Error(), target.Str)
	}
	return errors.Is(err, target)
}
