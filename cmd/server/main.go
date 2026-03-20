package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/internal/data/db"
	"github.com/yz626/edu-chain/pkg/logger"
)

func main() {
	// 加载配置
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志系统
	logConfig := conf.Logger.ToLoggerConfig()
	if err := logger.Init(logConfig); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 启动时记录日志
	logger.System().Info("Application starting",
		logger.StringField("mode", conf.Server.Mode),
		logger.StringField("host", conf.Server.Addr()),
	)

	// 初始化数据库连接
	if err := db.Init(conf); err != nil {
		logger.System().Fatal("Failed to initialize database",
			logger.StringField("error", err.Error()),
		)
	}

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 关闭时记录日志
	logger.System().Info("Application shutting down")

	// 同步日志
	if err := logger.Sync(); err != nil {
		log.Printf("Failed to sync logger: %v", err)
	}
}
