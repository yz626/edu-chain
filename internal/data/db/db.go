package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/google/wire"
	"github.com/yz626/edu-chain/config"
)

// ProviderSet 是 Wire 的 Provider 集合
//
// 使用方法：在 wire 包中导入此集合
//
//	import "github.com/yz626/edu-chain/internal/data/db"
//
//	var AppSet = wire.NewSet(
//	    db.ProviderSet,
//	)
var ProviderSet = wire.NewSet(NewDB, NewDBCloser)

var (
	// DB 全局数据库连接
	DB *gorm.DB
)

// NewDB 创建数据库连接
func NewDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	var err error

	// 连接数据库
	dsn := cfg.MySQLDSN()
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 配置连接池
	sqlDB, err := database.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// 设置连接最大存活时间
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Second)
	// 设置连接超时时间
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Timeout) * time.Second)

	fmt.Println("Database connected successfully")
	return database, nil
}

// NewDBCloser 创建数据库关闭函数
func NewDBCloser(db *gorm.DB) (func() error, error) {
	return func() error {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}, nil
}

// Init 初始化全局数据库连接（保留向后兼容）
func Init(cfg *config.DatabaseConfig) error {
	db, err := NewDB(cfg)
	if err != nil {
		return err
	}
	DB = db
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
