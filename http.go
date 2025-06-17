package gb

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

// InitHTTPServer 默认自动添加/前缀/healthz的any请求用于存活和就绪检查,没有配置前缀则默认前缀为/api
func InitHTTPServer(listenAddr string, opts ...GinRouterConfigOptionFunc) *HTTPServer {
	initRouter(opts...)
	return &HTTPServer{server: &http.Server{
		Addr:    listenAddr,
		Handler: engine,
	}}
}

// StartHTTPServer 启动http服务
func (h *HTTPServer) StartHTTPServer() {
	go func() {
		if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()
	h.setupGracefulShutdown()
}

// SetupGracefulShutdown 优雅关机
func (h *HTTPServer) setupGracefulShutdown() {
	NewHook().Close(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := h.server.Shutdown(ctx); err != nil {
			log.Printf("正常关机err: %s\n", err)
		} else {
			log.Printf("优雅关机成功\n")
		}
	})
}
