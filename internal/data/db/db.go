package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yz626/edu-chain/config"
)

// DB 全局数据库连接
var DB *gorm.DB

// Init 初始化数据库连接
func Init(cfg *config.Config) error {
	var err error

	// 配置GORM日志级别
	logLevel := logger.Info
	if cfg.Server.Mode == "release" {
		logLevel = logger.Warn
	}

	// 连接数据库
	dsn := cfg.Database.MySQLDSN()
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	// 设置连接最大存活时间
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.Timeout) * time.Second)

	fmt.Println("Database connected successfully")
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}
