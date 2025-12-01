// Package tracing provides basic distributed tracing utilities for DA Orchestrator.
//
// This is a simple implementation for MVP purposes. For production use, consider
// integrating with OpenTelemetry or similar distributed tracing systems.
//
// Example usage:
//
//	tracer := tracing.NewTracer("dago")
//	span, ctx := tracer.StartSpan(context.Background(), "execute-graph")
//	defer tracer.EndSpan(span)
//
//	span.SetTag("graph_id", "my-graph")
//	span.AddEvent("node-started", map[string]string{"node_id": "node-1"})
//
//	// Pass ctx to child operations to propagate trace context
//	childSpan, childCtx := tracer.StartSpan(ctx, "execute-node")
//	defer tracer.EndSpan(childSpan)
package tracing
