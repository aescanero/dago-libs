package graph

import (
	"testing"
)

func TestNewEdge(t *testing.T) {
	edge := NewEdge("node-1", "node-2")

	if edge.From != "node-1" {
		t.Errorf("expected From='node-1', got %q", edge.From)
	}
	if edge.To != "node-2" {
		t.Errorf("expected To='node-2', got %q", edge.To)
	}
}

func TestEdge_WithCondition(t *testing.T) {
	edge := NewEdge("node-1", "node-2").WithCondition("state.score > 0.5")

	if edge.Condition != "state.score > 0.5" {
		t.Errorf("expected condition 'state.score > 0.5', got %q", edge.Condition)
	}
}

func TestEdge_WithLabel(t *testing.T) {
	edge := NewEdge("node-1", "node-2").WithLabel("high score")

	if edge.Label != "high score" {
		t.Errorf("expected label 'high score', got %q", edge.Label)
	}
}

func TestEdge_Validate(t *testing.T) {
	tests := []struct {
		name        string
		edge        *Edge
		expectError bool
	}{
		{
			name:        "valid edge",
			edge:        &Edge{From: "node-1", To: "node-2"},
			expectError: false,
		},
		{
			name:        "empty from",
			edge:        &Edge{From: "", To: "node-2"},
			expectError: true,
		},
		{
			name:        "empty to",
			edge:        &Edge{From: "node-1", To: ""},
			expectError: true,
		},
		{
			name:        "self-loop",
			edge:        &Edge{From: "node-1", To: "node-1"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.edge.Validate()
			if tt.expectError && err == nil {
				t.Error("expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestEdge_ChainedMethods(t *testing.T) {
	edge := NewEdge("node-1", "node-2").
		WithCondition("state.ready == true").
		WithLabel("when ready")

	if edge.Condition != "state.ready == true" {
		t.Errorf("expected condition 'state.ready == true', got %q", edge.Condition)
	}
	if edge.Label != "when ready" {
		t.Errorf("expected label 'when ready', got %q", edge.Label)
	}
}
