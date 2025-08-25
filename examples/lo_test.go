package examples

import (
	"testing"

	"github.com/samber/lo"
)

func TestLoT(t *testing.T) {
	arr := []int64{1, 2, 3, 4, 5}
	for i := 0; i < 10; i++ {
		t.Log(lo.TernaryF(len(arr) > i, func() int64 {
			return arr[i]
		}, func() int64 {
			return -1
		}))
	}
}
