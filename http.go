package gb

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

// CreateHTTPServer 默认自动添加/前缀/healthz的any请求用于存活和就绪检查,没有配置前缀则默认前缀为/api
func CreateHTTPServer(listenAddr string, opts ...GinRouterConfigOptionFunc) *http.Server {
	initRouter(opts...)
	return &http.Server{
		Addr:    listenAddr,
		Handler: engine,
	}
}

// StartHTTPServer 启动http服务
func StartHTTPServer(server *http.Server) {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("listen: %s\n", err)
	}
}

// SetupGracefulShutdown 优雅关机
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
