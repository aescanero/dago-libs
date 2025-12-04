package config

import (
	"os"
	"testing"
	"time"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_ENV_VAR"
	value := "test_value"
	defaultValue := "default"

	// Test with environment variable set
	_ = os.Setenv(key, value)
	defer func() { _ = os.Unsetenv(key) }()

	result := GetEnv(key, defaultValue)
	if result != value {
		t.Errorf("expected %q, got %q", value, result)
	}

	// Test with environment variable not set
	_ = os.Unsetenv(key)
	result = GetEnv(key, defaultValue)
	if result != defaultValue {
		t.Errorf("expected default %q, got %q", defaultValue, result)
	}
}

func TestGetEnvInt(t *testing.T) {
	key := "TEST_INT_VAR"
	defaultValue := 42

	tests := []struct {
		name     string
		envValue string
		expected int
	}{
		{"valid int", "123", 123},
		{"invalid int", "not_a_number", defaultValue},
		{"empty string", "", defaultValue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				_ = os.Setenv(key, tt.envValue)
				defer func() { _ = os.Unsetenv(key) }()
			} else {
				_ = os.Unsetenv(key)
			}

			result := GetEnvInt(key, defaultValue)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGetEnvBool(t *testing.T) {
	key := "TEST_BOOL_VAR"

	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{"true lowercase", "true", false, true},
		{"TRUE uppercase", "TRUE", false, true},
		{"1", "1", false, true},
		{"yes", "yes", false, true},
		{"on", "on", false, true},
		{"false lowercase", "false", true, false},
		{"FALSE uppercase", "FALSE", true, false},
		{"0", "0", true, false},
		{"no", "no", true, false},
		{"off", "off", true, false},
		{"invalid value uses default", "invalid", true, true},
		{"empty uses default", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				_ = os.Setenv(key, tt.envValue)
				defer func() { _ = os.Unsetenv(key) }()
			} else {
				_ = os.Unsetenv(key)
			}

			result := GetEnvBool(key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetEnvDuration(t *testing.T) {
	key := "TEST_DURATION_VAR"
	defaultValue := 5 * time.Minute

	tests := []struct {
		name     string
		envValue string
		expected time.Duration
	}{
		{"valid duration seconds", "30s", 30 * time.Second},
		{"valid duration minutes", "10m", 10 * time.Minute},
		{"valid duration hours", "2h", 2 * time.Hour},
		{"invalid duration", "invalid", defaultValue},
		{"empty string", "", defaultValue},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				_ = os.Setenv(key, tt.envValue)
				defer func() { _ = os.Unsetenv(key) }()
			} else {
				_ = os.Unsetenv(key)
			}

			result := GetEnvDuration(key, defaultValue)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRequireEnv(t *testing.T) {
	key := "TEST_REQUIRED_VAR"
	value := "required_value"

	// Test with variable set
	_ = os.Setenv(key, value)
	defer func() { _ = os.Unsetenv(key) }()

	result, err := RequireEnv(key)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if result != value {
		t.Errorf("expected %q, got %q", value, result)
	}

	// Test with variable not set
	_ = os.Unsetenv(key)
	result, err = RequireEnv(key)
	if err == nil {
		t.Error("expected error for missing required variable")
	}
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestLoadFromEnv(t *testing.T) {
	// Set up test environment variables
	_ = os.Setenv("REDIS_ADDR", "redis.example.com:6379")
	_ = os.Setenv("REDIS_PASSWORD", "secret")
	_ = os.Setenv("REDIS_DB", "1")
	_ = os.Setenv("LOG_LEVEL", "debug")
	_ = os.Setenv("LOG_FORMAT", "json")
	_ = os.Setenv("METRICS_ENABLED", "true")
	_ = os.Setenv("METRICS_PORT", "9091")
	_ = os.Setenv("SERVICE_NAME", "test-service")
	_ = os.Setenv("SERVICE_PORT", "8081")
	_ = os.Setenv("DEFAULT_TIMEOUT", "10m")
	_ = os.Setenv("LLM_TIMEOUT", "3m")
	_ = os.Setenv("TOOL_TIMEOUT", "15m")

	defer func() {
		_ = os.Unsetenv("REDIS_ADDR")
		_ = os.Unsetenv("REDIS_PASSWORD")
		_ = os.Unsetenv("REDIS_DB")
		_ = os.Unsetenv("LOG_LEVEL")
		_ = os.Unsetenv("LOG_FORMAT")
		_ = os.Unsetenv("METRICS_ENABLED")
		_ = os.Unsetenv("METRICS_PORT")
		_ = os.Unsetenv("SERVICE_NAME")
		_ = os.Unsetenv("SERVICE_PORT")
		_ = os.Unsetenv("DEFAULT_TIMEOUT")
		_ = os.Unsetenv("LLM_TIMEOUT")
		_ = os.Unsetenv("TOOL_TIMEOUT")
	}()

	cfg := LoadFromEnv()

	if cfg.RedisAddr != "redis.example.com:6379" {
		t.Errorf("expected RedisAddr 'redis.example.com:6379', got %q", cfg.RedisAddr)
	}
	if cfg.RedisPassword != "secret" {
		t.Errorf("expected RedisPassword 'secret', got %q", cfg.RedisPassword)
	}
	if cfg.RedisDB != 1 {
		t.Errorf("expected RedisDB 1, got %d", cfg.RedisDB)
	}
	if cfg.LogLevel != "debug" {
		t.Errorf("expected LogLevel 'debug', got %q", cfg.LogLevel)
	}
	if cfg.LogFormat != "json" {
		t.Errorf("expected LogFormat 'json', got %q", cfg.LogFormat)
	}
	if !cfg.MetricsEnabled {
		t.Error("expected MetricsEnabled to be true")
	}
	if cfg.MetricsPort != 9091 {
		t.Errorf("expected MetricsPort 9091, got %d", cfg.MetricsPort)
	}
	if cfg.ServiceName != "test-service" {
		t.Errorf("expected ServiceName 'test-service', got %q", cfg.ServiceName)
	}
	if cfg.ServicePort != 8081 {
		t.Errorf("expected ServicePort 8081, got %d", cfg.ServicePort)
	}
	if cfg.DefaultTimeout != 10*time.Minute {
		t.Errorf("expected DefaultTimeout 10m, got %v", cfg.DefaultTimeout)
	}
	if cfg.LLMTimeout != 3*time.Minute {
		t.Errorf("expected LLMTimeout 3m, got %v", cfg.LLMTimeout)
	}
	if cfg.ToolTimeout != 15*time.Minute {
		t.Errorf("expected ToolTimeout 15m, got %v", cfg.ToolTimeout)
	}
}

func TestLoadFromEnv_Defaults(t *testing.T) {
	// Clear all relevant environment variables
	vars := []string{
		"REDIS_ADDR", "REDIS_PASSWORD", "REDIS_DB",
		"LOG_LEVEL", "LOG_FORMAT",
		"METRICS_ENABLED", "METRICS_PORT",
		"SERVICE_NAME", "SERVICE_PORT",
		"DEFAULT_TIMEOUT", "LLM_TIMEOUT", "TOOL_TIMEOUT",
	}
	for _, v := range vars {
		_ = os.Unsetenv(v)
	}

	cfg := LoadFromEnv()

	// Verify defaults
	if cfg.RedisAddr != "localhost:6379" {
		t.Errorf("expected default RedisAddr 'localhost:6379', got %q", cfg.RedisAddr)
	}
	if cfg.RedisDB != 0 {
		t.Errorf("expected default RedisDB 0, got %d", cfg.RedisDB)
	}
	if cfg.LogLevel != "info" {
		t.Errorf("expected default LogLevel 'info', got %q", cfg.LogLevel)
	}
	if cfg.LogFormat != "text" {
		t.Errorf("expected default LogFormat 'text', got %q", cfg.LogFormat)
	}
	if !cfg.MetricsEnabled {
		t.Error("expected default MetricsEnabled to be true")
	}
	if cfg.MetricsPort != 9090 {
		t.Errorf("expected default MetricsPort 9090, got %d", cfg.MetricsPort)
	}
	if cfg.ServiceName != "dago" {
		t.Errorf("expected default ServiceName 'dago', got %q", cfg.ServiceName)
	}
	if cfg.ServicePort != 8080 {
		t.Errorf("expected default ServicePort 8080, got %d", cfg.ServicePort)
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "valid config",
			config: Config{
				RedisAddr:   "localhost:6379",
				ServicePort: 8080,
				MetricsPort: 9090,
			},
			expectError: false,
		},
		{
			name: "empty redis addr",
			config: Config{
				RedisAddr:   "",
				ServicePort: 8080,
				MetricsPort: 9090,
			},
			expectError: true,
		},
		{
			name: "invalid service port (too low)",
			config: Config{
				RedisAddr:   "localhost:6379",
				ServicePort: 0,
				MetricsPort: 9090,
			},
			expectError: true,
		},
		{
			name: "invalid service port (too high)",
			config: Config{
				RedisAddr:   "localhost:6379",
				ServicePort: 70000,
				MetricsPort: 9090,
			},
			expectError: true,
		},
		{
			name: "invalid metrics port (too low)",
			config: Config{
				RedisAddr:   "localhost:6379",
				ServicePort: 8080,
				MetricsPort: -1,
			},
			expectError: true,
		},
		{
			name: "invalid metrics port (too high)",
			config: Config{
				RedisAddr:   "localhost:6379",
				ServicePort: 8080,
				MetricsPort: 100000,
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError && err == nil {
				t.Error("expected validation error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected validation error: %v", err)
			}
		})
	}
}
