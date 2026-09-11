package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterNotFound(t *testing.T) {
	mux := NewRouter(NewHandler())
	req := httptest.NewRequest(http.MethodGet, "/api/v1/nope", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	var got errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error body: %v (body = %s)", err, rec.Body.String())
	}
	if len(got.Errors) != 1 || got.Errors[0].Code != "NOT_FOUND" {
		t.Errorf("errors = %v, want a single NOT_FOUND entry", got.Errors)
	}
}

func TestRouterMethodNotAllowed(t *testing.T) {
	cases := []struct {
		name    string
		method  string
		path    string
		allowed string
	}{
		{"calculate wrong method", http.MethodGet, "/api/v1/calculate", "POST"},
		{"health wrong method", http.MethodPost, "/api/v1/health", "GET"},
		{"operations wrong method", http.MethodPost, "/api/v1/operations", "GET"},
	}
	mux := NewRouter(NewHandler())
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusMethodNotAllowed, rec.Body.String())
			}
			if got := rec.Header().Get("Allow"); got != tc.allowed {
				t.Errorf("Allow header = %q, want %q", got, tc.allowed)
			}
			var got errorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("unmarshal error body: %v (body = %s)", err, rec.Body.String())
			}
			if len(got.Errors) != 1 || got.Errors[0].Code != "METHOD_NOT_ALLOWED" {
				t.Errorf("errors = %v, want a single METHOD_NOT_ALLOWED entry", got.Errors)
			}
		})
	}
}
