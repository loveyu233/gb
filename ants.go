package gb

import (
	"github.com/panjf2000/ants/v2"
)

// AntsSubmit 函数用于处理AntsSubmit相关逻辑。
func AntsSubmit(task func()) error {
	return ants.Submit(task)
}

// AntsRelease 函数用于处理AntsRelease相关逻辑。
func AntsRelease() {
	ants.Release()
}

// AntsReboot 函数用于处理AntsReboot相关逻辑。
func AntsReboot() {
	ants.Reboot()
}

// AntsNewPool 函数用于处理AntsNewPool相关逻辑。
func AntsNewPool(size int, options ...ants.Option) (*ants.Pool, error) {
	return ants.NewPool(size, options...)
}
