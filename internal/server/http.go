package server

import (
	"github.com/gin-gonic/gin"
	"github.com/yz626/edu-chain/config"
	middleware "github.com/yz626/edu-chain/internal/routes/middleware/casbin"
)

// HTTP HTTP服务器
type HTTP struct {
	router *gin.Engine
	config *config.Config
}

// NewHTTPServer 创建HTTP服务器
func NewHTTPServer(cfg *config.Config) *HTTP {
	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	// 初始化 Casbin
	if err := middleware.InitCasbin("internal/casbin/model.conf", "internal/casbin/policy.csv"); err != nil {
		panic("Failed to initialize Casbin: " + err.Error())
	}

	// 设置全局 Casbin 中间件（可选）
	// router.Use(middleware.RequirePermission("*", "read"))

	return &HTTP{
		router: router,
		config: cfg,
	}
}

// Start 启动服务器
func (s *HTTP) Start() error {
	return s.router.Run(s.config.Server.Addr())
}

// GetRouter 获取路由实例（供外部使用）
func (s *HTTP) GetRouter() *gin.Engine {
	return s.router
}
