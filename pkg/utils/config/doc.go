// Package config provides utilities for loading configuration from environment variables.
//
// This package includes helper functions for reading environment variables with type
// conversion and default values, as well as a standard Config struct for common
// DA Orchestrator configuration options.
//
// Example usage:
//
//	cfg := config.LoadFromEnv()
//	if err := cfg.Validate(); err != nil {
//		log.Fatal(err)
//	}
//
//	// Use individual helpers
//	timeout := config.GetEnvDuration("TIMEOUT", 30*time.Second)
//	enabled := config.GetEnvBool("FEATURE_ENABLED", false)
package config
