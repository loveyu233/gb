package gb

import "github.com/gin-gonic/gin"

// MiddlewareRequestTime 注入请求时间
func MiddlewareRequestTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("request_time", Now())
	}
}
