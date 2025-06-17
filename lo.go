package gb

import (
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
