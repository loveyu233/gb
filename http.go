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

// InitHTTPServerAndStart 默认自动添加/前缀/healthz的any请求用于存活和就绪检查,没有配置前缀则默认前缀为/api
func InitHTTPServerAndStart(listenAddr string, opts ...GinRouterConfigOptionFunc) *HTTPServer {
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
	initPrivateRouter(config)
	server := &HTTPServer{server: &http.Server{
		Addr:    listenAddr,
		Handler: engine,
	}}
	if config.readTimeout > 0 {
		server.server.ReadTimeout = config.readTimeout
	}
	if config.writeTimeout > 0 {
		server.server.WriteTimeout = config.writeTimeout
	}
	if config.idleTimeout > 0 {
		server.server.IdleTimeout = config.idleTimeout
	}
	if config.maxHeaderBytes > 0 {
		server.server.MaxHeaderBytes = config.maxHeaderBytes
	}
	go server.startHTTPServer()
	server.setupGracefulShutdown()
	return server
}

// StartHTTPServer 启动http服务
func (h *HTTPServer) startHTTPServer() {
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
		}
	})
}
