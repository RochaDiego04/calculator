// Package e2e drives the full router + middleware + handler stack the way a
// real client does, asserting the whole API error contract end to end.
package e2e

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi/handler"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi/middleware"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi/respond"
)

func fullStack(maxBodyBytes int64) http.Handler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	mux := httpapi.NewRouter(handler.New(maxBodyBytes))
	return middleware.Chain(mux,
		middleware.Recovery(logger),
		middleware.RequestID,
		middleware.Logging(logger),
		middleware.CORS("http://localhost:5173"),
	)
}

func TestErrorContractEndToEnd(t *testing.T) {
	stack := fullStack(1 << 20)

	cases := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
		wantField  string
	}{
		{"malformed json", http.MethodPost, "/api/v1/calculate", `not json`,
			http.StatusBadRequest, "INVALID_JSON", ""},
		{"unknown field", http.MethodPost, "/api/v1/calculate", `{"operation":"add","a":1,"b":2,"x":9}`,
			http.StatusBadRequest, "INVALID_JSON", ""},
		{"missing a", http.MethodPost, "/api/v1/calculate", `{"operation":"add","b":2}`,
			http.StatusUnprocessableEntity, "MISSING_OPERAND", "a"},
		{"missing b for binary op", http.MethodPost, "/api/v1/calculate", `{"operation":"divide","a":12}`,
			http.StatusUnprocessableEntity, "MISSING_OPERAND", "b"},
		{"unknown operation", http.MethodPost, "/api/v1/calculate", `{"operation":"cos","a":1,"b":2}`,
			http.StatusUnprocessableEntity, "UNKNOWN_OPERATION", "operation"},
		{"division by zero", http.MethodPost, "/api/v1/calculate", `{"operation":"divide","a":12,"b":0}`,
			http.StatusUnprocessableEntity, "DIVISION_BY_ZERO", "b"},
		{"negative square root", http.MethodPost, "/api/v1/calculate", `{"operation":"sqrt","a":-4}`,
			http.StatusUnprocessableEntity, "NEGATIVE_SQUARE_ROOT", "a"},
		{"result not representable", http.MethodPost, "/api/v1/calculate", `{"operation":"power","a":1e308,"b":2}`,
			http.StatusUnprocessableEntity, "RESULT_NOT_REPRESENTABLE", ""},
		{"unknown path", http.MethodGet, "/api/v1/nope", "",
			http.StatusNotFound, "NOT_FOUND", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var body io.Reader
			if tc.body != "" {
				body = strings.NewReader(tc.body)
			}
			req := httptest.NewRequest(tc.method, tc.path, body)
			rec := httptest.NewRecorder()
			stack.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body = %s", rec.Code, tc.wantStatus, rec.Body.String())
			}

			var got respond.ErrorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("unmarshal error body: %v (body = %s)", err, rec.Body.String())
			}
			if len(got.Errors) != 1 {
				t.Fatalf("errors = %v, want exactly one", got.Errors)
			}
			// Assert on the code and field only, never on message text
			if got.Errors[0].Code != tc.wantCode {
				t.Errorf("code = %q, want %q", got.Errors[0].Code, tc.wantCode)
			}
			if got.Errors[0].Field != tc.wantField {
				t.Errorf("field = %q, want %q", got.Errors[0].Field, tc.wantField)
			}
		})
	}
}

func TestErrorContractMethodNotAllowed(t *testing.T) {
	stack := fullStack(1 << 20)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
	if got := rec.Header().Get("Allow"); got != http.MethodPost {
		t.Errorf("Allow header = %q, want %q", got, http.MethodPost)
	}
}

func TestErrorContractRequestTooLarge(t *testing.T) {
	stack := fullStack(16)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate",
		strings.NewReader(`{"operation":"add","a":1,"b":2}`))
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusRequestEntityTooLarge, rec.Body.String())
	}
	var got respond.ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error body: %v", err)
	}
	if len(got.Errors) != 1 || got.Errors[0].Code != "REQUEST_TOO_LARGE" {
		t.Errorf("errors = %v, want a single REQUEST_TOO_LARGE entry", got.Errors)
	}
}

func TestErrorContractSuccess(t *testing.T) {
	stack := fullStack(1 << 20)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate",
		strings.NewReader(`{"operation":"divide","a":12,"b":3}`))
	rec := httptest.NewRecorder()
	stack.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}
