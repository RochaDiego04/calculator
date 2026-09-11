package respond

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	JSON(rec, http.StatusCreated, map[string]string{"key": "value"})

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}
	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got["key"] != "value" {
		t.Errorf("body = %v, want key=value", got)
	}
}

func TestErrors(t *testing.T) {
	rec := httptest.NewRecorder()
	Errors(rec, http.StatusUnprocessableEntity, Error{Code: "TEST", Message: "test message", Field: "x"})

	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
	var got ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Errors) != 1 || got.Errors[0].Code != "TEST" {
		t.Errorf("errors = %v, want a single TEST entry", got.Errors)
	}
}

func TestErrorsMultiple(t *testing.T) {
	rec := httptest.NewRecorder()
	Errors(rec, http.StatusBadRequest, Error{Code: "A"}, Error{Code: "B"})

	var got ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(got.Errors) != 2 {
		t.Fatalf("errors = %v, want 2 entries", got.Errors)
	}
}
