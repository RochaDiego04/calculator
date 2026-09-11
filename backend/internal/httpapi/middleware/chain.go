package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

// Chain applies middleware in the order given: the first one listed is
// outermost, seeing the request before anything else.
func Chain(h http.Handler, mw ...Middleware) http.Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h = mw[i](h)
	}
	return h
}
