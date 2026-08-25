package util

import "reflect"

func IsLast(i int, a any) bool {
	switch reflect.TypeOf(a).Kind() {
	case reflect.Array, reflect.Slice:
		return i == reflect.ValueOf(a).Len()-1
	default:
		return false
	}
}
