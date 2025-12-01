package graph

// Edge represents a directed connection between two nodes in the graph.
type Edge struct {
	// ID is a unique identifier for this edge.
	ID string `json:"id,omitempty"`

	// From is the ID of the source node.
	From string `json:"from"`

	// To is the ID of the target node.
	To string `json:"to"`

	// Condition is an optional expression that must evaluate to true for this edge to be traversed.
	// If empty, the edge is always traversable.
	// The expression syntax is implementation-specific.
	Condition string `json:"condition,omitempty"`

	// Label provides a human-readable description of this edge.
	Label string `json:"label,omitempty"`

	// Metadata stores additional edge-specific data.
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// NewEdge creates a new edge with the given source and target.
func NewEdge(from, to string) *Edge {
	return &Edge{
		From: from,
		To:   to,
	}
}

// WithCondition sets a condition on the edge.
func (e *Edge) WithCondition(condition string) *Edge {
	e.Condition = condition
	return e
}

// WithLabel sets a label on the edge.
func (e *Edge) WithLabel(label string) *Edge {
	e.Label = label
	return e
}

// Validate checks if the edge configuration is valid.
func (e *Edge) Validate() error {
	if e.From == "" {
		return &ValidationError{Field: "from", Message: "edge source node ID cannot be empty"}
	}
	if e.To == "" {
		return &ValidationError{Field: "to", Message: "edge target node ID cannot be empty"}
	}
	if e.From == e.To {
		return &ValidationError{Field: "from/to", Message: "edge cannot connect a node to itself"}
	}
	return nil
}
