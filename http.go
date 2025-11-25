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

// InitHTTPServerAndStart 用来根据路由配置启动 HTTP 服务并注册钩子。
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
	engine := initPrivateRouter(config)
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

// startHTTPServer 用来启动底层 http.Server。
func (h *HTTPServer) startHTTPServer() {
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

// setupGracefulShutdown 用来注册系统信号以优雅关闭服务。
func (h *HTTPServer) setupGracefulShutdown() {
	NewHook().Close(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := h.server.Shutdown(ctx); err != nil {
			log.Printf("setup graceful shutdown err: %s\n", err)
		}
	})
}
