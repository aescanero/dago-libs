package errors

import (
	"errors"
	"testing"
)

func TestValidationError(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		message  string
		expected string
	}{
		{
			name:     "with field",
			field:    "username",
			message:  "cannot be empty",
			expected: "validation error on field 'username': cannot be empty",
		},
		{
			name:     "without field",
			field:    "",
			message:  "invalid input",
			expected: "validation error: invalid input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.field, tt.message)
			if err.Error() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, err.Error())
			}
			if err.Field != tt.field {
				t.Errorf("expected field %q, got %q", tt.field, err.Field)
			}
			if err.Message != tt.message {
				t.Errorf("expected message %q, got %q", tt.message, err.Message)
			}
		})
	}
}

func TestExecutionError(t *testing.T) {
	causeErr := errors.New("connection timeout")

	tests := []struct {
		name     string
		nodeID   string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "with cause",
			nodeID:   "node-123",
			message:  "failed to execute",
			cause:    causeErr,
			expected: "execution error at node 'node-123': failed to execute (caused by: connection timeout)",
		},
		{
			name:     "without cause",
			nodeID:   "node-456",
			message:  "timeout",
			cause:    nil,
			expected: "execution error at node 'node-456': timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewExecutionError(tt.nodeID, tt.message, tt.cause)
			if err.Error() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, err.Error())
			}

			// Test Unwrap
			if unwrapped := errors.Unwrap(err); unwrapped != tt.cause {
				t.Errorf("expected unwrapped error %v, got %v", tt.cause, unwrapped)
			}
		})
	}
}

func TestStateError(t *testing.T) {
	causeErr := errors.New("serialization failed")

	tests := []struct {
		name     string
		key      string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "with cause",
			key:      "user_data",
			message:  "cannot serialize",
			cause:    causeErr,
			expected: "state error for key 'user_data': cannot serialize (caused by: serialization failed)",
		},
		{
			name:     "without cause",
			key:      "config",
			message:  "not found",
			cause:    nil,
			expected: "state error for key 'config': not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewStateError(tt.key, tt.message, tt.cause)
			if err.Error() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, err.Error())
			}

			// Test Unwrap
			if unwrapped := errors.Unwrap(err); unwrapped != tt.cause {
				t.Errorf("expected unwrapped error %v, got %v", tt.cause, unwrapped)
			}
		})
	}
}

func TestToolError(t *testing.T) {
	causeErr := errors.New("command not found")

	tests := []struct {
		name     string
		toolName string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "with cause",
			toolName: "python",
			message:  "execution failed",
			cause:    causeErr,
			expected: "tool error 'python': execution failed (caused by: command not found)",
		},
		{
			name:     "without cause",
			toolName: "bash",
			message:  "timeout",
			cause:    nil,
			expected: "tool error 'bash': timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewToolError(tt.toolName, tt.message, tt.cause)
			if err.Error() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, err.Error())
			}

			// Test Unwrap
			if unwrapped := errors.Unwrap(err); unwrapped != tt.cause {
				t.Errorf("expected unwrapped error %v, got %v", tt.cause, unwrapped)
			}
		})
	}
}
