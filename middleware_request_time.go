package gb

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// MiddlewareRequestTime 函数用于处理MiddlewareRequestTime相关逻辑。
func MiddlewareRequestTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := Now()
		c.Set("request_time", startTime)
		c.Next()
		c.Header("response_time", fmt.Sprintf("%dms", Now().Sub(startTime).Milliseconds()))
	}
}
