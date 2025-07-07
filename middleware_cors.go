package gb

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许的请求来源
		c.Header("Access-Control-Allow-Origin", "*")
		// 允许的请求方法
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		// 允许的请求头
		c.Header("Access-Control-Allow-Headers", "*")
		// 允许暴露的响应头
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		// 是否允许携带 cookie
		c.Header("Access-Control-Allow-Credentials", "true")
		// 设置预检请求的缓存时间 1天
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
