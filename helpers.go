package testing

import "reflect"

func structVal(tc interface{}, fieldName string) (interface{}, bool) {
	fieldNotFound := reflect.Value{}
	rv := reflect.ValueOf(tc)
	expErr_ := rv.FieldByName(fieldName)
	if expErr_ == fieldNotFound {
		return nil, false
	} else {
		return expErr_.Interface(), true
	}
}
