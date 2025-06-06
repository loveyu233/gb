package gb

import (
	"github.com/gin-gonic/gin"
)

var (
	engine        *gin.Engine
	PublicRoutes  = make([]func(*gin.RouterGroup), 0) // 存储无需认证的公开路由处理函数
	PrivateRoutes = make([]func(*gin.RouterGroup), 0) // 存储需要认证的私有路由处理函数
)

func init() {
	PublicRoutes = append(PublicRoutes, func(group *gin.RouterGroup) {
		group.Any("/healthz", func(c *gin.Context) {
			c.Status(200)
		})
	})
}

type RouterConfig struct {
	model            string            // gin启动模式
	prefix           string            // api前缀
	authMiddleware   []gin.HandlerFunc // 认证api的中间件
	globalMiddleware []gin.HandlerFunc // 全局中间件
}

type GinRouterConfigOptionFunc func(*RouterConfig)

// WithGinRouterModel 设置gin的工作模式,不设置默认为debug
func WithGinRouterModel(model string) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.model = model
	}
}

// WithGinRouterPrefix 添加前缀不设置默认添加/api
func WithGinRouterPrefix(prefix string) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.prefix = prefix
	}
}

// WithGinRouterAuthHandler 用于对私有路由(PrivateRoutes)内的请求做校验
func WithGinRouterAuthHandler(handlers ...gin.HandlerFunc) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.authMiddleware = handlers
	}
}

// WithGinRouterGlobalMiddleware 用于对全局请求做校验
func WithGinRouterGlobalMiddleware(handlers ...gin.HandlerFunc) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.globalMiddleware = handlers
	}
}

// initRouter model默认为debug,prefix默认为/api,authMiddleware,globalMiddleware默认添加AddTraceID,AddRequestTime,ResponseLogger,GinRecovery
func initRouter(opts ...GinRouterConfigOptionFunc) {
	var config RouterConfig
	for _, opt := range opts {
		opt(&config)
	}
	if config.model == "" {
		config.model = "debug"
	}
	if config.prefix == "" {
		config.prefix = "/api"
	}
	if len(config.authMiddleware) == 0 {
		config.authMiddleware = []gin.HandlerFunc{GinAuth(map[string]any{}, DefaultGinConfig)}
	}
	if len(config.globalMiddleware) == 0 {
		config.globalMiddleware = []gin.HandlerFunc{AddTraceID(), AddRequestTime(), ResponseLogger(MiddlewareLogConfig{
			HeaderKeys: []string{"Token", "Authorization"},
			SaveLog:    nil,
			IsSaveLog:  false,
		}), GinRecovery(true)}
	}
	engine = newGinRouter(config.model, config.globalMiddleware...)
	registerRoutes(engine, config.prefix, config.authMiddleware...)
}

func newGinRouter(mode string, globalMiddlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(mode)
	engine := gin.Default()

	// 添加中间件
	engine.Use(globalMiddlewares...)

	return engine
}

func registerRoutes(r *gin.Engine, baseRouterPrefix string, authMiddlewares ...gin.HandlerFunc) {
	baseRouter := r.Group(baseRouterPrefix)

	// 注册公开路由
	for _, route := range PublicRoutes {
		route(baseRouter)
	}

	// 注册私有路由
	priRoute := baseRouter.Group("", authMiddlewares...)
	for _, route := range PrivateRoutes {
		route(priRoute)
	}

}
