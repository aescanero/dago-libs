// Package schema provides JSON schema definitions and validation for DA Orchestrator entities.
//
// This package includes JSON schemas for:
//   - Graph definitions (graph.schema.json)
//   - Executor node configurations (executor-node.schema.json)
//   - Router node configurations (router-node.schema.json)
//
// The Validator type provides methods to validate JSON data against these schemas,
// ensuring that graph definitions and node configurations conform to the expected structure.
//
// Schemas are embedded in the binary using go:embed, so they are always available
// at runtime without requiring external files.
package schema
