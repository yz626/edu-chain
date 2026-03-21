package db

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"

	"github.com/google/wire"
	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/pkg/logger"
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
	globalDB *gorm.DB
)

type DateDB struct {
	DB  *gorm.DB
	cfg *config.DatabaseConfig
	log *logger.Logger
}

// NewDB 创建数据库连接
func NewDB(cfg *config.DatabaseConfig, log *logger.Logger) (*DateDB, error) {
	var err error
	log = log.Named("db")

	// 连接数据库
	dsn := cfg.MySQLDSN()
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Info),
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

	log.With(logger.String("moudel", "data/db")).Info("数据库连接成功")
	return &DateDB{
		DB:  database,
		cfg: cfg,
		log: log,
	}, nil
}

// NewDBCloser 创建数据库关闭函数
func NewDBCloser(db *DateDB) (func() error, error) {
	return func() error {
		sqlDB, err := db.DB.DB()
		if err != nil {
			return err
		}
		db.log.With(logger.String("moudel", "data/db")).Info("关闭数据库连接")
		return sqlDB.Close()
	}, nil
}

// Init 初始化全局数据库连接（保留向后兼容）
func Init(cfg *config.DatabaseConfig) error {
	db, err := NewDB(cfg, logger.GetLogger())
	if err != nil {
		return err
	}
	globalDB = db.DB
	return nil
}

// Close 关闭数据库连接
func Close() error {
	if globalDB != nil {
		sqlDB, err := globalDB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB 获取数据库连接实例
func GetDB() *DateDB {
	return &DateDB{
		DB:  globalDB,
		cfg: nil,
		log: logger.GetLogger().Named("globalDB"),
	}
}

func (*DateDB) GetGormDB() *gorm.DB {
	return globalDB
}
