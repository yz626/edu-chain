package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/pkg/logger"
)

// GRPCServer gRPC服务器
type GRPCServer struct {
	server   *grpc.Server
	wg       sync.WaitGroup
	config   *config.Config
	log      *logger.Logger
	listener net.Listener
}

// NewGRPCServer 创建gRPC服务器
func NewGRPCServer(cfg *config.Config, log *logger.Logger) *GRPCServer {
	// 配置gRPC服务器选项
	serverOptions := []grpc.ServerOption{
		// 保持连接
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 300, // 5分钟
			MaxConnectionAge:  600, // 10分钟
			Timeout:           20,  // 20秒
		}),
		// 最大消息大小
		grpc.MaxConcurrentStreams(100),
	}

	server := grpc.NewServer(serverOptions...)

	return &GRPCServer{
		server: server,
		config: cfg,
		log:    log,
	}
}

// GetServer 获取gRPC服务器实例
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// RegisterService 注册服务
func (s *GRPCServer) RegisterService(fn func(server *grpc.Server)) {
	fn(s.server)
}

// Start 启动gRPC服务器
func (s *GRPCServer) Start() error {
	// 创建监听器
	addr := fmt.Sprintf("%s:%d", s.config.GRPC.Host, s.config.GRPC.Port)
	var err error
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听地址失败: %w", err)
	}

	// 启动服务器
	go func() {
		s.log.Info("启动gRPC服务器", logger.String("addr", addr))
		if err := s.server.Serve(s.listener); err != nil {
			s.log.Error("gRPC服务器启动失败", logger.Any("error", err))
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	s.log.Info("正在关闭gRPC服务器...")

	// 优雅关闭
	s.server.GracefulStop()

	s.log.Info("gRPC服务器已关闭")
	return nil
}

// StartWithListener 使用指定的监听器启动
func (s *GRPCServer) StartWithListener(listener net.Listener) error {
	s.listener = listener
	s.log.Info("启动gRPC服务器", logger.String("addr", listener.Addr().String()))
	return s.server.Serve(listener)
}

// Stop 停止服务器
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

// ==================== 拦截器 ====================

// unaryInterceptor  unary拦截器
func (s *GRPCServer) unaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		s.log.Info("gRPC调用", logger.String("method", info.FullMethod))

		// 可以在这里添加认证、日志等逻辑
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			s.log.Debug("元数据", logger.Any("metadata", md))
		}

		return handler(ctx, req)
	}
}

// streamInterceptor 流拦截器
func (s *GRPCServer) streamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		s.log.Info("gRPC流调用", logger.String("method", info.FullMethod))
		return handler(srv, ss)
	}
}

// ContextWithMetadata 将元数据添加到上下文
func ContextWithMetadata(ctx context.Context, key, value string) context.Context {
	md := metadata.Pairs(key, value)
	return metadata.NewOutgoingContext(ctx, md)
}
