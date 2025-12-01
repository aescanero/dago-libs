package ports

import (
	"context"
	"time"

	"github.com/aescanero/dago-libs/pkg/domain/state"
)

// StateStorage defines the interface for persisting execution state.
// For MVP, this is implemented using Redis.
type StateStorage interface {
	// Save persists the state for an execution.
	Save(ctx context.Context, executionID string, state state.State) error

	// Load retrieves the state for an execution.
	Load(ctx context.Context, executionID string) (state.State, error)

	// Delete removes the state for an execution.
	Delete(ctx context.Context, executionID string) error

	// Exists checks if state exists for an execution.
	Exists(ctx context.Context, executionID string) (bool, error)

	// SetTTL sets a time-to-live for state data.
	// After the TTL expires, the state will be automatically deleted.
	SetTTL(ctx context.Context, executionID string, ttl time.Duration) error

	// List returns all execution IDs that have stored state.
	List(ctx context.Context) ([]string, error)
}

// GraphStorage defines the interface for persisting graph definitions.
type GraphStorage interface {
	// Save persists a graph definition.
	Save(ctx context.Context, graphID string, graphData []byte) error

	// Load retrieves a graph definition.
	Load(ctx context.Context, graphID string) ([]byte, error)

	// Delete removes a graph definition.
	Delete(ctx context.Context, graphID string) error

	// Exists checks if a graph definition exists.
	Exists(ctx context.Context, graphID string) (bool, error)

	// List returns all stored graph IDs.
	List(ctx context.Context) ([]string, error)

	// ListVersions returns all versions of a graph (if versioning is supported).
	ListVersions(ctx context.Context, graphName string) ([]string, error)
}

// ExecutionMetadata contains metadata about a graph execution.
type ExecutionMetadata struct {
	// ExecutionID is the unique identifier for this execution.
	ExecutionID string `json:"execution_id"`

	// GraphID is the ID of the graph being executed.
	GraphID string `json:"graph_id"`

	// Status is the current status of the execution.
	Status ExecutionStatus `json:"status"`

	// StartedAt is when the execution started.
	StartedAt time.Time `json:"started_at"`

	// CompletedAt is when the execution completed (if finished).
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	// CurrentNodeID is the ID of the currently executing node.
	CurrentNodeID string `json:"current_node_id,omitempty"`

	// Error contains error information if the execution failed.
	Error string `json:"error,omitempty"`

	// Metadata contains additional execution-specific data.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ExecutionStatus represents the status of a graph execution.
type ExecutionStatus string

const (
	// ExecutionStatusPending indicates the execution is queued but not started.
	ExecutionStatusPending ExecutionStatus = "pending"

	// ExecutionStatusRunning indicates the execution is in progress.
	ExecutionStatusRunning ExecutionStatus = "running"

	// ExecutionStatusCompleted indicates the execution finished successfully.
	ExecutionStatusCompleted ExecutionStatus = "completed"

	// ExecutionStatusFailed indicates the execution failed with an error.
	ExecutionStatusFailed ExecutionStatus = "failed"

	// ExecutionStatusCancelled indicates the execution was cancelled.
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
)

// ExecutionStorage defines the interface for persisting execution metadata.
type ExecutionStorage interface {
	// Save persists execution metadata.
	Save(ctx context.Context, metadata ExecutionMetadata) error

	// Load retrieves execution metadata.
	Load(ctx context.Context, executionID string) (*ExecutionMetadata, error)

	// UpdateStatus updates the status of an execution.
	UpdateStatus(ctx context.Context, executionID string, status ExecutionStatus) error

	// List returns all execution metadata, optionally filtered by status.
	List(ctx context.Context, status *ExecutionStatus) ([]ExecutionMetadata, error)

	// Delete removes execution metadata.
	Delete(ctx context.Context, executionID string) error
}
