package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	var gotFromContext string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotFromContext = requestIDFromContext(r.Context())
	})
	wrapped := RequestID(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	header := rec.Header().Get("X-Request-Id")
	if header == "" {
		t.Fatal("X-Request-Id header not set")
	}
	if gotFromContext != header {
		t.Errorf("context request id = %q, want %q (header value)", gotFromContext, header)
	}
}

func TestRequestIDUnique(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	wrapped := RequestID(next)

	seen := map[string]bool{}
	for range 10 {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		id := rec.Header().Get("X-Request-Id")
		if seen[id] {
			t.Fatalf("duplicate request id: %s", id)
		}
		seen[id] = true
	}
}
