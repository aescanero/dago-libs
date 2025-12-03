package domain

import (
	"time"

	"github.com/aescanero/dago-libs/pkg/domain/graph"
	"github.com/aescanero/dago-libs/pkg/domain/state"
)

// Re-export types from sub-packages for compatibility
type (
	Graph    = graph.Graph
	Node     = graph.Node
	NodeType = graph.NodeType
	State    = state.State
)

// Node type constants
const (
	NodeTypeAgent      NodeType = "agent"
	NodeTypeParallel   NodeType = "parallel"
	NodeTypeConditional NodeType = "conditional"
	NodeTypeLoop       NodeType = "loop"
	NodeTypeMap        NodeType = "map"
	NodeTypeReduce     NodeType = "reduce"
)

// ExecutionStatus represents the status of graph or node execution
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"
	ExecutionStatusRunning   ExecutionStatus = "running"
	ExecutionStatusCompleted ExecutionStatus = "completed"
	ExecutionStatusFailed    ExecutionStatus = "failed"
	ExecutionStatusCancelled ExecutionStatus = "cancelled"
	ExecutionStatusSubmitted ExecutionStatus = "submitted"
)

// GraphState represents the state of a graph execution
type GraphState struct {
	GraphID      string                  `json:"graph_id"`
	Graph        *Graph                  `json:"graph"`
	Status       ExecutionStatus         `json:"status"`
	Inputs       map[string]interface{}  `json:"inputs"`
	NodeStates   map[string]*NodeState   `json:"node_states"`
	SubmittedAt  time.Time               `json:"submitted_at"`
	StartedAt    *time.Time              `json:"started_at,omitempty"`
	CompletedAt  *time.Time              `json:"completed_at,omitempty"`
	Error        string                  `json:"error,omitempty"`
}

// NodeState represents the state of a node execution
type NodeState struct {
	NodeID      string                 `json:"node_id"`
	Status      ExecutionStatus        `json:"status"`
	Output      interface{}            `json:"output,omitempty"`
	Error       string                 `json:"error,omitempty"`
	StartedAt   *time.Time             `json:"started_at,omitempty"`
	CompletedAt *time.Time             `json:"completed_at,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// EventType represents the type of an event
type EventType string

const (
	EventTypeGraphSubmitted  EventType = "graph.submitted"
	EventTypeGraphStarted    EventType = "graph.started"
	EventTypeGraphCompleted  EventType = "graph.completed"
	EventTypeGraphFailed     EventType = "graph.failed"
	EventTypeGraphCancelled  EventType = "graph.cancelled"
	EventTypeNodeStarted     EventType = "node.started"
	EventTypeNodeCompleted   EventType = "node.completed"
	EventTypeNodeFailed      EventType = "node.failed"
)

// Event represents an event in the system
type Event struct {
	ID        string                 `json:"id"`
	Type      EventType              `json:"type"`
	GraphID   string                 `json:"graph_id"`
	NodeID    string                 `json:"node_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// LLMRequest represents a request to an LLM
type LLMRequest struct {
	Model       string              `json:"model"`
	Messages    []Message           `json:"messages"`
	System      string              `json:"system,omitempty"`
	MaxTokens   int                 `json:"max_tokens"`
	Temperature float64             `json:"temperature,omitempty"`
	Tools       []Tool              `json:"tools,omitempty"`
}

// LLMResponse represents a response from an LLM
type LLMResponse struct {
	Content   string   `json:"content"`
	Model     string   `json:"model"`
	Usage     Usage    `json:"usage"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

// Message represents a message in an LLM conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool represents a tool that can be used by an LLM
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// ToolCall represents a tool invocation by an LLM
type ToolCall struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Input    map[string]interface{} `json:"input"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolCallID string      `json:"tool_call_id"`
	Output     interface{} `json:"output"`
	Error      string      `json:"error,omitempty"`
}

// Usage represents token usage statistics
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
