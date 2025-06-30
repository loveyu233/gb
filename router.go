package gb

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
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
	skipApiMap       map[string]struct{} // 日志输出跳过的api
	skipHealthz      bool                // 是否跳过健康检查请求的日志输出
	model            GinModel            // gin启动模式
	prefix           string              // api前缀
	authMiddleware   []gin.HandlerFunc   // 认证api的中间件
	globalMiddleware []gin.HandlerFunc   // 全局中间件
	recordHeaderKeys []string            // 需要记录的请求头
	saveLog          func(ReqLog)        // 保存请求日志
	tokenData        any                 // token携带的信息
}

type GinModel string

func (m GinModel) String() string {
	return string(m)
}

var (
	GinModelRelease GinModel = "release"
	GinModelDebug   GinModel = "debug"
	GinModelTest    GinModel = "test"
)

type GinRouterConfigOptionFunc func(*RouterConfig)

// WithGinRouterTokenData token携带的信息,必须是指针,尽量使用: new(结构体) 方式
func WithGinRouterTokenData(data any) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		if !IsPtr(data) {
			panic("data必须是指针类型")
		}
		config.tokenData = data
	}
}

// WithGinRouterModel 设置gin的工作模式,不设置默认为debug
func WithGinRouterModel(model GinModel) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.model = model
	}
}

// WithGinRouterSkipHealthzLog 是否跳过健康检查请求的日志输出
func WithGinRouterSkipHealthzLog() GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.skipHealthz = true
	}
}

func WithGinRouterSkipApiMap(skipApis ...string) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.skipApiMap = make(map[string]struct{})
		for _, item := range skipApis {
			config.skipApiMap[item] = struct{}{}
		}
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

// WithGinRouterLogRecordHeaderKeys 需要被记录的请求头
func WithGinRouterLogRecordHeaderKeys(keys []string) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.recordHeaderKeys = keys
	}
}

// WithGinRouterLogSaveLog 持久化日志可以在这里做
func WithGinRouterLogSaveLog(f func(ReqLog)) GinRouterConfigOptionFunc {
	return func(config *RouterConfig) {
		config.saveLog = f
	}
}

// initRouter model默认为debug,prefix默认为/api,authMiddleware,globalMiddleware默认添加AddTraceID,MiddlewareRequestTime,ResponseLogger,MiddlewareRecovery
func initRouter[T any](authConfig *GinAuthConfig[T], opts ...GinRouterConfigOptionFunc) {
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

	if config.skipApiMap == nil {
		config.skipApiMap = make(map[string]struct{})
	}

	if config.skipHealthz {
		config.skipApiMap[fmt.Sprintf("%s/healthz", config.prefix)] = struct{}{}
	}

	if config.model == GinModelDebug {
		var skips []string
		for k, _ := range config.skipApiMap {
			skips = append(skips, k)
		}
		fmt.Printf("跳过日志收集api:[%s]\n", strings.Join(skips, ";"))
	}

	if config.tokenData == nil {
		config.tokenData = new(TokenDefaultData)
	}

	if len(config.authMiddleware) == 0 {
		config.authMiddleware = []gin.HandlerFunc{GinAuth(authConfig)}
	}

	if len(config.globalMiddleware) == 0 {
		config.globalMiddleware = []gin.HandlerFunc{MiddlewareTraceID(), MiddlewareRequestTime(), MiddlewareLogger(MiddlewareLogConfig{
			HeaderKeys: config.recordHeaderKeys,
			SaveLog:    config.saveLog,
		}, config.skipApiMap), MiddlewareRecovery()}
	}
	engine = newGinRouter(config.model, config.globalMiddleware...)

	registerRoutes(engine, config.prefix, config.authMiddleware...)
}

func newGinRouter(mode GinModel, globalMiddlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(mode.String())
	engine := gin.New()

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
