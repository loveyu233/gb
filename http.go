package gb

import (
	"context"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

// InitHTTPServerAndStart 默认自动添加/前缀/healthz的any请求用于存活和就绪检查,没有配置前缀则默认前缀为/api
func InitHTTPServerAndStart(listenAddr string, opts ...GinRouterConfigOptionFunc) *HTTPServer {
	initRouter(opts...)
	server := &HTTPServer{server: &http.Server{
		Addr:    listenAddr,
		Handler: engine,
	}}
	go server.startHTTPServer()
	server.setupGracefulShutdown()
	return server
}

// StartHTTPServer 启动http服务
func (h *HTTPServer) startHTTPServer() {
	if err := h.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// setupGracefulShutdown 优雅关机
func (h *HTTPServer) setupGracefulShutdown() {
	NewHook().Close(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := h.server.Shutdown(ctx); err != nil {
			log.Printf("setup graceful shutdown err: %s\n", err)
		} else {
			log.Printf("setup graceful shutdown success \n")
		}
	})
}
