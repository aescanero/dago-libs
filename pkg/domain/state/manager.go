package state

import "context"

// Manager defines the interface for managing state evolution during graph execution.
// Implementations handle state transitions, persistence, and versioning.
type Manager interface {
	// Initialize creates a new execution state with the given initial data.
	Initialize(ctx context.Context, executionID string, initialState State) error

	// GetState retrieves the current state for an execution.
	GetState(ctx context.Context, executionID string) (State, error)

	// UpdateState updates the state for an execution.
	// The update function receives the current state and should return the modified state.
	UpdateState(ctx context.Context, executionID string, updateFn func(State) (State, error)) error

	// DeleteState removes all state data for an execution.
	DeleteState(ctx context.Context, executionID string) error

	// SaveSnapshot creates a named snapshot of the current state.
	// This is useful for checkpointing during long-running executions.
	SaveSnapshot(ctx context.Context, executionID string, snapshotName string) error

	// LoadSnapshot restores state from a named snapshot.
	LoadSnapshot(ctx context.Context, executionID string, snapshotName string) (State, error)

	// ListSnapshots returns all snapshot names for an execution.
	ListSnapshots(ctx context.Context, executionID string) ([]string, error)
}

// Transition represents a state transition event.
type Transition struct {
	ExecutionID string
	NodeID      string
	FromState   State
	ToState     State
	Timestamp   int64
}

// TransitionLogger defines the interface for logging state transitions.
// This is useful for debugging, auditing, and replay functionality.
type TransitionLogger interface {
	// LogTransition records a state transition.
	LogTransition(ctx context.Context, transition Transition) error

	// GetTransitions retrieves all transitions for an execution.
	GetTransitions(ctx context.Context, executionID string) ([]Transition, error)

	// GetTransitionsSince retrieves transitions after a given timestamp.
	GetTransitionsSince(ctx context.Context, executionID string, since int64) ([]Transition, error)
}
