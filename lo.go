package gb

import (
	"reflect"

	"github.com/samber/lo"
)

// LoMap 用来对切片执行映射转换。
func LoMap[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	return lo.Map(collection, iteratee)
}

// LoSliceToMap 用来把切片转换为键值映射。
func LoSliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	return lo.SliceToMap(collection, transform)
}

// LoTernary 用来模拟条件表达式并返回两个值之一。
func LoTernary[T any](condition bool, ifOutput T, elseOutput T) T {
	return lo.Ternary(condition, ifOutput, elseOutput)
}

// LoTernaryFunc 用来在条件成立时延迟执行对应函数。
func LoTernaryFunc[T any](condition bool, ifFunc func() T, elseFunc func() T) T {
	return lo.TernaryF(condition, ifFunc, elseFunc)
}

// LoWithout 用来返回排除了指定元素的新切片。
func LoWithout[T comparable, Slice interface{ ~[]T }](collection Slice, exclude ...T) Slice {
	return lo.Without(collection, exclude...)
}

// LoUniq 用来去除切片中的重复值。
func LoUniq[T comparable, Slice interface{ ~[]T }](collection Slice) Slice {
	return lo.Uniq(collection)
}

// LoContains 用来判断切片是否包含指定元素。
func LoContains[T comparable](collection []T, element T) bool {
	return lo.Contains(collection, element)
}

// IsPtr 用来检查传入对象是否为非空指针。
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

// LoToPtr 用来返回值的指针表示。
func LoToPtr[T any](x T) *T {
	return lo.ToPtr(x)
}

// LoFromPtr 用来在指针为空时返回零值，否则返回其内容。
func LoFromPtr[T any](x *T) T {
	return lo.FromPtr(x)
}
