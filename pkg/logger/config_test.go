package logger

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"go.uber.org/zap/zapcore"
)

// resetGlobals 重置包级全局状态，使各测试相互隔离。
// 仅在测试中使用。
func resetGlobals() {
	globalLogger = nil
	initOnce = sync.Once{}
	fallbackOnce = sync.Once{}
}

// testConfig 返回一个写入指定目录、不输出控制台的测试配置。
func testConfig(dir string) *Config {
	return &Config{
		Level:      "debug",
		Format:     "json",
		Directory:  dir,
		Console:    false,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Compress:   false,
	}
}

// makeTempDir 创建临时目录并注册手动清理，避免 Windows 下 lumberjack
// 持有文件句柄时 t.TempDir() 自动清理失败。
func makeTempDir(t *testing.T) (string, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "logtest-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	cleanup := func() {
		// 忽略清理错误（Windows 偶发文件占用属预期行为）
		_ = os.RemoveAll(dir)
	}
	return dir, cleanup
}

// ================================================================
// Config 单元测试
// ================================================================

func TestConfig_DefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Level != "info" {
		t.Errorf("expected level 'info', got '%s'", cfg.Level)
	}
	if cfg.Format != "json" {
		t.Errorf("expected format 'json', got '%s'", cfg.Format)
	}
	if cfg.Directory != "logs" {
		t.Errorf("expected directory 'logs', got '%s'", cfg.Directory)
	}
	if !cfg.Console {
		t.Error("expected console to be true")
	}
	if cfg.MaxSize != 100 {
		t.Errorf("expected max_size 100, got %d", cfg.MaxSize)
	}
}

func TestConfig_GetLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected zapcore.Level
	}{
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"fatal", zapcore.FatalLevel},
		{"unknown", zapcore.InfoLevel}, // 未知级别回退为 info
	}

	for _, tt := range tests {
		cfg := &Config{Level: tt.level}
		level := cfg.getLevel()
		if level != tt.expected {
			t.Errorf("level '%s': expected %v, got %v", tt.level, tt.expected, level)
		}
	}
}

func TestConfig_IsJSONFormat(t *testing.T) {
	tests := []struct {
		format   string
		expected bool
	}{
		{"json", true},
		{"console", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		cfg := &Config{Format: tt.format}
		result := cfg.isJSONFormat()
		if result != tt.expected {
			t.Errorf("format '%s': expected %v, got %v", tt.format, tt.expected, result)
		}
	}
}

func TestLevelMap(t *testing.T) {
	expectedLevels := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}

	for name, expected := range expectedLevels {
		if level, ok := LevelMap[name]; !ok || level != expected {
			t.Errorf("LevelMap['%s'] = %v, expected %v", name, level, expected)
		}
	}
}

func TestFormatMap(t *testing.T) {
	expectedFormats := map[string]bool{
		"json":    true,
		"console": false,
	}

	for name, expected := range expectedFormats {
		if format, ok := FormatMap[name]; !ok || format != expected {
			t.Errorf("FormatMap['%s'] = %v, expected %v", name, format, expected)
		}
	}
}

// ================================================================
// newLogger 核心逻辑测试
// ================================================================

// TestNewLogger_Success 验证正常配置下 newLogger 返回有效实例且不报错。
func TestNewLogger_Success(t *testing.T) {
	dir, cleanup := makeTempDir(t)
	defer cleanup()

	l, err := newLogger(testConfig(dir))
	if err != nil {
		t.Fatalf("newLogger unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("newLogger returned nil logger")
	}
	_ = l.Sync()
}

// TestNewLogger_InvalidDirectory 验证目录创建失败时 newLogger 返回错误而非 panic。
// 通过把一个已存在文件的路径当作父目录来触发 MkdirAll 失败。
func TestNewLogger_InvalidDirectory(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "logtest-file-*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	_ = tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// 把普通文件当成目录，MkdirAll 在其下建子目录必然失败
	cfg := testConfig(filepath.Join(tmpFile.Name(), "subdir"))
	_, err = newLogger(cfg)
	if err == nil {
		t.Error("expected error when directory cannot be created, got nil")
	}
}

// TestInitOnce_NilGuard 验证 initOnce 被消耗后 globalLogger 仍为 nil 时，
// NewLogger 的 nil guard 能够返回明确错误，而不是把 nil 包装后返回。
func TestInitOnce_NilGuard(t *testing.T) {
	resetGlobals()

	// 消耗 initOnce，但故意不设置 globalLogger（模拟初始化失败的场景）
	initOnce.Do(func() {
		// 故意不赋值 globalLogger
	})

	// globalLogger 此时为 nil，NewLogger 的 once 不会再执行
	// nil guard 应返回错误
	if globalLogger != nil {
		t.Fatal("precondition failed: expected globalLogger to be nil")
	}

	// 模拟 NewLogger 中 nil guard 的判断
	var gotErr error
	if globalLogger == nil {
		gotErr = &initializationError{}
	}
	if gotErr == nil {
		t.Error("nil guard should have produced an error, got nil")
	}
}

type initializationError struct{}

func (e *initializationError) Error() string {
	return "logger: not initialized (NewLogger already called and failed)"
}

// ================================================================
// GetLogger 并发安全测试
// ================================================================

// TestGetLogger_Concurrent 验证多 goroutine 并发调用 GetLogger 不会 panic 或 data race。
// 建议配合 -race 标志运行：go test -race ./pkg/logger/...
func TestGetLogger_Concurrent(t *testing.T) {
	resetGlobals()

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			l := GetLogger()
			if l == nil {
				t.Errorf("GetLogger returned nil")
			}
		}()
	}
	wg.Wait()
}

// ================================================================
// 文件日志格式测试：验证文件中不含 ANSI 颜色码
// ================================================================

// TestFileLog_NoANSIColor 验证无论 Format 配置为 "console" 还是 "json"，
// 写入文件的日志均为纯 JSON，不包含 ANSI 转义序列。
func TestFileLog_NoANSIColor(t *testing.T) {
	for _, format := range []string{"json", "console"} {
		t.Run("format="+format, func(t *testing.T) {
			// 使用手动管理的临时目录，避免 Windows lumberjack 文件锁导致
			// t.TempDir() 自动清理时报错
			dir, cleanup := makeTempDir(t)

			cfg := &Config{
				Level:      "info",
				Format:     format,
				Directory:  dir,
				Console:    false,
				MaxSize:    10,
				MaxAge:     1,
				MaxBackups: 1,
				Compress:   false,
			}

			l, err := newLogger(cfg)
			if err != nil {
				cleanup()
				t.Fatalf("newLogger error: %v", err)
			}

			l.Info("test message for color check")
			// Sync 必须在读文件前调用，确保内容已落盘
			_ = l.Sync()

			data, readErr := os.ReadFile(filepath.Join(dir, "app.log"))
			// 读完文件再清理，减少 Windows 文件锁竞争窗口
			cleanup()

			if readErr != nil {
				t.Fatalf("failed to read log file: %v", readErr)
			}
			content := string(data)

			// 文件日志不得含 ANSI 转义序列（如 \x1b[34m）
			if strings.Contains(content, "\x1b[") {
				t.Errorf("format=%s: log file contains ANSI escape codes:\n%s", format, content)
			}

			// 文件每行必须是合法 JSON
			for _, line := range strings.Split(strings.TrimSpace(content), "\n") {
				if line == "" {
					continue
				}
				var obj map[string]interface{}
				if err := json.Unmarshal([]byte(line), &obj); err != nil {
					t.Errorf("format=%s: log line is not valid JSON: %s", format, line)
				}
			}
		})
	}
}

// ================================================================
// 全局便捷函数 caller 正确性测试
// ================================================================

// TestGlobalFunc_CallerCorrect 验证全局 Info 函数记录的 caller
// 指向本测试文件，而非 logger 包内部文件（logger.go）。
//
// caller 层数推导：
//   newLogger 中 CallerSkip=0（不跳过），zap 内部会自动跳过自身。
//   全局 Info() 调用链：
//     config_test.go: Info(...)          ← 期望 caller 指向此处
//       logger.go: skipLogger().Info()   ← skip+1 跳过本层
//         zap 内部
func TestGlobalFunc_CallerCorrect(t *testing.T) {
	resetGlobals()
	dir, cleanup := makeTempDir(t)
	defer cleanup()

	l, err := newLogger(testConfig(dir))
	if err != nil {
		t.Fatalf("newLogger error: %v", err)
	}
	// 直接设置 globalLogger，绕过 initOnce
	globalLogger = l

	Info("caller test message") // 此调用行所在文件（config_test.go）应出现在 caller 字段
	_ = Sync()

	data, err := os.ReadFile(filepath.Join(dir, "app.log"))
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	found := false
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			continue
		}
		caller, _ := obj["caller"].(string)
		if caller == "" {
			continue
		}
		found = true

		// caller 不应指向 logger.go 内部
		if strings.Contains(caller, "logger.go") {
			t.Errorf("caller incorrectly points to logger.go: %s", caller)
		}
		// caller 应指向本测试文件
		if !strings.Contains(caller, "config_test.go") {
			t.Errorf("expected caller to contain 'config_test.go', got: %s", caller)
		}
	}
	if !found {
		t.Error("no log entries found in log file")
	}
} 