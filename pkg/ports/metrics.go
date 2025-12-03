package ports

import (
	"context"
	"time"
)

// MetricsCollector defines the interface for collecting system metrics.
// For MVP, this is implemented using Prometheus.
type MetricsCollector interface {
	// Counter metrics - these increment over time

	// IncGraphsSubmitted increments the count of submitted graphs.
	IncGraphsSubmitted(labels map[string]string)

	// IncGraphsCompleted increments the count of completed graphs.
	IncGraphsCompleted(labels map[string]string)

	// IncGraphsFailed increments the count of failed graphs.
	IncGraphsFailed(labels map[string]string)

	// IncNodesExecuted increments the count of executed nodes.
	IncNodesExecuted(nodeType string, labels map[string]string)

	// IncNodesFailed increments the count of failed nodes.
	IncNodesFailed(nodeType string, labels map[string]string)

	// IncToolExecutions increments the count of tool executions.
	IncToolExecutions(toolName string, labels map[string]string)

	// IncToolFailures increments the count of tool failures.
	IncToolFailures(toolName string, labels map[string]string)

	// IncLLMCalls increments the count of LLM API calls.
	IncLLMCalls(model string, labels map[string]string)

	// IncLLMTokens increments the count of LLM tokens used.
	IncLLMTokens(model string, tokenType string, count int, labels map[string]string)

	// Gauge metrics - these can go up or down

	// SetWorkerCount sets the current number of workers for a node type.
	SetWorkerCount(nodeType string, count int)

	// SetQueueDepth sets the current depth of the execution queue.
	SetQueueDepth(queueName string, depth int)

	// SetActiveExecutions sets the number of currently active executions.
	SetActiveExecutions(count int)

	// Histogram metrics - these track distributions

	// ObserveGraphDuration records the duration of a graph execution.
	ObserveGraphDuration(duration time.Duration, labels map[string]string)

	// ObserveNodeDuration records the duration of a node execution.
	ObserveNodeDuration(nodeType string, duration time.Duration, labels map[string]string)

	// ObserveToolDuration records the duration of a tool execution.
	ObserveToolDuration(toolName string, duration time.Duration, labels map[string]string)

	// ObserveLLMLatency records the latency of an LLM API call.
	ObserveLLMLatency(model string, duration time.Duration, labels map[string]string)

	// ObserveQueueWaitTime records how long an execution waited in the queue.
	ObserveQueueWaitTime(duration time.Duration, labels map[string]string)

	// RecordGraphSubmitted records a graph submission (compatibility method).
	RecordGraphSubmitted(status string)

	// RecordGraphCompleted records a graph completion (compatibility method).
	RecordGraphCompleted(status string, duration time.Duration)

	// RecordNodeExecuted records a node execution (compatibility method).
	RecordNodeExecuted(status string, duration time.Duration)

	// RecordWorkerPoolStatus records worker pool status (compatibility method).
	RecordWorkerPoolStatus(idle, busy, stopped int)
}

// MetricsConfig contains configuration for metrics collection.
type MetricsConfig struct {
	// Enabled controls whether metrics collection is active.
	Enabled bool `json:"enabled"`

	// Port is the HTTP port for the metrics endpoint.
	Port int `json:"port"`

	// Path is the HTTP path for the metrics endpoint (e.g., "/metrics").
	Path string `json:"path"`

	// Namespace is a prefix for all metric names.
	Namespace string `json:"namespace"`

	// Subsystem is a secondary prefix for metric names.
	Subsystem string `json:"subsystem"`
}

// HealthCheck represents a health check result.
type HealthCheck struct {
	// Name is the identifier for this health check.
	Name string `json:"name"`

	// Status indicates whether the component is healthy.
	Status HealthStatus `json:"status"`

	// Message provides additional details.
	Message string `json:"message,omitempty"`

	// LastChecked is when this health check was last performed.
	LastChecked time.Time `json:"last_checked"`
}

// HealthStatus represents the health status of a component.
type HealthStatus string

const (
	// HealthStatusHealthy indicates the component is functioning normally.
	HealthStatusHealthy HealthStatus = "healthy"

	// HealthStatusDegraded indicates the component is functioning but with issues.
	HealthStatusDegraded HealthStatus = "degraded"

	// HealthStatusUnhealthy indicates the component is not functioning.
	HealthStatusUnhealthy HealthStatus = "unhealthy"
)

// HealthChecker defines the interface for health checking.
type HealthChecker interface {
	// Check performs a health check and returns the result.
	Check(ctx context.Context) HealthCheck

	// Name returns the name of this health check.
	Name() string
}

// HealthRegistry manages multiple health checkers.
type HealthRegistry interface {
	// Register adds a health checker.
	Register(checker HealthChecker) error

	// Unregister removes a health checker.
	Unregister(name string) error

	// CheckAll runs all registered health checks.
	CheckAll(ctx context.Context) []HealthCheck

	// Check runs a specific health check by name.
	Check(ctx context.Context, name string) (*HealthCheck, error)
}
