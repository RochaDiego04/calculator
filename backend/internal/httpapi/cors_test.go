package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	wrapped := CORS("http://localhost:5173")(next)

	t.Run("reflects configured origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
			t.Errorf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
		}
		if rec.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
		}
	})

	t.Run("handles preflight", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		rec := httptest.NewRecorder()
		wrapped.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
		}
	})
}
