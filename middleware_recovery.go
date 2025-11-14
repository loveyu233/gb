package gb

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// MiddlewareRecovery 函数用于处理MiddlewareRecovery相关逻辑。
func MiddlewareRecovery(log ...GBLog) gin.HandlerFunc {
	var loclLog GBLog

	if len(log) > 0 {
		loclLog = log[0]
	} else {
		loclLog = new(GbDefaultlogger)
	}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				loclLog.Errorf("panic:%s;stack:%s", err, string(debug.Stack()))
				ResponseError(c, ErrServerBusy.WithMessage("panic:%v", err))
				c.AbortWithStatus(http.StatusOK)
			}
		}()
		c.Next()
	}
}
