package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi/handler"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi/respond"
)

func newTestRouter() *http.ServeMux {
	return NewRouter(handler.New(1 << 20))
}

func assertSingleErrorCode(t *testing.T, rec *httptest.ResponseRecorder, wantCode string) {
	t.Helper()
	var got respond.ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error body: %v (body = %s)", err, rec.Body.String())
	}
	if len(got.Errors) != 1 || got.Errors[0].Code != wantCode {
		t.Errorf("errors = %v, want a single %s entry", got.Errors, wantCode)
	}
}

func TestRouterNotFound(t *testing.T) {
	rec := httptest.NewRecorder()
	newTestRouter().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/nope", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	assertSingleErrorCode(t, rec, "NOT_FOUND")
}

func TestRouterMethodNotAllowed(t *testing.T) {
	cases := []struct {
		name    string
		method  string
		path    string
		allowed string
	}{
		{"calculate wrong method", http.MethodGet, "/api/v1/calculate", http.MethodPost},
		{"health wrong method", http.MethodPost, "/api/v1/health", http.MethodGet},
		{"operations wrong method", http.MethodPost, "/api/v1/operations", http.MethodGet},
	}
	mux := newTestRouter()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))

			if rec.Code != http.StatusMethodNotAllowed {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusMethodNotAllowed, rec.Body.String())
			}
			if got := rec.Header().Get("Allow"); got != tc.allowed {
				t.Errorf("Allow header = %q, want %q", got, tc.allowed)
			}
			assertSingleErrorCode(t, rec, "METHOD_NOT_ALLOWED")
		})
	}
}

func TestRouterRoutesToHandlers(t *testing.T) {
	mux := newTestRouter()
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/health", nil))

	if rec.Code != http.StatusOK {
		t.Errorf("health status = %d, want %d", rec.Code, http.StatusOK)
	}
}
