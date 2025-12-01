// Package schema provides JSON schema validation for graph definitions.
package schema

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

//go:embed graph.schema.json
var graphSchemaJSON string

//go:embed executor-node.schema.json
var executorNodeSchemaJSON string

//go:embed router-node.schema.json
var routerNodeSchemaJSON string

// Validator provides JSON schema validation for DA Orchestrator entities.
type Validator struct {
	compiler       *jsonschema.Compiler
	graphSchema    *jsonschema.Schema
	executorSchema *jsonschema.Schema
	routerSchema   *jsonschema.Schema
}

// NewValidator creates a new validator with all schemas loaded.
func NewValidator() (*Validator, error) {
	compiler := jsonschema.NewCompiler()
	compiler.Draft = jsonschema.Draft7

	// Add schemas to the compiler
	if err := compiler.AddResource("graph.schema.json", strings.NewReader(graphSchemaJSON)); err != nil {
		return nil, fmt.Errorf("failed to add graph schema: %w", err)
	}
	if err := compiler.AddResource("executor-node.schema.json", strings.NewReader(executorNodeSchemaJSON)); err != nil {
		return nil, fmt.Errorf("failed to add executor node schema: %w", err)
	}
	if err := compiler.AddResource("router-node.schema.json", strings.NewReader(routerNodeSchemaJSON)); err != nil {
		return nil, fmt.Errorf("failed to add router node schema: %w", err)
	}

	// Compile the schemas
	graphSchema, err := compiler.Compile("graph.schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile graph schema: %w", err)
	}
	executorSchema, err := compiler.Compile("executor-node.schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile executor node schema: %w", err)
	}
	routerSchema, err := compiler.Compile("router-node.schema.json")
	if err != nil {
		return nil, fmt.Errorf("failed to compile router node schema: %w", err)
	}

	return &Validator{
		compiler:       compiler,
		graphSchema:    graphSchema,
		executorSchema: executorSchema,
		routerSchema:   routerSchema,
	}, nil
}

// ValidateGraph validates a graph definition against the graph schema.
func (v *Validator) ValidateGraph(graphJSON []byte) error {
	var data interface{}
	if err := json.Unmarshal(graphJSON, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := v.graphSchema.Validate(data); err != nil {
		return fmt.Errorf("graph validation failed: %w", err)
	}

	return nil
}

// ValidateExecutorNode validates an executor node configuration.
func (v *Validator) ValidateExecutorNode(nodeJSON []byte) error {
	var data interface{}
	if err := json.Unmarshal(nodeJSON, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := v.executorSchema.Validate(data); err != nil {
		return fmt.Errorf("executor node validation failed: %w", err)
	}

	return nil
}

// ValidateRouterNode validates a router node configuration.
func (v *Validator) ValidateRouterNode(nodeJSON []byte) error {
	var data interface{}
	if err := json.Unmarshal(nodeJSON, &data); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	if err := v.routerSchema.Validate(data); err != nil {
		return fmt.Errorf("router node validation failed: %w", err)
	}

	return nil
}

// ValidationError wraps validation errors with additional context.
type ValidationError struct {
	SchemaType string
	Message    string
	Cause      error
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s schema validation failed: %s (caused by: %v)", e.SchemaType, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s schema validation failed: %s", e.SchemaType, e.Message)
}

// Unwrap implements the errors.Unwrap interface.
func (e *ValidationError) Unwrap() error {
	return e.Cause
}
