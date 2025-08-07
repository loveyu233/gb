package gb

import (
	"time"

	"golang.org/x/net/context"
)

func Context(ttl ...int64) context.Context {
	var sec int64 = 3
	if len(ttl) > 0 {
		sec = ttl[0]
	}
	timeout, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(sec))
	return timeout
}
