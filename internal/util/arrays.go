package util

import "reflect"

func IsLast(i int, a interface{}) bool {
	switch reflect.TypeOf(a).Kind() {
	case reflect.Array:
	case reflect.Slice:
		return i == reflect.ValueOf(a).Len()-1

	}
	return false
}
