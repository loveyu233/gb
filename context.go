package gb

import (
	"time"

	"golang.org/x/net/context"
)

// Context 函数用于处理Context相关逻辑。
func Context(ttl ...int64) (context.Context, context.CancelFunc) {
	var sec int64 = 3
	if len(ttl) > 0 {
		sec = ttl[0]
	}
	return context.WithTimeout(context.Background(), time.Second*time.Duration(sec))
}

// DurationSecond 函数用于处理DurationSecond相关逻辑。
func DurationSecond(Second int) time.Duration {
	return time.Duration(Second) * time.Second
}
