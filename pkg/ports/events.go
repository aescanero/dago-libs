package ports

import (
	"context"
	"time"
)

// EventType represents the type of event.
type EventType string

const (
	// EventTypeGraphStarted is emitted when graph execution begins.
	EventTypeGraphStarted EventType = "graph.started"

	// EventTypeGraphCompleted is emitted when graph execution completes successfully.
	EventTypeGraphCompleted EventType = "graph.completed"

	// EventTypeGraphFailed is emitted when graph execution fails.
	EventTypeGraphFailed EventType = "graph.failed"

	// EventTypeNodeStarted is emitted when a node begins execution.
	EventTypeNodeStarted EventType = "node.started"

	// EventTypeNodeCompleted is emitted when a node completes successfully.
	EventTypeNodeCompleted EventType = "node.completed"

	// EventTypeNodeFailed is emitted when a node execution fails.
	EventTypeNodeFailed EventType = "node.failed"

	// EventTypeStateChanged is emitted when execution state is modified.
	EventTypeStateChanged EventType = "state.changed"

	// EventTypeToolExecuted is emitted when a tool is executed.
	EventTypeToolExecuted EventType = "tool.executed"
)

// Event represents a system event.
type Event struct {
	// ID is a unique identifier for this event.
	ID string `json:"id"`

	// Type is the type of event.
	Type EventType `json:"type"`

	// Timestamp is when the event occurred.
	Timestamp time.Time `json:"timestamp"`

	// ExecutionID is the ID of the graph execution this event relates to.
	ExecutionID string `json:"execution_id"`

	// NodeID is the ID of the node this event relates to (if applicable).
	NodeID string `json:"node_id,omitempty"`

	// Data contains event-specific payload data.
	Data map[string]interface{} `json:"data,omitempty"`

	// Metadata contains additional event metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// EventHandler is a function that processes events.
type EventHandler func(ctx context.Context, event Event) error

// EventBus defines the interface for event publishing and subscription.
// For MVP, this is implemented using Redis Streams.
type EventBus interface {
	// Publish sends an event to a topic.
	Publish(ctx context.Context, topic string, event Event) error

	// Subscribe registers a handler for events on a topic.
	// The handler will be called for each event received.
	Subscribe(ctx context.Context, topic string, handler EventHandler) error

	// Unsubscribe removes a subscription from a topic.
	Unsubscribe(ctx context.Context, topic string) error

	// Close closes the event bus and cleans up resources.
	Close() error
}

// EventFilter defines criteria for filtering events.
type EventFilter struct {
	// Types filters events by type. If empty, all types are included.
	Types []EventType `json:"types,omitempty"`

	// ExecutionID filters events by execution ID.
	ExecutionID string `json:"execution_id,omitempty"`

	// NodeID filters events by node ID.
	NodeID string `json:"node_id,omitempty"`

	// Since filters events after this timestamp.
	Since time.Time `json:"since,omitempty"`

	// Until filters events before this timestamp.
	Until time.Time `json:"until,omitempty"`
}

// EventStore defines the interface for persisting and querying events.
type EventStore interface {
	// Store persists an event.
	Store(ctx context.Context, event Event) error

	// Query retrieves events matching the filter criteria.
	Query(ctx context.Context, filter EventFilter) ([]Event, error)

	// GetByID retrieves an event by its ID.
	GetByID(ctx context.Context, id string) (*Event, error)

	// GetByExecutionID retrieves all events for a specific execution.
	GetByExecutionID(ctx context.Context, executionID string) ([]Event, error)
}
