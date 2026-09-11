package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RochaDiego04/calculator/backend/internal/config"
)

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestNewLogger(t *testing.T) {
	cases := []struct {
		name        string
		appEnv      string
		wantHandler string
	}{
		{"production uses json handler", config.EnvProduction, "*slog.JSONHandler"},
		{"development uses text handler", config.EnvDevelopment, "*slog.TextHandler"},
		{"unknown env defaults to text handler", "staging", "*slog.TextHandler"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			logger := newLogger(tc.appEnv)
			if got := fmt.Sprintf("%T", logger.Handler()); got != tc.wantHandler {
				t.Errorf("newLogger(%q) handler = %s, want %s", tc.appEnv, got, tc.wantHandler)
			}
		})
	}
}

func testConfig() config.Config {
	return config.Config{
		Port:         "9999",
		AppEnv:       config.EnvDevelopment,
		CORSOrigin:   "http://localhost:5173",
		MaxBodyBytes: 1 << 20,
	}
}

func TestNewServer(t *testing.T) {
	srv := newServer(testConfig(), discardLogger())

	if srv.Addr != ":9999" {
		t.Errorf("Addr = %q, want %q", srv.Addr, ":9999")
	}
	if srv.ReadHeaderTimeout != readHeaderTimeout {
		t.Errorf("ReadHeaderTimeout = %v, want %v", srv.ReadHeaderTimeout, readHeaderTimeout)
	}
	if srv.ReadTimeout != readTimeout {
		t.Errorf("ReadTimeout = %v, want %v", srv.ReadTimeout, readTimeout)
	}
	if srv.WriteTimeout != writeTimeout {
		t.Errorf("WriteTimeout = %v, want %v", srv.WriteTimeout, writeTimeout)
	}
	if srv.IdleTimeout != idleTimeout {
		t.Errorf("IdleTimeout = %v, want %v", srv.IdleTimeout, idleTimeout)
	}
	if srv.Handler == nil {
		t.Fatal("Handler is nil")
	}
}

func TestNewHandlerRoutesToRealEndpoints(t *testing.T) {
	h := newHandler(testConfig(), discardLogger())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Errorf("Access-Control-Allow-Origin = %q, want %q (CORS middleware should be wired)", got, "http://localhost:5173")
	}
}
