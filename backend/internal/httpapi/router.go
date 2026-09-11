package httpapi

import "net/http"

func NewRouter(h *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/calculate", h.Calculate)
	mux.HandleFunc("/api/v1/calculate", h.methodNotAllowed("POST"))

	mux.HandleFunc("GET /api/v1/health", h.Health)
	mux.HandleFunc("/api/v1/health", h.methodNotAllowed("GET"))

	mux.HandleFunc("GET /api/v1/operations", h.Operations)
	mux.HandleFunc("/api/v1/operations", h.methodNotAllowed("GET"))

	mux.HandleFunc("/", h.NotFound)
	return mux
}
