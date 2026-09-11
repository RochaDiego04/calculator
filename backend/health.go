package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RochaDiego04/calculator/backend/internal/config"
)

const healthProbeTimeout = 2 * time.Second

// probeHealth checks the server running inside this same container. The
// distroless image has no shell and no curl, so the binary is its own probe.
func probeHealth() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	return probeHealthAt("http://localhost:" + cfg.Port)
}

func probeHealthAt(baseURL string) error {
	client := &http.Client{Timeout: healthProbeTimeout}

	resp, err := client.Get(baseURL + "/api/v1/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %s", resp.Status)
	}
	return nil
}
