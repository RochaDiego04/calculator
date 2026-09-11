package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbeHealthAt(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		wantErr bool
	}{
		{"ok is healthy", http.StatusOK, false},
		{"service unavailable is unhealthy", http.StatusServiceUnavailable, true},
		{"server error is unhealthy", http.StatusInternalServerError, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path != "/api/v1/health" {
					t.Errorf("path = %q, want %q", r.URL.Path, "/api/v1/health")
				}
				w.WriteHeader(tc.status)
			}))
			defer srv.Close()

			err := probeHealthAt(srv.URL)

			if tc.wantErr && err == nil {
				t.Error("got nil, want an error")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("got %v, want nil", err)
			}
		})
	}
}

func TestProbeHealthAtUnreachable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	addr := srv.URL
	srv.Close()

	if err := probeHealthAt(addr); err == nil {
		t.Error("got nil for a closed server, want an error")
	}
}
