package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mr-filatik/go-password-keeper/internal/server/http/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func run(
	t *testing.T,
	mdlwrFn func(http.Handler) http.Handler,
	req *http.Request,
) (string, string, string) {
	t.Helper()

	var gotCtxVal, gotReqHdr string

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v, ok := r.Context().Value(middleware.CtxKeyXRequestID).(string); ok {
			gotCtxVal = v
		}

		gotReqHdr = r.Header.Get(middleware.HeaderRequestID)

		w.WriteHeader(http.StatusOK)
	})

	rr := httptest.NewRecorder()
	mdlwrFn(next).ServeHTTP(rr, req)

	return gotCtxVal, gotReqHdr, rr.Header().Get(middleware.HeaderRequestID)
}

func TestRequestID_PreservesProvidedID(t *testing.T) {
	t.Parallel()

	mdlwr := middleware.RequestID()

	const wantID = "provided-id-123"

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
	req.Header.Set(middleware.HeaderRequestID, wantID)

	ctxVal, seenReqHdr, respHdr := run(t, mdlwr, req)

	require.Equalf(t, wantID, seenReqHdr, "next handler saw wrong request header")

	assert.Equalf(t, wantID, respHdr, "response header mismatch")

	if ctxVal != wantID {
		t.Fatalf("context value = %q, want %q", ctxVal, wantID)
	}
}

func TestRequestID_GeneratesWhenMissing(t *testing.T) {
	t.Parallel()

	mdlwr := middleware.RequestID()

	req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)

	ctxValue, reqHeader, respHeader := run(t, mdlwr, req)

	require.NotEmptyf(t, reqHeader,
		"next handler saw empty %s header; want generated UUID", middleware.HeaderRequestID)

	_, err := uuid.Parse(reqHeader)
	require.NoErrorf(t, err, "next handler saw non-UUID %q", reqHeader)

	assert.Equalf(t, reqHeader, respHeader, "response header mismatch")

	assert.Equalf(t, reqHeader, ctxValue, "context value mismatch")
}
