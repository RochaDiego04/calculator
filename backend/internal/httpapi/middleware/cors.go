package middleware

import (
	"net/http"
	"strings"
)

var (
	corsAllowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodOptions}
	corsAllowedHeaders = []string{"Content-Type"}
)

func CORS(origin string) Middleware {
	allowMethods := strings.Join(corsAllowedMethods, ", ")
	allowHeaders := strings.Join(corsAllowedHeaders, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", allowMethods)
			w.Header().Set("Access-Control-Allow-Headers", allowHeaders)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
