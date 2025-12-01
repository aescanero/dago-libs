package graph

import (
	"testing"
)

func TestBaseNode_GetID(t *testing.T) {
	node := &BaseNode{
		ID:   "test-node",
		Type: NodeTypeExecutor,
	}

	if node.GetID() != "test-node" {
		t.Errorf("expected ID 'test-node', got %q", node.GetID())
	}
}

func TestBaseNode_GetType(t *testing.T) {
	node := &BaseNode{
		ID:   "test-node",
		Type: NodeTypeRouter,
	}

	if node.GetType() != NodeTypeRouter {
		t.Errorf("expected type %q, got %q", NodeTypeRouter, node.GetType())
	}
}

func TestExecutorNode_Validate(t *testing.T) {
	tests := []struct {
		name        string
		node        *ExecutorNode
		expectError bool
	}{
		{
			name: "valid executor node",
			node: &ExecutorNode{
				BaseNode: BaseNode{
					ID:   "exec-1",
					Type: NodeTypeExecutor,
				},
				ExecutorType: "llm",
				Config:       map[string]interface{}{"model": "gpt-4"},
			},
			expectError: false,
		},
		{
			name: "empty ID",
			node: &ExecutorNode{
				BaseNode: BaseNode{
					ID:   "",
					Type: NodeTypeExecutor,
				},
				ExecutorType: "llm",
			},
			expectError: true,
		},
		{
			name: "empty executor type",
			node: &ExecutorNode{
				BaseNode: BaseNode{
					ID:   "exec-1",
					Type: NodeTypeExecutor,
				},
				ExecutorType: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Validate()
			if tt.expectError && err == nil {
				t.Error("expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestRouterNode_Validate(t *testing.T) {
	tests := []struct {
		name        string
		node        *RouterNode
		expectError bool
	}{
		{
			name: "valid router with routes",
			node: &RouterNode{
				BaseNode: BaseNode{
					ID:   "router-1",
					Type: NodeTypeRouter,
				},
				Routes: []Route{
					{Condition: "state.score > 0.5", Target: "node-1"},
				},
			},
			expectError: false,
		},
		{
			name: "valid router with default route",
			node: &RouterNode{
				BaseNode: BaseNode{
					ID:   "router-1",
					Type: NodeTypeRouter,
				},
				DefaultRoute: "node-1",
			},
			expectError: false,
		},
		{
			name: "empty ID",
			node: &RouterNode{
				BaseNode: BaseNode{
					ID:   "",
					Type: NodeTypeRouter,
				},
				Routes: []Route{{Target: "node-1"}},
			},
			expectError: true,
		},
		{
			name: "no routes and no default",
			node: &RouterNode{
				BaseNode: BaseNode{
					ID:   "router-1",
					Type: NodeTypeRouter,
				},
			},
			expectError: true,
		},
		{
			name: "route with empty target",
			node: &RouterNode{
				BaseNode: BaseNode{
					ID:   "router-1",
					Type: NodeTypeRouter,
				},
				Routes: []Route{
					{Condition: "true", Target: ""},
				},
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.node.Validate()
			if tt.expectError && err == nil {
				t.Error("expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "test_field",
		Message: "is invalid",
	}

	expected := "test_field: is invalid"
	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}
