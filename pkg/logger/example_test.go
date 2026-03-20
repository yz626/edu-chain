package logger

import (
	"context"
	"testing"
	"time"
)

// 示例：如何在业务代码中使用日志系统

func ExampleLogger() {
	// 初始化日志系统（通常在main函数中调用）
	cfg := DefaultConfig()
	// cfg := DevelopmentConfig()
	// cfg := ProductionConfig()
	_ = Init(cfg)

	// ===== 基础日志使用 =====

	// 使用默认日志器
	Debug("Debug message")
	Info("Info message")
	Warn("Warn message")
	Error("Error message")
	Fatal("Fatal message")

	// 使用格式化日志
	Debugf("Debug: %s", "message")
	Infof("Info: %s", "message")
	Warnf("Warn: %s", "message")
	Errorf("Error: %s", "message")
	Fatalf("Fatal: %s", "message")

	// ===== 使用带字段的日志 =====

	// 使用字段记录结构化日志
	Info("User login",
		StringField("username", "john"),
		StringField("ip", "192.168.1.1"),
		IntField("age", 25),
	)

	// 使用AnyField记录任意类型
	Info("Request details",
		AnyField("headers", map[string]string{
			"Content-Type": "application/json",
		}),
	)

	// ===== 使用专用日志器 =====

	// 系统日志 - 记录系统运行状态
	System().Info("Server started",
		StringField("host", "localhost"),
		IntField("port", 8080),
	)

	// 业务日志 - 记录业务操作
	Business().Info("Certificate issued",
		StringField("cert_id", "cert_12345"),
		StringField("user_id", "user_001"),
	)

	// 审计日志 - 记录用户操作（合规要求）
	Audit().Info("User action",
		StringField("user_id", "user_001"),
		StringField("action", "certificate_issue"),
		StringField("resource", "certificate/cert_12345"),
	)

	// 访问日志 - 记录HTTP请求
	Access().Info("HTTP Request",
		StringField("method", "POST"),
		StringField("path", "/api/v1/certificates"),
		StringField("client_ip", "192.168.1.100"),
		IntField("status_code", 200),
		Int64Field("duration_ms", 150),
	)

	// ===== 使用带上下文的日志器 =====

	// 创建带上下文的日志器
	ctxLogger := With(
		StringField("request_id", "req_12345"),
		StringField("user_id", "user_001"),
	)

	ctxLogger.Info("Request processed")

	// ===== 使用Context记录日志 =====

	// 创建上下文
	ctx := ContextWithRequestID(context.Background(), "req_12345")
	ctx = ContextWithUserID(ctx, "user_001")

	// 使用上下文记录日志
	LogWithContext(ctx, Default(), "info", "Request completed")

	// ===== 使用计时器 =====

	// 开始计时
	timer := StartTimer(Default(), "Database query")
	time.Sleep(100 * time.Millisecond)
	timer.End() // 自动记录耗时

	// ===== 使用辅助函数 =====

	// 记录用户登录
	LogUserLogin(Audit(), "user_001", "john", "192.168.1.1", true)

	// 记录证书颁发
	LogCertificateIssue(Business(), "user_001", "cert_12345", "org_001")

	// 记录证书验证
	LogCertificateVerify(Audit(), "user_001", "cert_12345", "verifier_001", true)

	// 记录区块链交易
	LogBlockchainTx(System(), "0x123...", "from_address", "to_address", "issueCertificate", true)

	// 同步日志（通常在程序退出时调用）
	_ = Sync()
}

// 示例：测试日志功能
func TestLogger(t *testing.T) {
	// 初始化日志系统
	cfg := DefaultConfig()
	cfg.Console = true
	cfg.Development = true
	if err := Init(cfg); err != nil {
		t.Fatalf("Failed to init logger: %v", err)
	}

	// 测试基础日志
	Debug("Debug message")
	Info("Info message")
	Warn("Warn message")
	Error("Error message")

	// 测试带字段的日志
	Info("Test with fields",
		StringField("key", "value"),
		IntField("number", 42),
		BoolField("flag", true),
	)

	// 测试专用日志器
	System().Info("System log")
	Business().Info("Business log")
	Audit().Info("Audit log")
	Access().Info("Access log")

	// 测试格式化日志
	Infof("Formatted: %s, %d", "test", 123)

	// 测试带上下文的日志器
	logger := With(StringField("context", "value"))
	logger.Info("Context log")

	// 测试Context日志
	ctx := ContextWithRequestID(context.Background(), "test_req_id")
	LogWithContext(ctx, Default(), "info", "Context log")

	// 测试计时器
	timer := StartTimer(Default(), "Test timer")
	time.Sleep(10 * time.Millisecond)
	timer.End()

	// 同步日志
	if err := Sync(); err != nil {
		t.Logf("Sync error: %v", err)
	}
}

// 示例：业务代码中使用日志
func businessLogicExample() {
	_ = Init(DefaultConfig())

	// 模拟业务逻辑
	userID := "user_123"
	certID := "cert_456"

	// 1. 记录业务开始
	Business().Info("Starting certificate issue",
		StringField("user_id", userID),
		StringField("cert_id", certID),
	)

	// 2. 记录处理步骤
	Business().Info("Validating certificate data",
		StringField("user_id", userID),
		StringField("cert_id", certID),
		StringField("step", "validation"),
	)

	// 3. 记录区块链操作
	txHash := "0xabc123"
	LogBlockchainTx(Business(), txHash, "contract_address", "user_address", "issue", true)

	// 4. 记录业务完成
	LogCertificateIssue(Business(), userID, certID, "org_001")
}

// 示例：HTTP中间件中使用日志
func httpLoggingExample() {
	_ = Init(DefaultConfig())

	// 模拟HTTP请求日志
	method := "POST"
	path := "/api/v1/certificates"
	clientIP := "192.168.1.100"
	userAgent := "Mozilla/5.0"
	duration := 150 * time.Millisecond
	statusCode := 200

	// 记录访问日志
	LogRequest(Access(), method, path, clientIP, userAgent, duration, statusCode)

	// 记录业务日志（请求处理）
	Business().Info("Request processed",
		StringField("method", method),
		StringField("path", path),
		StringField("client_ip", clientIP),
		Int64Field("duration_ms", duration.Milliseconds()),
		IntField("status_code", statusCode),
	)
}

// 示例：审计日志
func auditLoggingExample() {
	_ = Init(DefaultConfig())

	// 记录用户登录
	LogUserLogin(Audit(), "user_001", "john.doe@example.com", "192.168.1.100", true)

	// 记录失败的登录
	LogUserLogin(Audit(), "", "john.doe@example.com", "192.168.1.100", false)

	// 记录用户登出
	LogUserLogout(Audit(), "user_001", "john.doe@example.com")

	// 记录敏感操作
	Audit().Info("Sensitive action",
		StringField("user_id", "user_001"),
		StringField("action", "delete"),
		StringField("resource", "certificate/cert_123"),
		AnyField("details", map[string]interface{}{
			"reason":      "User requested deletion",
			"approved_by": "admin_001",
		}),
	)
}
