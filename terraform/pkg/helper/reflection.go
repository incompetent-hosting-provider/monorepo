package helper

import (
	"reflect"
)

func GetType(i interface{}) reflect.Type {
	return reflect.TypeOf(i)
}
