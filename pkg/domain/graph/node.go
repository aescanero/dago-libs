package graph

import (
	"context"

	"github.com/aescanero/dago-libs/pkg/domain/state"
)

// NodeType represents the type of a node in the graph.
type NodeType string

const (
	// NodeTypeExecutor represents a node that executes a task (calls LLM, runs tools, etc).
	NodeTypeExecutor NodeType = "executor"

	// NodeTypeRouter represents a node that makes routing decisions based on state.
	NodeTypeRouter NodeType = "router"

	// NodeTypeStart represents the entry point of the graph.
	NodeTypeStart NodeType = "start"

	// NodeTypeEnd represents an exit point of the graph.
	NodeTypeEnd NodeType = "end"
)

// Node defines the interface that all graph nodes must implement.
type Node interface {
	// GetID returns the unique identifier for this node.
	GetID() string

	// GetType returns the type of this node.
	GetType() NodeType

	// Execute runs the node's logic with the given context and state.
	// It returns the updated state and any error that occurred.
	// Implementations should be provided in the main dago repository, not here.
	Execute(ctx context.Context, state state.State) (state.State, error)

	// Validate checks if the node configuration is valid.
	Validate() error
}

// BaseNode provides common fields for all node types.
// This is a helper struct that can be embedded in concrete node implementations.
type BaseNode struct {
	ID          string                 `json:"id"`
	Type        NodeType               `json:"type"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// GetID returns the node's ID.
func (n *BaseNode) GetID() string {
	return n.ID
}

// GetType returns the node's type.
func (n *BaseNode) GetType() NodeType {
	return n.Type
}

// ExecutorNode represents a node that executes tasks like LLM calls or tool invocations.
type ExecutorNode struct {
	BaseNode
	// ExecutorType specifies what kind of executor this is (e.g., "llm", "tool", "python").
	ExecutorType string `json:"executor_type"`

	// Config contains executor-specific configuration.
	// The structure depends on the ExecutorType.
	Config map[string]interface{} `json:"config"`

	// InputMapping defines how to map state values to executor inputs.
	InputMapping map[string]string `json:"input_mapping,omitempty"`

	// OutputMapping defines how to map executor outputs back to state.
	OutputMapping map[string]string `json:"output_mapping,omitempty"`
}

// Execute is a placeholder that should be implemented in the main repository.
func (n *ExecutorNode) Execute(ctx context.Context, s state.State) (state.State, error) {
	// TODO: Implementation should be in dago repository
	panic("ExecutorNode.Execute must be implemented in the main dago repository")
}

// Validate checks if the executor node configuration is valid.
func (n *ExecutorNode) Validate() error {
	if n.ID == "" {
		return &ValidationError{Field: "id", Message: "executor node ID cannot be empty"}
	}
	if n.ExecutorType == "" {
		return &ValidationError{Field: "executor_type", Message: "executor type cannot be empty"}
	}
	return nil
}

// RouterNode represents a node that makes routing decisions based on state.
type RouterNode struct {
	BaseNode
	// Routes defines the routing logic.
	// Each route has a condition and target node ID.
	Routes []Route `json:"routes"`

	// DefaultRoute is the fallback route if no conditions match.
	DefaultRoute string `json:"default_route,omitempty"`
}

// Route represents a conditional routing rule.
type Route struct {
	// Condition is an expression evaluated against the state.
	// The expression syntax is implementation-specific (e.g., JSONPath, simple comparisons).
	Condition string `json:"condition"`

	// Target is the ID of the node to route to if the condition is true.
	Target string `json:"target"`

	// Description provides human-readable context for this route.
	Description string `json:"description,omitempty"`
}

// Execute is a placeholder that should be implemented in the main repository.
func (n *RouterNode) Execute(ctx context.Context, s state.State) (state.State, error) {
	// TODO: Implementation should be in dago repository
	panic("RouterNode.Execute must be implemented in the main dago repository")
}

// Validate checks if the router node configuration is valid.
func (n *RouterNode) Validate() error {
	if n.ID == "" {
		return &ValidationError{Field: "id", Message: "router node ID cannot be empty"}
	}
	if len(n.Routes) == 0 && n.DefaultRoute == "" {
		return &ValidationError{Field: "routes", Message: "router must have at least one route or a default route"}
	}
	for i, route := range n.Routes {
		if route.Target == "" {
			return &ValidationError{Field: "routes", Message: "route target cannot be empty at index " + string(rune(i))}
		}
	}
	return nil
}

// ValidationError represents a node validation error.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
