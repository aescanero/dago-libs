package state

import (
	"testing"
)

func TestNewState(t *testing.T) {
	s := NewState()
	if s == nil {
		t.Fatal("NewState returned nil")
	}
	if len(s) != 0 {
		t.Errorf("expected empty state, got %d items", len(s))
	}
}

func TestStateSetGet(t *testing.T) {
	s := NewState()

	// Test basic set/get
	s.Set("key1", "value1")
	val := s.Get("key1")
	if val != "value1" {
		t.Errorf("expected 'value1', got %v", val)
	}

	// Test non-existent key
	val = s.Get("nonexistent")
	if val != nil {
		t.Errorf("expected nil for non-existent key, got %v", val)
	}
}

func TestStateGetString(t *testing.T) {
	s := NewState()

	s.Set("str", "hello")
	s.Set("int", 123)
	s.Set("bool", true)

	tests := []struct {
		name        string
		key         string
		expectValue string
		expectOk    bool
	}{
		{"existing string", "str", "hello", true},
		{"int value", "int", "", false},
		{"bool value", "bool", "", false},
		{"non-existent", "missing", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := s.GetString(tt.key)
			if ok != tt.expectOk {
				t.Errorf("expected ok=%v, got %v", tt.expectOk, ok)
			}
			if val != tt.expectValue {
				t.Errorf("expected value %q, got %q", tt.expectValue, val)
			}
		})
	}
}

func TestStateGetInt(t *testing.T) {
	s := NewState()

	s.Set("int", 42)
	s.Set("float", 3.14)
	s.Set("int64", int64(100))
	s.Set("str", "hello")

	tests := []struct {
		name        string
		key         string
		expectValue int
		expectOk    bool
	}{
		{"int value", "int", 42, true},
		{"float64 value", "float", 3, true},
		{"int64 value", "int64", 100, true},
		{"string value", "str", 0, false},
		{"non-existent", "missing", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := s.GetInt(tt.key)
			if ok != tt.expectOk {
				t.Errorf("expected ok=%v, got %v", tt.expectOk, ok)
			}
			if val != tt.expectValue {
				t.Errorf("expected value %d, got %d", tt.expectValue, val)
			}
		})
	}
}

func TestStateGetBool(t *testing.T) {
	s := NewState()

	s.Set("bool", true)
	s.Set("int", 1)
	s.Set("str", "true")

	tests := []struct {
		name        string
		key         string
		expectValue bool
		expectOk    bool
	}{
		{"bool true", "bool", true, true},
		{"int value", "int", false, false},
		{"string value", "str", false, false},
		{"non-existent", "missing", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := s.GetBool(tt.key)
			if ok != tt.expectOk {
				t.Errorf("expected ok=%v, got %v", tt.expectOk, ok)
			}
			if val != tt.expectValue {
				t.Errorf("expected value %v, got %v", tt.expectValue, val)
			}
		})
	}
}

func TestStateDelete(t *testing.T) {
	s := NewState()
	s.Set("key1", "value1")

	if !s.Has("key1") {
		t.Error("expected key1 to exist")
	}

	s.Delete("key1")

	if s.Has("key1") {
		t.Error("expected key1 to be deleted")
	}
}

func TestStateHas(t *testing.T) {
	s := NewState()
	s.Set("existing", "value")

	if !s.Has("existing") {
		t.Error("expected 'existing' key to be present")
	}

	if s.Has("missing") {
		t.Error("expected 'missing' key to be absent")
	}
}

func TestStateCopy(t *testing.T) {
	s := NewState()
	s.Set("key1", "value1")
	s.Set("key2", map[string]interface{}{"nested": "data"})

	copy, err := s.Copy()
	if err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	// Verify copy has same data
	if val := copy.Get("key1"); val != "value1" {
		t.Errorf("expected 'value1', got %v", val)
	}

	// Modify original - copy should not change
	s.Set("key1", "modified")

	if val := copy.Get("key1"); val != "value1" {
		t.Errorf("copy was modified when original changed")
	}
}

func TestStateMerge(t *testing.T) {
	s1 := NewState()
	s1.Set("key1", "value1")
	s1.Set("key2", "value2")

	s2 := NewState()
	s2.Set("key2", "overwritten")
	s2.Set("key3", "value3")

	s1.Merge(s2)

	// Check merged values
	if val := s1.Get("key1"); val != "value1" {
		t.Errorf("expected key1='value1', got %v", val)
	}
	if val := s1.Get("key2"); val != "overwritten" {
		t.Errorf("expected key2='overwritten', got %v", val)
	}
	if val := s1.Get("key3"); val != "value3" {
		t.Errorf("expected key3='value3', got %v", val)
	}
}

func TestStateToJSON(t *testing.T) {
	s := NewState()
	s.Set("name", "test")
	s.Set("count", 42)

	jsonStr, err := s.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON failed: %v", err)
	}

	if jsonStr == "" {
		t.Error("expected non-empty JSON string")
	}

	// Verify it's valid JSON by parsing it back
	s2 := NewState()
	if err := s2.FromJSON(jsonStr); err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if val := s2.Get("name"); val != "test" {
		t.Errorf("expected name='test', got %v", val)
	}
}

func TestStateFromJSON(t *testing.T) {
	jsonStr := `{"name":"test","count":42,"active":true}`

	s := NewState()
	if err := s.FromJSON(jsonStr); err != nil {
		t.Fatalf("FromJSON failed: %v", err)
	}

	if val, ok := s.GetString("name"); !ok || val != "test" {
		t.Errorf("expected name='test', got %v (ok=%v)", val, ok)
	}

	if val, ok := s.GetInt("count"); !ok || val != 42 {
		t.Errorf("expected count=42, got %v (ok=%v)", val, ok)
	}

	if val, ok := s.GetBool("active"); !ok || !val {
		t.Errorf("expected active=true, got %v (ok=%v)", val, ok)
	}
}

func TestStateFromJSONInvalid(t *testing.T) {
	s := NewState()
	err := s.FromJSON("invalid json")
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestStateKeys(t *testing.T) {
	s := NewState()
	s.Set("key1", "value1")
	s.Set("key2", "value2")
	s.Set("key3", "value3")

	keys := s.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}

	keyMap := make(map[string]bool)
	for _, k := range keys {
		keyMap[k] = true
	}

	if !keyMap["key1"] || !keyMap["key2"] || !keyMap["key3"] {
		t.Errorf("missing expected keys, got %v", keys)
	}
}

func TestStateSize(t *testing.T) {
	s := NewState()
	if s.Size() != 0 {
		t.Errorf("expected size 0, got %d", s.Size())
	}

	s.Set("key1", "value1")
	s.Set("key2", "value2")

	if s.Size() != 2 {
		t.Errorf("expected size 2, got %d", s.Size())
	}
}

func TestStateClear(t *testing.T) {
	s := NewState()
	s.Set("key1", "value1")
	s.Set("key2", "value2")

	if s.Size() != 2 {
		t.Errorf("expected size 2 before clear, got %d", s.Size())
	}

	s.Clear()

	if s.Size() != 0 {
		t.Errorf("expected size 0 after clear, got %d", s.Size())
	}

	if s.Has("key1") || s.Has("key2") {
		t.Error("expected all keys to be cleared")
	}
}
