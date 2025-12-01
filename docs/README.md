# dago-libs Documentation

## Overview

**dago-libs** is the foundational shared library for the DA Orchestrator (Deep Agent in Go) multi-repository architecture. This library contains domain models, interfaces (ports), and common utilities that are shared across all DA Orchestrator components.

## Purpose

This library follows hexagonal architecture principles and provides:

- **Domain Models**: Core business entities (Graph, Node, State, etc.)
- **Ports (Interfaces)**: Contracts for external dependencies (LLM clients, storage, event bus, etc.)
- **JSON Schemas**: Validation schemas for graph definitions
- **Utilities**: Common helpers for logging, configuration, and tracing

**Important**: This library contains **NO implementations**. All concrete implementations belong in the main `dago` repository or node-specific repositories (`dago-node-*`).

## Architecture

```
dago-libs/
├── pkg/
│   ├── domain/         # Domain entities (no external deps)
│   │   ├── graph/      # Graph, Node, Edge definitions
│   │   ├── state/      # State management types
│   │   └── errors/     # Common error types
│   ├── ports/          # Interfaces for external dependencies
│   │   ├── llm.go      # LLM client interface
│   │   ├── tools.go    # Tool executor interface
│   │   ├── events.go   # Event bus interface
│   │   ├── storage.go  # Storage interfaces
│   │   └── metrics.go  # Metrics collector interface
│   ├── schema/         # JSON schemas + validator
│   └── utils/          # Common utilities
│       ├── logging/    # Structured logging
│       ├── config/     # Configuration loading
│       └── tracing/    # Distributed tracing
```

## Installation

```bash
go get github.com/aescanero/dago-libs@latest
```

## Usage Examples

### Working with Domain Models

```go
package main

import (
    "github.com/aescanero/dago-libs/pkg/domain/graph"
    "github.com/aescanero/dago-libs/pkg/domain/state"
)

func main() {
    // Create a new graph
    g := graph.NewGraph("my-workflow")

    // Add an executor node
    node := &graph.ExecutorNode{
        BaseNode: graph.BaseNode{
            ID:   "llm-node-1",
            Type: graph.NodeTypeExecutor,
            Name: "Generate Response",
        },
        ExecutorType: "llm",
        Config: map[string]interface{}{
            "model": "gpt-4",
            "temperature": 0.7,
        },
    }

    g.AddNode(node)
    g.EntryNode = "llm-node-1"

    // Validate the graph
    if err := g.Validate(); err != nil {
        panic(err)
    }

    // Create and manipulate state
    s := state.NewState()
    s.Set("user_input", "Hello, world!")
    s.Set("step", 1)

    if val, ok := s.GetString("user_input"); ok {
        println(val)
    }
}
```

### Using Ports (Interfaces)

```go
package main

import (
    "context"
    "github.com/aescanero/dago-libs/pkg/ports"
)

// Implement the LLMClient interface in your application
type MyLLMClient struct {
    // Your implementation details
}

func (c *MyLLMClient) Complete(ctx context.Context, req ports.CompletionRequest) (*ports.CompletionResponse, error) {
    // Your implementation
    return nil, nil
}

func (c *MyLLMClient) CompleteWithTools(ctx context.Context, req ports.CompletionRequest, tools []ports.Tool) (*ports.CompletionResponse, error) {
    // Your implementation
    return nil, nil
}

func (c *MyLLMClient) CompleteStructured(ctx context.Context, req ports.CompletionRequest, schema ports.JSONSchema) (*ports.StructuredResponse, error) {
    // Your implementation
    return nil, nil
}
```

### JSON Schema Validation

```go
package main

import (
    "github.com/aescanero/dago-libs/pkg/schema"
)

func main() {
    validator, err := schema.NewValidator()
    if err != nil {
        panic(err)
    }

    graphJSON := []byte(`{
        "id": "graph-1",
        "nodes": {
            "start": {
                "id": "start",
                "type": "executor",
                "executor_type": "llm"
            }
        },
        "entry_node": "start"
    }`)

    if err := validator.ValidateGraph(graphJSON); err != nil {
        println("Validation failed:", err.Error())
    }
}
```

### Logging

```go
package main

import (
    "github.com/aescanero/dago-libs/pkg/utils/logging"
)

func main() {
    logger := logging.NewLogger(logging.LevelInfo, "json")

    logger.Info("Starting execution", "graph_id", "my-graph")

    // Add contextual fields
    execLogger := logger.WithExecutionID("exec-123").WithGraphID("graph-1")
    execLogger.Info("Node started", "node_id", "node-1")
}
```

### Configuration

```go
package main

import (
    "github.com/aescanero/dago-libs/pkg/utils/config"
)

func main() {
    cfg := config.LoadFromEnv()

    if err := cfg.Validate(); err != nil {
        panic(err)
    }

    println("Redis address:", cfg.RedisAddr)
    println("Log level:", cfg.LogLevel)
}
```

## Key Concepts

### Graph

A **Graph** represents an execution workflow composed of nodes and edges. It defines the flow of execution through the system.

### Node

A **Node** is a unit of work within a graph. There are two main node types:
- **ExecutorNode**: Executes tasks (LLM calls, tool invocations, code execution)
- **RouterNode**: Makes routing decisions based on state conditions

### State

**State** is a flexible key-value map that flows through the graph execution. Nodes read from and write to the state as they execute.

### Ports

**Ports** are interfaces that define contracts for external dependencies. This allows for:
- Easy testing with mocks
- Swapping implementations without changing core logic
- Clear boundaries between domain and infrastructure

## Design Principles

1. **Zero Dependencies on Other DA Orchestrator Repos**: This library is the foundation and must not depend on `dago` or `dago-node-*` repositories.

2. **Minimal External Dependencies**: Only essential dependencies (UUID generation, JSON schema validation) are included.

3. **Backward Compatibility**: Changes to this library must maintain backward compatibility as it's used by multiple services.

4. **Interfaces Over Implementations**: The library provides contracts (interfaces), not implementations.

5. **Documentation First**: All exported types and functions have comprehensive godoc comments.

## For Implementers

If you're building a DA Orchestrator component:

1. **Import this library** for domain models and interfaces
2. **Implement the ports** (interfaces) for your specific needs
3. **Use the schemas** to validate graph definitions
4. **Leverage the utilities** for logging, config, and tracing

## Contributing

When adding new features to this library:

1. Ensure no dependencies on other DA Orchestrator repositories
2. Add comprehensive godoc comments
3. Include JSON schemas for new domain entities
4. Update this documentation
5. Maintain backward compatibility

## Version History

See [CHANGELOG.md](CHANGELOG.md) for version history and breaking changes.

## License

Apache 2.0 - See [LICENSE](../LICENSE) for details.
