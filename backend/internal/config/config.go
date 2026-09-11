package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	EnvDevelopment = "development"
	EnvProduction  = "production"
)

const (
	defaultPort            = "8080"
	defaultAppEnv          = EnvDevelopment
	defaultCORSOrigin      = "http://localhost:5173"
	defaultMaxBodyBytes    = 1 << 20
	defaultShutdownTimeout = 10 * time.Second
)

type Config struct {
	Port            string
	AppEnv          string
	CORSOrigin      string
	MaxBodyBytes    int64
	ShutdownTimeout time.Duration
}

// Load reads configuration from .env (when present) and the environment,
// returning an error rather than a default whenever a value is set but invalid.
func Load() (Config, error) {
	// A missing .env is normal outside local development.
	_ = godotenv.Load()

	cfg := Config{
		Port:            lookup("PORT", defaultPort),
		AppEnv:          lookup("APP_ENV", defaultAppEnv),
		CORSOrigin:      lookup("CORS_ORIGIN", defaultCORSOrigin),
		MaxBodyBytes:    defaultMaxBodyBytes,
		ShutdownTimeout: defaultShutdownTimeout,
	}

	if raw, ok := os.LookupEnv("MAX_BODY_BYTES"); ok {
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || n <= 0 {
			return Config{}, fmt.Errorf("MAX_BODY_BYTES must be a positive integer, got %q", raw)
		}
		cfg.MaxBodyBytes = n
	}

	if raw, ok := os.LookupEnv("SHUTDOWN_TIMEOUT"); ok {
		d, err := time.ParseDuration(raw)
		if err != nil || d <= 0 {
			return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT must be a positive duration such as 10s, got %q", raw)
		}
		cfg.ShutdownTimeout = d
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) validate() error {
	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("PORT must be a number between 1 and 65535, got %q", c.Port)
	}

	if c.AppEnv != EnvDevelopment && c.AppEnv != EnvProduction {
		return fmt.Errorf("APP_ENV must be %q or %q, got %q", EnvDevelopment, EnvProduction, c.AppEnv)
	}

	u, err := url.Parse(c.CORSOrigin)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("CORS_ORIGIN must be an absolute URL such as http://localhost:5173, got %q", c.CORSOrigin)
	}

	return nil
}

func lookup(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
