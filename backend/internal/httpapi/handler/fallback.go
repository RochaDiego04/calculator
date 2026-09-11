package handler

import (
	"net/http"
	"strings"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi/respond"
)

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	respond.Errors(w, http.StatusNotFound,
		respond.Error{Code: "NOT_FOUND", Message: "resource not found"})
}

// MethodNotAllowed returns a handler for a path whose method did not match,
// since ServeMux's built-in 405 never fires once a catch-all route is registered.
func (h *Handler) MethodNotAllowed(allowed ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		respond.Errors(w, http.StatusMethodNotAllowed,
			respond.Error{Code: "METHOD_NOT_ALLOWED", Message: "method not allowed for this resource"})
	}
}
