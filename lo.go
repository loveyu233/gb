package gb

import (
	"reflect"

	"github.com/samber/lo"
)

// LoMap 函数用于处理LoMap相关逻辑。
func LoMap[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	return lo.Map(collection, iteratee)
}

// LoSliceToMap 函数用于处理LoSliceToMap相关逻辑。
func LoSliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	return lo.SliceToMap(collection, transform)
}

// LoTernary 函数用于处理LoTernary相关逻辑。
func LoTernary[T any](condition bool, ifOutput T, elseOutput T) T {
	return lo.Ternary(condition, ifOutput, elseOutput)
}

// LoTernaryFunc 函数用于处理LoTernaryFunc相关逻辑。
func LoTernaryFunc[T any](condition bool, ifFunc func() T, elseFunc func() T) T {
	return lo.TernaryF(condition, ifFunc, elseFunc)
}

// LoWithout 函数用于处理LoWithout相关逻辑。
func LoWithout[T comparable, Slice interface{ ~[]T }](collection Slice, exclude ...T) Slice {
	return lo.Without(collection, exclude...)
}

// LoUniq 函数用于处理LoUniq相关逻辑。
func LoUniq[T comparable, Slice interface{ ~[]T }](collection Slice) Slice {
	return lo.Uniq(collection)
}

// LoContains 函数用于处理LoContains相关逻辑。
func LoContains[T comparable](collection []T, element T) bool {
	return lo.Contains(collection, element)
}

// IsPtr 函数用于处理IsPtr相关逻辑。
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

// LoToPtr 函数用于处理LoToPtr相关逻辑。
func LoToPtr[T any](x T) *T {
	return lo.ToPtr(x)
}

// LoFromPtr 函数用于处理LoFromPtr相关逻辑。
func LoFromPtr[T any](x *T) T {
	return lo.FromPtr(x)
}
