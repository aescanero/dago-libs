// Package state provides types and interfaces for managing execution state in the DA Orchestrator.
//
// The State type is a flexible key-value map that can store any JSON-serializable data.
// The Manager interface defines operations for state lifecycle management including
// initialization, updates, snapshots, and cleanup.
//
// State is the fundamental data structure that flows through the graph execution,
// being read and modified by nodes as the execution progresses.
package state
