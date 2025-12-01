// Package state provides types and interfaces for managing execution state.
package state

import (
	"encoding/json"
	"fmt"
)

// State represents the execution state as a flexible key-value map.
// It can store any JSON-serializable data.
type State map[string]interface{}

// NewState creates a new empty State.
func NewState() State {
	return make(State)
}

// Get retrieves a value from the state by key.
// Returns nil if the key doesn't exist.
func (s State) Get(key string) interface{} {
	return s[key]
}

// GetString retrieves a string value from the state.
// Returns empty string and false if the key doesn't exist or value is not a string.
func (s State) GetString(key string) (string, bool) {
	val, ok := s[key]
	if !ok {
		return "", false
	}
	str, ok := val.(string)
	return str, ok
}

// GetInt retrieves an int value from the state.
// Returns 0 and false if the key doesn't exist or value is not convertible to int.
func (s State) GetInt(key string) (int, bool) {
	val, ok := s[key]
	if !ok {
		return 0, false
	}

	switch v := val.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true
	case int64:
		return int(v), true
	default:
		return 0, false
	}
}

// GetBool retrieves a boolean value from the state.
// Returns false and false if the key doesn't exist or value is not a bool.
func (s State) GetBool(key string) (bool, bool) {
	val, ok := s[key]
	if !ok {
		return false, false
	}
	b, ok := val.(bool)
	return b, ok
}

// Set stores a value in the state.
func (s State) Set(key string, value interface{}) {
	s[key] = value
}

// Delete removes a key from the state.
func (s State) Delete(key string) {
	delete(s, key)
}

// Has checks if a key exists in the state.
func (s State) Has(key string) bool {
	_, ok := s[key]
	return ok
}

// Copy creates a deep copy of the state.
// Note: This uses JSON marshaling/unmarshaling for deep copy.
// TODO: Consider more efficient deep copy implementation for performance-critical paths.
func (s State) Copy() (State, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal state: %w", err)
	}

	var copy State
	if err := json.Unmarshal(data, &copy); err != nil {
		return nil, fmt.Errorf("failed to unmarshal state: %w", err)
	}

	return copy, nil
}

// Merge merges another state into this state.
// Values from the other state will overwrite existing values.
func (s State) Merge(other State) {
	for k, v := range other {
		s[k] = v
	}
}

// ToJSON converts the state to a JSON string.
func (s State) ToJSON() (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal state to JSON: %w", err)
	}
	return string(data), nil
}

// FromJSON populates the state from a JSON string.
func (s State) FromJSON(jsonStr string) error {
	if err := json.Unmarshal([]byte(jsonStr), &s); err != nil {
		return fmt.Errorf("failed to unmarshal JSON to state: %w", err)
	}
	return nil
}

// Keys returns all keys in the state.
func (s State) Keys() []string {
	keys := make([]string, 0, len(s))
	for k := range s {
		keys = append(keys, k)
	}
	return keys
}

// Size returns the number of entries in the state.
func (s State) Size() int {
	return len(s)
}

// Clear removes all entries from the state.
func (s State) Clear() {
	for k := range s {
		delete(s, k)
	}
}
