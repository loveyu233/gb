package gb

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	Engine        *gin.Engine
	PublicRoutes  = make([]func(*gin.RouterGroup), 0) // 存储无需认证的公开路由处理函数
	PrivateRoutes = make([]func(*gin.RouterGroup), 0) // 存储需要认证的私有路由处理函数
)

type RouterConfig struct {
	model            string            // gin启动模式
	prefix           string            // api前缀
	authMiddleware   []gin.HandlerFunc // 认证api的中间件
	globalMiddleware []gin.HandlerFunc // 全局中间件
}

type RouterConfigOption func(*RouterConfig)

func WithModel(model string) RouterConfigOption {
	return func(config *RouterConfig) {
		config.model = model
	}
}

func WithPrefix(prefix string) RouterConfigOption {
	return func(config *RouterConfig) {
		config.prefix = prefix
	}
}

func WithAuthHandler(handlers ...gin.HandlerFunc) RouterConfigOption {
	return func(config *RouterConfig) {
		config.authMiddleware = handlers
	}
}
func WithGlobalMiddleware(handlers ...gin.HandlerFunc) RouterConfigOption {
	return func(config *RouterConfig) {
		config.globalMiddleware = handlers
	}
}

// InitRouter model默认为debug,prefix默认为/api,authMiddleware,globalMiddleware默认添加AddTraceID,AddRequestTime,ResponseLogger,GinRecovery
func InitRouter(opts ...RouterConfigOption) {
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
	Engine = newGinRouter(config.model, config.globalMiddleware...)
	registerRoutes(Engine, config.prefix, config.authMiddleware...)
}

func newGinRouter(mode string, globalMiddlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(mode)
	engine := gin.Default()

	// 添加中间件
	engine.Use(globalMiddlewares...)

	return engine
}

func registerRoutes(r *gin.Engine, baseRouterPrefix string, authMiddlewares ...gin.HandlerFunc) {
	r.Any("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

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

func CreateHTTPServer(listenAddr string) *http.Server {
	return &http.Server{
		Addr:    listenAddr,
		Handler: Engine,
	}
}

func StartHTTPServer(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %s\n", err)
	}
}

func SetupGracefulShutdown(server *http.Server) {
	NewHook().Close(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("正常关机err: %s\n", err)
		} else {
			log.Printf("优雅关机成功\n")
		}
	})
}
