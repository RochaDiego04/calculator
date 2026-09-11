package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doCalculate(t *testing.T, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := NewHandler()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h.Calculate(rec, req)
	return rec
}

func TestCalculate(t *testing.T) {
	t.Run("successful divide", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"divide","a":12,"b":3}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
		}
		var got calculateResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Fatalf("unmarshal response: %v", err)
		}
		if got.Result != 4 {
			t.Errorf("Result = %v, want 4", got.Result)
		}
	})

	t.Run("successful sqrt omits b", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"sqrt","a":9}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
		}
		if strings.Contains(rec.Body.String(), `"b"`) {
			t.Errorf("response contains b field for a unary op: %s", rec.Body.String())
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		rec := doCalculate(t, `not json`)
		assertErrorCode(t, rec, http.StatusBadRequest, "INVALID_JSON")
	})

	t.Run("unknown field", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"add","a":1,"b":2,"x":9}`)
		assertErrorCode(t, rec, http.StatusBadRequest, "INVALID_JSON")
	})

	t.Run("missing a", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"add","b":2}`)
		assertErrorCode(t, rec, http.StatusUnprocessableEntity, "MISSING_OPERAND")
	})

	t.Run("missing b for binary op", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"divide","a":12}`)
		assertErrorCode(t, rec, http.StatusUnprocessableEntity, "MISSING_OPERAND")
	})

	t.Run("unknown operation", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"cos","a":1,"b":2}`)
		assertErrorCode(t, rec, http.StatusUnprocessableEntity, "UNKNOWN_OPERATION")
	})

	t.Run("division by zero", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"divide","a":12,"b":0}`)
		assertErrorCode(t, rec, http.StatusUnprocessableEntity, "DIVISION_BY_ZERO")
	})

	t.Run("negative square root", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"sqrt","a":-4}`)
		assertErrorCode(t, rec, http.StatusUnprocessableEntity, "NEGATIVE_SQUARE_ROOT")
	})

	t.Run("result not representable", func(t *testing.T) {
		rec := doCalculate(t, `{"operation":"power","a":1e308,"b":2}`)
		assertErrorCode(t, rec, http.StatusUnprocessableEntity, "RESULT_NOT_REPRESENTABLE")
	})
}

func assertErrorCode(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantCode string) {
	t.Helper()
	if rec.Code != wantStatus {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, wantStatus, rec.Body.String())
	}
	var got errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal error body: %v", err)
	}
	if len(got.Errors) != 1 {
		t.Fatalf("errors = %v, want exactly one", got.Errors)
	}
	if got.Errors[0].Code != wantCode {
		t.Errorf("code = %q, want %q", got.Errors[0].Code, wantCode)
	}
}
