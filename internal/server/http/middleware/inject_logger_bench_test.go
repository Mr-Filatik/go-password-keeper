package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/server/http/middleware"
)

// BenchmarkInjectLogger conducts benchmark testing of the InjectLogger middleware.
//
// The first allocation occurs when a logger is added to the current context,
// a new structure is created with the logger and a reference to the parent context.
//
// The second allocation occurs when r.WithContext(ctx) is called,
// when a shallow copy of the request is created, replacing the context in it.
//
// Results:
//   - Time per operation:        70.69 ns/op
//   - Bytes per operation:       368 B/op
//   - Allocations per operation: 2 allocs/op
//
// The results were obtained with the following settings:
//   - goos:   windows
//   - goarch: amd64
//   - cpu:    13th Gen Intel(R) Core(TM) i7-13700K
func BenchmarkInjectLogger(b *testing.B) {
	mockLogger := &MockLogger{}
	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		// Empty handler
	})

	mdlwr := middleware.InjectLogger(mockLogger)
	handler := mdlwr(next)

	req := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
	recorder := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		handler.ServeHTTP(recorder, req)
	}
}

// Results:
//   - Time per operation:        17.15 ns/op
//   - Bytes per operation:       48 B/op
//   - Allocations per operation: 1 allocs/op
//
// The results were obtained with the following settings:
//   - goos:   windows
//   - goarch: amd64
//   - cpu:    13th Gen Intel(R) Core(TM) i7-13700K
func BenchmarkInjectLoggerAction(b *testing.B) {
	mockLogger := &MockLogger{}
	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for range b.N {
		_ = middleware.InjectLoggerAction(ctx, mockLogger)
	}
}
