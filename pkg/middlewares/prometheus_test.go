package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMetrics(t *testing.T) {
	metrics := NewMetrics("data_pipeline")
	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.reg)
}

func TestMetricsHandler(t *testing.T) {
	metrics := NewMetrics("data_pipeline")
	handler := metrics.Handler()
	assert.NotNil(t, handler)

	// Test that handler is actually an http.Handler
	req := httptest.NewRequest("GET", "/metrics", http.NoBody)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	// Should return 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSourceMetrics(t *testing.T) {
	metrics := NewMetrics("data_pipeline")
	streamID := "stream-1"
	sourceID := "kafka-1"

	// Test all source metrics methods - should not panic
	assert.NotPanics(t, func() {
		metrics.OnSourceInput(streamID, sourceID)
		metrics.OnSourceInputSuccess(streamID, sourceID)
		metrics.OnSourceOutput(streamID, sourceID)
		metrics.OnSourceOutputSuccess(streamID, sourceID)
	})
}

func TestTransformMetrics(t *testing.T) {
	metrics := NewMetrics("data_pipeline")
	streamID := "stream-1"

	// Test all transform metrics methods - should not panic
	assert.NotPanics(t, func() {
		metrics.OnTransformInput(streamID)
		metrics.OnTransformInputSuccess(streamID)
		metrics.OnTransformConvert(streamID)
		metrics.OnTransformConvertSuccess(streamID)
		metrics.OnTransformOutput(streamID)
		metrics.OnTransformOutputSuccess(streamID)
	})
}

func TestSinkMetrics(t *testing.T) {
	metrics := NewMetrics("data_pipeline")
	streamID := "stream-1"
	sinkID := "http-1"

	// Test all sink metrics methods - should not panic
	assert.NotPanics(t, func() {
		metrics.OnSinkInput(streamID, sinkID)
		metrics.OnSinkInputSuccess(streamID, sinkID)
		metrics.OnSinkOutput(streamID, sinkID)
		metrics.OnSinkOutputSuccess(streamID, sinkID)
	})
}

func BenchmarkNewMetrics(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewMetrics("data_pipeline")
	}
}

func BenchmarkOnSourceInput(b *testing.B) {
	metrics := NewMetrics("data_pipeline")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.OnSourceInput("stream-1", "kafka-1")
	}
}

func BenchmarkOnSourceInputSuccess(b *testing.B) {
	metrics := NewMetrics("data_pipeline")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.OnSourceInputSuccess("stream-1", "kafka-1")
	}
}

func BenchmarkOnTransformInput(b *testing.B) {
	metrics := NewMetrics("data_pipeline")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.OnTransformInput("stream-1")
	}
}

func BenchmarkOnTransformConvert(b *testing.B) {
	metrics := NewMetrics("data_pipeline")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.OnTransformConvert("stream-1")
	}
}

func BenchmarkOnSinkInput(b *testing.B) {
	metrics := NewMetrics("data_pipeline")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.OnSinkInput("stream-1", "http-1")
	}
}

func BenchmarkOnSinkOutput(b *testing.B) {
	metrics := NewMetrics("data_pipeline")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.OnSinkOutput("stream-1", "http-1")
	}
}
