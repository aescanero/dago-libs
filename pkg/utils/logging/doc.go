// Package logging provides structured logging utilities using Go's standard slog package.
//
// The Logger type wraps slog.Logger with convenience methods for adding common fields
// like execution_id, node_id, and graph_id. It supports both text and JSON output formats.
//
// Example usage:
//
//	logger := logging.NewLogger(logging.LevelInfo, "json")
//	logger.Info("Starting execution", "graph_id", "my-graph")
//
//	// Add contextual fields
//	execLogger := logger.WithExecutionID("exec-123")
//	execLogger.Info("Node started", "node_id", "node-1")
package logging
