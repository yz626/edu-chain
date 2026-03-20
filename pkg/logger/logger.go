package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogLevel 日志级别
type LogLevel string

const (
	// LevelDebug 调试级别
	LevelDebug LogLevel = "debug"
	// LevelInfo 信息级别
	LevelInfo LogLevel = "info"
	// LevelWarn 警告级别
	LevelWarn LogLevel = "warn"
	// LevelError 错误级别
	LevelError LogLevel = "error"
	// LevelFatal 致命错误级别
	LevelFatal LogLevel = "fatal"
)

// LogType 日志类型
type LogType string

const (
	// TypeSystem 系统日志
	TypeSystem LogType = "system"
	// TypeBusiness 业务日志
	TypeBusiness LogType = "business"
	// TypeAudit 审计日志
	TypeAudit LogType = "audit"
	// TypeAccess 访问日志
	TypeAccess LogType = "access"
)

// Logger 日志器接口
type Logger interface {
	// Debug 调试日志
	Debug(msg string, fields ...Field)
	// Info 信息日志
	Info(msg string, fields ...Field)
	// Warn 警告日志
	Warn(msg string, fields ...Field)
	// Error 错误日志
	Error(msg string, fields ...Field)
	// Fatal 致命错误日志
	Fatal(msg string, fields ...Field)

	// Debugf 格式化调试日志
	Debugf(format string, args ...interface{})
	// Infof 格式化信息日志
	Infof(format string, args ...interface{})
	// Warnf 格式化警告日志
	Warnf(format string, args ...interface{})
	// Errorf 格式化错误日志
	Errorf(format string, args ...interface{})
	// Fatalf 格式化致命错误日志
	Fatalf(format string, args ...interface{})

	// With 创建带有上下文的日志器
	With(fields ...Field) Logger

	// Sync 同步日志缓冲区
	Sync() error
}

// Field 日志字段
type Field struct {
	Key   string
	Value interface{}
}

// StringField 创建字符串字段
func StringField(key, value string) Field {
	return Field{Key: key, Value: value}
}

// IntField 创建整数字段
func IntField(key string, value int) Field {
	return Field{Key: key, Value: value}
}

// Int64Field 创建64位整数字段
func Int64Field(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

// BoolField 创建布尔字段
func BoolField(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

// AnyField 创建任意类型字段
func AnyField(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}

// logger Zap日志器实现
type logger struct {
	zap         *zap.Logger
	logType     LogType
	fields      []Field
	level       zapcore.Level
	development bool
}

// 全局日志器
var (
	defaultLogger  Logger
	systemLogger   Logger
	businessLogger Logger
	auditLogger    Logger
	accessLogger   Logger

	loggerMutex sync.RWMutex
)

// Init 初始化日志系统
func Init(cfg *Config) error {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	// 创建默认日志器
	defaultLogger = newLogger(cfg, TypeSystem)

	// 创建各类专用日志器
	systemLogger = newLogger(cfg, TypeSystem)
	businessLogger = newLogger(cfg, TypeBusiness)
	auditLogger = newLogger(cfg, TypeAudit)
	accessLogger = newLogger(cfg, TypeAccess)

	return nil
}

// newLogger 创建新的日志器
func newLogger(cfg *Config, logType LogType) *logger {
	// 获取日志级别
	level := getZapLevel(cfg.GetLevel())

	// 创建编码器
	encoder := newEncoder(cfg, logType)

	// 创建写入器
	writer := newWriter(cfg, logType)

	// 创建核心
	core := zapcore.NewCore(
		encoder,
		writer,
		level,
	)

	// 创建Zap日志器
	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	if cfg.Development {
		zapLogger = zap.New(core, zap.Development(), zap.AddCaller(), zap.AddCallerSkip(2))
	}

	return &logger{
		zap:         zapLogger,
		logType:     logType,
		level:       level,
		development: cfg.Development,
	}
}

// newEncoder 创建编码器
func newEncoder(cfg *Config, logType LogType) zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 审计日志使用JSON格式
	if cfg.Format == "json" || logType == TypeAudit {
		return zapcore.NewJSONEncoder(encoderConfig)
	}

	// 开发环境使用Console格式
	if cfg.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewJSONEncoder(encoderConfig)
}

// newWriter 创建写入器
func newWriter(cfg *Config, logType LogType) zapcore.WriteSyncer {
	// 确定日志文件路径
	logPath := cfg.Directory
	if logPath == "" {
		logPath = "logs"
	}

	// 确保目录存在
	if err := os.MkdirAll(logPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create log directory: %v\n", err)
	}

	// 根据日志类型确定文件名
	var filename string
	switch logType {
	case TypeSystem:
		filename = "system.log"
	case TypeBusiness:
		filename = "business.log"
	case TypeAudit:
		filename = "audit.log"
	case TypeAccess:
		filename = "access.log"
	default:
		filename = "app.log"
	}

	// 创建日志轮转写入器
	writer := &lumberjack.Logger{
		Filename:   filepath.Join(logPath, filename),
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	// 同时输出到控制台（开发环境）
	if cfg.Console || cfg.Development {
		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(writer),
		)
	}

	return zapcore.AddSync(writer)
}

// getZapLevel 转换日志级别
func getZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	case LevelFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Debug 调试日志
func (l *logger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, l.toZapFields(fields...)...)
}

// Info 信息日志
func (l *logger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, l.toZapFields(fields...)...)
}

// Warn 警告日志
func (l *logger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, l.toZapFields(fields...)...)
}

// Error 错误日志
func (l *logger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, l.toZapFields(fields...)...)
}

// Fatal 致命错误日志
func (l *logger) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, l.toZapFields(fields...)...)
}

// Debugf 格式化调试日志
func (l *logger) Debugf(format string, args ...interface{}) {
	l.Debug(fmt.Sprintf(format, args...))
}

// Infof 格式化信息日志
func (l *logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Warnf 格式化警告日志
func (l *logger) Warnf(format string, args ...interface{}) {
	l.Warn(fmt.Sprintf(format, args...))
}

// Errorf 格式化错误日志
func (l *logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// Fatalf 格式化致命错误日志
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.Fatal(fmt.Sprintf(format, args...))
}

// With 创建带有上下文的日志器
func (l *logger) With(fields ...Field) Logger {
	newFields := append(l.fields, fields...)
	newZap := l.zap.With(l.toZapFields(newFields...)...)
	return &logger{
		zap:         newZap,
		logType:     l.logType,
		fields:      newFields,
		level:       l.level,
		development: l.development,
	}
}

// Sync 同步日志缓冲区
func (l *logger) Sync() error {
	return l.zap.Sync()
}

// toZapFields 转换字段
func (l *logger) toZapFields(fields ...Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)+len(l.fields))

	// 添加默认字段
	zapFields = append(zapFields, zap.String("type", string(l.logType)))

	// 添加上下文字段
	for _, field := range l.fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	// 添加额外字段
	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	return zapFields
}

// ===== 全局日志函数 =====

// Default 获取默认日志器
func Default() Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	if defaultLogger == nil {
		return &logger{zap: zap.NewNop()}
	}
	return defaultLogger
}

// System 获取系统日志器
func System() Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	if systemLogger == nil {
		return &logger{zap: zap.NewNop()}
	}
	return systemLogger
}

// Business 获取业务日志器
func Business() Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	if businessLogger == nil {
		return &logger{zap: zap.NewNop()}
	}
	return businessLogger
}

// Audit 获取审计日志器
func Audit() Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	if auditLogger == nil {
		return &logger{zap: zap.NewNop()}
	}
	return auditLogger
}

// Access 获取访问日志器
func Access() Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	if accessLogger == nil {
		return &logger{zap: zap.NewNop()}
	}
	return accessLogger
}

// Debug 调试日志
func Debug(msg string, fields ...Field) {
	Default().Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...Field) {
	Default().Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...Field) {
	Default().Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...Field) {
	Default().Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...Field) {
	Default().Fatal(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	Default().Debugf(format, args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Default().Infof(format, args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	Default().Warnf(format, args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	Default().Errorf(format, args...)
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, args ...interface{}) {
	Default().Fatalf(format, args...)
}

// With 创建带有上下文的日志器
func With(fields ...Field) Logger {
	return Default().With(fields...)
}

// Sync 同步所有日志器
func Sync() error {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()

	var errs []error
	if defaultLogger != nil {
		if err := defaultLogger.Sync(); err != nil {
			errs = append(errs, err)
		}
	}
	if systemLogger != nil {
		if err := systemLogger.Sync(); err != nil {
			errs = append(errs, err)
		}
	}
	if businessLogger != nil {
		if err := businessLogger.Sync(); err != nil {
			errs = append(errs, err)
		}
	}
	if auditLogger != nil {
		if err := auditLogger.Sync(); err != nil {
			errs = append(errs, err)
		}
	}
	if accessLogger != nil {
		if err := accessLogger.Sync(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

// GetCaller 获取调用者信息
func GetCaller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// LogRequest 记录HTTP请求
func LogRequest(logger Logger, method, path, clientIP, userAgent string, duration time.Duration, statusCode int) {
	logger.Info("HTTP Request",
		StringField("method", method),
		StringField("path", path),
		StringField("client_ip", clientIP),
		StringField("user_agent", userAgent),
		Int64Field("duration_ms", duration.Milliseconds()),
		IntField("status_code", statusCode),
	)
}

// LogError 记录错误
func LogError(logger Logger, err error, context string) {
	logger.Error(context,
		StringField("error", err.Error()),
	)
}

// ===== 审计日志辅助函数 =====

// LogUserLogin 记录用户登录
func LogUserLogin(logger Logger, userID, username, ip string, success bool) {
	level := "success"
	if !success {
		level = "failed"
	}
	logger.Info("User login",
		StringField("user_id", userID),
		StringField("username", username),
		StringField("ip", ip),
		StringField("result", level),
		StringField("action", "login"),
	)
}

// LogUserLogout 记录用户登出
func LogUserLogout(logger Logger, userID, username string) {
	logger.Info("User logout",
		StringField("user_id", userID),
		StringField("username", username),
		StringField("action", "logout"),
	)
}

// LogCertificateIssue 记录证书颁发
func LogCertificateIssue(logger Logger, userID, certID, orgID string) {
	logger.Info("Certificate issued",
		StringField("user_id", userID),
		StringField("cert_id", certID),
		StringField("org_id", orgID),
		StringField("action", "certificate_issue"),
	)
}

// LogCertificateVerify 记录证书验证
func LogCertificateVerify(logger Logger, userID, certID, verifierID string, result bool) {
	logger.Info("Certificate verified",
		StringField("user_id", userID),
		StringField("cert_id", certID),
		StringField("verifier_id", verifierID),
		BoolField("result", result),
		StringField("action", "certificate_verify"),
	)
}

// LogBlockchainTx 记录区块链交易
func LogBlockchainTx(logger Logger, txHash, from, to, method string, success bool) {
	level := "success"
	if !success {
		level = "failed"
	}
	logger.Info("Blockchain transaction",
		StringField("tx_hash", txHash),
		StringField("from", from),
		StringField("to", to),
		StringField("method", method),
		StringField("result", level),
		StringField("action", "blockchain_tx"),
	)
}

// ===== Context日志辅助 =====

// ContextWithRequestID 创建带有请求ID的上下文
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "request_id", requestID)
}

// ContextWithUserID 创建带有用户ID的上下文
func ContextWithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, "user_id", userID)
}

// ContextWithTraceID 创建带有追踪ID的上下文
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, "trace_id", traceID)
}

// GetRequestID 从上下文获取请求ID
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value("request_id").(string); ok {
		return id
	}
	return ""
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx context.Context) string {
	if id, ok := ctx.Value("user_id").(string); ok {
		return id
	}
	return ""
}

// GetTraceID 从上下文获取追踪ID
func GetTraceID(ctx context.Context) string {
	if id, ok := ctx.Value("trace_id").(string); ok {
		return id
	}
	return ""
}

// LogWithContext 记录带上下文的日志
func LogWithContext(ctx context.Context, logger Logger, level string, msg string, fields ...Field) {
	// 添加上下文字段
	contextFields := fields

	// 添加请求ID
	if requestID := GetRequestID(ctx); requestID != "" {
		contextFields = append(contextFields, StringField("request_id", requestID))
	}

	// 添加用户ID
	if userID := GetUserID(ctx); userID != "" {
		contextFields = append(contextFields, StringField("user_id", userID))
	}

	// 添加追踪ID
	if traceID := GetTraceID(ctx); traceID != "" {
		contextFields = append(contextFields, StringField("trace_id", traceID))
	}

	// 根据级别记录日志
	switch level {
	case "debug":
		logger.Debug(msg, contextFields...)
	case "info":
		logger.Info(msg, contextFields...)
	case "warn":
		logger.Warn(msg, contextFields...)
	case "error":
		logger.Error(msg, contextFields...)
	case "fatal":
		logger.Fatal(msg, contextFields...)
	}
}

// ===== JSON日志格式 =====

// JSONLog JSON格式日志
type JSONLog struct {
	Time      string      `json:"time"`
	Level     string      `json:"level"`
	Type      string      `json:"type"`
	Message   string      `json:"msg"`
	Caller    string      `json:"caller,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	UserID    string      `json:"user_id,omitempty"`
	TraceID   string      `json:"trace_id,omitempty"`
	Fields    interface{} `json:"fields,omitempty"`
}

// ToJSON 转换为JSON字符串
func (j *JSONLog) ToJSON() (string, error) {
	data, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ===== 性能监控 =====

// Timer 日志计时器
type Timer struct {
	start   time.Time
	logger  Logger
	message string
	fields  []Field
}

// StartTimer 开始计时
func StartTimer(logger Logger, message string, fields ...Field) *Timer {
	return &Timer{
		start:   time.Now(),
		logger:  logger,
		message: message,
		fields:  fields,
	}
}

// End 结束计时并记录日志
func (t *Timer) End() {
	duration := time.Since(t.start)
	t.logger.Info(t.message,
		append(t.fields, Int64Field("duration_ms", duration.Milliseconds()))...,
	)
}

// EndWithLevel 按指定级别结束计时
func (t *Timer) EndWithLevel(level string) {
	duration := time.Since(t.start)

	fields := append(t.fields, Int64Field("duration_ms", duration.Milliseconds()))

	switch level {
	case "debug":
		t.logger.Debug(t.message, fields...)
	case "info":
		t.logger.Info(t.message, fields...)
	case "warn":
		t.logger.Warn(t.message, fields...)
	case "error":
		t.logger.Error(t.message, fields...)
	}
}
