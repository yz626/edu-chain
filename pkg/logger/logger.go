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
	// initOnce 保证 NewLogger 只初始化一次
	initOnce sync.Once
	// fallbackOnce 保证 GetLogger 懒初始化只执行一次
	fallbackOnce sync.Once
)

// Logger 日志封装结构体
type Logger struct {
	*zap.Logger
}

// NewLogger 初始化全局日志（由 Wire 注入，应在程序启动时调用一次）。
// 若初始化失败，错误将被正确返回，不会将 nil 包装后返回给调用方。
func NewLogger(cfg *config.LoggerConfig) (*Logger, error) {
	var initErr error
	initOnce.Do(func() {
		l, err := newLogger(FromConfig(cfg))
		if err != nil {
			initErr = err
			return
		}
		globalLogger = l
	})
	if initErr != nil {
		return nil, initErr
	}
	if globalLogger == nil {
		return nil, fmt.Errorf("logger: not initialized (NewLogger already called and failed)")
	}
	return &Logger{globalLogger}, nil
}

// newLogger 创建新的 zap.Logger 实例。
//
// 核心设计：
//   - 文件 writer 始终使用 JSON encoder，保证日志文件可被机器解析，不含 ANSI 颜色码。
//   - 控制台 writer 使用彩色 console encoder（仅当 cfg.Console=true 且格式为 console 时）；
//     若 cfg.Format=json，控制台也使用 JSON encoder。
//   - 两路 writer 通过 zapcore.NewTee 并行写入，互不干扰。
//   - CallerSkip=1：跳过 Logger 方法本身这一层，使 caller 指向业务调用方。
func newLogger(cfg *Config) (*zap.Logger, error) {
	// 确保日志目录存在
	if err := os.MkdirAll(cfg.Directory, 0755); err != nil {
		return nil, fmt.Errorf("logger: failed to create log directory: %w", err)
	}

	// 基础编码器配置（文件与控制台共用时间格式、caller 等设置）
	baseEncoderCfg := zapcore.EncoderConfig{
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

	level := cfg.getLevel()

	// ---------- 文件 Core：始终 JSON，无颜色 ----------
	fileEncoderCfg := baseEncoderCfg // 值拷贝，单独修改不影响控制台
	fileEncoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderCfg)

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Directory, "app.log"),
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
	})

	fileCore := zapcore.NewCore(fileEncoder, fileWriter, level)

	// ---------- 控制台 Core：仅当启用时添加 ----------
	cores := []zapcore.Core{fileCore}

	if cfg.Console {
		var consoleEncoder zapcore.Encoder
		if cfg.isJSONFormat() {
			// 配置要求 JSON 格式，控制台同样输出 JSON（无颜色）
			consoleEncoderCfg := baseEncoderCfg
			consoleEncoderCfg.EncodeLevel = zapcore.LowercaseLevelEncoder
			consoleEncoder = zapcore.NewJSONEncoder(consoleEncoderCfg)
		} else {
			// console 格式：使用彩色级别标识，方便开发调试
			consoleEncoderCfg := baseEncoderCfg
			consoleEncoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
			consoleEncoder = zapcore.NewConsoleEncoder(consoleEncoderCfg)
		}
		consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
		cores = append(cores, consoleCore)
	}

	// 用 Tee 合并所有 core，各路独立写入
	core := zapcore.NewTee(cores...)

	// CallerSkip 不在此处设置：由调用方根据封装层数自行追加。
	// - 直接使用 *Logger（l.Info 等）：skip=0，zap 内部已正确定位到调用方。
	// - 全局便捷函数（Info/Error 等）：通过 skipLogger() 额外 +1 跳过函数本身。
	logger := zap.New(core, zap.AddCaller())

	if cfg.EnableStacktrace {
		logger = logger.WithOptions(zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return logger, nil
}

// GetLogger 获取全局日志实例。
//
// 若 NewLogger 尚未被调用，则以默认配置做一次懒初始化（使用 sync.Once 保证并发安全）。
// 程序正常启动时应先调用 NewLogger；此函数主要用于测试和工具命令等简单场景。
func GetLogger() *Logger {
	fallbackOnce.Do(func() {
		if globalLogger == nil {
			l, _ := newLogger(DefaultConfig())
			globalLogger = l
		}
	})
	return &Logger{globalLogger}
}

// With 添加上下文字段，返回新的 Logger 实例（不修改原实例）
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

// Named 创建带名称的子 Logger（名称会追加到 logger 字段，以 "." 分隔）
func (l *Logger) Named(name string) *Logger {
	return &Logger{l.Logger.Named(name)}
}

// ========== 全局便捷方法 ==========
//
// 说明：这些函数内部调用 GetLogger()，因此调用栈比直接使用 *Logger 多一层。
// 通过 WithOptions(zap.AddCallerSkip(1)) 额外跳过一层，使 caller 正确指向业务调用方。

// skipLogger 返回额外跳过一层 caller 的 zap.Logger，供全局便捷函数使用。
func skipLogger() *zap.Logger {
	return GetLogger().Logger.WithOptions(zap.AddCallerSkip(1))
}

// Debug 调试级别日志
func Debug(msg string, fields ...zap.Field) {
	skipLogger().Debug(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	skipLogger().Sugar().Debugf(format, args...)
}

// Info 信息级别日志
func Info(msg string, fields ...zap.Field) {
	skipLogger().Info(msg, fields...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	skipLogger().Sugar().Infof(format, args...)
}

// Warn 警告级别日志
func Warn(msg string, fields ...zap.Field) {
	skipLogger().Warn(msg, fields...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	skipLogger().Sugar().Warnf(format, args...)
}

// Error 错误级别日志
func Error(msg string, fields ...zap.Field) {
	skipLogger().Error(msg, fields...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	skipLogger().Sugar().Errorf(format, args...)
}

// Fatal 致命错误级别日志（会调用 os.Exit(1)）
func Fatal(msg string, fields ...zap.Field) {
	skipLogger().Fatal(msg, fields...)
}

// Fatalf 格式化致命错误日志（会调用 os.Exit(1)）
func Fatalf(format string, args ...interface{}) {
	skipLogger().Sugar().Fatalf(format, args...)
}

// Sync 刷新所有缓冲日志（程序退出前应调用）
func Sync() error {
	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}

// ========== 便捷字段构造函数 ==========

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

// Err 错误字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Any 任意类型字段
func Any(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}
