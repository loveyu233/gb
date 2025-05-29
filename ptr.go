package gb

import (
	"reflect"
)

// IsPtr 检查target是否为指针
func IsPtr(target any) bool {
	objValue := reflect.ValueOf(target)
	if objValue.Kind() != reflect.Ptr {
		return false
	}

	if objValue.IsNil() {
		return false
	}

	return true
}
