package schema

import (
	"errors"
	"testing"
)

func TestNewValidator(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	if validator == nil {
		t.Fatal("expected non-nil validator")
	}

	if validator.graphSchema == nil {
		t.Error("graph schema not loaded")
	}
	if validator.executorSchema == nil {
		t.Error("executor schema not loaded")
	}
	if validator.routerSchema == nil {
		t.Error("router schema not loaded")
	}
}

func TestValidateGraph_Valid(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	validGraph := []byte(`{
		"id": "graph-1",
		"name": "Test Graph",
		"version": "1.0",
		"nodes": {
			"start": {
				"id": "start",
				"type": "executor",
				"executor_type": "llm",
				"config": {
					"model": "gpt-4"
				}
			}
		},
		"edges": [],
		"entry_node": "start"
	}`)

	err = validator.ValidateGraph(validGraph)
	if err != nil {
		t.Errorf("validation failed for valid graph: %v", err)
	}
}

func TestValidateGraph_MissingID(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	invalidGraph := []byte(`{
		"name": "Test Graph",
		"nodes": {
			"start": {
				"id": "start",
				"type": "executor",
				"executor_type": "llm"
			}
		},
		"entry_node": "start"
	}`)

	err = validator.ValidateGraph(invalidGraph)
	if err == nil {
		t.Error("expected validation error for graph without ID")
	}
}

func TestValidateGraph_MissingEntryNode(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	invalidGraph := []byte(`{
		"id": "graph-1",
		"nodes": {
			"start": {
				"id": "start",
				"type": "executor",
				"executor_type": "llm"
			}
		}
	}`)

	err = validator.ValidateGraph(invalidGraph)
	if err == nil {
		t.Error("expected validation error for graph without entry_node")
	}
}

func TestValidateGraph_EmptyNodes(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	invalidGraph := []byte(`{
		"id": "graph-1",
		"nodes": {},
		"entry_node": "start"
	}`)

	err = validator.ValidateGraph(invalidGraph)
	if err == nil {
		t.Error("expected validation error for graph with empty nodes")
	}
}

func TestValidateGraph_InvalidJSON(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	invalidJSON := []byte(`{invalid json}`)

	err = validator.ValidateGraph(invalidJSON)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestValidateExecutorNode_Valid(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	validNode := []byte(`{
		"executor_type": "llm",
		"config": {
			"model": "gpt-4",
			"temperature": 0.7,
			"max_tokens": 2000
		}
	}`)

	err = validator.ValidateExecutorNode(validNode)
	if err != nil {
		t.Errorf("validation failed for valid executor node: %v", err)
	}
}

func TestValidateExecutorNode_MissingExecutorType(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	invalidNode := []byte(`{
		"config": {
			"model": "gpt-4"
		}
	}`)

	err = validator.ValidateExecutorNode(invalidNode)
	if err == nil {
		t.Error("expected validation error for node without executor_type")
	}
}

func TestValidateExecutorNode_MissingConfig(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	invalidNode := []byte(`{
		"executor_type": "llm"
	}`)

	err = validator.ValidateExecutorNode(invalidNode)
	if err == nil {
		t.Error("expected validation error for node without config")
	}
}

func TestValidateExecutorNode_LLMConfig(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	validLLMNode := []byte(`{
		"executor_type": "llm",
		"config": {
			"model": "gpt-4",
			"temperature": 0.7,
			"max_tokens": 1000,
			"system_prompt": "You are a helpful assistant"
		}
	}`)

	err = validator.ValidateExecutorNode(validLLMNode)
	if err != nil {
		t.Errorf("validation failed for valid LLM node: %v", err)
	}
}

func TestValidateExecutorNode_ToolConfig(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	validToolNode := []byte(`{
		"executor_type": "tool",
		"config": {
			"tool_name": "python",
			"parameters": {
				"script": "print('hello')"
			},
			"timeout": 300
		}
	}`)

	err = validator.ValidateExecutorNode(validToolNode)
	if err != nil {
		t.Errorf("validation failed for valid tool node: %v", err)
	}
}

func TestValidateRouterNode_Valid(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	validRouter := []byte(`{
		"routes": [
			{
				"condition": "state.score > 0.8",
				"target": "high-score-node",
				"description": "Route for high scores"
			},
			{
				"condition": "state.score <= 0.8",
				"target": "low-score-node"
			}
		],
		"default_route": "fallback-node"
	}`)

	err = validator.ValidateRouterNode(validRouter)
	if err != nil {
		t.Errorf("validation failed for valid router node: %v", err)
	}
}

func TestValidateRouterNode_OnlyDefaultRoute(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	validRouter := []byte(`{
		"default_route": "fallback-node"
	}`)

	err = validator.ValidateRouterNode(validRouter)
	if err != nil {
		t.Errorf("validation failed for router with only default route: %v", err)
	}
}

func TestValidateRouterNode_MissingTarget(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	invalidRouter := []byte(`{
		"routes": [
			{
				"condition": "state.score > 0.8"
			}
		]
	}`)

	err = validator.ValidateRouterNode(invalidRouter)
	if err == nil {
		t.Error("expected validation error for route without target")
	}
}

func TestValidateRouterNode_NoRoutesOrDefault(t *testing.T) {
	validator, err := NewValidator()
	if err != nil {
		t.Fatalf("NewValidator failed: %v", err)
	}

	invalidRouter := []byte(`{}`)

	err = validator.ValidateRouterNode(invalidRouter)
	if err == nil {
		t.Error("expected validation error for router without routes or default")
	}
}

func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *ValidationError
		expected string
	}{
		{
			name: "with cause",
			err: &ValidationError{
				SchemaType: "graph",
				Message:    "invalid structure",
				Cause:      errors.New("field validation error"),
			},
			expected: "graph schema validation failed: invalid structure (caused by: field validation error)",
		},
		{
			name: "without cause",
			err: &ValidationError{
				SchemaType: "node",
				Message:    "missing required field",
				Cause:      nil,
			},
			expected: "node schema validation failed: missing required field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.err.Error())
			}
		})
	}
}

func TestValidationError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &ValidationError{
		SchemaType: "test",
		Message:    "test error",
		Cause:      cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("expected unwrapped error to be %v, got %v", cause, unwrapped)
	}
}
