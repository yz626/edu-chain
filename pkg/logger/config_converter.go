package logger

import (
	"github.com/yz626/edu-chain/config"
)

// FromConfig 从应用配置转换为日志配置
func FromConfig(cfg *config.LoggerConfig) *Config {
	return &Config{
		Level:            cfg.Level,
		Format:           cfg.Format,
		Directory:        cfg.Directory,
		Console:          cfg.Console,
		MaxSize:          cfg.MaxSize,
		MaxAge:           cfg.MaxAge,
		MaxBackups:       cfg.MaxBackups,
		Compress:         cfg.Compress,
		EnableStacktrace: cfg.EnableStacktrace,
	}
}
