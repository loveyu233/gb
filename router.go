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

func InitRouter(ginMode, prefix string, authMiddlewares []gin.HandlerFunc, globalMiddlewares ...gin.HandlerFunc) {
	Engine = newGinRouter(ginMode, globalMiddlewares...)
	registerRoutes(Engine, prefix, authMiddlewares...)
}

func newGinRouter(mode string, globalMiddlewares ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(mode)
	engine := gin.Default()

	// 添加中间件
	engine.Use(globalMiddlewares...)

	return engine
}

func registerRoutes(r *gin.Engine, prefix string, authMiddlewares ...gin.HandlerFunc) {
	baseRouter := r.Group(prefix)

	// 注册公开路由
	for _, route := range PublicRoutes {
		route(baseRouter)
	}

	// 注册私有路由
	priRoute := baseRouter.Group("/auth", authMiddlewares...)
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
			log.Printf("graceful shutdown err: %s\n", err)
		} else {
			log.Printf("graceful shutdown success\n")
		}
	})
}
