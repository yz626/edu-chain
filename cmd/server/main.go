package main

import (
	"log"

	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/internal/data/db"
	"github.com/yz626/edu-chain/internal/data/redis"
	"github.com/yz626/edu-chain/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志器
	logger, err := logger.NewLogger(&cfg.Logger)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// 初始化数据库连接
	DB, err := db.NewDB(&cfg.Database, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 初始化 Redis 连接
	redisClient, err := redis.NewRedisClient(&cfg.Redis, logger)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	testMain(DB, redisClient)
}

func testMain(DB *db.DateDB, redis *redis.RedisClient) {
	fn, err := db.NewDBCloser(DB)
	if err != nil {
		logger.GetLogger().Named("数据库").Error("创建数据库关闭函数失败", logger.String("err", err.Error()))
	}
	defer fn()

	err = redis.Close()
	if err != nil {
		logger.GetLogger().Named("Redis").Error("关闭 Redis 连接失败", logger.String("err", err.Error()))
	}
}
