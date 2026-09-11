package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	h := httpapi.NewHandler()
	mux := httpapi.NewRouter(h)

	handler := httpapi.Chain(mux,
		httpapi.Recovery(logger),
		httpapi.RequestID,
		httpapi.Logging(logger),
		httpapi.CORS("http://localhost:5173"),
	)

	logger.Info("listening", "addr", ":8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
