package httpapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecovery(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	panicky := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})
	wrapped := Recovery(logger)(panicky)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
	var got errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error body: %v", err)
	}
	if len(got.Errors) != 1 || got.Errors[0].Code != "INTERNAL" {
		t.Errorf("errors = %v, want a single INTERNAL entry", got.Errors)
	}
}
