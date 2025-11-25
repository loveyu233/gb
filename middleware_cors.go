package gb

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Cors 用来为 gin 路由添加通用的跨域响应头。
func Cors() gin.HandlerFunc {
	const (
		allowMethods  = "POST, GET, OPTIONS, PUT, DELETE"
		allowHeaders  = "*"
		exposeHeaders = "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers"
		maxAge        = "86400"
	)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "false")
		}

		c.Header("Access-Control-Allow-Methods", allowMethods)
		c.Header("Access-Control-Allow-Headers", allowHeaders)
		c.Header("Access-Control-Expose-Headers", exposeHeaders)
		c.Header("Access-Control-Max-Age", maxAge)

		if strings.EqualFold(c.Request.Method, http.MethodOptions) {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
