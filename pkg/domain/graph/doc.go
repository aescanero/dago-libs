// Package graph provides the core domain types for representing execution graphs.
//
// A graph consists of nodes (execution units) connected by edges (transitions).
// The graph structure defines the flow of execution through the system.
//
// Node Types:
//   - ExecutorNode: Executes tasks like LLM calls, tool invocations, or code execution
//   - RouterNode: Makes routing decisions based on state conditions
//   - Start/End: Special nodes for graph entry and exit points
//
// This package defines only the domain models and interfaces. Actual implementations
// of node execution logic should be in the main dago repository.
package graph
