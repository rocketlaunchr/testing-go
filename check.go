package testing

type Comparator func(x, y interface{}) bool

type errComparator func(err, target error) bool
