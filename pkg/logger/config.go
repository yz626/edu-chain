package logger

import (
	"go.uber.org/zap/zapcore"
)

// LevelMap 日志级别映射
var LevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

// FormatMap 日志格式映射
var FormatMap = map[string]bool{
	"json":    true,
	"console": false,
}

// Config 日志配置
type Config struct {
	Level            string // 日志级别 (debug, info, warn, error, fatal)
	Format           string // 日志格式 (json, console)
	Directory        string // 日志目录
	Console          bool   // 是否输出到控制台
	MaxSize          int    // 单个日志文件最大大小(MB)
	MaxAge           int    // 日志文件保留天数
	MaxBackups       int    // 保留的日志文件数量
	Compress         bool   // 是否压缩旧日志
	EnableStacktrace bool   // 是否启用堆栈跟踪
}

// getLevel 获取日志级别
func (c *Config) getLevel() zapcore.Level {
	if level, ok := LevelMap[c.Level]; ok {
		return level
	}
	return zapcore.InfoLevel
}

// isJSONFormat 是否JSON格式
func (c *Config) isJSONFormat() bool {
	return FormatMap[c.Format]
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:            "info",
		Format:           "json",
		Directory:        "logs",
		Console:          true,
		MaxSize:          100,
		MaxAge:           30,
		MaxBackups:       10,
		Compress:         true,
		EnableStacktrace: false,
	}
}
