package graph

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// Graph represents a directed graph of nodes and edges that defines the execution flow.
type Graph struct {
	// ID is a unique identifier for this graph.
	ID string `json:"id"`

	// Name is a human-readable name for the graph.
	Name string `json:"name,omitempty"`

	// Description provides context about what this graph does.
	Description string `json:"description,omitempty"`

	// Nodes is a map of node ID to node instance.
	// Using a map allows O(1) lookups by ID.
	Nodes map[string]Node `json:"nodes"`

	// Edges defines the connections between nodes.
	Edges []*Edge `json:"edges"`

	// EntryNode is the ID of the node where execution begins.
	EntryNode string `json:"entry_node"`

	// Metadata stores additional graph-level data.
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// Version is the schema version of this graph definition.
	Version string `json:"version,omitempty"`
}

// NewGraph creates a new graph with a generated UUID.
func NewGraph(name string) *Graph {
	return &Graph{
		ID:       uuid.New().String(),
		Name:     name,
		Nodes:    make(map[string]Node),
		Edges:    make([]*Edge, 0),
		Metadata: make(map[string]interface{}),
		Version:  "1.0",
	}
}

// AddNode adds a node to the graph.
// Returns an error if a node with the same ID already exists.
func (g *Graph) AddNode(node Node) error {
	if node == nil {
		return fmt.Errorf("cannot add nil node")
	}

	nodeID := node.GetID()
	if nodeID == "" {
		return &ValidationError{Field: "node.id", Message: "node ID cannot be empty"}
	}

	if _, exists := g.Nodes[nodeID]; exists {
		return &ValidationError{Field: "node.id", Message: fmt.Sprintf("node with ID '%s' already exists", nodeID)}
	}

	if err := node.Validate(); err != nil {
		return fmt.Errorf("node validation failed: %w", err)
	}

	g.Nodes[nodeID] = node
	return nil
}

// GetNode retrieves a node by its ID.
// Returns nil if the node doesn't exist.
func (g *Graph) GetNode(nodeID string) Node {
	return g.Nodes[nodeID]
}

// RemoveNode removes a node from the graph.
// Also removes any edges connected to this node.
func (g *Graph) RemoveNode(nodeID string) {
	delete(g.Nodes, nodeID)

	// Remove edges connected to this node
	filteredEdges := make([]*Edge, 0)
	for _, edge := range g.Edges {
		if edge.From != nodeID && edge.To != nodeID {
			filteredEdges = append(filteredEdges, edge)
		}
	}
	g.Edges = filteredEdges
}

// AddEdge adds an edge to the graph.
// Validates that both source and target nodes exist.
func (g *Graph) AddEdge(edge *Edge) error {
	if edge == nil {
		return fmt.Errorf("cannot add nil edge")
	}

	if err := edge.Validate(); err != nil {
		return fmt.Errorf("edge validation failed: %w", err)
	}

	// Verify that both nodes exist
	if g.GetNode(edge.From) == nil {
		return &ValidationError{Field: "edge.from", Message: fmt.Sprintf("source node '%s' does not exist", edge.From)}
	}
	if g.GetNode(edge.To) == nil {
		return &ValidationError{Field: "edge.to", Message: fmt.Sprintf("target node '%s' does not exist", edge.To)}
	}

	g.Edges = append(g.Edges, edge)
	return nil
}

// GetOutgoingEdges returns all edges originating from the given node.
func (g *Graph) GetOutgoingEdges(nodeID string) []*Edge {
	edges := make([]*Edge, 0)
	for _, edge := range g.Edges {
		if edge.From == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// GetIncomingEdges returns all edges targeting the given node.
func (g *Graph) GetIncomingEdges(nodeID string) []*Edge {
	edges := make([]*Edge, 0)
	for _, edge := range g.Edges {
		if edge.To == nodeID {
			edges = append(edges, edge)
		}
	}
	return edges
}

// Validate performs comprehensive validation of the graph structure.
func (g *Graph) Validate() error {
	if g.ID == "" {
		return &ValidationError{Field: "id", Message: "graph ID cannot be empty"}
	}

	if len(g.Nodes) == 0 {
		return &ValidationError{Field: "nodes", Message: "graph must have at least one node"}
	}

	if g.EntryNode == "" {
		return &ValidationError{Field: "entry_node", Message: "graph must have an entry node"}
	}

	// Verify entry node exists
	if g.GetNode(g.EntryNode) == nil {
		return &ValidationError{Field: "entry_node", Message: fmt.Sprintf("entry node '%s' does not exist", g.EntryNode)}
	}

	// Validate all nodes
	for _, node := range g.Nodes {
		if err := node.Validate(); err != nil {
			return fmt.Errorf("node '%s' validation failed: %w", node.GetID(), err)
		}
	}

	// Validate all edges
	for i, edge := range g.Edges {
		if err := edge.Validate(); err != nil {
			return fmt.Errorf("edge %d validation failed: %w", i, err)
		}

		// Verify both nodes exist
		if g.GetNode(edge.From) == nil {
			return &ValidationError{Field: "edge.from", Message: fmt.Sprintf("edge references non-existent source node '%s'", edge.From)}
		}
		if g.GetNode(edge.To) == nil {
			return &ValidationError{Field: "edge.to", Message: fmt.Sprintf("edge references non-existent target node '%s'", edge.To)}
		}
	}

	// TODO: Add cycle detection for graphs that shouldn't have cycles
	// TODO: Add reachability analysis to detect orphaned nodes

	return nil
}

// ToJSON serializes the graph to JSON.
func (g *Graph) ToJSON() (string, error) {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal graph to JSON: %w", err)
	}
	return string(data), nil
}

// FromJSON deserializes a graph from JSON.
// Note: This is a basic implementation. Full deserialization with proper node type
// handling should be implemented in the main dago repository.
func FromJSON(jsonStr string) (*Graph, error) {
	var g Graph
	if err := json.Unmarshal([]byte(jsonStr), &g); err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph from JSON: %w", err)
	}
	return &g, nil
}

// Clone creates a deep copy of the graph.
// Note: This uses JSON serialization for simplicity.
// TODO: Consider more efficient cloning for performance-critical paths.
func (g *Graph) Clone() (*Graph, error) {
	jsonStr, err := g.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to serialize graph for cloning: %w", err)
	}
	return FromJSON(jsonStr)
}

// NodeCount returns the number of nodes in the graph.
func (g *Graph) NodeCount() int {
	return len(g.Nodes)
}

// EdgeCount returns the number of edges in the graph.
func (g *Graph) EdgeCount() int {
	return len(g.Edges)
}
