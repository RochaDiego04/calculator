package httpapi

import (
	"log/slog"
	"net/http"
)

func Recovery(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"error", rec,
						"request_id", requestIDFromContext(r.Context()))
					writeErrors(w, http.StatusInternalServerError,
						apiError{Code: "INTERNAL", Message: "internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
