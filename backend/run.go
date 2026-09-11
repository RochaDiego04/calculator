package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RochaDiego04/calculator/backend/internal/config"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi/handler"
	"github.com/RochaDiego04/calculator/backend/internal/httpapi/middleware"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 15 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

func run() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	logger := newLogger(cfg.AppEnv)
	srv := newServer(cfg, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("listening", "addr", srv.Addr, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down", "timeout", cfg.ShutdownTimeout)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("stopped")
}

func newLogger(appEnv string) *slog.Logger {
	if appEnv == config.EnvProduction {
		return slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}

func newHandler(cfg config.Config, logger *slog.Logger) http.Handler {
	mux := httpapi.NewRouter(handler.New(cfg.MaxBodyBytes))
	return middleware.Chain(mux,
		middleware.Recovery(logger),
		middleware.RequestID,
		middleware.Logging(logger),
		middleware.CORS(cfg.CORSOrigin),
	)
}

func newServer(cfg config.Config, logger *slog.Logger) *http.Server {
	return &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           newHandler(cfg, logger),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}
