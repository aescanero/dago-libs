// Package config provides configuration loading utilities.
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// GetEnv retrieves an environment variable with a default fallback.
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt retrieves an integer environment variable with a default fallback.
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

// GetEnvBool retrieves a boolean environment variable with a default fallback.
// Accepts: "true", "1", "yes", "on" for true (case-insensitive).
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		switch value {
		case "true", "1", "yes", "on", "TRUE", "YES", "ON":
			return true
		case "false", "0", "no", "off", "FALSE", "NO", "OFF":
			return false
		}
	}
	return defaultValue
}

// GetEnvDuration retrieves a duration environment variable with a default fallback.
// Duration should be specified as a string parseable by time.ParseDuration (e.g., "5s", "10m").
func GetEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// RequireEnv retrieves a required environment variable.
// Returns an error if the variable is not set.
func RequireEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return value, nil
}

// Config represents common configuration options for DA Orchestrator components.
type Config struct {
	// Redis configuration
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Logging configuration
	LogLevel  string
	LogFormat string

	// Metrics configuration
	MetricsEnabled bool
	MetricsPort    int

	// Service configuration
	ServiceName string
	ServicePort int

	// Timeouts
	DefaultTimeout time.Duration
	LLMTimeout     time.Duration
	ToolTimeout    time.Duration
}

// LoadFromEnv loads configuration from environment variables.
func LoadFromEnv() Config {
	return Config{
		// Redis
		RedisAddr:     GetEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: GetEnv("REDIS_PASSWORD", ""),
		RedisDB:       GetEnvInt("REDIS_DB", 0),

		// Logging
		LogLevel:  GetEnv("LOG_LEVEL", "info"),
		LogFormat: GetEnv("LOG_FORMAT", "text"),

		// Metrics
		MetricsEnabled: GetEnvBool("METRICS_ENABLED", true),
		MetricsPort:    GetEnvInt("METRICS_PORT", 9090),

		// Service
		ServiceName: GetEnv("SERVICE_NAME", "dago"),
		ServicePort: GetEnvInt("SERVICE_PORT", 8080),

		// Timeouts
		DefaultTimeout: GetEnvDuration("DEFAULT_TIMEOUT", 5*time.Minute),
		LLMTimeout:     GetEnvDuration("LLM_TIMEOUT", 2*time.Minute),
		ToolTimeout:    GetEnvDuration("TOOL_TIMEOUT", 5*time.Minute),
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.RedisAddr == "" {
		return fmt.Errorf("redis address cannot be empty")
	}
	if c.ServicePort <= 0 || c.ServicePort > 65535 {
		return fmt.Errorf("service port must be between 1 and 65535")
	}
	if c.MetricsPort <= 0 || c.MetricsPort > 65535 {
		return fmt.Errorf("metrics port must be between 1 and 65535")
	}
	return nil
}
