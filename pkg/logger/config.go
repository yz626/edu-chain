package logger

// Config 日志配置
type Config struct {
	// Level 日志级别 (debug, info, warn, error, fatal)
	Level string `mapstructure:"level" json:"level"`
	// Format 日志格式 (json, console)
	Format string `mapstructure:"format" json:"format"`
	// Directory 日志目录
	Directory string `mapstructure:"directory" json:"directory"`
	// Console 是否输出到控制台
	Console bool `mapstructure:"console" json:"console"`
	// Development 是否为开发模式
	Development bool `mapstructure:"development" json:"development"`
	// MaxSize 单个日志文件最大大小(MB)
	MaxSize int `mapstructure:"max_size" json:"max_size"`
	// MaxAge 日志文件保留天数
	MaxAge int `mapstructure:"max_age" json:"max_age"`
	// MaxBackups 保留的日志文件数量
	MaxBackups int `mapstructure:"max_backups" json:"max_backups"`
	// Compress 是否压缩日志文件
	Compress bool `mapstructure:"compress" json:"compress"`
	// EnableStacktrace 是否启用堆栈跟踪
	EnableStacktrace bool `mapstructure:"enable_stacktrace" json:"enable_stacktrace"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:            "info",
		Format:           "json",
		Directory:        "logs",
		Console:          true,
		Development:      false,
		MaxSize:          100,
		MaxAge:           30,
		MaxBackups:       10,
		Compress:         true,
		EnableStacktrace: false,
	}
}

// DevelopmentConfig 返回开发环境配置
func DevelopmentConfig() *Config {
	return &Config{
		Level:            "debug",
		Format:           "console",
		Directory:        "logs",
		Console:          true,
		Development:      true,
		MaxSize:          10,
		MaxAge:           7,
		MaxBackups:       3,
		Compress:         false,
		EnableStacktrace: true,
	}
}

// ProductionConfig 返回生产环境配置
func ProductionConfig() *Config {
	return &Config{
		Level:            "info",
		Format:           "json",
		Directory:        "logs",
		Console:          false,
		Development:      false,
		MaxSize:          100,
		MaxAge:           30,
		MaxBackups:       10,
		Compress:         true,
		EnableStacktrace: false,
	}
}

// GetLevel 获取日志级别
func (c *Config) GetLevel() LogLevel {
	switch c.Level {
	case "debug":
		return LevelDebug
	case "info":
		return LevelInfo
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "fatal":
		return LevelFatal
	default:
		return LevelInfo
	}
}

// IsDevelopment 判断是否为开发模式
func (c *Config) IsDevelopment() bool {
	return c.Development
}

// IsJSON 判断是否为JSON格式
func (c *Config) IsJSON() bool {
	return c.Format == "json"
}
