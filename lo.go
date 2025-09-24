package gb

import (
	"reflect"

	"github.com/samber/lo"
)

func LoMap[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	return lo.Map(collection, iteratee)
}

func LoSliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	return lo.SliceToMap(collection, transform)
}

func LoTernary[T any](condition bool, ifOutput T, elseOutput T) T {
	return lo.Ternary(condition, ifOutput, elseOutput)
}

func LoTernaryFunc[T any](condition bool, ifFunc func() T, elseFunc func() T) T {
	return lo.TernaryF(condition, ifFunc, elseFunc)
}

func LoWithout[T comparable, Slice interface{ ~[]T }](collection Slice, exclude ...T) Slice {
	return lo.Without(collection, exclude...)
}

func LoUniq[T comparable, Slice interface{ ~[]T }](collection Slice) Slice {
	return lo.Uniq(collection)
}

func LoContains[T comparable](collection []T, element T) bool {
	return lo.Contains(collection, element)
}

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

func LoToPtr[T any](x T) *T {
	return lo.ToPtr(x)
}

func LoFromPtr[T any](x *T) T {
	return lo.FromPtr(x)
}
