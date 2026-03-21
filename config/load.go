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
