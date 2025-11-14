package gb

import (
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	defaultRestyClient *resty.Client
	restyOnce          sync.Once
)

// initRestyClient 函数用于处理initRestyClient相关逻辑。
func initRestyClient() {
	defaultRestyClient = resty.New()
	defaultRestyClient.
		SetTimeout(30 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(2 * time.Second)
}

// RestyClient 函数用于处理RestyClient相关逻辑。
func RestyClient() *resty.Client {
	restyOnce.Do(initRestyClient)
	return defaultRestyClient
}

// R 函数用于处理R相关逻辑。
func R() *resty.Request {
	return RestyClient().R()
}
