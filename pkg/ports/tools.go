package ports

import (
	"context"
	"time"
)

// ToolType represents the type of tool executor.
type ToolType string

const (
	// ToolTypePython executes Python code.
	ToolTypePython ToolType = "python"

	// ToolTypeBash executes bash commands.
	ToolTypeBash ToolType = "bash"

	// ToolTypeHTTP makes HTTP requests.
	ToolTypeHTTP ToolType = "http"

	// ToolTypeCustom is a user-defined tool type.
	ToolTypeCustom ToolType = "custom"
)

// ToolSchema defines the schema for a tool's inputs.
type ToolSchema struct {
	// Name is the unique identifier for the tool.
	Name string `json:"name"`

	// Description explains what the tool does.
	Description string `json:"description"`

	// InputSchema is a JSON schema defining the tool's input parameters.
	InputSchema map[string]interface{} `json:"input_schema"`

	// OutputSchema is a JSON schema defining the tool's output format.
	OutputSchema map[string]interface{} `json:"output_schema,omitempty"`
}

// ToolResult represents the result of a tool execution.
type ToolResult struct {
	// Success indicates whether the tool executed successfully.
	Success bool `json:"success"`

	// Output contains the tool's output data.
	Output map[string]interface{} `json:"output,omitempty"`

	// Error contains error information if the tool failed.
	Error string `json:"error,omitempty"`

	// ExecutionTime is the duration of the tool execution.
	ExecutionTime time.Duration `json:"execution_time"`

	// Metadata contains additional execution metadata.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// ToolExecutor defines the interface for executing tools.
// Each tool type (Python, Bash, HTTP, etc.) implements this interface.
type ToolExecutor interface {
	// Execute runs the tool with the given parameters.
	Execute(ctx context.Context, params map[string]interface{}) (*ToolResult, error)

	// Schema returns the tool's schema definition.
	Schema() *ToolSchema

	// Type returns the type of this tool.
	Type() ToolType

	// Validate checks if the given parameters are valid for this tool.
	Validate(params map[string]interface{}) error
}

// ToolRegistry manages available tools and their executors.
type ToolRegistry interface {
	// Register adds a tool executor to the registry.
	Register(name string, executor ToolExecutor) error

	// Get retrieves a tool executor by name.
	Get(name string) (ToolExecutor, error)

	// List returns all registered tool names.
	List() []string

	// Unregister removes a tool executor from the registry.
	Unregister(name string) error

	// GetByType returns all tools of a specific type.
	GetByType(toolType ToolType) []ToolExecutor
}

// ToolConfig represents configuration for tool execution.
type ToolConfig struct {
	// Timeout is the maximum execution time for tools.
	Timeout time.Duration `json:"timeout"`

	// MaxRetries is the number of retry attempts on failure.
	MaxRetries int `json:"max_retries"`

	// RetryDelay is the delay between retry attempts.
	RetryDelay time.Duration `json:"retry_delay"`

	// Sandbox enables sandboxed execution for security.
	Sandbox bool `json:"sandbox"`

	// Environment contains environment variables for tool execution.
	Environment map[string]string `json:"environment,omitempty"`
}
