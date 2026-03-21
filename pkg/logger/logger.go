package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/wire"
	"github.com/yz626/edu-chain/config"
	"gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ProviderSet = wire.NewSet(NewLogger)

var (
	// globalLogger 全局日志实例
	globalLogger *zap.Logger
	// once 用于保证全局日志只初始化一次
	once sync.Once
)

// Logger 日志封装结构体
type Logger struct {
	*zap.Logger
}

// NewLogger 初始化全局日志
func NewLogger(cfg *config.LoggerConfig) (*Logger, error) {
	once.Do(func() {
		logger, err := newLogger(FromConfig(cfg))
		if err != nil {
			return
		}
		globalLogger = logger
	})
	return &Logger{globalLogger}, nil
}

// newLogger 创建新的日志实例
func newLogger(cfg *Config) (*zap.Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(cfg.Directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 根据格式选择编码器
	var encoder zapcore.Encoder
	if cfg.isJSONFormat() {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 日志级别
	level := cfg.getLevel()

	// 创建写入器
	writers := []zapcore.WriteSyncer{}

	// 控制台输出
	if cfg.Console {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}

	// 文件输出
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Directory, "app.log"),
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	})
	writers = append(writers, fileWriter)

	// 创建核心
	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		level,
	)

	// 创建Logger
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 如果启用堆栈跟踪
	if cfg.EnableStacktrace {
		logger = logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return logger, nil
}

// GetLogger 获取全局日志实例
func GetLogger() *Logger {
	if globalLogger == nil {
		globalLogger, _ = newLogger(DefaultConfig())
	}
	return &Logger{globalLogger}
}

// With 添加上下文字段
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

// Named 创建带名称的日志实例
func (l *Logger) Named(name string) *Logger {
	return &Logger{l.Logger.Named(name)}
}

// ========== 全局便捷方法 ==========

// Debug 调试级别日志
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	GetLogger().Sugar().Debugf(format, args...)
}

// Info 信息级别日志
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	GetLogger().Sugar().Infof(format, args...)
}

// Warn 警告级别日志
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	GetLogger().Sugar().Warnf(format, args...)
}

// Error 错误级别日志
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	GetLogger().Sugar().Errorf(format, args...)
}

// Fatal 致命错误级别日志
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	GetLogger().Sugar().Fatalf(format, args...)
}

// Sync 同步日志缓冲
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// ========== 便捷字段方法 ==========

// String 字符串字段
func String(key, value string) zap.Field {
	return zap.String(key, value)
}

// Int 整数字段
func Int(key string, value int) zap.Field {
	return zap.Int(key, value)
}

// Int64 64位整数字段
func Int64(key string, value int64) zap.Field {
	return zap.Int64(key, value)
}

// Bool 布尔字段
func Bool(key string, value bool) zap.Field {
	return zap.Bool(key, value)
}

// Error 错误字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Any 任意类型字段
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
