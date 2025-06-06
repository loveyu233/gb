package gb

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MiddlewareTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("trace_id")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		c.Set("trace_id", traceID)
	}
}
