package ports

import (
	"context"
	"time"
)

// WorkerType represents the type of worker.
type WorkerType string

const (
	// WorkerTypeExecutor is a worker that executes executor nodes.
	WorkerTypeExecutor WorkerType = "executor"

	// WorkerTypeRouter is a worker that executes router nodes.
	WorkerTypeRouter WorkerType = "router"
)

// WorkerStatus represents the current status of a worker.
type WorkerStatus string

const (
	// WorkerStatusIdle means the worker is running but not processing any task.
	WorkerStatusIdle WorkerStatus = "idle"

	// WorkerStatusBusy means the worker is currently processing a task.
	WorkerStatusBusy WorkerStatus = "busy"

	// WorkerStatusUnhealthy means the worker missed heartbeat(s).
	WorkerStatusUnhealthy WorkerStatus = "unhealthy"

	// WorkerStatusStopped means the worker has been explicitly stopped.
	WorkerStatusStopped WorkerStatus = "stopped"
)

// WorkerInfo contains information about a registered worker.
type WorkerInfo struct {
	// ID is the unique identifier for this worker.
	ID string `json:"id"`

	// Type is the type of worker (executor or router).
	Type WorkerType `json:"type"`

	// Status is the current status of the worker.
	Status WorkerStatus `json:"status"`

	// RegisteredAt is when the worker was first registered.
	RegisteredAt time.Time `json:"registered_at"`

	// LastHeartbeat is the timestamp of the last heartbeat received.
	LastHeartbeat time.Time `json:"last_heartbeat"`

	// CurrentTask is the ID of the task currently being processed (if any).
	CurrentTask string `json:"current_task,omitempty"`

	// PendingTasks is the number of tasks pending in this worker's queue.
	PendingTasks int `json:"pending_tasks"`

	// Version is the worker's software version.
	Version string `json:"version,omitempty"`

	// Metadata contains additional worker-specific information.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// WorkerFilter defines criteria for filtering workers.
type WorkerFilter struct {
	// Types filters workers by type. If empty, all types are included.
	Types []WorkerType `json:"types,omitempty"`

	// Statuses filters workers by status. If empty, all statuses are included.
	Statuses []WorkerStatus `json:"statuses,omitempty"`

	// HealthyOnly if true, only returns workers with recent heartbeats.
	HealthyOnly bool `json:"healthy_only,omitempty"`
}

// WorkerRegistry defines the interface for managing worker registration and heartbeats.
// This interface is transport-agnostic and can be implemented using Redis, Kafka,
// WebSockets, database, or any other mechanism.
type WorkerRegistry interface {
	// Register registers a new worker in the system.
	// This should be called when a worker starts up.
	Register(ctx context.Context, worker WorkerInfo) error

	// Unregister removes a worker from the registry.
	// This should be called when a worker shuts down gracefully.
	Unregister(ctx context.Context, workerID string) error

	// Heartbeat updates the last heartbeat timestamp for a worker.
	// This should be called periodically (e.g., every 5-10 seconds) to indicate
	// that the worker is still alive.
	Heartbeat(ctx context.Context, workerID string, status WorkerStatus, currentTask string) error

	// GetWorker retrieves information about a specific worker.
	GetWorker(ctx context.Context, workerID string) (*WorkerInfo, error)

	// ListWorkers retrieves all workers matching the filter criteria.
	ListWorkers(ctx context.Context, filter WorkerFilter) ([]WorkerInfo, error)

	// GetWorkerStats returns aggregate statistics about workers.
	GetWorkerStats(ctx context.Context, workerType WorkerType) (*WorkerStats, error)

	// CleanupStaleWorkers removes workers that haven't sent a heartbeat
	// within the specified timeout duration.
	// This is typically called periodically by the orchestrator.
	CleanupStaleWorkers(ctx context.Context, timeout time.Duration) (int, error)
}

// WorkerStats contains aggregate statistics about workers of a specific type.
type WorkerStats struct {
	// Type is the worker type these stats refer to.
	Type WorkerType `json:"type"`

	// TotalWorkers is the total number of registered workers.
	TotalWorkers int `json:"total_workers"`

	// IdleWorkers is the number of workers in idle status.
	IdleWorkers int `json:"idle_workers"`

	// BusyWorkers is the number of workers in busy status.
	BusyWorkers int `json:"busy_workers"`

	// UnhealthyWorkers is the number of workers in unhealthy status.
	UnhealthyWorkers int `json:"unhealthy_workers"`

	// TotalPendingTasks is the total number of pending tasks across all workers.
	TotalPendingTasks int `json:"total_pending_tasks"`
}
