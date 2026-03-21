package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	// ConfigPath 配置文件路径
	ConfigPath = "config/config.yaml"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `mapstructure:"host"`          // 服务器地址
	Port         int    `mapstructure:"port"`          // 服务器端口
	Mode         string `mapstructure:"mode"`          // 运行模式 (debug, release)
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 读取超时(秒)
	WriteTimeout int    `mapstructure:"write_timeout"` // 写入超时(秒)
}

// GRPCConfig gRPC服务器配置
type GRPCConfig struct {
	Host string `mapstructure:"host"` // gRPC服务器地址
	Port int    `mapstructure:"port"` // gRPC服务器端口
}

// Addr 获取服务器地址
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Name         string `mapstructure:"name"`           // 数据库名称
	Host         string `mapstructure:"host"`           // 数据库地址
	Port         int    `mapstructure:"port"`           // 数据库端口
	Username     string `mapstructure:"username"`       // 数据库用户名
	Password     string `mapstructure:"password"`       // 数据库密码
	Database     string `mapstructure:"database"`       // 数据库名称
	SSLMode      string `mapstructure:"sslmode"`        // SSL模式
	MaxOpenConns int    `mapstructure:"max_open_conns"` // 最大连接数
	MaxIdleConns int    `mapstructure:"max_idle_conns"` // 最大空闲连接数
	MaxLifetime  int    `mapstructure:"max_lifetime"`   // 连接超时时间
	Timeout      int    `mapstructure:"timeout"`        // 连接超时时间
}

// DSN 获取数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

// MySQLDSN 获取MySQL连接字符串
func (c *DatabaseConfig) MySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database,
	)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level            string `mapstructure:"level"`             // 日志级别
	Format           string `mapstructure:"format"`            // 日志格式
	Directory        string `mapstructure:"directory"`         // 日志目录
	Console          bool   `mapstructure:"console"`           // 是否输出到控制台
	MaxSize          int    `mapstructure:"max_size"`          // 单个日志文件最大大小(MB)
	MaxAge           int    `mapstructure:"max_age"`           // 日志文件保留天数
	MaxBackups       int    `mapstructure:"max_backups"`       // 保留的日志文件数量
	Compress         bool   `mapstructure:"compress"`          // 是否压缩旧日志
	EnableStacktrace bool   `mapstructure:"enable_stacktrace"` // 是否启用堆栈跟踪
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string `mapstructure:"secret"`         // JWT密钥
	Expire        int    `mapstructure:"expire"`         // 访问令牌过期时间（秒）
	RefreshExpire int    `mapstructure:"refresh_expire"` // 刷新令牌过期时间（秒）
	Issuer        string `mapstructure:"issuer"`         // JWT签发者
}

// Load 加载配置
func Load(path ...string) (*Config, error) {
	if len(path) > 1 {
		return nil, fmt.Errorf("too many config paths")
	}

	v := viper.New()
	setDefaults(v)
	v.AutomaticEnv()

	// 配置文件路径，先尝试读取环境变量中的配置文件路径
	// 如果未设置，尝试使用传入的参数
	// 如果未设置，使用默认配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		if len(path) > 0 {
			configPath = path[0]
		} else {
			configPath = ConfigPath
		}

	}

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", 30)  // 读取超时30秒
	v.SetDefault("server.write_timeout", 30) // 写入超时30秒

	v.SetDefault("grpc.host", "0.0.0.0")
	v.SetDefault("grpc.port", 9090)

	v.SetDefault("database.name", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.username", "root")
	v.SetDefault("database.password", "123456")
	v.SetDefault("database.name", "edu_chain")
	v.SetDefault("database.sslmode", "disable")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)

	// 日志配置默认值
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("logger.directory", "logs")
	v.SetDefault("logger.console", true)
	v.SetDefault("logger.max_size", 100)
	v.SetDefault("logger.max_age", 30)
	v.SetDefault("logger.max_backups", 10)
	v.SetDefault("logger.compress", true)
	v.SetDefault("logger.enable_stacktrace", false)

	// JWT配置默认值
	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.expire", 3600)           // 访问令牌1小时
	v.SetDefault("jwt.refresh_expire", 604800) // 刷新令牌7天
	v.SetDefault("jwt.issuer", "edu-chain")
}
