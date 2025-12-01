package tracing

import (
	"context"
	"testing"
	"time"
)

func TestNewTracer(t *testing.T) {
	tracer := NewTracer("test-service")
	if tracer == nil {
		t.Fatal("NewTracer returned nil")
	}
	if tracer.serviceName != "test-service" {
		t.Errorf("expected service name 'test-service', got %q", tracer.serviceName)
	}
}

func TestTracer_StartSpan(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	span, newCtx := tracer.StartSpan(ctx, "test-operation")

	if span == nil {
		t.Fatal("StartSpan returned nil span")
	}
	if newCtx == nil {
		t.Fatal("StartSpan returned nil context")
	}

	// Verify span fields
	if span.Name != "test-operation" {
		t.Errorf("expected span name 'test-operation', got %q", span.Name)
	}
	if span.Context.TraceID == "" {
		t.Error("span should have a trace ID")
	}
	if span.Context.SpanID == "" {
		t.Error("span should have a span ID")
	}
	if span.Context.ParentSpanID != "" {
		t.Error("root span should not have a parent span ID")
	}
	if span.Status != SpanStatusUnset {
		t.Errorf("expected initial status %q, got %q", SpanStatusUnset, span.Status)
	}

	// Verify service name tag was added
	if span.Tags["service.name"] != "test-service" {
		t.Errorf("expected service.name tag 'test-service', got %q", span.Tags["service.name"])
	}
}

func TestTracer_StartSpan_WithParent(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	// Create parent span
	parentSpan, parentCtx := tracer.StartSpan(ctx, "parent-operation")

	// Create child span
	childSpan, _ := tracer.StartSpan(parentCtx, "child-operation")

	// Verify child has parent relationship
	if childSpan.Context.TraceID != parentSpan.Context.TraceID {
		t.Error("child span should have same trace ID as parent")
	}
	if childSpan.Context.ParentSpanID != parentSpan.Context.SpanID {
		t.Error("child span parent ID should match parent span ID")
	}
}

func TestTracer_EndSpan(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	span, _ := tracer.StartSpan(ctx, "test-operation")

	// Wait a bit to ensure duration is non-zero
	time.Sleep(10 * time.Millisecond)

	tracer.EndSpan(span)

	if span.EndTime.IsZero() {
		t.Error("span should have end time set")
	}
	if span.Status == SpanStatusUnset {
		t.Error("span status should be set after ending")
	}
	if span.Status != SpanStatusOK {
		t.Errorf("expected status %q, got %q", SpanStatusOK, span.Status)
	}
}

func TestSpan_SetTag(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	span, _ := tracer.StartSpan(ctx, "test-operation")

	span.SetTag("key1", "value1")
	span.SetTag("key2", "value2")

	if span.Tags["key1"] != "value1" {
		t.Errorf("expected tag 'key1'='value1', got %q", span.Tags["key1"])
	}
	if span.Tags["key2"] != "value2" {
		t.Errorf("expected tag 'key2'='value2', got %q", span.Tags["key2"])
	}
}

func TestSpan_AddEvent(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	span, _ := tracer.StartSpan(ctx, "test-operation")

	attrs := map[string]string{
		"attr1": "value1",
		"attr2": "value2",
	}
	span.AddEvent("test-event", attrs)

	if len(span.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(span.Events))
	}

	event := span.Events[0]
	if event.Name != "test-event" {
		t.Errorf("expected event name 'test-event', got %q", event.Name)
	}
	if event.Attributes["attr1"] != "value1" {
		t.Errorf("expected attribute 'attr1'='value1', got %q", event.Attributes["attr1"])
	}
	if event.Timestamp.IsZero() {
		t.Error("event should have a timestamp")
	}
}

func TestSpan_SetStatus(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	span, _ := tracer.StartSpan(ctx, "test-operation")

	span.SetStatus(SpanStatusError)

	if span.Status != SpanStatusError {
		t.Errorf("expected status %q, got %q", SpanStatusError, span.Status)
	}
}

func TestSpan_SetError(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	span, _ := tracer.StartSpan(ctx, "test-operation")

	testErr := context.DeadlineExceeded
	span.SetError(testErr)

	if span.Status != SpanStatusError {
		t.Errorf("expected status %q, got %q", SpanStatusError, span.Status)
	}
	if span.Tags["error"] != "true" {
		t.Errorf("expected error tag 'true', got %q", span.Tags["error"])
	}
	if span.Tags["error.message"] != testErr.Error() {
		t.Errorf("expected error.message %q, got %q", testErr.Error(), span.Tags["error.message"])
	}
}

func TestSpan_Duration(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	span, _ := tracer.StartSpan(ctx, "test-operation")

	// Test duration before end
	time.Sleep(10 * time.Millisecond)
	duration1 := span.Duration()
	if duration1 < 10*time.Millisecond {
		t.Error("duration should be at least 10ms")
	}

	// Test duration after end
	tracer.EndSpan(span)
	duration2 := span.Duration()
	if duration2.Milliseconds() <= 0 {
		t.Error("duration after end should be positive")
	}
}

func TestSpanFromContext(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	// Test with no span in context
	spanCtx := SpanFromContext(ctx)
	if spanCtx != nil {
		t.Error("expected nil span context from empty context")
	}

	// Test with span in context
	_, newCtx := tracer.StartSpan(ctx, "test-operation")
	spanCtx = SpanFromContext(newCtx)
	if spanCtx == nil {
		t.Fatal("expected non-nil span context")
	}
	if spanCtx.TraceID == "" {
		t.Error("span context should have trace ID")
	}
	if spanCtx.SpanID == "" {
		t.Error("span context should have span ID")
	}
}

func TestExtractTraceID(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	// Test with no trace
	traceID := ExtractTraceID(ctx)
	if traceID != "" {
		t.Errorf("expected empty trace ID, got %q", traceID)
	}

	// Test with trace
	span, newCtx := tracer.StartSpan(ctx, "test-operation")
	traceID = ExtractTraceID(newCtx)
	if traceID != span.Context.TraceID {
		t.Errorf("expected trace ID %q, got %q", span.Context.TraceID, traceID)
	}
}

func TestExtractSpanID(t *testing.T) {
	tracer := NewTracer("test-service")
	ctx := context.Background()

	// Test with no span
	spanID := ExtractSpanID(ctx)
	if spanID != "" {
		t.Errorf("expected empty span ID, got %q", spanID)
	}

	// Test with span
	span, newCtx := tracer.StartSpan(ctx, "test-operation")
	spanID = ExtractSpanID(newCtx)
	if spanID != span.Context.SpanID {
		t.Errorf("expected span ID %q, got %q", span.Context.SpanID, spanID)
	}
}

func TestSpanStatus_Constants(t *testing.T) {
	// Verify constants are defined
	if SpanStatusUnset == "" {
		t.Error("SpanStatusUnset should be defined")
	}
	if SpanStatusOK == "" {
		t.Error("SpanStatusOK should be defined")
	}
	if SpanStatusError == "" {
		t.Error("SpanStatusError should be defined")
	}

	// Verify they're different
	statuses := []SpanStatus{SpanStatusUnset, SpanStatusOK, SpanStatusError}
	seen := make(map[SpanStatus]bool)
	for _, s := range statuses {
		if seen[s] {
			t.Errorf("duplicate status value: %q", s)
		}
		seen[s] = true
	}
}
