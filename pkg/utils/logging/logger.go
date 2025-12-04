// Package logging provides structured logging utilities for DA Orchestrator.
package logging

import (
	"context"
	"log/slog"
	"os"
)

// LogLevel represents the severity level of a log message.
type LogLevel string

const (
	// LevelDebug is for detailed debugging information.
	LevelDebug LogLevel = "debug"

	// LevelInfo is for general informational messages.
	LevelInfo LogLevel = "info"

	// LevelWarn is for warning messages.
	LevelWarn LogLevel = "warn"

	// LevelError is for error messages.
	LevelError LogLevel = "error"
)

// Logger wraps slog.Logger with additional convenience methods.
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger with the specified level and format.
func NewLogger(level LogLevel, format string) *Logger {
	var slogLevel slog.Level
	switch level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: slogLevel,
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}

// NewDefaultLogger creates a logger with INFO level and text format.
func NewDefaultLogger() *Logger {
	return NewLogger(LevelInfo, "text")
}

// WithContext returns a logger with context values added.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		Logger: l.With(),
	}
}

// WithFields returns a logger with additional fields.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.With(args...),
	}
}

// WithField returns a logger with an additional field.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{
		Logger: l.With(key, value),
	}
}

// WithExecutionID returns a logger with the execution ID field.
func (l *Logger) WithExecutionID(executionID string) *Logger {
	return l.WithField("execution_id", executionID)
}

// WithNodeID returns a logger with the node ID field.
func (l *Logger) WithNodeID(nodeID string) *Logger {
	return l.WithField("node_id", nodeID)
}

// WithGraphID returns a logger with the graph ID field.
func (l *Logger) WithGraphID(graphID string) *Logger {
	return l.WithField("graph_id", graphID)
}

// LoggerConfig contains configuration for the logger.
type LoggerConfig struct {
	// Level is the minimum log level to output.
	Level LogLevel `json:"level"`

	// Format is the output format ("text" or "json").
	Format string `json:"format"`

	// AddSource adds source file and line number to log entries.
	AddSource bool `json:"add_source"`
}

// DefaultConfig returns a default logger configuration.
func DefaultConfig() LoggerConfig {
	return LoggerConfig{
		Level:     LevelInfo,
		Format:    "text",
		AddSource: false,
	}
}

// NewLoggerFromConfig creates a logger from a configuration.
func NewLoggerFromConfig(cfg LoggerConfig) *Logger {
	var slogLevel slog.Level
	switch cfg.Level {
	case LevelDebug:
		slogLevel = slog.LevelDebug
	case LevelInfo:
		slogLevel = slog.LevelInfo
	case LevelWarn:
		slogLevel = slog.LevelWarn
	case LevelError:
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     slogLevel,
		AddSource: cfg.AddSource,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return &Logger{
		Logger: slog.New(handler),
	}
}
