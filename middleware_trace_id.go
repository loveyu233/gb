package gb

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceIDHeader = "Trace-Id"

// MiddlewareTraceID 函数用于处理MiddlewareTraceID相关逻辑。
func MiddlewareTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get(TraceIDHeader)
		if traceID == "" {
			traceID = c.Request.Header.Get("trace_id")
		}
		if traceID == "" {
			traceID = c.Request.Header.Get("X-Request-Id")
		}
		if traceID == "" {
			traceID = uuid.NewString()
		}
		c.Header(TraceIDHeader, traceID)
		c.Header("X-Request-Id", traceID)
		c.Set("trace_id", traceID)
		c.Next()
	}
}
