package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RochaDiego04/calculator/backend/internal/calc"
)

func TestHealth(t *testing.T) {
	h := New(testMaxBodyBytes)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	h.Health(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got healthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if got.Status != "ok" {
		t.Errorf("Status = %q, want %q", got.Status, "ok")
	}
}

func TestOperations(t *testing.T) {
	h := New(testMaxBodyBytes)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/operations", nil)
	rec := httptest.NewRecorder()
	h.Operations(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got []calc.Descriptor
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if len(got) != len(calc.Operations()) {
		t.Errorf("returned %d operations, want %d", len(got), len(calc.Operations()))
	}
}

func TestNotFound(t *testing.T) {
	h := New(testMaxBodyBytes)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/nope", nil)
	rec := httptest.NewRecorder()
	h.NotFound(rec, req)

	assertErrorCode(t, rec, http.StatusNotFound, "NOT_FOUND")
}

func TestMethodNotAllowed(t *testing.T) {
	h := New(testMaxBodyBytes)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()
	h.MethodNotAllowed(http.MethodPost)(rec, req)

	assertErrorCode(t, rec, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED")
	if got := rec.Header().Get("Allow"); got != http.MethodPost {
		t.Errorf("Allow header = %q, want %q", got, http.MethodPost)
	}
}
