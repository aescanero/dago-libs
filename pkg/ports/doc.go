// Package ports defines the interfaces (ports) for external dependencies.
//
// Following hexagonal architecture principles, this package contains only
// interface definitions. Concrete implementations should be in separate
// repositories (dago, dago-node-*) that depend on this library.
//
// Ports include:
//   - LLMClient: Interface for Large Language Model providers
//   - ToolExecutor: Interface for executing tools (Python, Bash, HTTP, etc.)
//   - EventBus: Interface for event publishing and subscription (Redis Streams)
//   - StateStorage: Interface for persisting execution state (Redis)
//   - MetricsCollector: Interface for collecting system metrics (Prometheus)
//
// This design allows for:
//   - Easy testing with mock implementations
//   - Swapping implementations without changing core logic
//   - Clear boundaries between domain logic and infrastructure
package ports
