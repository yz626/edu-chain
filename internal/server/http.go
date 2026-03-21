package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/pkg/logger"
)

// HTTPServer HTTP服务器
type HTTPServer struct {
	router *gin.Engine
	server *http.Server
	log    *logger.Logger
}

// HTTPServerSet Wire Provider Set
var HTTPServerSet = wire.NewSet(
	NewHTTPServer,
	wire.Bind(new(IHTTPServer), new(*HTTPServer)),
)

// IHTTPServer HTTP服务器接口
type IHTTPServer interface {
	Start() error
	StartWithAddr(addr string) error
	GetRouter() *gin.Engine
	InitMiddlewares()
}

// NewHTTPServer 创建HTTP服务器（Wire 注入点）
func NewHTTPServer(cfg *config.Config, router *gin.Engine, log *logger.Logger) *HTTPServer {
	gin.SetMode(cfg.Server.Mode)

	server := &http.Server{
		Addr:           cfg.Server.Addr(),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &HTTPServer{
		router: router,
		server: server,
		log:    log,
	}
}

// GetRouter 获取路由实例
func (s *HTTPServer) GetRouter() *gin.Engine {
	return s.router
}

// Start 启动服务器（支持优雅关闭）
func (s *HTTPServer) Start() error {
	go func() {
		s.log.Info("启动HTTP服务器", logger.String("addr", s.server.Addr))
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("HTTP服务器启动失败", logger.Any("error", err))
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.log.Info("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务器关闭失败: %w", err)
	}

	s.log.Info("服务器已优雅关闭")
	return nil
}

// StartWithAddr 使用指定地址启动（不处理信号）
func (s *HTTPServer) StartWithAddr(addr string) error {
	if addr != "" {
		s.server.Addr = addr
	}
	s.log.Info("启动HTTP服务器", logger.String("addr", s.server.Addr))
	return s.server.ListenAndServe()
}
