package gb

import (
	"time"

	"golang.org/x/net/context"
)

// Context 用来创建一个带超时的 context，默认 3 秒。
func Context(ttl ...int64) (context.Context, context.CancelFunc) {
	var sec int64 = 3
	if len(ttl) > 0 {
		sec = ttl[0]
	}
	return context.WithTimeout(context.Background(), time.Second*time.Duration(sec))
}

// DurationSecond 用来把秒值转换为 time.Duration。
func DurationSecond(Second int) time.Duration {
	return time.Duration(Second) * time.Second
}
