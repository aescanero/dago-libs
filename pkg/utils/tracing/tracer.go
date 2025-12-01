// Package tracing provides distributed tracing utilities.
package tracing

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// SpanContext represents the context of a trace span.
type SpanContext struct {
	// TraceID uniquely identifies the entire trace.
	TraceID string

	// SpanID uniquely identifies this span within the trace.
	SpanID string

	// ParentSpanID is the ID of the parent span (empty for root spans).
	ParentSpanID string
}

// Span represents a unit of work in a distributed trace.
type Span struct {
	Context   SpanContext
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Tags      map[string]string
	Events    []SpanEvent
	Status    SpanStatus
}

// SpanEvent represents a point-in-time event within a span.
type SpanEvent struct {
	Name       string
	Timestamp  time.Time
	Attributes map[string]string
}

// SpanStatus represents the status of a span.
type SpanStatus string

const (
	// SpanStatusUnset indicates the span status is not set.
	SpanStatusUnset SpanStatus = "unset"

	// SpanStatusOK indicates the span completed successfully.
	SpanStatusOK SpanStatus = "ok"

	// SpanStatusError indicates the span encountered an error.
	SpanStatusError SpanStatus = "error"
)

// Tracer provides basic tracing functionality.
// For MVP, this is a simple implementation. Production use should integrate
// with OpenTelemetry or similar tracing systems.
type Tracer struct {
	serviceName string
}

// NewTracer creates a new tracer.
func NewTracer(serviceName string) *Tracer {
	return &Tracer{
		serviceName: serviceName,
	}
}

// StartSpan creates a new span.
func (t *Tracer) StartSpan(ctx context.Context, name string) (*Span, context.Context) {
	spanID := uuid.New().String()

	// Check if there's a parent span in the context
	parentSpan, _ := ctx.Value(spanContextKey).(*SpanContext)

	var traceID, parentSpanID string
	if parentSpan != nil {
		traceID = parentSpan.TraceID
		parentSpanID = parentSpan.SpanID
	} else {
		// Root span - generate new trace ID
		traceID = uuid.New().String()
	}

	spanCtx := SpanContext{
		TraceID:      traceID,
		SpanID:       spanID,
		ParentSpanID: parentSpanID,
	}

	span := &Span{
		Context:   spanCtx,
		Name:      name,
		StartTime: time.Now(),
		Tags:      make(map[string]string),
		Events:    make([]SpanEvent, 0),
		Status:    SpanStatusUnset,
	}

	// Add service name tag
	span.Tags["service.name"] = t.serviceName

	// Store span context in the returned context
	newCtx := context.WithValue(ctx, spanContextKey, &spanCtx)

	return span, newCtx
}

// EndSpan marks a span as complete.
func (t *Tracer) EndSpan(span *Span) {
	span.EndTime = time.Now()
	if span.Status == SpanStatusUnset {
		span.Status = SpanStatusOK
	}
	// TODO: In production, export span to tracing backend here
}

// SetTag adds a tag to the span.
func (s *Span) SetTag(key, value string) {
	s.Tags[key] = value
}

// AddEvent adds an event to the span.
func (s *Span) AddEvent(name string, attributes map[string]string) {
	s.Events = append(s.Events, SpanEvent{
		Name:       name,
		Timestamp:  time.Now(),
		Attributes: attributes,
	})
}

// SetStatus sets the status of the span.
func (s *Span) SetStatus(status SpanStatus) {
	s.Status = status
}

// SetError marks the span as errored and records the error message.
func (s *Span) SetError(err error) {
	s.Status = SpanStatusError
	s.Tags["error"] = "true"
	s.Tags["error.message"] = err.Error()
}

// Duration returns the duration of the span.
func (s *Span) Duration() time.Duration {
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// Context key for storing span context
type contextKey string

const spanContextKey contextKey = "span_context"

// SpanFromContext retrieves the span context from a context.
func SpanFromContext(ctx context.Context) *SpanContext {
	spanCtx, _ := ctx.Value(spanContextKey).(*SpanContext)
	return spanCtx
}

// ExtractTraceID extracts the trace ID from a context.
func ExtractTraceID(ctx context.Context) string {
	if spanCtx := SpanFromContext(ctx); spanCtx != nil {
		return spanCtx.TraceID
	}
	return ""
}

// ExtractSpanID extracts the span ID from a context.
func ExtractSpanID(ctx context.Context) string {
	if spanCtx := SpanFromContext(ctx); spanCtx != nil {
		return spanCtx.SpanID
	}
	return ""
}
