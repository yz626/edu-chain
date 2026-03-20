package logger

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

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
		{"unknown", zapcore.InfoLevel}, // 默认值
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
