package ports

import (
	"context"
	"time"
)

// Message represents a single message in an LLM conversation.
type Message struct {
	// Role is the role of the message sender (e.g., "system", "user", "assistant").
	Role string `json:"role"`

	// Content is the text content of the message.
	Content string `json:"content"`

	// Name is an optional name for the message sender.
	Name string `json:"name,omitempty"`
}

// Tool represents a tool that can be called by the LLM.
type Tool struct {
	// Name is the unique identifier for the tool.
	Name string `json:"name"`

	// Description explains what the tool does.
	Description string `json:"description"`

	// Parameters is a JSON schema defining the tool's input parameters.
	Parameters map[string]interface{} `json:"parameters"`
}

// ToolCall represents a request from the LLM to call a tool.
type ToolCall struct {
	// ID is a unique identifier for this tool call.
	ID string `json:"id"`

	// Name is the name of the tool to call.
	Name string `json:"name"`

	// Arguments contains the tool call arguments as a JSON object.
	Arguments map[string]interface{} `json:"arguments"`
}

// CompletionRequest represents a request for LLM text completion.
type CompletionRequest struct {
	// Messages is the conversation history.
	Messages []Message `json:"messages"`

	// Model is the identifier of the LLM model to use.
	Model string `json:"model"`

	// Temperature controls randomness in the response (0.0 to 1.0).
	Temperature float64 `json:"temperature,omitempty"`

	// MaxTokens is the maximum number of tokens to generate.
	MaxTokens int `json:"max_tokens,omitempty"`

	// TopP controls nucleus sampling (0.0 to 1.0).
	TopP float64 `json:"top_p,omitempty"`

	// Stop is a list of sequences where the LLM should stop generating.
	Stop []string `json:"stop,omitempty"`

	// PresencePenalty penalizes new tokens based on whether they appear in the text so far.
	PresencePenalty float64 `json:"presence_penalty,omitempty"`

	// FrequencyPenalty penalizes new tokens based on their frequency in the text so far.
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`

	// User is an optional identifier representing the end-user.
	User string `json:"user,omitempty"`
}

// CompletionResponse represents the response from an LLM completion.
type CompletionResponse struct {
	// ID is a unique identifier for this completion.
	ID string `json:"id"`

	// Model is the model that generated this completion.
	Model string `json:"model"`

	// Message is the generated message.
	Message Message `json:"message"`

	// ToolCalls contains any tool calls requested by the LLM.
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`

	// FinishReason indicates why the generation stopped.
	FinishReason string `json:"finish_reason"`

	// Usage contains token usage information.
	Usage UsageInfo `json:"usage"`

	// CreatedAt is the timestamp when this completion was created.
	CreatedAt time.Time `json:"created_at"`
}

// UsageInfo contains token usage statistics.
type UsageInfo struct {
	// PromptTokens is the number of tokens in the prompt.
	PromptTokens int `json:"prompt_tokens"`

	// CompletionTokens is the number of tokens in the completion.
	CompletionTokens int `json:"completion_tokens"`

	// TotalTokens is the total number of tokens used.
	TotalTokens int `json:"total_tokens"`
}

// JSONSchema represents a JSON schema for structured output.
type JSONSchema map[string]interface{}

// StructuredResponse represents a response with structured JSON output.
type StructuredResponse struct {
	// Data contains the structured data conforming to the provided schema.
	Data map[string]interface{} `json:"data"`

	// Usage contains token usage information.
	Usage UsageInfo `json:"usage"`

	// CreatedAt is the timestamp when this response was created.
	CreatedAt time.Time `json:"created_at"`
}

// LLMClient defines the interface for interacting with Large Language Models.
// Implementations should handle provider-specific details (OpenAI, Anthropic, etc).
type LLMClient interface {
	// Complete performs a standard text completion.
	Complete(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

	// CompleteWithTools performs a completion with tool calling support.
	// The LLM can request tool executions via ToolCall objects in the response.
	CompleteWithTools(ctx context.Context, req CompletionRequest, tools []Tool) (*CompletionResponse, error)

	// CompleteStructured performs a completion with guaranteed JSON schema conformance.
	// The response will be validated against the provided schema.
	CompleteStructured(ctx context.Context, req CompletionRequest, schema JSONSchema) (*StructuredResponse, error)

	// StreamComplete performs a streaming completion (optional for MVP).
	// Returns a channel that yields completion chunks as they arrive.
	// StreamComplete(ctx context.Context, req CompletionRequest) (<-chan CompletionChunk, error)
}

// CompletionChunk represents a chunk of a streaming completion (for future use).
type CompletionChunk struct {
	Delta   string `json:"delta"`
	IsFinal bool   `json:"is_final"`
}
