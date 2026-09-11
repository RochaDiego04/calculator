package middleware

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi/respond"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestRecovery(t *testing.T) {
	panicky := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})
	wrapped := Recovery(discardLogger())(panicky)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	var got respond.ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error body: %v", err)
	}
	if len(got.Errors) != 1 || got.Errors[0].Code != "INTERNAL" {
		t.Errorf("errors = %v, want a single INTERNAL entry", got.Errors)
	}
}

func TestRequestID(t *testing.T) {
	var gotFromContext string
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotFromContext = RequestIDFromContext(r.Context())
	})
	wrapped := RequestID(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	header := rec.Header().Get(HeaderRequestID)
	if header == "" {
		t.Fatalf("%s header not set", HeaderRequestID)
	}
	if gotFromContext != header {
		t.Errorf("context request id = %q, want %q (header value)", gotFromContext, header)
	}
}

func TestRequestIDUnique(t *testing.T) {
	wrapped := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	seen := map[string]bool{}
	for range 10 {
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

		id := rec.Header().Get(HeaderRequestID)
		if seen[id] {
			t.Fatalf("duplicate request id: %s", id)
		}
		seen[id] = true
	}
}

func TestCORS(t *testing.T) {
	const origin = "http://localhost:5173"
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	wrapped := CORS(origin)(next)

	t.Run("reflects configured origin", func(t *testing.T) {
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != origin {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, origin)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("handles preflight", func(t *testing.T) {
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, httptest.NewRequest(http.MethodOptions, "/", nil))

		if rec.Code != http.StatusNoContent {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
		}
	})
}

func TestChainOrder(t *testing.T) {
	var order []string
	record := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}
	final := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	})

	Chain(final, record("outer"), record("inner")).
		ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))

	want := []string{"outer", "inner", "handler"}
	for i, w := range want {
		if order[i] != w {
			t.Fatalf("order = %v, want %v", order, want)
		}
	}
}

func TestLoggingCapturesStatus(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	wrapped := Logging(discardLogger())(next)

	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	if rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusTeapot)
	}
}
