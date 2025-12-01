// Package errors defines common error types used across the DA Orchestrator system.
package errors

import "fmt"

// ValidationError represents an error that occurs during validation of domain entities.
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ExecutionError represents an error that occurs during graph or node execution.
type ExecutionError struct {
	NodeID  string
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *ExecutionError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("execution error at node '%s': %s (caused by: %v)", e.NodeID, e.Message, e.Cause)
	}
	return fmt.Sprintf("execution error at node '%s': %s", e.NodeID, e.Message)
}

// Unwrap implements the errors.Unwrap interface.
func (e *ExecutionError) Unwrap() error {
	return e.Cause
}

// NewExecutionError creates a new ExecutionError.
func NewExecutionError(nodeID, message string, cause error) *ExecutionError {
	return &ExecutionError{
		NodeID:  nodeID,
		Message: message,
		Cause:   cause,
	}
}

// StateError represents an error related to state management.
type StateError struct {
	Key     string
	Message string
	Cause   error
}

// Error implements the error interface.
func (e *StateError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("state error for key '%s': %s (caused by: %v)", e.Key, e.Message, e.Cause)
	}
	return fmt.Sprintf("state error for key '%s': %s", e.Key, e.Message)
}

// Unwrap implements the errors.Unwrap interface.
func (e *StateError) Unwrap() error {
	return e.Cause
}

// NewStateError creates a new StateError.
func NewStateError(key, message string, cause error) *StateError {
	return &StateError{
		Key:     key,
		Message: message,
		Cause:   cause,
	}
}

// ToolError represents an error that occurs during tool execution.
type ToolError struct {
	ToolName string
	Message  string
	Cause    error
}

// Error implements the error interface.
func (e *ToolError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("tool error '%s': %s (caused by: %v)", e.ToolName, e.Message, e.Cause)
	}
	return fmt.Sprintf("tool error '%s': %s", e.ToolName, e.Message)
}

// Unwrap implements the errors.Unwrap interface.
func (e *ToolError) Unwrap() error {
	return e.Cause
}

// NewToolError creates a new ToolError.
func NewToolError(toolName, message string, cause error) *ToolError {
	return &ToolError{
		ToolName: toolName,
		Message:  message,
		Cause:    cause,
	}
}
