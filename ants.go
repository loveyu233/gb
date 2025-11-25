package gb

import (
	"github.com/panjf2000/ants/v2"
)

// AntsSubmit 用来把任务提交到全局 ants 协程池。
func AntsSubmit(task func()) error {
	return ants.Submit(task)
}

// AntsRelease 用来释放默认的 ants 协程池。
func AntsRelease() {
	ants.Release()
}

// AntsReboot 用来重启默认的 ants 协程池。
func AntsReboot() {
	ants.Reboot()
}

// AntsNewPool 用来创建指定大小的 ants 协程池。
func AntsNewPool(size int, options ...ants.Option) (*ants.Pool, error) {
	return ants.NewPool(size, options...)
}
