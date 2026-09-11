package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("CORS_ORIGIN", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if cfg.Port != defaultPort {
		t.Errorf("Port = %q, want %q", cfg.Port, defaultPort)
	}
	if cfg.AppEnv != EnvDevelopment {
		t.Errorf("AppEnv = %q, want %q", cfg.AppEnv, EnvDevelopment)
	}
	if cfg.MaxBodyBytes != defaultMaxBodyBytes {
		t.Errorf("MaxBodyBytes = %d, want %d", cfg.MaxBodyBytes, defaultMaxBodyBytes)
	}
	if cfg.ShutdownTimeout != defaultShutdownTimeout {
		t.Errorf("ShutdownTimeout = %v, want %v", cfg.ShutdownTimeout, defaultShutdownTimeout)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("APP_ENV", EnvProduction)
	t.Setenv("CORS_ORIGIN", "https://calc.example.com")
	t.Setenv("MAX_BODY_BYTES", "2048")
	t.Setenv("SHUTDOWN_TIMEOUT", "30s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}
	if cfg.Port != "9090" {
		t.Errorf("Port = %q, want %q", cfg.Port, "9090")
	}
	if cfg.AppEnv != EnvProduction {
		t.Errorf("AppEnv = %q, want %q", cfg.AppEnv, EnvProduction)
	}
	if cfg.CORSOrigin != "https://calc.example.com" {
		t.Errorf("CORSOrigin = %q, want %q", cfg.CORSOrigin, "https://calc.example.com")
	}
	if cfg.MaxBodyBytes != 2048 {
		t.Errorf("MaxBodyBytes = %d, want 2048", cfg.MaxBodyBytes)
	}
	if cfg.ShutdownTimeout != 30*time.Second {
		t.Errorf("ShutdownTimeout = %v, want %v", cfg.ShutdownTimeout, 30*time.Second)
	}
}

func TestLoadRejectsInvalidValues(t *testing.T) {
	cases := []struct {
		name  string
		key   string
		value string
	}{
		{"non numeric port", "PORT", "abc"},
		{"port out of range", "PORT", "70000"},
		{"unknown app env", "APP_ENV", "staging"},
		{"relative cors origin", "CORS_ORIGIN", "localhost:5173"},
		{"negative max body bytes", "MAX_BODY_BYTES", "-1"},
		{"non numeric max body bytes", "MAX_BODY_BYTES", "big"},
		{"unparseable shutdown timeout", "SHUTDOWN_TIMEOUT", "soon"},
		{"zero shutdown timeout", "SHUTDOWN_TIMEOUT", "0s"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv(tc.key, tc.value)

			if _, err := Load(); err == nil {
				t.Fatalf("Load() with %s=%q returned nil error, want a validation failure", tc.key, tc.value)
			}
		})
	}
}
