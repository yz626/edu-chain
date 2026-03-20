package logger

import (
	"testing"
)

// TestNamed_Chaining 演示 Named 方法的链式调用
func TestNamed_Chaining(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	// 多次 Named 调用不会覆盖，而是继承
	log1 := GetLogger().Named("module1")
	log2 := log1.Named("module2")
	log3 := log2.Named("module3")

	log1.Info("log1 message") // 输出: module1
	log2.Info("log2 message") // 输出: module1.module2
	log3.Info("log3 message") // 输出: module1.module2.module3
}

// TestWith_Chaining 演示 With 方法的链式调用
func TestWith_Chaining(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	// 多次 With 调用不会覆盖，而是累加字段
	log1 := GetLogger().With(String("field1", "value1"))
	log2 := log1.With(String("field2", "value2"))
	log3 := log2.With(String("field3", "value3"))

	log1.Info("log1 message") // 包含 field1
	log2.Info("log2 message") // 包含 field1, field2
	log3.Info("log3 message") // 包含 field1, field2, field3
}

// TestNamed_With_Combination 演示 Named 和 With 组合使用
func TestNamed_With_Combination(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	// 组合使用：先 Named 后 With
	log := GetLogger().
		Named("user").
		With(String("user_id", "100")).
		With(String("action", "login"))

	log.Info("用户登录")
}

// TestWith_Override 演示 With 相同字段名的覆盖行为
func TestWith_Override(t *testing.T) {
	setupTestDir()
	defer cleanupTestDir()

	cfg := DefaultConfig()
	cfg.Directory = testLogDir
	cfg.Console = true

	Init(cfg)

	// 相同字段名：后面的值会覆盖前面的值
	log := GetLogger().
		With(String("name", "first")).
		With(String("name", "second"))

	log.Info("相同字段名后面的值会覆盖前面的值")
}
