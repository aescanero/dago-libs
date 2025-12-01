package graph

import (
	"context"
	"testing"

	"github.com/aescanero/dago-libs/pkg/domain/state"
)

// Mock node for testing
type mockNode struct {
	id       string
	nodeType NodeType
}

func (n *mockNode) GetID() string {
	return n.id
}

func (n *mockNode) GetType() NodeType {
	return n.nodeType
}

func (n *mockNode) Execute(ctx context.Context, s state.State) (state.State, error) {
	return s, nil
}

func (n *mockNode) Validate() error {
	if n.id == "" {
		return &ValidationError{Field: "id", Message: "cannot be empty"}
	}
	return nil
}

func TestNewGraph(t *testing.T) {
	g := NewGraph("test-graph")

	if g.ID == "" {
		t.Error("expected graph ID to be generated")
	}
	if g.Name != "test-graph" {
		t.Errorf("expected name 'test-graph', got %q", g.Name)
	}
	if g.Nodes == nil {
		t.Error("expected nodes map to be initialized")
	}
	if g.Edges == nil {
		t.Error("expected edges slice to be initialized")
	}
	if g.Version != "1.0" {
		t.Errorf("expected version '1.0', got %q", g.Version)
	}
}

func TestGraphAddNode(t *testing.T) {
	g := NewGraph("test")

	node := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	err := g.AddNode(node)
	if err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}

	// Verify node was added
	if g.GetNode("node-1") != node {
		t.Error("node was not added to graph")
	}

	// Test adding duplicate
	err = g.AddNode(node)
	if err == nil {
		t.Error("expected error when adding duplicate node")
	}
}

func TestGraphAddNodeNil(t *testing.T) {
	g := NewGraph("test")

	err := g.AddNode(nil)
	if err == nil {
		t.Error("expected error when adding nil node")
	}
}

func TestGraphAddNodeInvalid(t *testing.T) {
	g := NewGraph("test")

	// Node with empty ID should fail validation
	node := &mockNode{id: "", nodeType: NodeTypeExecutor}
	err := g.AddNode(node)
	if err == nil {
		t.Error("expected error when adding node with empty ID")
	}
}

func TestGraphGetNode(t *testing.T) {
	g := NewGraph("test")

	node := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	g.AddNode(node)

	retrieved := g.GetNode("node-1")
	if retrieved != node {
		t.Error("GetNode returned wrong node")
	}

	// Test non-existent node
	retrieved = g.GetNode("non-existent")
	if retrieved != nil {
		t.Error("expected nil for non-existent node")
	}
}

func TestGraphRemoveNode(t *testing.T) {
	g := NewGraph("test")

	node1 := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	node2 := &mockNode{id: "node-2", nodeType: NodeTypeExecutor}
	g.AddNode(node1)
	g.AddNode(node2)

	edge := NewEdge("node-1", "node-2")
	g.AddEdge(edge)

	// Remove node-1
	g.RemoveNode("node-1")

	if g.GetNode("node-1") != nil {
		t.Error("node-1 should be removed")
	}

	// Check that edge was also removed
	if len(g.Edges) != 0 {
		t.Errorf("expected 0 edges after removing node, got %d", len(g.Edges))
	}
}

func TestGraphAddEdge(t *testing.T) {
	g := NewGraph("test")

	node1 := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	node2 := &mockNode{id: "node-2", nodeType: NodeTypeExecutor}
	g.AddNode(node1)
	g.AddNode(node2)

	edge := NewEdge("node-1", "node-2")
	err := g.AddEdge(edge)
	if err != nil {
		t.Fatalf("AddEdge failed: %v", err)
	}

	if len(g.Edges) != 1 {
		t.Errorf("expected 1 edge, got %d", len(g.Edges))
	}
}

func TestGraphAddEdgeNil(t *testing.T) {
	g := NewGraph("test")

	err := g.AddEdge(nil)
	if err == nil {
		t.Error("expected error when adding nil edge")
	}
}

func TestGraphAddEdgeNonExistentNode(t *testing.T) {
	g := NewGraph("test")

	node1 := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	g.AddNode(node1)

	// Edge to non-existent node
	edge := NewEdge("node-1", "non-existent")
	err := g.AddEdge(edge)
	if err == nil {
		t.Error("expected error when adding edge to non-existent node")
	}
}

func TestGraphGetOutgoingEdges(t *testing.T) {
	g := NewGraph("test")

	node1 := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	node2 := &mockNode{id: "node-2", nodeType: NodeTypeExecutor}
	node3 := &mockNode{id: "node-3", nodeType: NodeTypeExecutor}
	g.AddNode(node1)
	g.AddNode(node2)
	g.AddNode(node3)

	g.AddEdge(NewEdge("node-1", "node-2"))
	g.AddEdge(NewEdge("node-1", "node-3"))
	g.AddEdge(NewEdge("node-2", "node-3"))

	outgoing := g.GetOutgoingEdges("node-1")
	if len(outgoing) != 2 {
		t.Errorf("expected 2 outgoing edges from node-1, got %d", len(outgoing))
	}
}

func TestGraphGetIncomingEdges(t *testing.T) {
	g := NewGraph("test")

	node1 := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	node2 := &mockNode{id: "node-2", nodeType: NodeTypeExecutor}
	node3 := &mockNode{id: "node-3", nodeType: NodeTypeExecutor}
	g.AddNode(node1)
	g.AddNode(node2)
	g.AddNode(node3)

	g.AddEdge(NewEdge("node-1", "node-3"))
	g.AddEdge(NewEdge("node-2", "node-3"))

	incoming := g.GetIncomingEdges("node-3")
	if len(incoming) != 2 {
		t.Errorf("expected 2 incoming edges to node-3, got %d", len(incoming))
	}
}

func TestGraphValidate(t *testing.T) {
	tests := []struct {
		name        string
		setupGraph  func() *Graph
		expectError bool
	}{
		{
			name: "valid graph",
			setupGraph: func() *Graph {
				g := NewGraph("test")
				node := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
				g.AddNode(node)
				g.EntryNode = "node-1"
				return g
			},
			expectError: false,
		},
		{
			name: "empty graph ID",
			setupGraph: func() *Graph {
				g := NewGraph("test")
				g.ID = ""
				node := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
				g.AddNode(node)
				g.EntryNode = "node-1"
				return g
			},
			expectError: true,
		},
		{
			name: "no nodes",
			setupGraph: func() *Graph {
				g := NewGraph("test")
				g.EntryNode = "node-1"
				return g
			},
			expectError: true,
		},
		{
			name: "no entry node",
			setupGraph: func() *Graph {
				g := NewGraph("test")
				node := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
				g.AddNode(node)
				return g
			},
			expectError: true,
		},
		{
			name: "non-existent entry node",
			setupGraph: func() *Graph {
				g := NewGraph("test")
				node := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
				g.AddNode(node)
				g.EntryNode = "non-existent"
				return g
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := tt.setupGraph()
			err := g.Validate()
			if tt.expectError && err == nil {
				t.Error("expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestGraphToJSON(t *testing.T) {
	g := NewGraph("test")
	node := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	g.AddNode(node)
	g.EntryNode = "node-1"

	jsonStr, err := g.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if jsonStr == "" {
		t.Error("expected non-empty JSON string")
	}
}

func TestGraphNodeCount(t *testing.T) {
	g := NewGraph("test")

	if g.NodeCount() != 0 {
		t.Errorf("expected 0 nodes, got %d", g.NodeCount())
	}

	g.AddNode(&mockNode{id: "node-1", nodeType: NodeTypeExecutor})
	g.AddNode(&mockNode{id: "node-2", nodeType: NodeTypeExecutor})

	if g.NodeCount() != 2 {
		t.Errorf("expected 2 nodes, got %d", g.NodeCount())
	}
}

func TestGraphEdgeCount(t *testing.T) {
	g := NewGraph("test")

	node1 := &mockNode{id: "node-1", nodeType: NodeTypeExecutor}
	node2 := &mockNode{id: "node-2", nodeType: NodeTypeExecutor}
	g.AddNode(node1)
	g.AddNode(node2)

	if g.EdgeCount() != 0 {
		t.Errorf("expected 0 edges, got %d", g.EdgeCount())
	}

	g.AddEdge(NewEdge("node-1", "node-2"))

	if g.EdgeCount() != 1 {
		t.Errorf("expected 1 edge, got %d", g.EdgeCount())
	}
}
