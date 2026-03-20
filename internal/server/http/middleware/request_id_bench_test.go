package middleware_test

// import (
// 	"context"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/mr-filatik/go-password-keeper/internal/platform/logging"
// 	"github.com/mr-filatik/go-password-keeper/internal/server/http/middleware"
// )

// // BenchmarkRequestID conducts benchmark testing of the RequestID middleware.
// //
// // Results:
// //   - Time per operation:        483.8 ns/op
// //   - Bytes per operation:       1424 B/op
// //   - Allocations per operation: 13 allocs/op
// //
// // The results were obtained with the following settings:
// //   - goos:   windows
// //   - goarch: amd64
// //   - cpu:    13th Gen Intel(R) Core(TM) i7-13700K
// func BenchmarkRequestID(b *testing.B) {
// 	mockLogger := &MockLogger{}
// 	next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
// 		// Empty handler
// 	})

// 	mdlwr := middleware.RequestID()
// 	handler := mdlwr(next)

// 	ctx := logging.ToContext(context.Background(), mockLogger)
// 	baseReq := httptest.NewRequest(http.MethodGet, "/test", http.NoBody).WithContext(ctx)
// 	// baseReq.Header.Set(middleware.HeaderRequestID, "test-uuid-123")

// 	recorder := httptest.NewRecorder()

// 	// recorder := httptest.NewRecorder()

// 	// baseReq := httptest.NewRequest(http.MethodGet, "/test", http.NoBody)
// 	// //baseReq.Header.Set(middleware.HeaderRequestID, "aaaaaaa")

// 	// req := baseReq.Clone(logging.ToContext(context.Background(), mockLogger))

// 	b.ResetTimer()

// 	for range b.N {
// 		baseReq.Header.Del(middleware.HeaderRequestID)

// 		handler.ServeHTTP(recorder, baseReq)
// 	}
// }

// // go test -bench . -benchmem -memprofile mem.out
