package gb

import (
	"time"

	"github.com/gin-gonic/gin"
)

var (
	PublicRoutes  = make([]func(*gin.RouterGroup), 0) // 存储无需认证的公开路由处理函数
	PrivateRoutes = make([]func(*gin.RouterGroup), 0) // 存储需要认证的私有路由处理函数
)

type RouterConfig struct {
	outputHealthz    bool              // 是否输出健康检查请求的日志输出
	model            GinModel          // gin启动模式
	prefix           string            // api前缀
	authMiddleware   []gin.HandlerFunc // 认证api的中间件
	globalMiddleware []gin.HandlerFunc // 全局中间件
	recordHeaderKeys []string          // 需要记录的请求头
	saveLog          func(ReqLog)      // 保存请求日志
	readTimeout      time.Duration
	writeTimeout     time.Duration
	idleTimeout      time.Duration
	maxHeaderBytes   int
	skipLog          bool
}

type GinModel string

// String 方法用于处理String相关逻辑。
func (m GinModel) String() string {
	return string(m)
}

var (
	GinModelRelease GinModel = "release"
	GinModelDebug   GinModel = "debug"
	GinModelTest    GinModel = "test"
)

type GinRouterConfigOptionFunc func(*RouterConfig)

// WithGinSkipLog 函数用于处理WithGinSkipLog相关逻辑。
func WithGinSkipLog(skipLog bool) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.skipLog = skipLog
	}
}

// WithGinReadTimeout 函数用于处理WithGinReadTimeout相关逻辑。
func WithGinReadTimeout(d time.Duration) GinRouterConfigOptionFunc {
	return func(routerConfig *RouterConfig) {
		routerConfig.readTimeout = d
	}
}

// WithGinWriteTimeout 函数用于处理WithGinWriteTimeout相关逻辑。
func WithGinWriteTimeout(d time.Duration) GinRouterConfigOptionFunc {
	return func(routerConfig *RouterConfig) {
		routerConfig.writeTimeout = d
	}
}

// WithGinIdleTimeout 函数用于处理WithGinIdleTimeout相关逻辑。
func WithGinIdleTimeout(d time.Duration) GinRouterConfigOptionFunc {
	return func(routerConfig *RouterConfig) {
		routerConfig.idleTimeout = d
	}
}

// WithGinMaxHeaderBytes 函数用于处理WithGinMaxHeaderBytes相关逻辑。
func WithGinMaxHeaderBytes(d int) GinRouterConfigOptionFunc {
	return func(routerConfig *RouterConfig) {
		routerConfig.maxHeaderBytes = d
	}
}

// WithGinRouterModel 函数用于处理WithGinRouterModel相关逻辑。
func WithGinRouterModel(model GinModel) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.model = model
	}
}

// WithGinRouterOutputHealthzLog 函数用于处理WithGinRouterOutputHealthzLog相关逻辑。
func WithGinRouterOutputHealthzLog() GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.outputHealthz = true
	}
}

// WithGinRouterPrefix 函数用于处理WithGinRouterPrefix相关逻辑。
func WithGinRouterPrefix(prefix string) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.prefix = prefix
	}
}

// WithGinRouterAuthHandler 函数用于处理WithGinRouterAuthHandler相关逻辑。
func WithGinRouterAuthHandler(handlers ...gin.HandlerFunc) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.authMiddleware = handlers
	}
}

// WithGinRouterGlobalMiddleware 函数用于处理WithGinRouterGlobalMiddleware相关逻辑。
func WithGinRouterGlobalMiddleware(handlers ...gin.HandlerFunc) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.globalMiddleware = handlers
	}
}

// WithGinRouterLogRecordHeaderKeys 函数用于处理WithGinRouterLogRecordHeaderKeys相关逻辑。
func WithGinRouterLogRecordHeaderKeys(keys []string) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.recordHeaderKeys = keys
	}
}

// WithGinRouterLogSaveLog 函数用于处理WithGinRouterLogSaveLog相关逻辑。
func WithGinRouterLogSaveLog(f func(ReqLog)) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.saveLog = f
	}
}

// initPrivateRouter 函数用于处理initPrivateRouter相关逻辑。
func initPrivateRouter(config RouterConfig) *gin.Engine {
	publicRoutes := make([]func(*gin.RouterGroup), 0, len(PublicRoutes)+1)
	publicRoutes = append(publicRoutes, func(group *gin.RouterGroup) {
		if !config.outputHealthz {
			group.Any("/healthz", GinLogSetSkipLogFlag(), func(c *gin.Context) {
				c.Status(200)
			})
		} else {
			group.Any("/healthz", func(c *gin.Context) {
				c.Status(200)
			})
		}
	})
	publicRoutes = append(publicRoutes, PublicRoutes...)

	privateRoutes := append(make([]func(*gin.RouterGroup), 0, len(PrivateRoutes)), PrivateRoutes...)

	config.globalMiddleware = append(config.globalMiddleware, MiddlewareTraceID(), MiddlewareRequestTime(), MiddlewareRecovery())
	if !config.skipLog {
		config.globalMiddleware = append(config.globalMiddleware, MiddlewareLogger(MiddlewareLogConfig{
			HeaderKeys: config.recordHeaderKeys,
			SaveLog:    config.saveLog,
		}))
	}

	engine := newGinRouter(config.model, config.globalMiddleware...)
	registerRoutes(engine, config.prefix, publicRoutes, privateRoutes, config.authMiddleware...)
	return engine
}

// newGinRouter 函数用于处理newGinRouter相关逻辑。
func newGinRouter(mode GinModel, globalMiddlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(mode.String())
	engine := gin.New()

	// 添加中间件
	engine.Use(globalMiddlewares...)

	return engine
}

// registerRoutes 函数用于处理registerRoutes相关逻辑。
func registerRoutes(r *gin.Engine, baseRouterPrefix string, publicRoutes, privateRoutes []func(*gin.RouterGroup), authMiddlewares ...gin.HandlerFunc) {
	baseRouter := r.Group(baseRouterPrefix)
	for _, route := range publicRoutes {
		route(baseRouter)
	}

	priRoute := baseRouter.Group("", authMiddlewares...)
	for _, route := range privateRoutes {
		route(priRoute)
	}
}
