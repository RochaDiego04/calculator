package httpapi

import (
	"net/http"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi/handler"
)

const apiPrefix = "/api/v1"

func NewRouter(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Each route is registered twice: once for its method, once without one so
	// a wrong method gets a 405 instead of falling through to the 404 catch-all.
	mux.HandleFunc("POST "+apiPrefix+"/calculate", h.Calculate)
	mux.HandleFunc(apiPrefix+"/calculate", h.MethodNotAllowed(http.MethodPost))

	mux.HandleFunc("GET "+apiPrefix+"/health", h.Health)
	mux.HandleFunc(apiPrefix+"/health", h.MethodNotAllowed(http.MethodGet))

	mux.HandleFunc("GET "+apiPrefix+"/operations", h.Operations)
	mux.HandleFunc(apiPrefix+"/operations", h.MethodNotAllowed(http.MethodGet))

	mux.HandleFunc("/", h.NotFound)
	return mux
}
