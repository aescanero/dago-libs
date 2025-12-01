# dago-libs

[![Go Reference](https://pkg.go.dev/badge/github.com/aescanero/dago-libs.svg)](https://pkg.go.dev/github.com/aescanero/dago-libs)
[![Go Report Card](https://goreportcard.com/badge/github.com/aescanero/dago-libs)](https://goreportcard.com/report/github.com/aescanero/dago-libs)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

**Shared libraries foundation for DA Orchestrator (Deep Agent in Go)**

## Overview

`dago-libs` is the foundational library for the DA Orchestrator multi-repository architecture. It provides domain models, interface definitions (ports), and common utilities following hexagonal/clean architecture principles.

**Key Principle**: This library contains **only interfaces and domain models** - no implementations. All concrete implementations belong in the main [`dago`](https://github.com/aescanero/dago) repository or node-specific repositories.

## Features

- **Domain Models**: Graph, Node, Edge, State, and Error types
- **Ports (Interfaces)**: Contracts for LLM clients, tool executors, event bus, storage, and metrics
- **JSON Schemas**: Validation for graph definitions and node configurations
- **Utilities**: Structured logging, configuration loading, and distributed tracing helpers
- **Zero Implementation Dependencies**: Pure domain layer with minimal external dependencies

## Installation

```bash
go get github.com/aescanero/dago-libs@latest
```

## Quick Start

```go
package main

import (
    "github.com/aescanero/dago-libs/pkg/domain/graph"
    "github.com/aescanero/dago-libs/pkg/domain/state"
)

func main() {
    // Create a graph
    g := graph.NewGraph("my-workflow")

    // Add a node
    node := &graph.ExecutorNode{
        BaseNode: graph.BaseNode{
            ID:   "llm-1",
            Type: graph.NodeTypeExecutor,
        },
        ExecutorType: "llm",
        Config: map[string]interface{}{
            "model": "gpt-4",
        },
    }

    g.AddNode(node)
    g.EntryNode = "llm-1"

    // Validate
    if err := g.Validate(); err != nil {
        panic(err)
    }

    // Work with state
    s := state.NewState()
    s.Set("input", "Hello!")

    if val, ok := s.GetString("input"); ok {
        println(val)
    }
}
```

## Architecture

```
pkg/
├── domain/          # Pure domain models (no external deps)
│   ├── graph/       # Graph, Node, Edge
│   ├── state/       # State management
│   └── errors/      # Error types
├── ports/           # Interface definitions
│   ├── llm.go       # LLM client interface
│   ├── tools.go     # Tool executor interface
│   ├── events.go    # Event bus interface
│   ├── storage.go   # Storage interfaces
│   └── metrics.go   # Metrics collector interface
├── schema/          # JSON schemas + validator
└── utils/           # Common utilities
    ├── logging/     # Structured logging
    ├── config/      # Configuration
    └── tracing/     # Distributed tracing
```

## Documentation

- **[Full Documentation](docs/README.md)**: Complete usage guide with examples
- **[Changelog](docs/CHANGELOG.md)**: Version history and breaking changes
- **[GoDoc](https://pkg.go.dev/github.com/aescanero/dago-libs)**: API reference

## Design Principles

1. **Foundation First**: Zero dependencies on other DA Orchestrator repositories
2. **Interfaces Over Implementations**: Define contracts, not concrete code
3. **Minimal Dependencies**: Only essential external packages
4. **Backward Compatibility**: API stability is critical
5. **Documentation**: Comprehensive godoc for all exports

## For Implementers

When building DA Orchestrator components:

1. Import `dago-libs` for domain models and interfaces
2. Implement the ports (interfaces) for your specific use case
3. Use the JSON schemas to validate graph definitions
4. Leverage utilities for logging, config, and tracing

Example implementation:

```go
import "github.com/aescanero/dago-libs/pkg/ports"

type MyLLMClient struct {
    // Your implementation
}

func (c *MyLLMClient) Complete(ctx context.Context, req ports.CompletionRequest) (*ports.CompletionResponse, error) {
    // Your LLM integration code
}
```

## Development

```bash
# Install dependencies
make deps

# Run tests
make test

# Run tests with coverage
make test-coverage

# Lint code
make lint

# Format code
make fmt

# Create a release
make release VERSION=v1.0.0
```

## Key Components

### Domain Layer

- **Graph**: Execution workflow with nodes and edges
- **Node**: Execution units (Executor, Router, Start, End)
- **State**: Key-value map for execution state
- **Errors**: Structured error types

### Ports Layer

- **LLMClient**: Interface for LLM providers (OpenAI, Anthropic, etc.)
- **ToolExecutor**: Interface for executing tools (Python, Bash, HTTP, etc.)
- **EventBus**: Interface for pub/sub events (Redis Streams for MVP)
- **StateStorage**: Interface for state persistence (Redis for MVP)
- **MetricsCollector**: Interface for metrics (Prometheus for MVP)

### Schema Layer

- JSON schemas for graph, executor nodes, and router nodes
- Validator using `jsonschema/v5`
- Embedded schemas (no external files needed)

### Utils Layer

- **Logging**: Structured logging with `slog`
- **Config**: Environment variable loading with type conversion
- **Tracing**: Basic distributed tracing (OpenTelemetry-compatible)

## Contributing

Contributions are welcome! When adding features:

1. Maintain backward compatibility
2. Add comprehensive godoc comments
3. Include tests for new functionality
4. Update documentation
5. Ensure no dependencies on other DA Orchestrator repos

## Related Repositories

- [`dago`](https://github.com/aescanero/dago): Main orchestrator implementation
- `dago-node-executor`: Executor node workers (coming soon)
- `dago-node-router`: Router node workers (coming soon)

## License

Apache 2.0 - See [LICENSE](LICENSE) for details.

## Project

- **Domain**: disasterproject.com
- **Organization**: DA Orchestrator
- **Author**: aescanero
- **Repository**: [github.com/aescanero/dago-libs](https://github.com/aescanero/dago-libs)
