package logging

import (
	"testing"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name   string
		level  LogLevel
		format string
	}{
		{"debug text", LevelDebug, "text"},
		{"info json", LevelInfo, "json"},
		{"warn text", LevelWarn, "text"},
		{"error json", LevelError, "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger(tt.level, tt.format)
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}
			if logger.Logger == nil {
				t.Error("internal slog.Logger is nil")
			}
		})
	}
}

func TestNewDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()
	if logger == nil {
		t.Fatal("NewDefaultLogger returned nil")
	}
	if logger.Logger == nil {
		t.Error("internal slog.Logger is nil")
	}
}

func TestLogger_WithField(t *testing.T) {
	logger := NewDefaultLogger()
	newLogger := logger.WithField("key", "value")

	if newLogger == nil {
		t.Fatal("WithField returned nil")
	}
	if newLogger.Logger == nil {
		t.Error("WithField returned logger with nil internal logger")
	}
}

func TestLogger_WithFields(t *testing.T) {
	logger := NewDefaultLogger()

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
		"key3": true,
	}

	newLogger := logger.WithFields(fields)

	if newLogger == nil {
		t.Fatal("WithFields returned nil")
	}
	if newLogger.Logger == nil {
		t.Error("WithFields returned logger with nil internal logger")
	}
}

func TestLogger_WithExecutionID(t *testing.T) {
	logger := NewDefaultLogger()
	newLogger := logger.WithExecutionID("exec-123")

	if newLogger == nil {
		t.Fatal("WithExecutionID returned nil")
	}
	if newLogger.Logger == nil {
		t.Error("WithExecutionID returned logger with nil internal logger")
	}
}

func TestLogger_WithNodeID(t *testing.T) {
	logger := NewDefaultLogger()
	newLogger := logger.WithNodeID("node-456")

	if newLogger == nil {
		t.Fatal("WithNodeID returned nil")
	}
	if newLogger.Logger == nil {
		t.Error("WithNodeID returned logger with nil internal logger")
	}
}

func TestLogger_WithGraphID(t *testing.T) {
	logger := NewDefaultLogger()
	newLogger := logger.WithGraphID("graph-789")

	if newLogger == nil {
		t.Fatal("WithGraphID returned nil")
	}
	if newLogger.Logger == nil {
		t.Error("WithGraphID returned logger with nil internal logger")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Level != LevelInfo {
		t.Errorf("expected default level %q, got %q", LevelInfo, cfg.Level)
	}
	if cfg.Format != "text" {
		t.Errorf("expected default format 'text', got %q", cfg.Format)
	}
	if cfg.AddSource {
		t.Error("expected AddSource to be false by default")
	}
}

func TestNewLoggerFromConfig(t *testing.T) {
	tests := []struct {
		name   string
		config LoggerConfig
	}{
		{
			name: "debug json with source",
			config: LoggerConfig{
				Level:     LevelDebug,
				Format:    "json",
				AddSource: true,
			},
		},
		{
			name: "info text without source",
			config: LoggerConfig{
				Level:     LevelInfo,
				Format:    "text",
				AddSource: false,
			},
		},
		{
			name: "error json",
			config: LoggerConfig{
				Level:  LevelError,
				Format: "json",
			},
		},
		{
			name: "invalid level defaults to info",
			config: LoggerConfig{
				Level:  "invalid",
				Format: "text",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLoggerFromConfig(tt.config)
			if logger == nil {
				t.Fatal("NewLoggerFromConfig returned nil")
			}
			if logger.Logger == nil {
				t.Error("internal slog.Logger is nil")
			}
		})
	}
}

func TestLogger_Chaining(t *testing.T) {
	// Test that chaining multiple WithX methods works
	logger := NewDefaultLogger().
		WithExecutionID("exec-123").
		WithGraphID("graph-456").
		WithNodeID("node-789").
		WithField("custom", "value")

	if logger == nil {
		t.Fatal("chained logger is nil")
	}
	if logger.Logger == nil {
		t.Error("chained logger has nil internal logger")
	}
}
