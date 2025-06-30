package gb

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
)

// MiddlewareRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
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
