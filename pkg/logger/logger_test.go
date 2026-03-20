package logger

import (
	"os"
	"path/filepath"
	"testing"
)

// 测试目录
var testLogDir = "test_logs"

func setupTestDir() {
	os.MkdirAll(testLogDir, 0755)
}

func cleanupTestDir() {
	os.RemoveAll(testLogDir)
}

func TestLogger_Init(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := &Config{
		Level:      "info",
		Format:     "console",
		Directory:  testLogDir,
		Console:    true,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Compress:   false,
	}

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 验证日志可以正常工作
	log := GetLogger()
	if log == nil {
		t.Fatal("GetLogger returned nil")
	}

	// 测试日志输出
	log.Info("Test info message")
	log.Debug("Test debug message")
	log.Warn("Test warn message")
	log.Error("Test error message")
}

func TestLogger_InitJSONFormat(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := &Config{
		Level:      "debug",
		Format:     "json",
		Directory:  testLogDir,
		Console:    false,
		MaxSize:    10,
		MaxAge:     1,
		MaxBackups: 1,
		Compress:   false,
	}

	err := Init(cfg)
	if err != nil {
		t.Fatalf("Init with JSON format failed: %v", err)
	}

	log := GetLogger()
	log.Info("JSON format test")
}

func TestLogger_Named(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	log := GetLogger()
	namedLog := log.Named("test-module")

	if namedLog == nil {
		t.Fatal("Named returned nil")
	}

	namedLog.Info("Named logger test")
}

func TestLogger_With(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	log := GetLogger()
	logWithCtx := log.With(String("request_id", "12345"), Int("user_id", 100))

	logWithCtx.Info("Logger with context")
}

func TestLogger_ConvenienceMethods(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true
	cfg.Level = "debug"

	Init(cfg)

	// 测试格式化方法
	Debugf("Debug format: %s", "test")
	Infof("Info format: %d", 123)
	Warnf("Warn format: %v", true)
	Errorf("Error format: %f", 3.14)
}

func TestLogger_FieldHelpers(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	// 测试各种字段类型
	log := GetLogger()
	log.Info("Field types test",
		String("string_key", "value"),
		Int("int_key", 42),
		Int64("int64_key", 1234567890),
		Bool("bool_key", true),
	)
}

func TestLogger_ErrorField(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	testErr := os.ErrNotExist
	log := GetLogger()
	log.Error("Error with field", Err(testErr))
}

func TestLogger_Sync(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	err := Sync()
	if err != nil {
		t.Logf("Sync error (may be expected): %v", err)
	}
}

func TestLogger_FileRotation(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := &Config{
		Level:      "info",
		Format:     "json",
		Directory:  testLogDir,
		Console:    false,
		MaxSize:    1, // 1MB
		MaxAge:     1,
		MaxBackups: 3,
		Compress:   true,
	}

	Init(cfg)

	log := GetLogger()

	// 写入大量日志以触发轮转
	for i := 0; i < 1000; i++ {
		log.Info("Log message for rotation test",
			Int("index", i),
			String("data", "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"),
		)
	}

	Sync()

	// 检查是否创建了多个日志文件
	files, err := filepath.Glob(filepath.Join(testLogDir, "*.log*"))
	if err != nil {
		t.Logf("Error listing files: %v", err)
	}

	t.Logf("Log files found: %v", files)
}

func TestLogger_LevelFiltering(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := &Config{
		Level:     "warn",
		Format:    "console",
		Directory: testLogDir,
		Console:   true,
	}

	Init(cfg)

	log := GetLogger()

	// 这些应该被过滤掉
	log.Debug("This debug should be hidden")
	log.Info("This info should be hidden")

	// 这些应该被输出
	log.Warn("This warn should be visible")
	log.Error("This error should be visible")
}

func TestLogger_Concurrent(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = false // 减少并发输出干扰

	Init(cfg)

	done := make(chan bool, 10)

	// 并发写入日志
	for i := 0; i < 10; i++ {
		go func(idx int) {
			log := GetLogger().Named("goroutine")
			for j := 0; j < 100; j++ {
				log.Info("Concurrent log",
					Int("goroutine", idx),
					Int("iteration", j),
				)
			}
			done <- true
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < 10; i++ {
		<-done
	}

	Sync()
}

func TestLogger_EmptyFields(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	log := GetLogger()

	// 测试空字段
	log.Info("Empty fields test", String("", "empty key"))
	log.Info("Nil fields test")
}

func TestLogger_AnyField(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	log := GetLogger()

	// 测试 Any 类型字段
	mapVal := map[string]int{"a": 1, "b": 2}
	sliceVal := []int{1, 2, 3}

	log.Info("Any field test",
		Any("map", mapVal),
		Any("slice", sliceVal),
		Any("struct", struct{ Name string }{Name: "test"}),
	)
}
