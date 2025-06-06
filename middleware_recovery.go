package gb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

// MiddlewareRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func MiddlewareRecovery(stack bool, log ...GBLog) gin.HandlerFunc {
	var loclLog GBLog

	if len(log) > 0 {
		loclLog = log[0]
	} else {
		loclLog = new(GbDefaultlogger)
	}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					loclLog.Errorf("path: %s\nerror: %v\nreq:%s\n", c.Request.URL.Path, err, string(httpRequest))
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					loclLog.Errorf("[Recovery from panic]\nerr:%s\nreq:%s\nstack:%s\n", err, string(httpRequest), string(debug.Stack()))
				} else {
					loclLog.Errorf("[Recovery from panic]\nerr:%s\nreq:%s\n", err, string(httpRequest))

				}
				ResponseError(c, ErrServerBusy.WithMessage(fmt.Sprintf("[Recovery from panic]: %v", err)))
				c.AbortWithStatus(http.StatusOK)
			}
		}()
		c.Next()
	}
}
