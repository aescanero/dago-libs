# Changelog

All notable changes to dago-libs will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial library structure
- Domain layer with Graph, Node, State, and Error types
- Ports layer with interfaces for LLM, Tools, Events, Storage, and Metrics
- JSON schemas for graph validation
- Schema validator using jsonschema/v5
- Utilities for logging (slog-based), configuration, and tracing
- Comprehensive documentation

## [1.0.0] - TBD

### Added
- Initial release of dago-libs
- Domain models for graph execution (Graph, Node, Edge, State)
- Port interfaces following hexagonal architecture
  - LLMClient: Interface for LLM providers
  - ToolExecutor: Interface for tool execution
  - EventBus: Interface for event pub/sub (Redis Streams)
  - StateStorage: Interface for state persistence (Redis)
  - MetricsCollector: Interface for metrics collection (Prometheus)
- JSON schemas for validation
  - Graph schema
  - Executor node schema
  - Router node schema
- Utility packages
  - Structured logging with slog
  - Configuration from environment variables
  - Basic distributed tracing
- Complete documentation and examples

### Design Decisions
- Zero dependencies on other DA Orchestrator repositories
- Minimal external dependencies (uuid, jsonschema)
- Interface-first design for easy testing and swapping implementations
- Embedded JSON schemas for runtime validation
- MVP focus on Redis-based implementations for storage and events

### Notes
- This library contains only domain models and interfaces
- Concrete implementations should be in `dago` or `dago-node-*` repositories
- Backward compatibility will be maintained in all future releases
