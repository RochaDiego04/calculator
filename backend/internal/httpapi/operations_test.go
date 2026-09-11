package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RochaDiego04/calculator/backend/internal/calc"
)

func TestOperations(t *testing.T) {
	h := NewHandler()
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
