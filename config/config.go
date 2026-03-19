package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
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

// Load 加载配置
func Load() (*Config, error) {
	v := viper.New()
	setDefaults(v)
	v.AutomaticEnv()

	// 配置文件路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// 默认配置文件路径 (相对于项目根目录)
		configPath = "config/config.yaml"
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

	v.SetDefault("database.name", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.username", "root")
	v.SetDefault("database.password", "123456")
	v.SetDefault("database.name", "edu_chain")
	v.SetDefault("database.sslmode", "disable")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
}
