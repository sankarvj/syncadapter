package utils

import (
	"reflect"
	"strings"
)

func Tablename(in interface{}) string {
	return strings.ToLower(reflect.TypeOf(in).Elem().Name() + "s")
}
