package middleware

import (
	"log/slog"
	"net/http"

	"github.com/RochaDiego04/calculator/backend/internal/httpapi/respond"
)

func Recovery(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered",
						"error", rec,
						"request_id", RequestIDFromContext(r.Context()))
					respond.Errors(w, http.StatusInternalServerError,
						respond.Error{Code: "INTERNAL", Message: "internal server error"})
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
