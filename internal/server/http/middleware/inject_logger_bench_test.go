package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mr-filatik/go-password-keeper/internal/server/http/middleware"
)

// BenchmarkInjectLogger conducts benchmark testing of the InjectLogger middleware.
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

	for range b.N {
		handler.ServeHTTP(recorder, req)
	}
}
