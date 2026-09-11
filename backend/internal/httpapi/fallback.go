package httpapi

import (
	"net/http"
	"strings"
)

func (h *Handler) NotFound(w http.ResponseWriter, r *http.Request) {
	writeErrors(w, http.StatusNotFound, apiError{Code: "NOT_FOUND", Message: "resource not found"})
}

func (h *Handler) methodNotAllowed(allowed ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", strings.Join(allowed, ", "))
		writeErrors(w, http.StatusMethodNotAllowed,
			apiError{Code: "METHOD_NOT_ALLOWED", Message: "method not allowed for this resource"})
	}
}
