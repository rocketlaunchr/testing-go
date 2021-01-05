package testing

import "reflect"

func structVal(tc interface{}, fieldName string) (interface{}, bool) {
	fieldNotFound := reflect.Value{}
	rv := reflect.ValueOf(tc)
	_expErr := rv.FieldByName(fieldName)
	if _expErr == fieldNotFound {
		return nil, false
	}
	return _expErr.Interface(), true
}
