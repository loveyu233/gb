package gb

import (
	"github.com/panjf2000/ants/v2"
)

func AntsSubmit(task func()) error {
	return ants.Submit(task)
}

func AntsRelease() {
	ants.Release()
}

func AntsReboot() {
	ants.Reboot()
}

func AntsNewPool(size int, options ...ants.Option) (*ants.Pool, error) {
	return ants.NewPool(size)
}
