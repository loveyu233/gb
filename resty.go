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

// initRestyClient 用来按默认配置初始化 Resty 客户端。
func initRestyClient() {
	defaultRestyClient = resty.New()
	defaultRestyClient.
		SetTimeout(30 * time.Second).
		SetRetryCount(2).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(2 * time.Second)
}

// RestyClient 用来返回带懒加载的 Resty 单例。
func RestyClient() *resty.Client {
	restyOnce.Do(initRestyClient)
	return defaultRestyClient
}

// R 用来基于默认客户端创建请求。
func R() *resty.Request {
	return RestyClient().R()
}
